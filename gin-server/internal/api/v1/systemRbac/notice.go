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

type NoticeApi struct{}

// GetNoticePageHandler
// @Tags systemRbacnoticeApi
// @Summary GetNoticePageHandler 分页查询通知列表-后台使用
// @Description GetNoticePageHandler 分页查询通知列表-后台使用
// @Param data body req.GetNoticePageReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.GetNoticePageRes}
// @Router /sys/notice/page [GET]
func (s *NoticeApi) GetNoticePageHandler(c *gin.Context) {
	var req req.GetNoticePageReq
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
	{
		val := c.Query("title")
		if val != "" {
			req.Title = val

		}
	}
	{
		val := c.Query("noticeType")
		if val != "" {
			parsed, err := strconv.Atoi(val)
			if err != nil {
				c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
				return
			}
			req.NoticeType = parsed

		}
	}
	{
		val := c.Query("status")
		if val != "" {
			parsed, err := strconv.Atoi(val)
			if err != nil {
				c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
				return
			}
			req.Status = parsed

		}
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := noticeService.GetNoticePage(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// GetMyNoticesHandler
// @Tags systemRbacnoticeApi
// @Summary GetMyNoticesHandler 获取当前用户通知列表-前台使用
// @Description GetMyNoticesHandler 获取当前用户通知列表-前台使用
// @Param data body req.GetMyNoticesReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.GetMyNoticesRes}
// @Router /sys/notice/my [GET]
func (s *NoticeApi) GetMyNoticesHandler(c *gin.Context) {
	var req req.GetMyNoticesReq
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
	{
		val := c.Query("isRead")
		if val != "" {
			parsed, err := strconv.Atoi(val)
			if err != nil {
				c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
				return
			}
			req.IsRead = parsed

		}
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := noticeService.GetMyNotices(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// GetNoticeDetailHandler
// @Tags systemRbacnoticeApi
// @Summary GetNoticeDetailHandler 获取通知详情-前台/后台使用
// @Description GetNoticeDetailHandler 获取通知详情-前台/后台使用
// @Param data body req.GetNoticeDetailReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.GetNoticeDetailRes}
// @Router /sys/notice/:id [GET]
func (s *NoticeApi) GetNoticeDetailHandler(c *gin.Context) {
	var req req.GetNoticeDetailReq
	// path 参数
	{
		val := c.Param("id")
		parsed, err := strconv.Atoi(val)
		if err != nil {
			c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
			return
		}
		req.Id = parsed

	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := noticeService.GetNoticeDetail(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// CreateNoticeHandler
// @Tags systemRbacnoticeApi
// @Summary CreateNoticeHandler 创建通知-后台使用
// @Description CreateNoticeHandler 创建通知-后台使用
// @Param data body req.CreateNoticeReq true "请求参数"
// @Success 200 {object} vo.Result{}
// @Router /sys/notice [POST]
func (s *NoticeApi) CreateNoticeHandler(c *gin.Context) {
	var req req.CreateNoticeReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	err := noticeService.CreateNotice(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, "操作成功"))
}

// UpdateNoticeHandler
// @Tags systemRbacnoticeApi
// @Summary UpdateNoticeHandler 更新通知-后台使用
// @Description UpdateNoticeHandler 更新通知-后台使用
// @Param data body req.UpdateNoticeReq true "请求参数"
// @Success 200 {object} vo.Result{}
// @Router /sys/notice [PUT]
func (s *NoticeApi) UpdateNoticeHandler(c *gin.Context) {
	var req req.UpdateNoticeReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	err := noticeService.UpdateNotice(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, "操作成功"))
}

// DeleteNoticeHandler
// @Tags systemRbacnoticeApi
// @Summary DeleteNoticeHandler 删除通知-后台使用
// @Description DeleteNoticeHandler 删除通知-后台使用
// @Param data body req.DeleteNoticeReq true "请求参数"
// @Success 200 {object} vo.Result{}
// @Router /sys/notice/:id [DELETE]
func (s *NoticeApi) DeleteNoticeHandler(c *gin.Context) {
	var req req.DeleteNoticeReq
	// path 参数
	{
		val := c.Param("id")
		parsed, err := strconv.Atoi(val)
		if err != nil {
			c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
			return
		}
		req.Id = parsed

	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	err := noticeService.DeleteNotice(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, "操作成功"))
}

// PublishNoticeHandler
// @Tags systemRbacnoticeApi
// @Summary PublishNoticeHandler 发布通知-后台使用
// @Description PublishNoticeHandler 发布通知-后台使用
// @Param data body req.PublishNoticeReq true "请求参数"
// @Success 200 {object} vo.Result{}
// @Router /sys/notice/:id/publish [POST]
func (s *NoticeApi) PublishNoticeHandler(c *gin.Context) {
	var req req.PublishNoticeReq
	// path 参数
	{
		val := c.Param("id")
		parsed, err := strconv.Atoi(val)
		if err != nil {
			c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
			return
		}
		req.Id = parsed

	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	err := noticeService.PublishNotice(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, "操作成功"))
}

// MarkNoticeAsReadHandler
// @Tags systemRbacnoticeApi
// @Summary MarkNoticeAsReadHandler 标记通知为已读-前台使用
// @Description MarkNoticeAsReadHandler 标记通知为已读-前台使用
// @Param data body req.MarkNoticeAsReadReq true "请求参数"
// @Success 200 {object} vo.Result{}
// @Router /sys/notice/:id/read [POST]
func (s *NoticeApi) MarkNoticeAsReadHandler(c *gin.Context) {
	var req req.MarkNoticeAsReadReq
	// path 参数
	{
		val := c.Param("id")
		parsed, err := strconv.Atoi(val)
		if err != nil {
			c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
			return
		}
		req.Id = parsed

	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	err := noticeService.MarkNoticeAsRead(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, "操作成功"))
}

// MarkAllNoticesAsReadHandler
// @Tags systemRbacnoticeApi
// @Summary MarkAllNoticesAsReadHandler 标记所有通知为已读-前台使用
// @Description MarkAllNoticesAsReadHandler 标记所有通知为已读-前台使用
// @Success 200 {object} vo.Result{}
// @Router /sys/notice/read-all [POST]
func (s *NoticeApi) MarkAllNoticesAsReadHandler(c *gin.Context) {
	err := noticeService.MarkAllNoticesAsRead(c)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, "操作成功"))
}

// GetNoticeUnreadCountHandler
// @Tags systemRbacnoticeApi
// @Summary GetNoticeUnreadCountHandler 获取未读通知数量-前台使用
// @Description GetNoticeUnreadCountHandler 获取未读通知数量-前台使用
// @Success 200 {object} vo.Result{data=_.GetNoticeUnreadCountRes}
// @Router /sys/notice/unread-count [GET]
func (s *NoticeApi) GetNoticeUnreadCountHandler(c *gin.Context) {
	data, err := noticeService.GetNoticeUnreadCount(c)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// GetNoticeChannelsHandler
// @Tags systemRbacnoticeApi
// @Summary GetNoticeChannelsHandler 获取推送渠道-后台使用
// @Description GetNoticeChannelsHandler 获取推送渠道-后台使用
// @Success 200 {object} vo.Result{data=_.GetNoticeChannelsRes}
// @Router /sys/notice/channels [GET]
func (s *NoticeApi) GetNoticeChannelsHandler(c *gin.Context) {
	data, err := noticeService.GetNoticeChannels(c)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// GetNoticeSendLogsHandler
// @Tags systemRbacnoticeApi
// @Summary GetNoticeSendLogsHandler 获取通知推送记录-后台使用
// @Description GetNoticeSendLogsHandler 获取通知推送记录-后台使用
// @Param data body req.GetNoticeSendLogsReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.GetNoticeSendLogsRes}
// @Router /sys/notice/:id/send-logs [GET]
func (s *NoticeApi) GetNoticeSendLogsHandler(c *gin.Context) {
	var req req.GetNoticeSendLogsReq
	// path 参数
	{
		val := c.Param("id")
		parsed, err := strconv.Atoi(val)
		if err != nil {
			c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
			return
		}
		req.Id = parsed

	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := noticeService.GetNoticeSendLogs(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// RetryNoticeSendHandler
// @Tags systemRbacnoticeApi
// @Summary RetryNoticeSendHandler 重试失败推送-后台使用
// @Description RetryNoticeSendHandler 重试失败推送-后台使用
// @Param data body req.RetryNoticeSendReq true "请求参数"
// @Success 200 {object} vo.Result{}
// @Router /sys/notice/:id/retry [POST]
func (s *NoticeApi) RetryNoticeSendHandler(c *gin.Context) {
	var req req.RetryNoticeSendReq
	// path 参数
	{
		val := c.Param("id")
		parsed, err := strconv.Atoi(val)
		if err != nil {
			c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
			return
		}
		req.Id = parsed

	}
	// query 参数
	{
		val := c.Query("channel")
		if val != "" {
			req.Channel = val

		}
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	err := noticeService.RetryNoticeSend(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, "操作成功"))
}