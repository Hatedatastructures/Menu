package systemRbac

import (
	req "shack/internal/model/systemRbac/request"
	_ "shack/internal/model/systemRbac/response"
	"shack/internal/utils/validator"
	"shack/internal/vo"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ApiTestApi struct{}

// SaveApiTestLogHandler
// @Tags systemRbacapiTestApi
// @Summary SaveApiTestLogHandler 保存API测试请求记录-前台使用
// @Description SaveApiTestLogHandler 保存API测试请求记录-前台使用
// @Param data body req.SaveApiTestLogReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.SaveApiTestLogRes}
// @Router /api-test/log [POST]
func (s *ApiTestApi) SaveApiTestLogHandler(c *gin.Context) {
	var req req.SaveApiTestLogReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := apiTestService.SaveApiTestLog(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// GetApiTestLogListHandler
// @Tags systemRbacapiTestApi
// @Summary GetApiTestLogListHandler 获取API测试记录列表-前台使用
// @Description GetApiTestLogListHandler 获取API测试记录列表-前台使用
// @Param data body req.GetApiTestLogListReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.GetApiTestLogListRes}
// @Router /api-test/log/list [GET]
func (s *ApiTestApi) GetApiTestLogListHandler(c *gin.Context) {
	var req req.GetApiTestLogListReq
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
		val := c.Query("keyword")
		if val != "" {
			req.Keyword = val

		}
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := apiTestService.GetApiTestLogList(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// GetApiTestLogDetailHandler
// @Tags systemRbacapiTestApi
// @Summary GetApiTestLogDetailHandler 获取API测试记录详情-前台使用
// @Description GetApiTestLogDetailHandler 获取API测试记录详情-前台使用
// @Param data body req.GetApiTestLogDetailReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.GetApiTestLogDetailRes}
// @Router /api-test/log/:id [GET]
func (s *ApiTestApi) GetApiTestLogDetailHandler(c *gin.Context) {
	var req req.GetApiTestLogDetailReq
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
	data, err := apiTestService.GetApiTestLogDetail(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// DeleteApiTestLogHandler
// @Tags systemRbacapiTestApi
// @Summary DeleteApiTestLogHandler 删除API测试记录-前台使用
// @Description DeleteApiTestLogHandler 删除API测试记录-前台使用
// @Param data body req.DeleteApiTestLogReq true "请求参数"
// @Success 200 {object} vo.Result{}
// @Router /api-test/log/:id [DELETE]
func (s *ApiTestApi) DeleteApiTestLogHandler(c *gin.Context) {
	var req req.DeleteApiTestLogReq
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
	err := apiTestService.DeleteApiTestLog(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, "操作成功"))
}

// ClearApiTestLogHandler
// @Tags systemRbacapiTestApi
// @Summary ClearApiTestLogHandler 清空API测试记录-前台使用
// @Description ClearApiTestLogHandler 清空API测试记录-前台使用
// @Success 200 {object} vo.Result{}
// @Router /api-test/log/clear [DELETE]
func (s *ApiTestApi) ClearApiTestLogHandler(c *gin.Context) {
	err := apiTestService.ClearApiTestLog(c)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, "操作成功"))
}
