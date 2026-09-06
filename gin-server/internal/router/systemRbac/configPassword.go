// router 解析
package systemRbac

import (
	"shack/internal/middleware"
	"shack/internal/model/common"
	"shack/internal/utils"

	"github.com/gin-gonic/gin"
)

type ConfigPasswordRouter struct{}

func (s *ConfigPasswordRouter) InitConfigPasswordRouter(Router *gin.RouterGroup, apiSet map[string]struct{}) (R gin.IRoutes) {
	configPasswordRouter := Router.Group("").Use(middleware.OperationRecord())
	utils.RegisterApi(configPasswordRouter, apiSet, "",
		utils.NewRegisterApiParam(common.HttpGet, "/api/config/password", "获取密码配置-后台使用", configPasswordApi.GetPasswordConfigHandler),
		utils.NewRegisterApiParam(common.HttpPost, "/api/config/password", "保存密码配置-后台使用", configPasswordApi.SavePasswordConfigHandler),
	)
	return configPasswordRouter
}
