package systemRbac

import (
	"context"
	"fmt"
	"time"

	"shack/internal/global"
	"shack/internal/model/systemRbac"
	req "shack/internal/model/systemRbac/request"
	res "shack/internal/model/systemRbac/response"
	"shack/internal/utils"

	biz_err "shack/internal/error"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	// Redis key 前缀
	RedisKeyUserOnline     = "chat:online:"      // 用户在线状态
	RedisKeyUnreadCount    = "chat:unread:"      // 未读消息数
	RedisKeyConversation   = "chat:conversation:" // 会话信息
	RedisKeyBlacklist      = "chat:blacklist:"    // 黑名单缓存
	RedisKeyRecentContacts = "chat:recent:"       // 最近联系人

	// 过期时间
	RedisExpireOnline     = 300  // 5分钟
	RedisExpireUnread     = 600  // 10分钟
	RedisExpireConversation = 3600 // 1小时
	RedisExpireRecent     = 1800 // 30分钟
)

type ChatService struct{}

// SendChatMessage 发送私聊消息
// Auth: shack
// Time: 2026年04月03日
func (s *ChatService) SendChatMessage(
	c *gin.Context,
	r req.SendChatMessageReq,
) (rs res.SendChatMessageRes, err error) {
	// 获取当前用户
	currentUserID := utils.GetUserID(c)
	if currentUserID == 0 {
		return rs, biz_err.New(biz_err.AUTH_ERROR, "用户未登录")
	}

	// 参数验证
	if r.ReceiverId == 0 {
		return rs, biz_err.New(biz_err.PARAM_MISSING, "接收者ID不能为空")
	}
	if r.Content == "" {
		return rs, biz_err.New(biz_err.PARAM_MISSING, "消息内容不能为空")
	}
	if r.MsgType < 1 || r.MsgType > 3 {
		return rs, biz_err.New(biz_err.PARAM_ERROR, "消息类型不正确")
	}

	// 不能发送消息给自己
	if int64(currentUserID) == r.ReceiverId {
		return rs, biz_err.New(biz_err.PARAM_ERROR, "不能发送消息给自己")
	}

	// 检查是否被拉黑
	isBlocked, err := s.checkUserBlocked(context.Background(), int64(currentUserID), r.ReceiverId)
	if err != nil {
		return rs, err
	}
	if isBlocked {
		return rs, biz_err.New(biz_err.FORBIDDEN_OPERATION, "对方已将您拉黑")
	}

	// 查找接收者是否存在
	var receiver systemRbac.User
	if err := global.GVA_DB.First(&receiver, r.ReceiverId).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return rs, biz_err.New(biz_err.USER_NOT_FOUND, "接收者不存在")
		}
		return rs, biz_err.New(biz_err.DB_ERROR, "查询接收者失败")
	}

	// 开始事务
	tx := global.GVA_DB.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 获取或创建会话
	conversationID, err := s.getOrCreateConversation(tx, int64(currentUserID), r.ReceiverId)
	if err != nil {
		return rs, err
	}

	// 创建消息
	now := time.Now()
	message := systemRbac.ChatMessage{
		ConversationID: conversationID,
		SenderID:       int64(currentUserID),
		ReceiverID:     r.ReceiverId,
		Content:        r.Content,
		MsgType:        r.MsgType,
		Status:         1, // 发送成功
		IsRead:         false,
	}

	if err := tx.Create(&message).Error; err != nil {
		global.GVA_LOG.Error("创建聊天消息失败", zap.Error(err))
		return rs, biz_err.New(biz_err.DB_ERROR, "发送消息失败")
	}

	// 更新会话信息
	if err := s.updateConversationAfterSend(tx, conversationID, int64(currentUserID), r.ReceiverId, r.Content, &now); err != nil {
		return rs, err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR, "发送消息失败")
	}

	// 清除缓存
	s.clearChatCache(int64(currentUserID), r.ReceiverId)

	// TODO: 这里可以通过WebSocket推送实时消息

	// 获取发送者信息
	var sender systemRbac.User
	global.GVA_DB.First(&sender, currentUserID)

	rs = res.SendChatMessageRes{
		Id:         int64(message.ID),
		SenderId:   int64(currentUserID),
		SenderName: sender.NickName,
		ReceiverId: r.ReceiverId,
		Content:    r.Content,
		MsgType:    r.MsgType,
		IsRead:     0,
		SendTime:   now,
	}

	return rs, nil
}

