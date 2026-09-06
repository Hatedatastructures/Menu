//router 解析
package systemRbac

import (
	"github.com/gin-gonic/gin"
	"shack/internal/model/common"
	"shack/internal/utils"
)

type CaptchaSliderRouter struct{}

func (s *CaptchaSliderRouter) InitCaptchaSliderRouter(Router *gin.RouterGroup, apiSet map[string]struct{}) (R gin.IRoutes) {
	captchaSliderRouter := Router.Group("")
	utils.RegisterApi(captchaSliderRouter, apiSet, "",
		utils.NewRegisterApiParam(common.HttpPost, "/api/captcha/slider/generate", "生成滑块验证码-前台使用", captchaSliderApi.GenerateCaptchaSliderHandler),
		utils.NewRegisterApiParam(common.HttpPost, "/api/captcha/slider/verify", "验证滑块验证码-前台使用", captchaSliderApi.VerifyCaptchaSliderHandler),
		utils.NewRegisterApiParam(common.HttpPost, "/api/captcha/slider/config", "保存滑块验证码配置-后台使用", captchaSliderApi.SaveCaptchaSliderConfigHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/api/captcha/slider/config", "获取滑块验证码配置-后台使用", captchaSliderApi.GetCaptchaSliderConfigHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/api/captcha/slider/preview", "获取滑块验证码预览-后台使用(调试用)", captchaSliderApi.PreviewCaptchaSliderHandler),
	)
	return captchaSliderRouter
}
