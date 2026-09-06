// router 解析
package systemRbac

import (
	"shack/internal/middleware"
	"shack/internal/model/common"
	"shack/internal/utils"

	"github.com/gin-gonic/gin"
)

type ConfigEmailRouter struct{}

func (s *ConfigEmailRouter) InitConfigEmailRouter(Router *gin.RouterGroup, apiSet map[string]struct{}) (R gin.IRoutes) {
	configEmailRouter := Router.Group("").Use(middleware.OperationRecord())
	utils.RegisterApi(configEmailRouter, apiSet, "",
		utils.NewRegisterApiParam(common.HttpGet, "/api/config/email", "获取邮件配置-后台使用", configEmailApi.GetEmailConfigHandler),
		utils.NewRegisterApiParam(common.HttpPost, "/api/config/email", "保存邮件配置-后台使用", configEmailApi.SaveEmailConfigHandler),
		utils.NewRegisterApiParam(common.HttpPost, "/api/config/email/test", "测试邮件发送-后台使用", configEmailApi.TestEmailHandler),
	)
	return configEmailRouter
}
