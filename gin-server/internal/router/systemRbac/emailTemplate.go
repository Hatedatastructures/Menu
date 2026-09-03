//router 解析
package systemRbac

import (
	"github.com/gin-gonic/gin"
	"shack/internal/model/common"
	"shack/internal/utils"
)

type EmailTemplateRouter struct{}

func (s *EmailTemplateRouter) InitEmailTemplateRouter(Router *gin.RouterGroup, apiSet map[string]struct{}) (R gin.IRoutes) {
	emailTemplateRouter := Router.Group("")
	utils.RegisterApi(emailTemplateRouter, apiSet, "",
		utils.NewRegisterApiParam(common.HttpPost, "/email/template", "创建邮件模板-后台使用", emailTemplateApi.CreateEmailTemplateHandler),
		utils.NewRegisterApiParam(common.HttpPut, "/email/template/:id", "更新邮件模板-后台使用", emailTemplateApi.UpdateEmailTemplateHandler),
		utils.NewRegisterApiParam(common.HttpDelete, "/email/template/:id", "删除邮件模板-后台使用", emailTemplateApi.DeleteEmailTemplateHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/email/template/list", "获取邮件模板列表-后台使用", emailTemplateApi.GetEmailTemplateListHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/email/template/:id", "获取邮件模板详情-后台使用", emailTemplateApi.GetEmailTemplateHandler),
		utils.NewRegisterApiParam(common.HttpPost, "/email/template/preview", "预览邮件模板-后台使用", emailTemplateApi.PreviewEmailTemplateHandler),
		utils.NewRegisterApiParam(common.HttpPost, "/email/template/ai-generate", "AI生成邮件模板-后台使用", emailTemplateApi.AiGenerateEmailTemplateHandler),
	)
	return emailTemplateRouter
}
