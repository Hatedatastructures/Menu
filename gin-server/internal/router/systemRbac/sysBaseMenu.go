//router 解析
package systemRbac

import (
	"github.com/gin-gonic/gin"
	"shack/internal/model/common"
	"shack/internal/utils"
)

type SysBaseMenuRouter struct{}

func (s *SysBaseMenuRouter) InitSysBaseMenuRouter(Router *gin.RouterGroup, apiSet map[string]struct{}) (R gin.IRoutes) {
	sysBaseMenuRouter := Router.Group("")
	utils.RegisterApi(sysBaseMenuRouter, apiSet, "",
		utils.NewRegisterApiParam(common.HttpDelete, "/basemenu/:id", "删除基础菜单-后台使用", sysBaseMenuApi.DeleteBaseMenuHandler),
		utils.NewRegisterApiParam(common.HttpPut, "/basemenu/update", "更新基础菜单-后台使用", sysBaseMenuApi.UpdateBaseMenuHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/basemenu/:id", "获取基础菜单详情-后台使用", sysBaseMenuApi.GetBaseMenuByIdHandler),
	)
	return sysBaseMenuRouter
}
