package systemRbac

import (
	"context"
	"encoding/json"
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
	RedisKeyNoticeUnread   = "notice:unread:"    // 未读通知数
	RedisKeyNoticeDetail   = "notice:detail:"    // 通知详情缓存
	RedisKeyNoticeList     = "notice:list:"      // 通知列表缓存
	RedisKeyNoticeChannels = "notice:channels"   // 推送渠道配置

	// 过期时间
	RedisExpireNoticeUnread = 600  // 10分钟
	RedisExpireDetail       = 300  // 5分钟
	RedisExpireList         = 180  // 3分钟
	RedisExpireChannels     = 3600 // 1小时

	// 通知类型
	NoticeTypeNotification = 1 // 通知
	NoticeTypeAnnouncement = 2 // 公告

	// 通知状态
	NoticeStatusDraft     = 0 // 草稿
	NoticeStatusPublished = 1 // 已发布
	NoticeStatusRevoked   = 2 // 已撤回

	// 推送对象类型
	TargetTypeUsers = 1 // 指定用户
	TargetTypeDept  = 2 // 按部门
	TargetTypeAll   = 3 // 全部用户

	// 推送渠道
	ChannelSystem    = "system"    // 系统内
	ChannelEmail     = "email"     // 邮件
	ChannelSMS       = "sms"       // 短信
	ChannelWebSocket = "websocket" // WebSocket
)

type NoticeService struct{}

// GetNoticePage 分页查询通知列表-后台使用
// Auth: shack
// Time: 2026年04月03日
func (s *NoticeService) GetNoticePage(
	c *gin.Context,
	r req.GetNoticePageReq,
) (rs res.GetNoticePageRes, err error) {
	// 分页参数
	page := r.Page
	if page < 1 {
		page = 1
	}
	pageSize := r.PageSize
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	// 构建查询
	db := global.GVA_DB.Model(&systemRbac.Notice{})

	// 条件过滤
	if r.Title != "" {
		db = db.Where("title LIKE ?", "%"+r.Title+"%")
	}
	if r.NoticeType > 0 {
		db = db.Where("notice_type = ?", r.NoticeType)
	}
	if r.Status >= 0 {
		db = db.Where("status = ?", r.Status)
	}

	// 查询总数
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR, "查询通知列表失败")
	}

	// 分页查询
	var notices []systemRbac.Notice
	offset := (page - 1) * pageSize
	if err := db.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&notices).Error; err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR, "查询通知列表失败")
	}

	// 构建响应
	list := make([]res.GetNoticePageResList, 0, len(notices))
	for _, notice := range notices {
		list = append(list, res.GetNoticePageResList{
			Id:         int64(notice.ID),
			Title:      notice.Title,
			Content:    notice.Content,
			NoticeType: notice.NoticeType,
			Channels:   notice.Channels,
			TargetType: notice.TargetType,
			TargetIds:  notice.TargetIds,
			Status:     notice.Status,
			CreateBy:   notice.CreateBy,
			CreateName: notice.CreateName,
			CreateTime: notice.CreatedAt,
		})
	}

	rs = res.GetNoticePageRes{
		Title:      r.Title,
		NoticeType: r.NoticeType,
		Status:     r.Status,
		List:       list,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
	}

	return rs, nil
}

