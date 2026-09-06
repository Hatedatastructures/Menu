package systemRbac

import (
	"strconv"

	req "shack/internal/model/systemRbac/request"
	_ "shack/internal/model/systemRbac/response"
	"shack/internal/utils/validator"
	"shack/internal/vo"

	"github.com/gin-gonic/gin"
)

type EmailLogApi struct{}

// GetEmailLogListHandler
// @Tags systemRbacemailLogApi
// @Summary GetEmailLogListHandler 获取邮件日志列表-后台使用
// @Description GetEmailLogListHandler 获取邮件日志列表-后台使用
// @Param data body req.GetEmailLogListReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.GetEmailLogListRes}
// @Router /email/log/list [GET]
func (s *EmailLogApi) GetEmailLogListHandler(c *gin.Context) {
	var req req.GetEmailLogListReq
	// query 参数
	{
		val := c.Query("page")
		if val != "" {
			parsed, err := strconv.Atoi(val)
			if err != nil {
				c.JSON(201, vo.Fail(c, "", err))
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
				c.JSON(201, vo.Fail(c, "", err))
				return
			}
			req.PageSize = parsed

		}
	}
	{
		val := c.Query("toEmail")
		if val != "" {
			req.ToEmail = val

		}
	}
	{
		val := c.Query("status")
		if val != "" {
			req.Status = val

		}
	}
	{
		val := c.Query("bizType")
		if val != "" {
			req.BizType = val

		}
	}
	{
		val := c.Query("startTime")
		if val != "" {
			req.StartTime = val

		}
	}
	{
		val := c.Query("endTime")
		if val != "" {
			req.EndTime = val

		}
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := emailLogService.GetEmailLogList(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// GetEmailLogHandler
// @Tags systemRbacemailLogApi
// @Summary GetEmailLogHandler 获取邮件日志详情-后台使用
// @Description GetEmailLogHandler 获取邮件日志详情-后台使用
// @Param data body req.GetEmailLogReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.GetEmailLogRes}
// @Router /email/log/:id [GET]
func (s *EmailLogApi) GetEmailLogHandler(c *gin.Context) {
	var req req.GetEmailLogReq
	// path 参数
	{
		val := c.Param("id")
		parsed, err := strconv.Atoi(val)
		if err != nil {
			c.JSON(201, vo.Fail(c, "", err))
			return
		}
		req.Id = parsed

	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := emailLogService.GetEmailLog(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// ResendEmailHandler
// @Tags systemRbacemailLogApi
// @Summary ResendEmailHandler 重发邮件-后台使用
// @Description ResendEmailHandler 重发邮件-后台使用
// @Param data body req.ResendEmailReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.ResendEmailRes}
// @Router /email/log/resend/:id [POST]
func (s *EmailLogApi) ResendEmailHandler(c *gin.Context) {
	var req req.ResendEmailReq
	// path 参数
	{
		val := c.Param("id")
		parsed, err := strconv.Atoi(val)
		if err != nil {
			c.JSON(201, vo.Fail(c, "", err))
			return
		}
		req.Id = parsed

	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := emailLogService.ResendEmail(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}