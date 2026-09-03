//router 解析
package systemRbac

import (
	"github.com/gin-gonic/gin"
	"shack/internal/model/common"
	"shack/internal/utils"
)

type SysUserRouter struct{}

func (s *SysUserRouter) InitSysUserRouter(Router *gin.RouterGroup, apiSet map[string]struct{}) (R gin.IRoutes) {
	sysUserRouter := Router.Group("")
	utils.RegisterApi(sysUserRouter, apiSet, "",
		utils.NewRegisterApiParam(common.HttpPost, "/user/register", "用户注册-后台使用", sysUserApi.RegisterHandler),
		utils.NewRegisterApiParam(common.HttpPost, "/user/login", "用户登录-后台使用", sysUserApi.LoginHandler),
		utils.NewRegisterApiParam(common.HttpPut, "/user/changepassword", "修改用户密码-后台使用", sysUserApi.ChangePasswordHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/user/list", "获取用户列表-后台使用", sysUserApi.GetUserInfoListHandler),
		utils.NewRegisterApiParam(common.HttpPut, "/user/:id/authority", "设置用户角色(单一)-后台使用", sysUserApi.SetUserAuthorityHandler),
		utils.NewRegisterApiParam(common.HttpPut, "/user/authorities", "设置用户角色(多角色)-后台使用", sysUserApi.SetUserAuthoritiesHandler),
		utils.NewRegisterApiParam(common.HttpDelete, "/user/:id", "删除用户-后台使用", sysUserApi.DeleteUserHandler),
		utils.NewRegisterApiParam(common.HttpPut, "/user/info", "设置用户信息-后台使用", sysUserApi.SetUserInfoHandler),
		utils.NewRegisterApiParam(common.HttpPut, "/user/selfinfo", "设置自身信息-后台使用", sysUserApi.SetSelfInfoHandler),
		utils.NewRegisterApiParam(common.HttpPut, "/userselfsetting", "设置用户配置-后台使用", sysUserApi.SetSelfSettingHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/user/uuid/:uuid", "获取用户信息-后台使用", sysUserApi.GetUserInfHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/user/id/:id", "通过ID获取用户信息-后台使用", sysUserApi.FindUserByIdHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/user/byuuid/:uuid", "通过UUID获取用户信息-后台使用", sysUserApi.FindUserByUuidHandler),
		utils.NewRegisterApiParam(common.HttpPut, "/user/resetpassword", "重置用户密码-后台使用", sysUserApi.ResetPasswordHandler),
	)
	return sysUserRouter
}