// GetMyNotices 获取当前用户通知列表-前台使用
// Auth: shack
// Time: 2026年04月03日
func (s *NoticeService) GetMyNotices(
	c *gin.Context,
	r req.GetMyNoticesReq,
) (rs res.GetMyNoticesRes, err error) {
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
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	// 查询当前用户的通知记录
	db := global.GVA_DB.Table("sys_notices n").
		Select("DISTINCT n.*").
		Joins("INNER JOIN sys_notice_records nr ON n.id = nr.notice_id").
		Where("nr.user_id = ?", currentUserID)

	// 条件过滤
	if r.IsRead >= 0 {
		db = db.Where("nr.is_read = ?", r.IsRead == 1)
	}

	// 只显示已发布的通知
	db = db.Where("n.status = ?", NoticeStatusPublished)

	// 查询总数
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR, "查询通知列表失败")
	}

	// 分页查询
	var notices []systemRbac.Notice
	offset := (page - 1) * pageSize
	if err := db.Order("n.created_at DESC").Offset(offset).Limit(pageSize).Find(&notices).Error; err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR, "查询通知列表失败")
	}

	// 构建响应并获取阅读状态
	list := make([]res.GetMyNoticesResList, 0, len(notices))
	for _, notice := range notices {
		// 查询阅读状态
		var record systemRbac.NoticeRecord
		global.GVA_DB.Where(
			"notice_id = ? AND user_id = ?",
			notice.ID, currentUserID,
		).First(&record)

		list = append(list, res.GetMyNoticesResList{
			Id:         int64(notice.ID),
			Title:      notice.Title,
			Content:    notice.Content,
			NoticeType: notice.NoticeType,
			Status:     notice.Status,
			CreateTime: notice.CreatedAt,
		})
	}

	rs = res.GetMyNoticesRes{
		IsRead:   r.IsRead,
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}

	return rs, nil
}

// GetNoticeDetail 获取通知详情-前台/后台使用
// Auth: shack
// Time: 2026年04月03日
func (s *NoticeService) GetNoticeDetail(
	c *gin.Context,
	r req.GetNoticeDetailReq,
) (rs res.GetNoticeDetailRes, err error) {
	// 尝试从Redis获取缓存
	cacheKey := fmt.Sprintf("%s%d", RedisKeyNoticeDetail, r.Id)
	if global.GVA_REDIS != nil {
		cached, err := global.GVA_REDIS.Get(context.Background(), cacheKey).Result()
		if err == nil && cached != "" {
			// TODO: 反序列化缓存数据
			_ = cached
		}
	}

	// 从数据库查询
	var notice systemRbac.Notice
	if err := global.GVA_DB.First(&notice, r.Id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return rs, biz_err.New(biz_err.PARAM_ERROR, "通知不存在")
		}
		return rs, biz_err.New(biz_err.DB_ERROR, "查询通知详情失败")
	}

	rs = res.GetNoticeDetailRes{
		Id:         int64(notice.ID),
		Title:      notice.Title,
		Content:    notice.Content,
		NoticeType: notice.NoticeType,
		Channels:   notice.Channels,
		TargetType: notice.TargetType,
		TargetIds:  notice.TargetIds,
		Status:     notice.Status,
		CreateBy:   notice.CreateBy,
		CreateName: notice.CreateName,
		CreateTime: notice.CreatedAt,
	}

	// 缓存结果
	if global.GVA_REDIS != nil {
		// TODO: 序列化后缓存
		global.GVA_REDIS.Set(context.Background(), cacheKey, "1", RedisExpireDetail*time.Second)
	}

	return rs, nil
}

// CreateNotice 创建通知-后台使用
// Auth: shack
// Time: 2026年04月03日
func (s *NoticeService) CreateNotice(
	c *gin.Context,
	r req.CreateNoticeReq,
) error {
	currentUserID := utils.GetUserID(c)
	if currentUserID == 0 {
		return biz_err.New(biz_err.AUTH_ERROR, "用户未登录")
	}

	// 参数验证
	if r.Title == "" {
		return biz_err.New(biz_err.PARAM_MISSING, "通知标题不能为空")
	}
	if r.Content == "" {
		return biz_err.New(biz_err.PARAM_MISSING, "通知内容不能为空")
	}
	if r.NoticeType < 1 || r.NoticeType > 2 {
		return biz_err.New(biz_err.PARAM_ERROR, "通知类型不正确")
	}
	if r.TargetType < 1 || r.TargetType > 3 {
		return biz_err.New(biz_err.PARAM_ERROR, "推送对象类型不正确")
	}

	// 验证推送渠道JSON格式
	if r.Channels != "" {
		var channels []string
		if err := json.Unmarshal([]byte(r.Channels), &channels); err != nil {
			return biz_err.New(biz_err.PARAM_FORMAT, "推送渠道格式不正确")
		}
	}

	// 验证推送对象ID JSON格式
	if r.TargetType == TargetTypeUsers && r.TargetIds != "" {
		var targetIDs []int64
		if err := json.Unmarshal([]byte(r.TargetIds), &targetIDs); err != nil {
			return biz_err.New(biz_err.PARAM_FORMAT, "推送对象ID格式不正确")
		}
	}

	// 获取当前用户信息
	var user systemRbac.User
	if err := global.GVA_DB.First(&user, currentUserID).Error; err != nil {
		return biz_err.New(biz_err.DB_ERROR, "获取用户信息失败")
	}

	// 创建通知
	notice := systemRbac.Notice{
		Title:       r.Title,
		Content:     r.Content,
		NoticeType:  r.NoticeType,
		Channels:    r.Channels,
		TargetType:  r.TargetType,
		TargetIds:   r.TargetIds,
		Status:      NoticeStatusDraft,
		CreateBy:    int64(currentUserID),
		CreateName:  user.NickName,
	}

	if err := global.GVA_DB.Create(&notice).Error; err != nil {
		global.GVA_LOG.Error("创建通知失败", zap.Error(err))
		return biz_err.New(biz_err.DB_ERROR, "创建通知失败")
	}

	return nil
}

