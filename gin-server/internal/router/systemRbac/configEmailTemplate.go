// router 解析
package systemRbac

import (
	"shack/internal/middleware"
	"shack/internal/model/common"
	"shack/internal/utils"

	"github.com/gin-gonic/gin"
)

type ConfigEmailTemplateRouter struct{}

func (s *ConfigEmailTemplateRouter) InitConfigEmailTemplateRouter(Router *gin.RouterGroup, apiSet map[string]struct{}) (R gin.IRoutes) {
	configEmailTemplateRouter := Router.Group("").Use(middleware.OperationRecord())
	utils.RegisterApi(configEmailTemplateRouter, apiSet, "",
		utils.NewRegisterApiParam(common.HttpGet, "/api/config/emailTemplate", "获取邮件模板配置-后台使用", configEmailTemplateApi.GetEmailTemplateConfigHandler),
		utils.NewRegisterApiParam(common.HttpPost, "/api/config/emailTemplate", "保存邮件模板配置-后台使用", configEmailTemplateApi.SaveEmailTemplateConfigHandler),
	)
	return configEmailTemplateRouter
}
