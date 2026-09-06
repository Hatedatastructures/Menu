package systemRbac

import (
	"strconv"

	req "shack/internal/model/systemRbac/request"
	_ "shack/internal/model/systemRbac/response"
	"shack/internal/utils/validator"
	"shack/internal/vo"

	"github.com/gin-gonic/gin"
)

type EmailTemplateApi struct{}

// CreateEmailTemplateHandler
// @Tags systemRbacemailTemplateApi
// @Summary CreateEmailTemplateHandler 创建邮件模板-后台使用
// @Description CreateEmailTemplateHandler 创建邮件模板-后台使用
// @Param data body req.CreateEmailTemplateReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.CreateEmailTemplateRes}
// @Router /email/template [POST]
func (s *EmailTemplateApi) CreateEmailTemplateHandler(c *gin.Context) {
	var req req.CreateEmailTemplateReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := emailTemplateService.CreateEmailTemplate(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// UpdateEmailTemplateHandler
// @Tags systemRbacemailTemplateApi
// @Summary UpdateEmailTemplateHandler 更新邮件模板-后台使用
// @Description UpdateEmailTemplateHandler 更新邮件模板-后台使用
// @Param data body req.UpdateEmailTemplateReq true "请求参数"
// @Success 200 {object} vo.Result{}
// @Router /email/template/:id [PUT]
func (s *EmailTemplateApi) UpdateEmailTemplateHandler(c *gin.Context) {
	var req req.UpdateEmailTemplateReq
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
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	err := emailTemplateService.UpdateEmailTemplate(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, "操作成功"))
}

// DeleteEmailTemplateHandler
// @Tags systemRbacemailTemplateApi
// @Summary DeleteEmailTemplateHandler 删除邮件模板-后台使用
// @Description DeleteEmailTemplateHandler 删除邮件模板-后台使用
// @Param data body req.DeleteEmailTemplateReq true "请求参数"
// @Success 200 {object} vo.Result{}
// @Router /email/template/:id [DELETE]
func (s *EmailTemplateApi) DeleteEmailTemplateHandler(c *gin.Context) {
	var req req.DeleteEmailTemplateReq
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
	err := emailTemplateService.DeleteEmailTemplate(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, "操作成功"))
}

// GetEmailTemplateListHandler
// @Tags systemRbacemailTemplateApi
// @Summary GetEmailTemplateListHandler 获取邮件模板列表-后台使用
// @Description GetEmailTemplateListHandler 获取邮件模板列表-后台使用
// @Param data body req.GetEmailTemplateListReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.GetEmailTemplateListRes}
// @Router /email/template/list [GET]
func (s *EmailTemplateApi) GetEmailTemplateListHandler(c *gin.Context) {
	var req req.GetEmailTemplateListReq
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
		val := c.Query("type")
		if val != "" {
			req.Type = val

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
	data, err := emailTemplateService.GetEmailTemplateList(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// GetEmailTemplateHandler
// @Tags systemRbacemailTemplateApi
// @Summary GetEmailTemplateHandler 获取邮件模板详情-后台使用
// @Description GetEmailTemplateHandler 获取邮件模板详情-后台使用
// @Param data body req.GetEmailTemplateReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.GetEmailTemplateRes}
// @Router /email/template/:id [GET]
func (s *EmailTemplateApi) GetEmailTemplateHandler(c *gin.Context) {
	var req req.GetEmailTemplateReq
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
	data, err := emailTemplateService.GetEmailTemplate(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// PreviewEmailTemplateHandler
// @Tags systemRbacemailTemplateApi
// @Summary PreviewEmailTemplateHandler 预览邮件模板-后台使用
// @Description PreviewEmailTemplateHandler 预览邮件模板-后台使用
// @Param data body req.PreviewEmailTemplateReq true "请求参数"
// @Success 200 {object} vo.Result{}
// @Router /email/template/preview [POST]
func (s *EmailTemplateApi) PreviewEmailTemplateHandler(c *gin.Context) {
	var req req.PreviewEmailTemplateReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	err := emailTemplateService.PreviewEmailTemplate(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, "操作成功"))
}

// AiGenerateEmailTemplateHandler
// @Tags systemRbacemailTemplateApi
// @Summary AiGenerateEmailTemplateHandler AI生成邮件模板-后台使用
// @Description AiGenerateEmailTemplateHandler AI生成邮件模板-后台使用
// @Param data body req.AiGenerateEmailTemplateReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.AiGenerateEmailTemplateRes}
// @Router /email/template/ai-generate [POST]
func (s *EmailTemplateApi) AiGenerateEmailTemplateHandler(c *gin.Context) {
	var req req.AiGenerateEmailTemplateReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := emailTemplateService.AiGenerateEmailTemplate(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}