// GetChatHistory 获取聊天记录
// Auth: shack
// Time: 2026年04月03日
func (s *ChatService) GetChatHistory(
	c *gin.Context,
	r req.GetChatHistoryReq,
) (rs res.GetChatHistoryRes, err error) {
	currentUserID := utils.GetUserID(c)
	if currentUserID == 0 {
		return rs, biz_err.New(biz_err.AUTH_ERROR, "用户未登录")
	}

	// 分页参数
	page := r.Page
	if page < 1 {
		page = 1
	}
	pageSize := r.PageSize
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	// 查询条件：当前用户与目标用户的聊天记录
	db := global.GVA_DB.Model(&systemRbac.ChatMessage{}).Where(
		"(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
		currentUserID, r.TargetId, r.TargetId, currentUserID,
	)

	// 查询总数
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR, "查询聊天记录失败")
	}

	// 分页查询
	var messages []systemRbac.ChatMessage
	offset := (page - 1) * pageSize
	if err := db.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&messages).Error; err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR, "查询聊天记录失败")
	}

	// 反转消息顺序（最新的在后面）
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	// 构建响应
	list := make([]res.GetChatHistoryResList, 0, len(messages))
	userIDs := make(map[int64]bool)

	// 收集所有涉及的用ID
	for _, msg := range messages {
		userIDs[msg.SenderID] = true
	}

	// 批量查询用户信息
	var users []systemRbac.User
	global.GVA_DB.Where("id IN ?", userIDs).Find(&users)
	userMap := make(map[int64]systemRbac.User)
	for _, user := range users {
		userMap[int64(user.ID)] = user
	}

	// 组装数据
	for _, msg := range messages {
		sender, ok := userMap[msg.SenderID]
		senderName := "未知用户"
		senderAvatar := ""
		if ok {
			senderName = sender.NickName
			senderAvatar = sender.HeaderImg
		}

		isRead := 0
		if msg.IsRead {
			isRead = 1
		}

		list = append(list, res.GetChatHistoryResList{
			Id:           int64(msg.ID),
			SenderId:     msg.SenderID,
			SenderName:   senderName,
			SenderAvatar: senderAvatar,
			ReceiverId:   msg.ReceiverID,
			Content:      msg.Content,
			MsgType:      msg.MsgType,
			IsRead:       isRead,
			SendTime:     msg.CreatedAt,
		})
	}

	rs = res.GetChatHistoryRes{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}

	return rs, nil
}

// GetRecentContacts 获取最近联系人
// Auth: shack
// Time: 2026年04月03日
func (s *ChatService) GetRecentContacts(
	c *gin.Context,
) (rs res.GetRecentContactsRes, err error) {
	currentUserID := utils.GetUserID(c)
	if currentUserID == 0 {
		return rs, biz_err.New(biz_err.AUTH_ERROR, "用户未登录")
	}

	// 尝试从Redis获取缓存
	cacheKey := fmt.Sprintf("%s%d", RedisKeyRecentContacts, currentUserID)
	if global.GVA_REDIS != nil {
		cached, err := global.GVA_REDIS.Get(context.Background(), cacheKey).Result()
		if err == nil && cached != "" {
			// TODO: 反序列化缓存数据
			_ = cached
		}
	}

	// 查询最近联系人的最新消息
	type ContactMsg struct {
		ID           int64
		SenderID     int64
		ReceiverID   int64
		Content      string
		MsgType      int
		IsRead       bool
		CreatedAt    time.Time
	}

	var messages []ContactMsg
	err = global.GVA_DB.Raw(`
		SELECT id, sender_id, receiver_id, content, msg_type, is_read, created_at
		FROM chat_messages
		WHERE sender_id = ? OR receiver_id = ?
		ORDER BY created_at DESC
		LIMIT 50
	`, currentUserID, currentUserID).Scan(&messages).Error

	if err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR, "查询最近联系人失败")
	}

	// 收集用户ID
	userIDSet := make(map[int64]bool)
	for _, msg := range messages {
		if msg.SenderID != int64(currentUserID) {
			userIDSet[msg.SenderID] = true
		}
		if msg.ReceiverID != int64(currentUserID) {
			userIDSet[msg.ReceiverID] = true
		}
	}

	// 查询用户信息
	var users []systemRbac.User
	if len(userIDSet) > 0 {
		userIDs := make([]int64, 0, len(userIDSet))
		for id := range userIDSet {
			userIDs = append(userIDs, id)
		}
		global.GVA_DB.Where("id IN ?", userIDs).Find(&users)
	}

	userMap := make(map[int64]systemRbac.User)
	for _, user := range users {
		userMap[int64(user.ID)] = user
	}

	// 构建响应列表
	contactMap := make(map[int64]res.GetRecentContactsResList)
	for _, msg := range messages {
		var targetID int64
		if msg.SenderID == int64(currentUserID) {
			targetID = msg.ReceiverID
		} else {
			targetID = msg.SenderID
		}

		// 如果已经有这个联系人的记录，跳过（保留最新的）
		if _, exists := contactMap[targetID]; exists {
			continue
		}

		isRead := 0
		if msg.IsRead {
			isRead = 1
		}

		contactMap[targetID] = res.GetRecentContactsResList{
			Id:         int64(msg.ID),
			SenderId:   msg.SenderID,
			SenderName: userMap[targetID].NickName,
			ReceiverId: msg.ReceiverID,
			Content:    msg.Content,
			MsgType:    msg.MsgType,
			IsRead:     isRead,
			SendTime:   msg.CreatedAt,
		}
	}

	// 转换为列表
	list := make([]res.GetRecentContactsResList, 0, len(contactMap))
	for _, contact := range contactMap {
		list = append(list, contact)
	}

	rs = res.GetRecentContactsRes{
		List: list,
	}

	return rs, nil
}

