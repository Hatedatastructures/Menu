// router 解析
package systemRbac

import (
	"shack/internal/middleware"
	"shack/internal/model/common"
	"shack/internal/utils"

	"github.com/gin-gonic/gin"
)

type ConfigSecurityRouter struct{}

func (s *ConfigSecurityRouter) InitConfigSecurityRouter(Router *gin.RouterGroup, apiSet map[string]struct{}) (R gin.IRoutes) {
	configSecurityRouter := Router.Group("").Use(middleware.OperationRecord())
	utils.RegisterApi(configSecurityRouter, apiSet, "",
		utils.NewRegisterApiParam(common.HttpGet, "/api/config/security", "获取安全配置-后台使用", configSecurityApi.GetSecurityConfigHandler),
		utils.NewRegisterApiParam(common.HttpPost, "/api/config/security", "保存安全配置-后台使用", configSecurityApi.SaveSecurityConfigHandler),
		utils.NewRegisterApiParam(common.HttpPost, "/api/config/security/generate-keys", "生成RSA密钥对-后台使用", configSecurityApi.GenerateRSAKeysHandler),
	)
	return configSecurityRouter
}
