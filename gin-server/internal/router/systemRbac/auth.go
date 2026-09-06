//router 解析
package systemRbac

import (
	"github.com/gin-gonic/gin"
	"shack/internal/model/common"
	"shack/internal/utils"
)

type AuthRouter struct{}

func (s *AuthRouter) InitAuthRouter(Router *gin.RouterGroup, apiSet map[string]struct{}) (R gin.IRoutes) {
	authRouter := Router.Group("")
	utils.RegisterApi(authRouter, apiSet, "",
		utils.NewRegisterApiParam(common.HttpGet, "/auth/captcha", "获取验证码-前台使用", authApi.GetCaptchaHandler),
		utils.NewRegisterApiParam(common.HttpPost, "/auth/login", "用户登录-后台使用", authApi.AuthLoginHandler),
		utils.NewRegisterApiParam(common.HttpPost, "/auth/logout", "用户退出登录-后台使用", authApi.LogoutHandler),
		utils.NewRegisterApiParam(common.HttpPost, "/auth/password", "修改密码-后台使用", authApi.UpdatePasswordHandler),
		utils.NewRegisterApiParam(common.HttpPost, "/auth/register", "用户注册-前台使用", authApi.AuthRegisterHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/sys/config-group/public", "获取公开配置(含注册配置)-中台使用", authApi.GetPublicConfigHandler),
	)
	return authRouter
}
