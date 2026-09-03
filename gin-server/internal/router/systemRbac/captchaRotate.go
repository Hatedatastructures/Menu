//router 解析
package systemRbac

import (
	"github.com/gin-gonic/gin"
	"shack/internal/model/common"
	"shack/internal/utils"
)

type CaptchaRotateRouter struct{}

func (s *CaptchaRotateRouter) InitCaptchaRotateRouter(Router *gin.RouterGroup, apiSet map[string]struct{}) (R gin.IRoutes) {
	captchaRotateRouter := Router.Group("")
	utils.RegisterApi(captchaRotateRouter, apiSet, "",
		utils.NewRegisterApiParam(common.HttpPost, "/api/captcha/rotate/generate", "生成旋转验证码-前台使用", captchaRotateApi.GenerateCaptchaRotateHandler),
		utils.NewRegisterApiParam(common.HttpPost, "/api/captcha/rotate/verify", "验证旋转验证码-前台使用", captchaRotateApi.VerifyCaptchaRotateHandler),
		utils.NewRegisterApiParam(common.HttpPost, "/api/captcha/rotate/config", "保存旋转验证码配置-后台使用", captchaRotateApi.SaveCaptchaRotateConfigHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/api/captcha/rotate/config", "获取旋转验证码配置-后台使用", captchaRotateApi.GetCaptchaRotateConfigHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/api/captcha/rotate/preview", "获取旋转验证码预览-后台使用(调试用)", captchaRotateApi.PreviewCaptchaRotateHandler),
	)
	return captchaRotateRouter
}
