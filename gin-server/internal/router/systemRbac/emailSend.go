//router 解析
package systemRbac

import (
	"github.com/gin-gonic/gin"
	"shack/internal/model/common"
	"shack/internal/utils"
)

type EmailSendRouter struct{}

func (s *EmailSendRouter) InitEmailSendRouter(Router *gin.RouterGroup, apiSet map[string]struct{}) (R gin.IRoutes) {
	emailSendRouter := Router.Group("")
	utils.RegisterApi(emailSendRouter, apiSet, "",
		utils.NewRegisterApiParam(common.HttpPost, "/email/send", "发送普通邮件-后台使用", emailSendApi.SendEmailHandler),
		utils.NewRegisterApiParam(common.HttpPost, "/email/send-template", "发送模板邮件-后台使用", emailSendApi.SendTemplateEmailHandler),
		utils.NewRegisterApiParam(common.HttpPost, "/email/send-batch", "批量发送邮件-后台使用", emailSendApi.SendBatchEmailHandler),
	)
	return emailSendRouter
}