// UpdateNotice 更新通知-后台使用
// Auth: shack
// Time: 2026年04月03日
func (s *NoticeService) UpdateNotice(
	c *gin.Context,
	r req.UpdateNoticeReq,
) error {
	currentUserID := utils.GetUserID(c)
	if currentUserID == 0 {
		return biz_err.New(biz_err.AUTH_ERROR, "用户未登录")
	}

	// 参数验证
	if r.Id == 0 {
		return biz_err.New(biz_err.PARAM_MISSING, "通知ID不能为空")
	}
	if r.Title == "" {
		return biz_err.New(biz_err.PARAM_MISSING, "通知标题不能为空")
	}
	if r.Content == "" {
		return biz_err.New(biz_err.PARAM_MISSING, "通知内容不能为空")
	}

	// 查询通知是否存在
	var notice systemRbac.Notice
	if err := global.GVA_DB.First(&notice, r.Id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return biz_err.New(biz_err.PARAM_ERROR, "通知不存在")
		}
		return biz_err.New(biz_err.DB_ERROR, "查询通知失败")
	}

	// 已发布的通知不能修改
	if notice.Status == NoticeStatusPublished {
		return biz_err.New(biz_err.FORBIDDEN_OPERATION, "已发布的通知不能修改")
	}

	// 更新通知
	updates := map[string]interface{}{
		"title":       r.Title,
		"content":     r.Content,
		"notice_type": r.NoticeType,
		"channels":    r.Channels,
		"target_type": r.TargetType,
		"target_ids":  r.TargetIds,
	}

	if err := global.GVA_DB.Model(&notice).Updates(updates).Error; err != nil {
		return biz_err.New(biz_err.DB_ERROR, "更新通知失败")
	}

	// 清除缓存
	s.clearNoticeCache(int64(notice.ID))

	return nil
}

// DeleteNotice 删除通知-后台使用
// Auth: shack
// Time: 2026年04月03日
func (s *NoticeService) DeleteNotice(
	c *gin.Context,
	r req.DeleteNoticeReq,
) error {
	currentUserID := utils.GetUserID(c)
	if currentUserID == 0 {
		return biz_err.New(biz_err.AUTH_ERROR, "用户未登录")
	}

	// 查询通知是否存在
	var notice systemRbac.Notice
	if err := global.GVA_DB.First(&notice, r.Id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return biz_err.New(biz_err.PARAM_ERROR, "通知不存在")
		}
		return biz_err.New(biz_err.DB_ERROR, "查询通知失败")
	}

	// 已发布的通知不能删除
	if notice.Status == NoticeStatusPublished {
		return biz_err.New(biz_err.FORBIDDEN_OPERATION, "已发布的通知不能删除")
	}

	// 删除通知
	if err := global.GVA_DB.Delete(&notice).Error; err != nil {
		return biz_err.New(biz_err.DB_ERROR, "删除通知失败")
	}

	// 清除缓存
	s.clearNoticeCache(int64(r.Id))

	return nil
}

