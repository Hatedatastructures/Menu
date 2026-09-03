//router 解析
package systemRbac

import (
	"github.com/gin-gonic/gin"
	"shack/internal/model/common"
	"shack/internal/utils"
)

type CasbinRouter struct{}

func (s *CasbinRouter) InitCasbinRouter(Router *gin.RouterGroup, apiSet map[string]struct{}) (R gin.IRoutes) {
	casbinRouter := Router.Group("")
	utils.RegisterApi(casbinRouter, apiSet, "",
		utils.NewRegisterApiParam(common.HttpPost, "/api/casbin", "更新权限-后台使用", casbinApi.UpdateCasbinHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/api/casbin/:authorityId", "获取权限列表-后台使用", casbinApi.GetPolicyPathByAuthorityIdHandler),
		utils.NewRegisterApiParam(common.HttpDelete, "/api/casbin", "清除权限-后台使用", casbinApi.ClearCasbinHandler),
		utils.NewRegisterApiParam(common.HttpPost, "/api/casbin/refresh", "刷新Casbin-后台使用", casbinApi.FreshCasbinHandler),
	)
	return casbinRouter
}
