package systemRbac

import (
	"strconv"

	req "shack/internal/model/systemRbac/request"
	_ "shack/internal/model/systemRbac/response"
	"shack/internal/utils/validator"
	"shack/internal/vo"

	"github.com/gin-gonic/gin"
)

type EmailLimitApi struct{}

// GetEmailLimitLogHandler
// @Tags systemRbacemailLimitApi
// @Summary GetEmailLimitLogHandler 获取限流日志-后台使用
// @Description GetEmailLimitLogHandler 获取限流日志-后台使用
// @Param data body req.GetEmailLimitLogReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.GetEmailLimitLogRes}
// @Router /email/limit/log [GET]
func (s *EmailLimitApi) GetEmailLimitLogHandler(c *gin.Context) {
	var req req.GetEmailLimitLogReq
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
		val := c.Query("email")
		if val != "" {
			req.Email = val

		}
	}
	{
		val := c.Query("ip")
		if val != "" {
			req.Ip = val

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
	data, err := emailLimitService.GetEmailLimitLog(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// AddEmailBlacklistHandler
// @Tags systemRbacemailLimitApi
// @Summary AddEmailBlacklistHandler 添加黑名单-后台使用
// @Description AddEmailBlacklistHandler 添加黑名单-后台使用
// @Param data body req.AddEmailBlacklistReq true "请求参数"
// @Success 200 {object} vo.Result{}
// @Router /email/blacklist [POST]
func (s *EmailLimitApi) AddEmailBlacklistHandler(c *gin.Context) {
	var req req.AddEmailBlacklistReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	err := emailLimitService.AddEmailBlacklist(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, "操作成功"))
}

// DeleteEmailBlacklistHandler
// @Tags systemRbacemailLimitApi
// @Summary DeleteEmailBlacklistHandler 删除黑名单-后台使用
// @Description DeleteEmailBlacklistHandler 删除黑名单-后台使用
// @Param data body req.DeleteEmailBlacklistReq true "请求参数"
// @Success 200 {object} vo.Result{}
// @Router /email/blacklist/:id [DELETE]
func (s *EmailLimitApi) DeleteEmailBlacklistHandler(c *gin.Context) {
	var req req.DeleteEmailBlacklistReq
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
	err := emailLimitService.DeleteEmailBlacklist(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, "操作成功"))
}

// GetEmailBlacklistListHandler
// @Tags systemRbacemailLimitApi
// @Summary GetEmailBlacklistListHandler 获取黑名单列表-后台使用
// @Description GetEmailBlacklistListHandler 获取黑名单列表-后台使用
// @Param data body req.GetEmailBlacklistListReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.GetEmailBlacklistListRes}
// @Router /email/blacklist/list [GET]
func (s *EmailLimitApi) GetEmailBlacklistListHandler(c *gin.Context) {
	var req req.GetEmailBlacklistListReq
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

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := emailLimitService.GetEmailBlacklistList(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}