// PublishNotice 发布通知-后台使用
// Auth: shack
// Time: 2026年04月03日
func (s *NoticeService) PublishNotice(
	c *gin.Context,
	r req.PublishNoticeReq,
) error {
	currentUserID := utils.GetUserID(c)
	if currentUserID == 0 {
		return biz_err.New(biz_err.AUTH_ERROR, "用户未登录")
	}

	// 查询通知是否存在
	var notice systemRbac.Notice
	if err := global.GVA_DB.First(&notice, r.Id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return biz_err.New(biz_err.PARAM_ERROR, "通知不存在")
		}
		return biz_err.New(biz_err.DB_ERROR, "查询通知失败")
	}

	// 只有草稿状态可以发布
	if notice.Status != NoticeStatusDraft {
		return biz_err.New(biz_err.FORBIDDEN_OPERATION, "只有草稿状态的通知可以发布")
	}

	// 开始事务
	tx := global.GVA_DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 更新通知状态
	now := time.Now()
	if err := tx.Model(&notice).Updates(map[string]interface{}{
		"status":       NoticeStatusPublished,
		"publish_time": &now,
	}).Error; err != nil {
		return biz_err.New(biz_err.DB_ERROR, "发布通知失败")
	}

	// 创建通知记录
	if err := s.createNoticeRecords(tx, &notice); err != nil {
		return err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return biz_err.New(biz_err.DB_ERROR, "发布通知失败")
	}

	// 清除缓存
	s.clearNoticeCache(int64(r.Id))

	// TODO: 根据channels配置发送通知

	return nil
}

// MarkNoticeAsRead 标记通知为已读-前台使用
// Auth: shack
// Time: 2026年04月03日
func (s *NoticeService) MarkNoticeAsRead(
	c *gin.Context,
	r req.MarkNoticeAsReadReq,
) error {
	currentUserID := utils.GetUserID(c)
	if currentUserID == 0 {
		return biz_err.New(biz_err.AUTH_ERROR, "用户未登录")
	}

	// 查询通知记录是否存在
	var record systemRbac.NoticeRecord
	err := global.GVA_DB.Where(
		"notice_id = ? AND user_id = ?",
		r.Id, currentUserID,
	).First(&record).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return biz_err.New(biz_err.PARAM_ERROR, "通知记录不存在")
		}
		return biz_err.New(biz_err.DB_ERROR, "查询通知记录失败")
	}

	// 如果已读，直接返回
	if record.IsRead {
		return nil
	}

	// 更新为已读
	now := time.Now()
	if err := global.GVA_DB.Model(&record).Updates(map[string]interface{}{
		"is_read":   true,
		"read_time": &now,
	}).Error; err != nil {
		return biz_err.New(biz_err.DB_ERROR, "标记通知失败")
	}

	// 清除未读数缓存
	cacheKey := fmt.Sprintf("%s%d", RedisKeyNoticeUnread, currentUserID)
	if global.GVA_REDIS != nil {
		global.GVA_REDIS.Del(context.Background(), cacheKey)
	}

	return nil
}

// MarkAllNoticesAsRead 标记所有通知为已读-前台使用
// Auth: shack
// Time: 2026年04月03日
func (s *NoticeService) MarkAllNoticesAsRead(
	c *gin.Context,
) error {
	currentUserID := utils.GetUserID(c)
	if currentUserID == 0 {
		return biz_err.New(biz_err.AUTH_ERROR, "用户未登录")
	}

	// 查询所有未读的通知记录
	var records []systemRbac.NoticeRecord
	if err := global.GVA_DB.Where(
		"user_id = ? AND is_read = ?",
		currentUserID, false,
	).Find(&records).Error; err != nil {
		return biz_err.New(biz_err.DB_ERROR, "查询通知记录失败")
	}

	// 如果没有未读记录，直接返回
	if len(records) == 0 {
		return nil
	}

	// 批量更新为已读
	now := time.Now()
	if err := global.GVA_DB.Model(&systemRbac.NoticeRecord{}).
		Where("user_id = ? AND is_read = ?", currentUserID, false).
		Updates(map[string]interface{}{
			"is_read":   true,
			"read_time": &now,
		}).Error; err != nil {
		return biz_err.New(biz_err.DB_ERROR, "标记通知失败")
	}

	// 清除未读数缓存
	cacheKey := fmt.Sprintf("%s%d", RedisKeyNoticeUnread, currentUserID)
	if global.GVA_REDIS != nil {
		global.GVA_REDIS.Del(context.Background(), cacheKey)
	}

	return nil
}

