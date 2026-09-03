package systemRbac

import (
	req "shack/internal/model/systemRbac/request"
	_ "shack/internal/model/systemRbac/response"
	"shack/internal/utils/validator"
	"shack/internal/vo"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TemplateApi struct{}

// CreateTemplateHandler
// @Tags systemRbactemplateApi
// @Summary CreateTemplateHandler 创建模板-后台使用
// @Description CreateTemplateHandler 创建模板-后台使用
// @Param data body req.CreateTemplateReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.CreateTemplateRes}
// @Router /system/template [POST]
func (s *TemplateApi) CreateTemplateHandler(c *gin.Context) {
	var req req.CreateTemplateReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := templateService.CreateTemplate(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// DeleteTemplateHandler
// @Tags systemRbactemplateApi
// @Summary DeleteTemplateHandler 删除模板-后台使用
// @Description DeleteTemplateHandler 删除模板-后台使用
// @Param data body req.DeleteTemplateReq true "请求参数"
// @Success 200 {object} vo.Result{}
// @Router /system/template/:id [DELETE]
func (s *TemplateApi) DeleteTemplateHandler(c *gin.Context) {
	var req req.DeleteTemplateReq
	// path 参数
	{
		val := c.Param("id")
		parsed, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			c.JSON(201, vo.Fail(c, "", err))
			return
		}
		req.Id = uint(parsed)

	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	err := templateService.DeleteTemplate(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, "操作成功"))
}

// UpdateTemplateHandler
// @Tags systemRbactemplateApi
// @Summary UpdateTemplateHandler 更新模板-后台使用
// @Description UpdateTemplateHandler 更新模板-后台使用
// @Param data body req.UpdateTemplateReq true "请求参数"
// @Success 200 {object} vo.Result{}
// @Router /system/template/:id [PUT]
func (s *TemplateApi) UpdateTemplateHandler(c *gin.Context) {
	var req req.UpdateTemplateReq
	// path 参数
	{
		val := c.Param("id")
		parsed, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			c.JSON(201, vo.Fail(c, "", err))
			return
		}
		req.Id = uint(parsed)

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
	err := templateService.UpdateTemplate(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, "操作成功"))
}

// GetTemplateDetailHandler
// @Tags systemRbactemplateApi
// @Summary GetTemplateDetailHandler 获取模板详情-后台使用
// @Description GetTemplateDetailHandler 获取模板详情-后台使用
// @Param data body req.GetTemplateDetailReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.GetTemplateDetailRes}
// @Router /system/template/:id [GET]
func (s *TemplateApi) GetTemplateDetailHandler(c *gin.Context) {
	var req req.GetTemplateDetailReq
	// path 参数
	{
		val := c.Param("id")
		parsed, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			c.JSON(201, vo.Fail(c, "", err))
			return
		}
		req.Id = uint(parsed)

	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := templateService.GetTemplateDetail(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// GetTemplateListHandler
// @Tags systemRbactemplateApi
// @Summary GetTemplateListHandler 获取模板分页列表-后台使用
// @Description GetTemplateListHandler 获取模板分页列表-后台使用
// @Param data body req.GetTemplateListReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.GetTemplateListRes}
// @Router /system/template/list [GET]
func (s *TemplateApi) GetTemplateListHandler(c *gin.Context) {
	var req req.GetTemplateListReq
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
		val := c.Query("keyword")
		if val != "" {
			req.Keyword = val

		}
	}
	{
		val := c.Query("status")
		if val != "" {
			req.Status = val

		}
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := templateService.GetTemplateList(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}
