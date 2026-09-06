//router 解析
package menu

import (
	"github.com/gin-gonic/gin"
	"shack/internal/model/common"
	"shack/internal/utils"
)

type HealthRouter struct{}

func (s *HealthRouter) InitHealthRouter(Router *gin.RouterGroup, apiSet map[string]struct{}) (R gin.IRoutes) {
	healthRouter := Router.Group("")
	utils.RegisterApi(healthRouter, apiSet, "",
		utils.NewRegisterApiParam(common.HttpGet, "/healthz", "系统健康检查模块 对应 C++ server /healthz 和 /readyz 健康检查", healthApi.HealthzHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/readyz", "就绪检查", healthApi.ReadyzHandler),
	)
	return healthRouter
}
