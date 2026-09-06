//router 解析
package systemRbac

import (
	"github.com/gin-gonic/gin"
	"shack/internal/model/common"
	"shack/internal/utils"
)

type SysMenuRouter struct{}

func (s *SysMenuRouter) InitSysMenuRouter(Router *gin.RouterGroup, apiSet map[string]struct{}) (R gin.IRoutes) {
	sysMenuRouter := Router.Group("")
	utils.RegisterApi(sysMenuRouter, apiSet, "",
		utils.NewRegisterApiParam(common.HttpGet, "/menu/tree/:authorityId", "获取动态菜单树-后台使用", sysMenuApi.GetMenuTreeHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/menu/list/:authorityId", "获取基础菜单列表(分页)-后台使用", sysMenuApi.GetInfoListHandler),
		utils.NewRegisterApiParam(common.HttpPost, "/menu/create", "添加基础菜单-后台使用", sysMenuApi.AddBaseMenuHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/menu/base/tree/:authorityId", "获取基础菜单树-后台使用", sysMenuApi.GetBaseMenuTreeHandler),
		utils.NewRegisterApiParam(common.HttpPost, "/menu/authority", "为角色分配菜单-后台使用", sysMenuApi.AddMenuAuthorityHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/menu/authority/:authorityId", "获取角色菜单权限-后台使用", sysMenuApi.GetMenuAuthorityHandler),
	)
	return sysMenuRouter
}