// GetChatUsers 获取用户列表
// Auth: shack
// Time: 2026年04月03日
func (s *ChatService) GetChatUsers(
	c *gin.Context,
) (rs res.GetChatUsersRes, err error) {
	currentUserID := utils.GetUserID(c)
	if currentUserID == 0 {
		return rs, biz_err.New(biz_err.AUTH_ERROR, "用户未登录")
	}

	// 查询所有启用的用户
	var users []systemRbac.User
	if err := global.GVA_DB.Where("enable = ?", 1).Find(&users).Error; err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR, "查询用户列表失败")
	}

	list := make([]res.GetChatUsersResList, 0, len(users))
	currentUserIDInt64 := int64(currentUserID)

	for _, user := range users {
		// 排除自己
		if int64(user.ID) == currentUserIDInt64 {
			continue
		}

		// 检查是否被拉黑
		isBlocked, _ := s.checkUserBlocked(context.Background(), currentUserIDInt64, int64(user.ID))

		// 获取最后一条消息时间
		type LastMessageInfo struct {
			CreatedAt time.Time
			Content   string
		}
		var lastMsgInfo LastMessageInfo
		global.GVA_DB.Raw(`
			SELECT created_at, content
			FROM chat_messages
			WHERE ((sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?))
			ORDER BY created_at DESC
			LIMIT 1
		`, currentUserIDInt64, int64(user.ID), int64(user.ID), currentUserIDInt64).Scan(&lastMsgInfo)

		list = append(list, res.GetChatUsersResList{
			Id:             int64(user.ID),
			Username:       user.Username,
			Nickname:       user.NickName,
			Avatar:         user.HeaderImg,
			LastMessage:    lastMsgInfo.Content,
			LastMessageTime: lastMsgInfo.CreatedAt,
			IsBlocked:      isBlocked,
		})
	}

	rs = res.GetChatUsersRes{
		List: list,
	}

	return rs, nil
}

// MarkChatAsRead 标记消息为已读
// Auth: shack
// Time: 2026年04月03日
func (s *ChatService) MarkChatAsRead(
	c *gin.Context,
	r req.MarkChatAsReadReq,
) error {
	currentUserID := utils.GetUserID(c)
	if currentUserID == 0 {
		return biz_err.New(biz_err.AUTH_ERROR, "用户未登录")
	}

	// 标记来自指定用户的所有未读消息为已读
	if err := global.GVA_DB.Model(&systemRbac.ChatMessage{}).
		Where("sender_id = ? AND receiver_id = ? AND is_read = ?", r.SenderId, currentUserID, false).
		Update("is_read", true).Error; err != nil {
		return biz_err.New(biz_err.DB_ERROR, "标记消息失败")
	}

	// 清除未读数缓存
	cacheKey := fmt.Sprintf("%s%d", RedisKeyUnreadCount, currentUserID)
	if global.GVA_REDIS != nil {
		global.GVA_REDIS.Del(context.Background(), cacheKey)
	}

	return nil
}

