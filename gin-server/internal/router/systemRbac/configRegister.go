// router 解析
package systemRbac

import (
	"shack/internal/middleware"
	"shack/internal/model/common"
	"shack/internal/utils"

	"github.com/gin-gonic/gin"
)

type ConfigRegisterRouter struct{}

func (s *ConfigRegisterRouter) InitConfigRegisterRouter(Router *gin.RouterGroup, apiSet map[string]struct{}) (R gin.IRoutes) {
	configRegisterRouter := Router.Group("").Use(middleware.OperationRecord())
	utils.RegisterApi(configRegisterRouter, apiSet, "",
		utils.NewRegisterApiParam(common.HttpGet, "/api/config/register", "获取注册配置-后台使用", configRegisterApi.GetRegisterConfigHandler),
		utils.NewRegisterApiParam(common.HttpPost, "/api/config/register", "保存注册配置-后台使用", configRegisterApi.SaveRegisterConfigHandler),
	)
	return configRegisterRouter
}
