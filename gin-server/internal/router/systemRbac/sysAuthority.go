//router 解析
package systemRbac

import (
	"github.com/gin-gonic/gin"
	"shack/internal/model/common"
	"shack/internal/utils"
)

type SysAuthorityRouter struct{}

func (s *SysAuthorityRouter) InitSysAuthorityRouter(Router *gin.RouterGroup, apiSet map[string]struct{}) (R gin.IRoutes) {
	sysAuthorityRouter := Router.Group("")
	utils.RegisterApi(sysAuthorityRouter, apiSet, "",
		utils.NewRegisterApiParam(common.HttpPost, "/authority/create", "创建角色-后台使用", sysAuthorityApi.CreateAuthorityHandler),
		utils.NewRegisterApiParam(common.HttpPost, "/authority/copy", "复制角色-后台使用", sysAuthorityApi.CopyAuthorityHandler),
		utils.NewRegisterApiParam(common.HttpPut, "/authority/update", "更新角色-后台使用", sysAuthorityApi.UpdateAuthorityHandler),
		utils.NewRegisterApiParam(common.HttpDelete, "/authority/:authorityId", "删除角色-后台使用", sysAuthorityApi.DeleteAuthorityHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/authority/list", "获取角色列表-后台使用", sysAuthorityApi.GetAuthorityInfoListHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/authority/structlist/:authorityId", "获取角色结构列表-后台使用", sysAuthorityApi.GetStructAuthorityListHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/authority/info/:authorityId", "获取角色信息-后台使用", sysAuthorityApi.GetAuthorityInfoHandler),
		utils.NewRegisterApiParam(common.HttpPut, "/authority/data", "设置角色数据权限-后台使用", sysAuthorityApi.SetDataAuthorityHandler),
		utils.NewRegisterApiParam(common.HttpPut, "/authority/menu", "设置角色菜单权限-后台使用", sysAuthorityApi.SetMenuAuthorityHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/authority/parent/:authorityId", "获取父角色ID-后台使用", sysAuthorityApi.GetParentAuthorityIDHandler),
	)
	return sysAuthorityRouter
}