// GetNoticeUnreadCount 获取未读通知数量-前台使用
// Auth: shack
// Time: 2026年04月03日
func (s *NoticeService) GetNoticeUnreadCount(
	c *gin.Context,
) (rs res.GetNoticeUnreadCountRes, err error) {
	currentUserID := utils.GetUserID(c)
	if currentUserID == 0 {
		return rs, biz_err.New(biz_err.AUTH_ERROR, "用户未登录")
	}

	// 尝试从Redis获取缓存
	cacheKey := fmt.Sprintf("%s%d", RedisKeyNoticeUnread, currentUserID)
	if global.GVA_REDIS != nil {
		if count, err := global.GVA_REDIS.Get(context.Background(), cacheKey).Int(); err == nil {
			rs = res.GetNoticeUnreadCountRes{
				Count: count,
			}
			return rs, nil
		}
	}

	// 从数据库查询
	var count int64
	if err := global.GVA_DB.Model(&systemRbac.NoticeRecord{}).
		Joins("INNER JOIN sys_notices ON sys_notices.id = sys_notice_records.notice_id").
		Where("sys_notice_records.user_id = ? AND sys_notice_records.is_read = ? AND sys_notices.status = ?",
			currentUserID, false, NoticeStatusPublished).
		Count(&count).Error; err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR, "查询未读通知失败")
	}

	// 缓存结果
	if global.GVA_REDIS != nil {
		global.GVA_REDIS.Set(context.Background(), cacheKey, count, RedisExpireNoticeUnread*time.Second)
	}

	rs = res.GetNoticeUnreadCountRes{
		Count: int(count),
	}

	return rs, nil
}

// GetNoticeChannels 获取推送渠道-后台使用
// Auth: shack
// Time: 2026年04月03日
func (s *NoticeService) GetNoticeChannels(
	c *gin.Context,
) (rs res.GetNoticeChannelsRes, err error) {
	// 尝试从Redis获取缓存
	if global.GVA_REDIS != nil {
		cached, err := global.GVA_REDIS.Get(context.Background(), RedisKeyNoticeChannels).Result()
		if err == nil && cached != "" {
			// TODO: 反序列化缓存数据
			_ = cached
		}
	}

	// 返回推送渠道列表
	channels := []res.GetNoticeChannelsResList{
		{
			Channel: ChannelSystem,
			Name:    "系统内通知",
			Enabled: true,
		},
		{
			Channel: ChannelEmail,
			Name:    "邮件通知",
			Enabled: false, // 需要配置邮件服务后才能启用
		},
		{
			Channel: ChannelSMS,
			Name:    "短信通知",
			Enabled: false, // 需要配置短信服务后才能启用
		},
		{
			Channel: ChannelWebSocket,
			Name:    "实时推送",
			Enabled: true,
		},
	}

	rs = res.GetNoticeChannelsRes{
		List: channels,
	}

	// 缓存结果
	if global.GVA_REDIS != nil {
		// TODO: 序列化后缓存
		global.GVA_REDIS.Set(context.Background(), RedisKeyNoticeChannels, "1", RedisExpireChannels*time.Second)
	}

	return rs, nil
}

// GetNoticeSendLogs 获取通知推送记录-后台使用
// Auth: shack
// Time: 2026年04月03日
func (s *NoticeService) GetNoticeSendLogs(
	c *gin.Context,
	r req.GetNoticeSendLogsReq,
) (rs res.GetNoticeSendLogsRes, err error) {
	// 查询推送记录
	var records []systemRbac.NoticeRecord
	if err := global.GVA_DB.Where("notice_id = ?", r.Id).Find(&records).Error; err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR, "查询推送记录失败")
	}

	// 构建响应
	list := make([]res.GetNoticeSendLogsResList, 0, len(records))
	for _, record := range records {
		list = append(list, res.GetNoticeSendLogsResList{
			Id:       int64(record.ID),
			NoticeId: record.NoticeId,
			Channel:  record.Channel,
			TargetId: record.UserId,
			Status:   record.SendStatus,
			ErrorMsg: record.ErrorMsg,
			SendTime: record.CreatedAt,
		})
	}

	rs = res.GetNoticeSendLogsRes{
		List: list,
	}

	return rs, nil
}

