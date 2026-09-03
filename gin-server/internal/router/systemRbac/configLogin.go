// router 解析
package systemRbac

import (
	"shack/internal/middleware"
	"shack/internal/model/common"
	"shack/internal/utils"

	"github.com/gin-gonic/gin"
)

type ConfigLoginRouter struct{}

func (s *ConfigLoginRouter) InitConfigLoginRouter(Router *gin.RouterGroup, apiSet map[string]struct{}) (R gin.IRoutes) {
	configLoginRouter := Router.Group("").Use(middleware.OperationRecord())
	utils.RegisterApi(configLoginRouter, apiSet, "",
		utils.NewRegisterApiParam(common.HttpGet, "/api/config/login", "获取登录配置-后台使用", configLoginApi.GetLoginConfigHandler),
		utils.NewRegisterApiParam(common.HttpPost, "/api/config/login", "保存登录配置-后台使用", configLoginApi.SaveLoginConfigHandler),
	)
	return configLoginRouter
}
