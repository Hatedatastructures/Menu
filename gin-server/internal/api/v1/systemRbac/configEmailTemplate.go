package systemRbac

import (
	biz_err "shack/internal/error"
	req "shack/internal/model/systemRbac/request"
	_ "shack/internal/model/systemRbac/response"
	"shack/internal/utils/validator"
	"shack/internal/vo"

	"github.com/gin-gonic/gin"
)

type ConfigEmailTemplateApi struct{}

// GetEmailTemplateConfigHandler
// @Tags systemRbacconfigEmailTemplateApi
// @Summary GetEmailTemplateConfigHandler 获取邮件模板配置-后台使用
// @Description GetEmailTemplateConfigHandler 获取邮件模板配置-后台使用
// @Success 200 {object} vo.Result{data=_.GetEmailTemplateConfigRes}
// @Router /api/config/emailTemplate [GET]
func (s *ConfigEmailTemplateApi) GetEmailTemplateConfigHandler(c *gin.Context) {
	data, err := configEmailTemplateService.GetEmailTemplateConfig(c)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// SaveEmailTemplateConfigHandler
// @Tags systemRbacconfigEmailTemplateApi
// @Summary SaveEmailTemplateConfigHandler 保存邮件模板配置-后台使用
// @Description SaveEmailTemplateConfigHandler 保存邮件模板配置-后台使用
// @Param data body req.SaveEmailTemplateConfigReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.SaveEmailTemplateConfigRes}
// @Router /api/config/emailTemplate [POST]
func (s *ConfigEmailTemplateApi) SaveEmailTemplateConfigHandler(c *gin.Context) {
	var req req.SaveEmailTemplateConfigReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := configEmailTemplateService.SaveEmailTemplateConfig(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}
