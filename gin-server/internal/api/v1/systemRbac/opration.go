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

type OprationApi struct{}

// ClearOperationRecordHandler
// @Tags systemRbacoprationApi
// @Summary ClearOperationRecordHandler 清空所有操作记录-后台使用
// @Description ClearOperationRecordHandler 清空所有操作记录-后台使用
// @Success 200 {object} vo.Result{}
// @Router /operation/clear [DELETE]
func (s *OprationApi) ClearOperationRecordHandler(c *gin.Context) {
	err := oprationService.ClearOperationRecord(c)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, "操作成功"))
}

// GetOperationStatisticsHandler
// @Tags systemRbacoprationApi
// @Summary GetOperationStatisticsHandler 获取操作统计信息-后台使用
// @Description GetOperationStatisticsHandler 获取操作统计信息-后台使用
// @Success 200 {object} vo.Result{data=_.GetOperationStatisticsRes}
// @Router /operation/statistics [GET]
func (s *OprationApi) GetOperationStatisticsHandler(c *gin.Context) {
	data, err := oprationService.GetOperationStatistics(c)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// GetOperationRecordsByUserIdHandler
// @Tags systemRbacoprationApi
// @Summary GetOperationRecordsByUserIdHandler 根据用户ID获取操作记录-后台使用
// @Description GetOperationRecordsByUserIdHandler 根据用户ID获取操作记录-后台使用
// @Param data body req.GetOperationRecordsByUserIdReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.GetOperationRecordsByUserIdRes}
// @Router /operation/user/:userId [GET]
func (s *OprationApi) GetOperationRecordsByUserIdHandler(c *gin.Context) {
	var req req.GetOperationRecordsByUserIdReq
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
	data, err := oprationService.GetOperationRecordsByUserId(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// GetOperationRecordsByTimeRangeHandler
// @Tags systemRbacoprationApi
// @Summary GetOperationRecordsByTimeRangeHandler 按时间范围查询操作记录-后台使用
// @Description GetOperationRecordsByTimeRangeHandler 按时间范围查询操作记录-后台使用
// @Param data body req.GetOperationRecordsByTimeRangeReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.GetOperationRecordsByTimeRangeRes}
// @Router /operation/range [GET]
func (s *OprationApi) GetOperationRecordsByTimeRangeHandler(c *gin.Context) {
	var req req.GetOperationRecordsByTimeRangeReq
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
	data, err := oprationService.GetOperationRecordsByTimeRange(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// GetRecentOperationRecordsHandler
// @Tags systemRbacoprationApi
// @Summary GetRecentOperationRecordsHandler 获取最近N条记录-后台使用
// @Description GetRecentOperationRecordsHandler 获取最近N条记录-后台使用
// @Param data body req.GetRecentOperationRecordsReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.GetRecentOperationRecordsRes}
// @Router /operation/recent [GET]
func (s *OprationApi) GetRecentOperationRecordsHandler(c *gin.Context) {
	var req req.GetRecentOperationRecordsReq
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
	data, err := oprationService.GetRecentOperationRecords(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// GetErrorRecordsHandler
// @Tags systemRbacoprationApi
// @Summary GetErrorRecordsHandler 获取异常请求记录-后台使用
// @Description GetErrorRecordsHandler 获取异常请求记录-后台使用
// @Param data body req.GetErrorRecordsReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.GetErrorRecordsRes}
// @Router /operation/errors [GET]
func (s *OprationApi) GetErrorRecordsHandler(c *gin.Context) {
	var req req.GetErrorRecordsReq
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
	data, err := oprationService.GetErrorRecords(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// DeleteExpiredRecordsHandler
// @Tags systemRbacoprationApi
// @Summary DeleteExpiredRecordsHandler 批量删除过期记录-后台使用
// @Description DeleteExpiredRecordsHandler 批量删除过期记录-后台使用
// @Param data body req.DeleteExpiredRecordsReq true "请求参数"
// @Success 200 {object} vo.Result{}
// @Router /operation/expired/delete [POST]
func (s *OprationApi) DeleteExpiredRecordsHandler(c *gin.Context) {
	var req req.DeleteExpiredRecordsReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	err := oprationService.DeleteExpiredRecords(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, "操作成功"))
}

// GetOperationRecordsByMethodHandler
// @Tags systemRbacoprationApi
// @Summary GetOperationRecordsByMethodHandler 按请求方法查询记录-后台使用
// @Description GetOperationRecordsByMethodHandler 按请求方法查询记录-后台使用
// @Param data body req.GetOperationRecordsByMethodReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.GetOperationRecordsByMethodRes}
// @Router /operation/method/:method [GET]
func (s *OprationApi) GetOperationRecordsByMethodHandler(c *gin.Context) {
	var req req.GetOperationRecordsByMethodReq
	// path 参数
	{
		val := c.Param("method")
		parsed := val
		var err error
		if err != nil {
			c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
			return
		}
		req.Method = parsed
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
	data, err := oprationService.GetOperationRecordsByMethod(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// GetOperationRecordsByPathHandler
// @Tags systemRbacoprationApi
// @Summary GetOperationRecordsByPathHandler 按请求路径查询记录-后台使用
// @Description GetOperationRecordsByPathHandler 按请求路径查询记录-后台使用
// @Param data body req.GetOperationRecordsByPathReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.GetOperationRecordsByPathRes}
// @Router /operation/path [GET]
func (s *OprationApi) GetOperationRecordsByPathHandler(c *gin.Context) {
	var req req.GetOperationRecordsByPathReq
	// query 参数
	{
		val := c.Query("path")
		if val != "" {
			parsed := val
			var err error
			if err != nil {
				c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
				return
			}
			req.Path = parsed
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
	data, err := oprationService.GetOperationRecordsByPath(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// CountTodayOperationsHandler
// @Tags systemRbacoprationApi
// @Summary CountTodayOperationsHandler 统计今日操作数-后台使用
// @Description CountTodayOperationsHandler 统计今日操作数-后台使用
// @Success 200 {object} vo.Result{data=_.CountTodayOperationsRes}
// @Router /operation/today/count [GET]
func (s *OprationApi) CountTodayOperationsHandler(c *gin.Context) {
	data, err := oprationService.CountTodayOperations(c)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}