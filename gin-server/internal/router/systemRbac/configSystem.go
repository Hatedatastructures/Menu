// router 解析
package systemRbac

import (
	"shack/internal/middleware"
	"shack/internal/model/common"
	"shack/internal/utils"

	"github.com/gin-gonic/gin"
)

type ConfigSystemRouter struct{}

func (s *ConfigSystemRouter) InitConfigSystemRouter(Router *gin.RouterGroup, apiSet map[string]struct{}) (R gin.IRoutes) {
	configSystemRouter := Router.Group("").Use(middleware.OperationRecord())
	utils.RegisterApi(configSystemRouter, apiSet, "",
		utils.NewRegisterApiParam(common.HttpGet, "/api/config/system", "获取系统配置-后台使用", configSystemApi.GetSystemConfigHandler),
		utils.NewRegisterApiParam(common.HttpPost, "/api/config/system", "保存系统配置-后台使用", configSystemApi.SaveSystemConfigHandler),
	)
	return configSystemRouter
}
