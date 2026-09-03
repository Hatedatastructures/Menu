package systemRbac

import (
	"strconv"

	req "shack/internal/model/systemRbac/request"
	_ "shack/internal/model/systemRbac/response"
	"shack/internal/utils/validator"
	"shack/internal/vo"

	"github.com/gin-gonic/gin"
)

type ApiApi struct{}

// CreateApiHandler
// @Tags systemRbacapiApi
// @Summary CreateApiHandler 创建API-后台使用
// @Description CreateApiHandler 创建API-后台使用
// @Param data body req.CreateApiReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.CreateApiRes}
// @Router /api/apis [POST]
func (s *ApiApi) CreateApiHandler(c *gin.Context) {
	var req req.CreateApiReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := apiService.CreateApi(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// GetApiListHandler
// @Tags systemRbacapiApi
// @Summary GetApiListHandler API列表-后台使用
// @Description GetApiListHandler API列表-后台使用
// @Param data body req.GetApiListReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.GetApiListRes}
// @Router /api/apis [GET]
func (s *ApiApi) GetApiListHandler(c *gin.Context) {
	var req req.GetApiListReq
	// query 参数
	{
		val := c.Query("keyword")
		if val != "" {
			parsed := val
			var err error
			if err != nil {
				c.JSON(201, vo.Fail(c, "", err))
				return
			}
			req.Keyword = parsed
		}
	}
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
		val := c.Query("size")
		if val != "" {
			parsed, err := strconv.Atoi(val)
			if err != nil {
				c.JSON(201, vo.Fail(c, "", err))
				return
			}
			req.Size = parsed
		}
	}
	{
		val := c.Query("method")
		if val != "" {
			parsed := val
			var err error
			if err != nil {
				c.JSON(201, vo.Fail(c, "", err))
				return
			}
			req.Method = parsed
		}
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := apiService.GetApiList(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// GetApiDetailHandler
// @Tags systemRbacapiApi
// @Summary GetApiDetailHandler API详情-后台使用
// @Description GetApiDetailHandler API详情-后台使用
// @Param data body req.GetApiDetailReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.GetApiDetailRes}
// @Router /api/apis/:id [GET]
func (s *ApiApi) GetApiDetailHandler(c *gin.Context) {
	var req req.GetApiDetailReq
	// path 参数
	{
		val := c.Param("id")
		parsed := val
		var err error
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
	data, err := apiService.GetApiDetail(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// UpdateApiHandler
// @Tags systemRbacapiApi
// @Summary UpdateApiHandler 更新API-后台使用
// @Description UpdateApiHandler 更新API-后台使用
// @Param data body req.UpdateApiReq true "请求参数"
// @Success 200 {object} vo.Result{}
// @Router /api/apis/:id [PUT]
func (s *ApiApi) UpdateApiHandler(c *gin.Context) {
	var req req.UpdateApiReq
	// path 参数
	{
		val := c.Param("id")
		parsed := val
		var err error
		if err != nil {
			c.JSON(201, vo.Fail(c, "", err))
			return
		}
		req.Id = parsed
	}
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	err := apiService.UpdateApi(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, "操作成功"))
}

// DeleteApiHandler
// @Tags systemRbacapiApi
// @Summary DeleteApiHandler 删除API-后台使用
// @Description DeleteApiHandler 删除API-后台使用
// @Param data body req.DeleteApiReq true "请求参数"
// @Success 200 {object} vo.Result{}
// @Router /api/apis/:id [DELETE]
func (s *ApiApi) DeleteApiHandler(c *gin.Context) {
	var req req.DeleteApiReq
	// path 参数
	{
		val := c.Param("id")
		parsed := val
		var err error
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
	err := apiService.DeleteApi(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, "操作成功"))
}
