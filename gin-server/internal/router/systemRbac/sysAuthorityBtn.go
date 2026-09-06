//router 解析
package systemRbac

import (
	"github.com/gin-gonic/gin"
	"shack/internal/model/common"
	"shack/internal/utils"
)

type SysAuthorityBtnRouter struct{}

func (s *SysAuthorityBtnRouter) InitSysAuthorityBtnRouter(Router *gin.RouterGroup, apiSet map[string]struct{}) (R gin.IRoutes) {
	sysAuthorityBtnRouter := Router.Group("")
	utils.RegisterApi(sysAuthorityBtnRouter, apiSet, "",
		utils.NewRegisterApiParam(common.HttpGet, "/authoritybtn/:authorityId/:menuID", "获取角色按钮权限-后台使用", sysAuthorityBtnApi.GetAuthorityBtnHandler),
		utils.NewRegisterApiParam(common.HttpPut, "/authoritybtn", "设置角色按钮权限-后台使用", sysAuthorityBtnApi.SetAuthorityBtnHandler),
		utils.NewRegisterApiParam(common.HttpDelete, "/authoritybtn/check/:id", "检查角色按钮是否可以删除-后台使用", sysAuthorityBtnApi.CanRemoveAuthorityBtnHandler),
	)
	return sysAuthorityBtnRouter
}
