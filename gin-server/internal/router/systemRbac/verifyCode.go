//router 解析
package systemRbac

import (
	"github.com/gin-gonic/gin"
	"shack/internal/model/common"
	"shack/internal/utils"
)

type VerifyCodeRouter struct{}

func (s *VerifyCodeRouter) InitVerifyCodeRouter(Router *gin.RouterGroup, apiSet map[string]struct{}) (R gin.IRoutes) {
	verifyCodeRouter := Router.Group("")
	utils.RegisterApi(verifyCodeRouter, apiSet, "",
		utils.NewRegisterApiParam(common.HttpPost, "/email/send-code", "发送验证码-前台使用", verifyCodeApi.SendVerifyCodeHandler),
		utils.NewRegisterApiParam(common.HttpPost, "/email/verify-code", "验证验证码-前台使用", verifyCodeApi.VerifyCodeHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/email/code/status", "获取验证码状态-前台使用", verifyCodeApi.GetVerifyCodeStatusHandler),
	)
	return verifyCodeRouter
}