// GetChatUnreadCount 获取未读消息数量
// Auth: shack
// Time: 2026年04月03日
func (s *ChatService) GetChatUnreadCount(
	c *gin.Context,
) (rs res.GetChatUnreadCountRes, err error) {
	currentUserID := utils.GetUserID(c)
	if currentUserID == 0 {
		return rs, biz_err.New(biz_err.AUTH_ERROR, "用户未登录")
	}

	// 尝试从Redis获取缓存
	cacheKey := fmt.Sprintf("%s%d", RedisKeyUnreadCount, currentUserID)
	if global.GVA_REDIS != nil {
		if count, err := global.GVA_REDIS.Get(context.Background(), cacheKey).Int(); err == nil {
			rs = res.GetChatUnreadCountRes{
				Count: count,
			}
			return rs, nil
		}
	}

	// 从数据库查询
	var count int64
	if err := global.GVA_DB.Model(&systemRbac.ChatMessage{}).
		Where("receiver_id = ? AND is_read = ?", currentUserID, false).
		Count(&count).Error; err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR, "查询未读消息失败")
	}

	// 缓存结果
	if global.GVA_REDIS != nil {
		global.GVA_REDIS.Set(context.Background(), cacheKey, count, RedisExpireUnread*time.Second)
	}

	rs = res.GetChatUnreadCountRes{
		Count: int(count),
	}

	return rs, nil
}

// CheckUserOnline 检查用户是否在线
// Auth: shack
// Time: 2026年04月03日
func (s *ChatService) CheckUserOnline(
	c *gin.Context,
	r req.CheckUserOnlineReq,
) (rs res.CheckUserOnlineRes, err error) {
	currentUserID := utils.GetUserID(c)
	if currentUserID == 0 {
		return rs, biz_err.New(biz_err.AUTH_ERROR, "用户未登录")
	}

	// 从Redis检查在线状态
	cacheKey := fmt.Sprintf("%s%d", RedisKeyUserOnline, r.UserId)
	online := false

	if global.GVA_REDIS != nil {
		if exists, _ := global.GVA_REDIS.Exists(context.Background(), cacheKey).Result(); exists > 0 {
			online = true
		}
	}

	rs = res.CheckUserOnlineRes{
		Online: online,
	}

	return rs, nil
}

// ClearChatHistory 清空聊天记录
// Auth: shack
// Time: 2026年04月03日
func (s *ChatService) ClearChatHistory(
	c *gin.Context,
	r req.ClearChatHistoryReq,
) error {
	currentUserID := utils.GetUserID(c)
	if currentUserID == 0 {
		return biz_err.New(biz_err.AUTH_ERROR, "用户未登录")
	}

	// 删除聊天记录
	if err := global.GVA_DB.Where(
		"(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
		currentUserID, r.TargetId, r.TargetId, currentUserID,
	).Delete(&systemRbac.ChatMessage{}).Error; err != nil {
		return biz_err.New(biz_err.DB_ERROR, "清空聊天记录失败")
	}

	// 清除缓存
	s.clearChatCache(int64(currentUserID), int64(r.TargetId))

	return nil
}

// BlockUser 拉黑用户
// Auth: shack
// Time: 2026年04月03日
func (s *ChatService) BlockUser(
	c *gin.Context,
	r req.BlockUserReq,
) error {
	currentUserID := utils.GetUserID(c)
	if currentUserID == 0 {
		return biz_err.New(biz_err.AUTH_ERROR, "用户未登录")
	}

	// 检查是否已经拉黑
	var count int64
	global.GVA_DB.Model(&systemRbac.ChatBlacklist{}).
		Where("user_id = ? AND blocked_user_id = ?", currentUserID, r.TargetId).
		Count(&count)

	if count > 0 {
		return biz_err.New(biz_err.PARAM_ERROR, "已经拉黑该用户")
	}

	// 创建拉黑记录
	blacklist := systemRbac.ChatBlacklist{
		UserId:        int64(currentUserID),
		BlockedUserId: int64(r.TargetId),
	}

	if err := global.GVA_DB.Create(&blacklist).Error; err != nil {
		return biz_err.New(biz_err.DB_ERROR, "拉黑用户失败")
	}

	// 清除黑名单缓存
	cacheKey := fmt.Sprintf("%s%d", RedisKeyBlacklist, currentUserID)
	if global.GVA_REDIS != nil {
		global.GVA_REDIS.Del(context.Background(), cacheKey)
	}

	return nil
}

