//router 解析
package systemRbac

import (
	"github.com/gin-gonic/gin"
	"shack/internal/model/common"
	"shack/internal/utils"
)

type EmailLogRouter struct{}

func (s *EmailLogRouter) InitEmailLogRouter(Router *gin.RouterGroup, apiSet map[string]struct{}) (R gin.IRoutes) {
	emailLogRouter := Router.Group("")
	utils.RegisterApi(emailLogRouter, apiSet, "",
		utils.NewRegisterApiParam(common.HttpGet, "/email/log/list", "获取邮件日志列表-后台使用", emailLogApi.GetEmailLogListHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/email/log/:id", "获取邮件日志详情-后台使用", emailLogApi.GetEmailLogHandler),
		utils.NewRegisterApiParam(common.HttpPost, "/email/log/resend/:id", "重发邮件-后台使用", emailLogApi.ResendEmailHandler),
	)
	return emailLogRouter
}
