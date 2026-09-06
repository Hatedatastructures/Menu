//router 解析
package systemRbac

import (
	"github.com/gin-gonic/gin"
	"shack/internal/model/common"
	"shack/internal/utils"
)

type EmailLimitRouter struct{}

func (s *EmailLimitRouter) InitEmailLimitRouter(Router *gin.RouterGroup, apiSet map[string]struct{}) (R gin.IRoutes) {
	emailLimitRouter := Router.Group("")
	utils.RegisterApi(emailLimitRouter, apiSet, "",
		utils.NewRegisterApiParam(common.HttpGet, "/email/limit/log", "获取限流日志-后台使用", emailLimitApi.GetEmailLimitLogHandler),
		utils.NewRegisterApiParam(common.HttpPost, "/email/blacklist", "添加黑名单-后台使用", emailLimitApi.AddEmailBlacklistHandler),
		utils.NewRegisterApiParam(common.HttpDelete, "/email/blacklist/:id", "删除黑名单-后台使用", emailLimitApi.DeleteEmailBlacklistHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/email/blacklist/list", "获取黑名单列表-后台使用", emailLimitApi.GetEmailBlacklistListHandler),
	)
	return emailLimitRouter
}