// UnblockUser 取消拉黑
// Auth: shack
// Time: 2026年04月03日
func (s *ChatService) UnblockUser(
	c *gin.Context,
	r req.UnblockUserReq,
) error {
	currentUserID := utils.GetUserID(c)
	if currentUserID == 0 {
		return biz_err.New(biz_err.AUTH_ERROR, "用户未登录")
	}

	// 删除拉黑记录
	if err := global.GVA_DB.Where(
		"user_id = ? AND blocked_user_id = ?",
		currentUserID, int64(r.TargetId),
	).Delete(&systemRbac.ChatBlacklist{}).Error; err != nil {
		return biz_err.New(biz_err.DB_ERROR, "取消拉黑失败")
	}

	// 清除黑名单缓存
	cacheKey := fmt.Sprintf("%s%d", RedisKeyBlacklist, currentUserID)
	if global.GVA_REDIS != nil {
		global.GVA_REDIS.Del(context.Background(), cacheKey)
	}

	return nil
}

// GetBlacklist 获取黑名单列表
// Auth: shack
// Time: 2026年04月03日
func (s *ChatService) GetBlacklist(
	c *gin.Context,
) (rs res.GetBlacklistRes, err error) {
	currentUserID := utils.GetUserID(c)
	if currentUserID == 0 {
		return rs, biz_err.New(biz_err.AUTH_ERROR, "用户未登录")
	}

	// 尝试从Redis获取缓存
	cacheKey := fmt.Sprintf("%s%d", RedisKeyBlacklist, currentUserID)
	if global.GVA_REDIS != nil {
		cached, err := global.GVA_REDIS.Get(context.Background(), cacheKey).Result()
		if err == nil && cached != "" {
			// TODO: 反序列化缓存数据
			_ = cached
		}
	}

	// 查询黑名单
	var blacklists []systemRbac.ChatBlacklist
	if err := global.GVA_DB.Where("user_id = ?", currentUserID).Find(&blacklists).Error; err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR, "查询黑名单失败")
	}

	// 如果黑名单为空，直接返回
	if len(blacklists) == 0 {
		rs = res.GetBlacklistRes{
			List: []res.GetBlacklistResList{},
		}
		return rs, nil
	}

	// 收集被拉黑用户ID
	blockedUserIDs := make([]int64, len(blacklists))
	for i, bl := range blacklists {
		blockedUserIDs[i] = bl.BlockedUserId
	}

	// 查询被拉黑用户信息
	var blockedUsers []systemRbac.User
	global.GVA_DB.Where("id IN ?", blockedUserIDs).Find(&blockedUsers)

	userMap := make(map[int64]systemRbac.User)
	for _, user := range blockedUsers {
		userMap[int64(user.ID)] = user
	}

	// 构建响应
	list := make([]res.GetBlacklistResList, 0, len(blacklists))
	for _, bl := range blacklists {
		if user, ok := userMap[bl.BlockedUserId]; ok {
			list = append(list, res.GetBlacklistResList{
				Id:              int64(bl.ID),
				UserId:          bl.UserId,
				BlockedUserId:   bl.BlockedUserId,
				BlockedUserName: user.NickName,
				BlockedUserAvatar: user.HeaderImg,
				CreateTime:      bl.CreatedAt,
			})
		}
	}

	rs = res.GetBlacklistRes{
		List: list,
	}

	return rs, nil
}

// CheckIsBlocked 检查是否拉黑
// Auth: shack
// Time: 2026年04月03日
func (s *ChatService) CheckIsBlocked(
	c *gin.Context,
	r req.CheckIsBlockedReq,
) (rs res.CheckIsBlockedRes, err error) {
	currentUserID := utils.GetUserID(c)
	if currentUserID == 0 {
		return rs, biz_err.New(biz_err.AUTH_ERROR, "用户未登录")
	}

	blocked, err := s.checkUserBlocked(context.Background(), int64(currentUserID), int64(r.TargetId))
	if err != nil {
		return rs, err
	}

	rs = res.CheckIsBlockedRes{
		Blocked: blocked,
	}

	return rs, nil
}

