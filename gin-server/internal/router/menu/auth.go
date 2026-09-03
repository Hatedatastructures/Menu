//router 解析
package menu

import (
	"github.com/gin-gonic/gin"
	"shack/internal/model/common"
	"shack/internal/utils"
)

type AuthRouter struct{}

func (s *AuthRouter) InitAuthRouter(Router *gin.RouterGroup, apiSet map[string]struct{}) (R gin.IRoutes) {
	authRouter := Router.Group("")
	utils.RegisterApi(authRouter, apiSet, "",
		utils.NewRegisterApiParam(common.HttpPost, "/api/v1/auth/register", "用户认证模块 对应 C++ server AuthRoutes (register/login/refresh) 所有接口均为 POST,无需鉴权 注册新用户", authApi.RegisterHandler),
		utils.NewRegisterApiParam(common.HttpPost, "/api/v1/auth/login", "用户登录", authApi.LoginHandler),
		utils.NewRegisterApiParam(common.HttpPost, "/api/v1/auth/refresh", "刷新令牌", authApi.RefreshTokenHandler),
	)
	return authRouter
}