// RetryNoticeSend 重试失败推送-后台使用
// Auth: shack
// Time: 2026年04月03日
func (s *NoticeService) RetryNoticeSend(
	c *gin.Context,
	r req.RetryNoticeSendReq,
) error {
	// 查询通知是否存在
	var notice systemRbac.Notice
	if err := global.GVA_DB.First(&notice, r.Id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return biz_err.New(biz_err.PARAM_ERROR, "通知不存在")
		}
		return biz_err.New(biz_err.DB_ERROR, "查询通知失败")
	}

	// 查询失败的推送记录
	var records []systemRbac.NoticeRecord
	if err := global.GVA_DB.Where(
		"notice_id = ? AND channel = ? AND send_status = ?",
		r.Id, r.Channel, 0,
	).Find(&records).Error; err != nil {
		return biz_err.New(biz_err.DB_ERROR, "查询推送记录失败")
	}

	// 如果没有失败记录，直接返回
	if len(records) == 0 {
		return biz_err.New(biz_err.PARAM_ERROR, "没有需要重试的推送记录")
	}

	// TODO: 根据channel重新发送通知

	return nil
}

// ============ 辅助函数 ============

// createNoticeRecords 创建通知记录
func (s *NoticeService) createNoticeRecords(db *gorm.DB, notice *systemRbac.Notice) error {
	// 解析推送渠道
	var channels []string
	if notice.Channels != "" {
		if err := json.Unmarshal([]byte(notice.Channels), &channels); err != nil {
			channels = []string{ChannelSystem} // 默认系统内通知
		}
	} else {
		channels = []string{ChannelSystem}
	}

	// 解析推送对象ID
	var targetIDs []int64
	if notice.TargetType == TargetTypeUsers && notice.TargetIds != "" {
		if err := json.Unmarshal([]byte(notice.TargetIds), &targetIDs); err != nil {
			return biz_err.New(biz_err.PARAM_FORMAT, "推送对象ID格式不正确")
		}
	}

	// 获取目标用户列表
	userIDs := make([]int64, 0)
	switch notice.TargetType {
	case TargetTypeAll:
		// 全部用户
		var users []systemRbac.User
		if err := db.Where("enable = ?", 1).Find(&users).Error; err != nil {
			return biz_err.New(biz_err.DB_ERROR, "查询用户列表失败")
		}
		for _, user := range users {
			userIDs = append(userIDs, int64(user.ID))
		}
	case TargetTypeUsers:
		// 指定用户
		userIDs = targetIDs
	case TargetTypeDept:
		// 按部门（TODO: 需要关联部门表）
		return biz_err.New(biz_err.PARAM_ERROR, "按部门推送功能暂未实现")
	}

	// 创建通知记录
	records := make([]systemRbac.NoticeRecord, 0, len(userIDs)*len(channels))
	for _, userID := range userIDs {
		for _, channel := range channels {
			records = append(records, systemRbac.NoticeRecord{
				NoticeId:   int64(notice.ID),
				UserId:     userID,
				IsRead:     false,
				Channel:    channel,
				SendStatus: 1, // 默认成功
			})
		}
	}

	// 批量插入
	if len(records) > 0 {
		if err := db.Create(&records).Error; err != nil {
			global.GVA_LOG.Error("创建通知记录失败", zap.Error(err))
			return biz_err.New(biz_err.DB_ERROR, "创建通知记录失败")
		}
	}

	return nil
}

// clearNoticeCache 清除通知相关缓存
func (s *NoticeService) clearNoticeCache(noticeID int64) {
	if global.GVA_REDIS == nil {
		return
	}

	ctx := context.Background()

	// 清除详情缓存
	cacheKey := fmt.Sprintf("%s%d", RedisKeyNoticeDetail, noticeID)
	global.GVA_REDIS.Del(ctx, cacheKey)

	// 清除列表缓存
	cacheKey = fmt.Sprintf("%s*", RedisKeyNoticeList)
	// TODO: 使用SCAN命令删除匹配的key
}

