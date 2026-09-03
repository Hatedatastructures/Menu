package systemRbac

import (
	"strconv"

	"github.com/gin-gonic/gin"
	req "shack/internal/model/systemRbac/request"
	_ "shack/internal/model/systemRbac/response"
	"shack/internal/vo"
	"shack/internal/utils/validator"
	biz_err "shack/internal/error"
)

type ChatApi struct{}

// SendChatMessageHandler
// @Tags systemRbacchatApi
// @Summary SendChatMessageHandler 发送私聊消息-前台使用
// @Description SendChatMessageHandler 发送私聊消息-前台使用
// @Param data body req.SendChatMessageReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.SendChatMessageRes}
// @Router /sys/chat/send [POST]
func (s *ChatApi) SendChatMessageHandler(c *gin.Context) {
	var req req.SendChatMessageReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := chatService.SendChatMessage(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// GetChatHistoryHandler
// @Tags systemRbacchatApi
// @Summary GetChatHistoryHandler 获取聊天记录-前台使用
// @Description GetChatHistoryHandler 获取聊天记录-前台使用
// @Param data body req.GetChatHistoryReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.GetChatHistoryRes}
// @Router /sys/chat/history/:targetId [GET]
func (s *ChatApi) GetChatHistoryHandler(c *gin.Context) {
	var req req.GetChatHistoryReq
	// path 参数
	{
		val := c.Param("targetId")
		parsed, err := strconv.Atoi(val)
		if err != nil {
			c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
			return
		}
		req.TargetId = parsed

	}
	// query 参数
	{
		val := c.Query("page")
		if val != "" {
			parsed, err := strconv.Atoi(val)
			if err != nil {
				c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
				return
			}
			req.Page = parsed

		}
	}
	{
		val := c.Query("pageSize")
		if val != "" {
			parsed, err := strconv.Atoi(val)
			if err != nil {
				c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
				return
			}
			req.PageSize = parsed

		}
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := chatService.GetChatHistory(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// GetRecentContactsHandler
// @Tags systemRbacchatApi
// @Summary GetRecentContactsHandler 获取最近联系人-前台使用
// @Description GetRecentContactsHandler 获取最近联系人-前台使用
// @Success 200 {object} vo.Result{data=_.GetRecentContactsRes}
// @Router /sys/chat/contacts [GET]
func (s *ChatApi) GetRecentContactsHandler(c *gin.Context) {
	data, err := chatService.GetRecentContacts(c)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// GetChatUsersHandler
// @Tags systemRbacchatApi
// @Summary GetChatUsersHandler 获取用户列表(用于选择聊天对象)-前台使用
// @Description GetChatUsersHandler 获取用户列表(用于选择聊天对象)-前台使用
// @Success 200 {object} vo.Result{data=_.GetChatUsersRes}
// @Router /sys/chat/users [GET]
func (s *ChatApi) GetChatUsersHandler(c *gin.Context) {
	data, err := chatService.GetChatUsers(c)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// MarkChatAsReadHandler
// @Tags systemRbacchatApi
// @Summary MarkChatAsReadHandler 标记消息为已读-前台使用
// @Description MarkChatAsReadHandler 标记消息为已读-前台使用
// @Param data body req.MarkChatAsReadReq true "请求参数"
// @Success 200 {object} vo.Result{}
// @Router /sys/chat/read/:senderId [POST]
func (s *ChatApi) MarkChatAsReadHandler(c *gin.Context) {
	var req req.MarkChatAsReadReq
	// path 参数
	{
		val := c.Param("senderId")
		parsed, err := strconv.Atoi(val)
		if err != nil {
			c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
			return
		}
		req.SenderId = parsed

	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	err := chatService.MarkChatAsRead(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, "操作成功"))
}

// GetChatUnreadCountHandler
// @Tags systemRbacchatApi
// @Summary GetChatUnreadCountHandler 获取未读消息数量-前台使用
// @Description GetChatUnreadCountHandler 获取未读消息数量-前台使用
// @Success 200 {object} vo.Result{data=_.GetChatUnreadCountRes}
// @Router /sys/chat/unread-count [GET]
func (s *ChatApi) GetChatUnreadCountHandler(c *gin.Context) {
	data, err := chatService.GetChatUnreadCount(c)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// CheckUserOnlineHandler
// @Tags systemRbacchatApi
// @Summary CheckUserOnlineHandler 检查用户是否在线-前台使用
// @Description CheckUserOnlineHandler 检查用户是否在线-前台使用
// @Param data body req.CheckUserOnlineReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.CheckUserOnlineRes}
// @Router /sys/chat/online/:userId [GET]
func (s *ChatApi) CheckUserOnlineHandler(c *gin.Context) {
	var req req.CheckUserOnlineReq
	// path 参数
	{
		val := c.Param("userId")
		parsed, err := strconv.Atoi(val)
		if err != nil {
			c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
			return
		}
		req.UserId = parsed

	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := chatService.CheckUserOnline(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// ClearChatHistoryHandler
// @Tags systemRbacchatApi
// @Summary ClearChatHistoryHandler 清空聊天记录-前台使用
// @Description ClearChatHistoryHandler 清空聊天记录-前台使用
// @Param data body req.ClearChatHistoryReq true "请求参数"
// @Success 200 {object} vo.Result{}
// @Router /sys/chat/clear/:targetId [DELETE]
func (s *ChatApi) ClearChatHistoryHandler(c *gin.Context) {
	var req req.ClearChatHistoryReq
	// path 参数
	{
		val := c.Param("targetId")
		parsed, err := strconv.Atoi(val)
		if err != nil {
			c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
			return
		}
		req.TargetId = parsed

	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	err := chatService.ClearChatHistory(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, "操作成功"))
}

// BlockUserHandler
// @Tags systemRbacchatApi
// @Summary BlockUserHandler 拉黑用户-前台使用
// @Description BlockUserHandler 拉黑用户-前台使用
// @Param data body req.BlockUserReq true "请求参数"
// @Success 200 {object} vo.Result{}
// @Router /sys/chat/block/:targetId [POST]
func (s *ChatApi) BlockUserHandler(c *gin.Context) {
	var req req.BlockUserReq
	// path 参数
	{
		val := c.Param("targetId")
		parsed, err := strconv.Atoi(val)
		if err != nil {
			c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
			return
		}
		req.TargetId = parsed

	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	err := chatService.BlockUser(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, "操作成功"))
}

// UnblockUserHandler
// @Tags systemRbacchatApi
// @Summary UnblockUserHandler 取消拉黑-前台使用
// @Description UnblockUserHandler 取消拉黑-前台使用
// @Param data body req.UnblockUserReq true "请求参数"
// @Success 200 {object} vo.Result{}
// @Router /sys/chat/block/:targetId [DELETE]
func (s *ChatApi) UnblockUserHandler(c *gin.Context) {
	var req req.UnblockUserReq
	// path 参数
	{
		val := c.Param("targetId")
		parsed, err := strconv.Atoi(val)
		if err != nil {
			c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
			return
		}
		req.TargetId = parsed

	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	err := chatService.UnblockUser(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, "操作成功"))
}

// GetBlacklistHandler
// @Tags systemRbacchatApi
// @Summary GetBlacklistHandler 获取黑名单列表-前台使用
// @Description GetBlacklistHandler 获取黑名单列表-前台使用
// @Success 200 {object} vo.Result{data=_.GetBlacklistRes}
// @Router /sys/chat/blacklist [GET]
func (s *ChatApi) GetBlacklistHandler(c *gin.Context) {
	data, err := chatService.GetBlacklist(c)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// CheckIsBlockedHandler
// @Tags systemRbacchatApi
// @Summary CheckIsBlockedHandler 检查是否拉黑-前台使用
// @Description CheckIsBlockedHandler 检查是否拉黑-前台使用
// @Param data body req.CheckIsBlockedReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.CheckIsBlockedRes}
// @Router /sys/chat/blocked/:targetId [GET]
func (s *ChatApi) CheckIsBlockedHandler(c *gin.Context) {
	var req req.CheckIsBlockedReq
	// path 参数
	{
		val := c.Param("targetId")
		parsed, err := strconv.Atoi(val)
		if err != nil {
			c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
			return
		}
		req.TargetId = parsed

	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := chatService.CheckIsBlocked(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// GetMessageStatsHandler
// @Tags systemRbacchatApi
// @Summary GetMessageStatsHandler 获取消息统计-前台使用
// @Description GetMessageStatsHandler 获取消息统计-前台使用
// @Success 200 {object} vo.Result{data=_.GetMessageStatsRes}
// @Router /sys/chat/stats [GET]
func (s *ChatApi) GetMessageStatsHandler(c *gin.Context) {
	data, err := chatService.GetMessageStats(c)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}