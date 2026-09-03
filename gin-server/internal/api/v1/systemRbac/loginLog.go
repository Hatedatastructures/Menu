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

type LoginLogApi struct{}

// GetLoginLogListHandler
// @Tags systemRbacloginLogApi
// @Summary GetLoginLogListHandler 获取登录日志列表-后台使用
// @Description GetLoginLogListHandler 获取登录日志列表-后台使用
// @Param data body req.GetLoginLogListReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.GetLoginLogListRes}
// @Router /loginLog [GET]
func (s *LoginLogApi) GetLoginLogListHandler(c *gin.Context) {
	var req req.GetLoginLogListReq
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
		val := c.Query("size")
		if val != "" {
			parsed, err := strconv.Atoi(val)
			if err != nil {
				c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
				return
			}
			req.Size = parsed
		}
	}
	{
		val := c.Query("username")
		if val != "" {
			parsed := val
			var err error
			if err != nil {
				c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
				return
			}
			req.Username = parsed
		}
	}
	{
		val := c.Query("status")
		if val != "" {
			parsed := val
			var err error
			if err != nil {
				c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
				return
			}
			req.Status = parsed
		}
	}
	{
		val := c.Query("startTime")
		if val != "" {
			parsed := val
			var err error
			if err != nil {
				c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
				return
			}
			req.StartTime = parsed
		}
	}
	{
		val := c.Query("endTime")
		if val != "" {
			parsed := val
			var err error
			if err != nil {
				c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
				return
			}
			req.EndTime = parsed
		}
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := loginLogService.GetLoginLogList(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// GetLoginLogByUsernameHandler
// @Tags systemRbacloginLogApi
// @Summary GetLoginLogByUsernameHandler 根据用户名查询登录日志-后台使用
// @Description GetLoginLogByUsernameHandler 根据用户名查询登录日志-后台使用
// @Param data body req.GetLoginLogByUsernameReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.GetLoginLogByUsernameRes}
// @Router /loginLog/user/:username [GET]
func (s *LoginLogApi) GetLoginLogByUsernameHandler(c *gin.Context) {
	var req req.GetLoginLogByUsernameReq
	// path 参数
	{
		val := c.Param("username")
		parsed := val
		var err error
		if err != nil {
			c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
			return
		}
		req.Username = parsed
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
		val := c.Query("size")
		if val != "" {
			parsed, err := strconv.Atoi(val)
			if err != nil {
				c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
				return
			}
			req.Size = parsed
		}
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := loginLogService.GetLoginLogByUsername(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// GetLoginLogByStatusHandler
// @Tags systemRbacloginLogApi
// @Summary GetLoginLogByStatusHandler 根据状态查询登录日志-后台使用
// @Description GetLoginLogByStatusHandler 根据状态查询登录日志-后台使用
// @Param data body req.GetLoginLogByStatusReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.GetLoginLogByStatusRes}
// @Router /loginLog/status/:status [GET]
func (s *LoginLogApi) GetLoginLogByStatusHandler(c *gin.Context) {
	var req req.GetLoginLogByStatusReq
	// path 参数
	{
		val := c.Param("status")
		parsed := val
		var err error
		if err != nil {
			c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
			return
		}
		req.Status = parsed
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
		val := c.Query("size")
		if val != "" {
			parsed, err := strconv.Atoi(val)
			if err != nil {
				c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
				return
			}
			req.Size = parsed
		}
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := loginLogService.GetLoginLogByStatus(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// GetLoginLogByTimeRangeHandler
// @Tags systemRbacloginLogApi
// @Summary GetLoginLogByTimeRangeHandler 按时间范围查询登录日志-后台使用
// @Description GetLoginLogByTimeRangeHandler 按时间范围查询登录日志-后台使用
// @Param data body req.GetLoginLogByTimeRangeReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.GetLoginLogByTimeRangeRes}
// @Router /loginLog/range [GET]
func (s *LoginLogApi) GetLoginLogByTimeRangeHandler(c *gin.Context) {
	var req req.GetLoginLogByTimeRangeReq
	// query 参数
	{
		val := c.Query("startTime")
		if val != "" {
			parsed := val
			var err error
			if err != nil {
				c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
				return
			}
			req.StartTime = parsed
		}
	}
	{
		val := c.Query("endTime")
		if val != "" {
			parsed := val
			var err error
			if err != nil {
				c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
				return
			}
			req.EndTime = parsed
		}
	}
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
		val := c.Query("size")
		if val != "" {
			parsed, err := strconv.Atoi(val)
			if err != nil {
				c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
				return
			}
			req.Size = parsed
		}
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := loginLogService.GetLoginLogByTimeRange(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// GetRecentLoginLogHandler
// @Tags systemRbacloginLogApi
// @Summary GetRecentLoginLogHandler 获取最近N条登录记录-后台使用
// @Description GetRecentLoginLogHandler 获取最近N条登录记录-后台使用
// @Param data body req.GetRecentLoginLogReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.GetRecentLoginLogRes}
// @Router /loginLog/recent [GET]
func (s *LoginLogApi) GetRecentLoginLogHandler(c *gin.Context) {
	var req req.GetRecentLoginLogReq
	// query 参数
	{
		val := c.Query("limit")
		if val != "" {
			parsed, err := strconv.Atoi(val)
			if err != nil {
				c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
				return
			}
			req.Limit = parsed
		}
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := loginLogService.GetRecentLoginLog(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// GetLoginStatisticsHandler
// @Tags systemRbacloginLogApi
// @Summary GetLoginStatisticsHandler 获取登录统计信息-后台使用
// @Description GetLoginStatisticsHandler 获取登录统计信息-后台使用
// @Success 200 {object} vo.Result{data=_.GetLoginStatisticsRes}
// @Router /loginLog/statistics [GET]
func (s *LoginLogApi) GetLoginStatisticsHandler(c *gin.Context) {
	data, err := loginLogService.GetLoginStatistics(c)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// ClearLoginLogHandler
// @Tags systemRbacloginLogApi
// @Summary ClearLoginLogHandler 清空登录日志-后台使用
// @Description ClearLoginLogHandler 清空登录日志-后台使用
// @Success 200 {object} vo.Result{}
// @Router /loginLog/clear [DELETE]
func (s *LoginLogApi) ClearLoginLogHandler(c *gin.Context) {
	err := loginLogService.ClearLoginLog(c)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, "操作成功"))
}

// DeleteLoginLogHandler
// @Tags systemRbacloginLogApi
// @Summary DeleteLoginLogHandler 批量删除登录日志-后台使用
// @Description DeleteLoginLogHandler 批量删除登录日志-后台使用
// @Param data body req.DeleteLoginLogReq true "请求参数"
// @Success 200 {object} vo.Result{}
// @Router /loginLog/delete [POST]
func (s *LoginLogApi) DeleteLoginLogHandler(c *gin.Context) {
	var req req.DeleteLoginLogReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	err := loginLogService.DeleteLoginLog(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, "操作成功"))
}