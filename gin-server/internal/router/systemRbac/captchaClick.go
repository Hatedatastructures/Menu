//router 解析
package systemRbac

import (
	"github.com/gin-gonic/gin"
	"shack/internal/model/common"
	"shack/internal/utils"
)

type CaptchaClickRouter struct{}

func (s *CaptchaClickRouter) InitCaptchaClickRouter(Router *gin.RouterGroup, apiSet map[string]struct{}) (R gin.IRoutes) {
	captchaClickRouter := Router.Group("")
	utils.RegisterApi(captchaClickRouter, apiSet, "",
		utils.NewRegisterApiParam(common.HttpPost, "/api/captcha/click/generate", "生成点击验证码-前台使用", captchaClickApi.GenerateCaptchaClickHandler),
		utils.NewRegisterApiParam(common.HttpPost, "/api/captcha/click/verify", "验证点击验证码-前台使用", captchaClickApi.VerifyCaptchaClickHandler),
		utils.NewRegisterApiParam(common.HttpPost, "/api/captcha/click/config", "保存点击验证码配置-后台使用", captchaClickApi.SaveCaptchaClickConfigHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/api/captcha/click/config", "获取点击验证码配置-后台使用", captchaClickApi.GetCaptchaClickConfigHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/api/captcha/click/preview", "获取点击验证码预览-后台使用(调试用)", captchaClickApi.PreviewCaptchaClickHandler),
	)
	return captchaClickRouter
}