// GetMessageStats 获取消息统计
// Auth: shack
// Time: 2026年04月03日
func (s *ChatService) GetMessageStats(
	c *gin.Context,
) (rs res.GetMessageStatsRes, err error) {
	currentUserID := utils.GetUserID(c)
	if currentUserID == 0 {
		return rs, biz_err.New(biz_err.AUTH_ERROR, "用户未登录")
	}

	// 查询未读消息数
	var chatCount int64
	if err := global.GVA_DB.Model(&systemRbac.ChatMessage{}).
		Where("receiver_id = ? AND is_read = ?", currentUserID, false).
		Count(&chatCount).Error; err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR, "查询消息统计失败")
	}

	rs = res.GetMessageStatsRes{
		ChatCount: int(chatCount),
	}

	return rs, nil
}

// ============ 辅助函数 ============

// getOrCreateConversation 获取或创建会话
func (s *ChatService) getOrCreateConversation(db *gorm.DB, user1ID, user2ID int64) (int64, error) {
	// 确保 user1ID < user2ID，保证会话唯一性
	if user1ID > user2ID {
		user1ID, user2ID = user2ID, user1ID
	}

	// 查找已存在的会话
	var conversation systemRbac.ChatConversation
	err := db.Where(
		"user1_id = ? AND user2_id = ?",
		user1ID, user2ID,
	).First(&conversation).Error

	if err == nil {
		// 会话已存在
		return int64(conversation.ID), nil
	}

	if err != gorm.ErrRecordNotFound {
		// 数据库错误
		return 0, biz_err.New(biz_err.DB_ERROR, "查询会话失败")
	}

	// 创建新会话
	conversation = systemRbac.ChatConversation{
		User1ID:       user1ID,
		User2ID:       user2ID,
		User1Unread:   0,
		User2Unread:   0,
	}

	if err := db.Create(&conversation).Error; err != nil {
		return 0, biz_err.New(biz_err.DB_ERROR, "创建会话失败")
	}

	return int64(conversation.ID), nil
}

// updateConversationAfterSend 发送消息后更新会话
func (s *ChatService) updateConversationAfterSend(
	db *gorm.DB,
	conversationID int64,
	senderID, receiverID int64,
	content string,
	time *time.Time,
) error {
	// 获取会话信息
	var conversation systemRbac.ChatConversation
	if err := db.First(&conversation, conversationID).Error; err != nil {
		return biz_err.New(biz_err.DB_ERROR, "查询会话失败")
	}

	// 更新最后消息信息
	updates := map[string]interface{}{
		"last_msg":      content,
		"last_msg_time": time.Unix(),
	}

	// 更新未读数
	if senderID == conversation.User1ID {
		updates["user2_unread"] = conversation.User2Unread + 1
	} else {
		updates["user1_unread"] = conversation.User1Unread + 1
	}

	if err := db.Model(&conversation).Updates(updates).Error; err != nil {
		return biz_err.New(biz_err.DB_ERROR, "更新会话失败")
	}

	return nil
}

// checkUserBlocked 检查用户是否被拉黑
func (s *ChatService) checkUserBlocked(ctx context.Context, userID, targetID int64) (bool, error) {
	// 尝试从Redis获取缓存
	cacheKey := fmt.Sprintf("%s%d", RedisKeyBlacklist, userID)
	if global.GVA_REDIS != nil {
		cached, err := global.GVA_REDIS.Get(ctx, cacheKey).Result()
		if err == nil && cached != "" {
			// TODO: 反序列化缓存数据并检查
			_ = cached
		}
	}

	// 从数据库查询
	var count int64
	if err := global.GVA_DB.Model(&systemRbac.ChatBlacklist{}).
		Where("user_id = ? AND blocked_user_id = ?", targetID, userID).
		Count(&count).Error; err != nil {
		return false, biz_err.New(biz_err.DB_ERROR, "查询黑名单失败")
	}

	return count > 0, nil
}

// clearChatCache 清除聊天相关缓存
func (s *ChatService) clearChatCache(userIDs ...int64) {
	if global.GVA_REDIS == nil {
		return
	}

	ctx := context.Background()
	for _, userID := range userIDs {
		// 清除未读数缓存
		cacheKey := fmt.Sprintf("%s%d", RedisKeyUnreadCount, userID)
		global.GVA_REDIS.Del(ctx, cacheKey)

		// 清除最近联系人缓存
		cacheKey = fmt.Sprintf("%s%d", RedisKeyRecentContacts, userID)
		global.GVA_REDIS.Del(ctx, cacheKey)

		// 清除黑名单缓存
		cacheKey = fmt.Sprintf("%s%d", RedisKeyBlacklist, userID)
		global.GVA_REDIS.Del(ctx, cacheKey)
	}
}

