package systemRbac

import (
	biz_err "shack/internal/error"
	req "shack/internal/model/systemRbac/request"
	_ "shack/internal/model/systemRbac/response"
	"shack/internal/utils/validator"
	"shack/internal/vo"

	"github.com/gin-gonic/gin"
)

type CaptchaSliderApi struct{}

// GenerateCaptchaSliderHandler
// @Tags systemRbaccaptchaSliderApi
// @Summary GenerateCaptchaSliderHandler 生成滑块验证码-前台使用
// @Description GenerateCaptchaSliderHandler 生成滑块验证码-前台使用
// @Param data body req.GenerateCaptchaSliderReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.GenerateCaptchaSliderRes}
// @Router /api/captcha/slider/generate [POST]
func (s *CaptchaSliderApi) GenerateCaptchaSliderHandler(c *gin.Context) {
	var req req.GenerateCaptchaSliderReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := captchaSliderService.GenerateCaptchaSlider(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// VerifyCaptchaSliderHandler
// @Tags systemRbaccaptchaSliderApi
// @Summary VerifyCaptchaSliderHandler 验证滑块验证码-前台使用
// @Description VerifyCaptchaSliderHandler 验证滑块验证码-前台使用
// @Param data body req.VerifyCaptchaSliderReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.VerifyCaptchaSliderRes}
// @Router /api/captcha/slider/verify [POST]
func (s *CaptchaSliderApi) VerifyCaptchaSliderHandler(c *gin.Context) {
	var req req.VerifyCaptchaSliderReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := captchaSliderService.VerifyCaptchaSlider(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// SaveCaptchaSliderConfigHandler
// @Tags systemRbaccaptchaSliderApi
// @Summary SaveCaptchaSliderConfigHandler 保存滑块验证码配置-后台使用
// @Description SaveCaptchaSliderConfigHandler 保存滑块验证码配置-后台使用
// @Param data body req.SaveCaptchaSliderConfigReq true "请求参数"
// @Success 200 {object} vo.Result{}
// @Router /api/captcha/slider/config [POST]
func (s *CaptchaSliderApi) SaveCaptchaSliderConfigHandler(c *gin.Context) {
	var req req.SaveCaptchaSliderConfigReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := captchaSliderService.SaveCaptchaSliderConfig(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// GetCaptchaSliderConfigHandler
// @Tags systemRbaccaptchaSliderApi
// @Summary GetCaptchaSliderConfigHandler 获取滑块验证码配置-后台使用
// @Description GetCaptchaSliderConfigHandler 获取滑块验证码配置-后台使用
// @Success 200 {object} vo.Result{data=_.GetCaptchaSliderConfigRes}
// @Router /api/captcha/slider/config [GET]
func (s *CaptchaSliderApi) GetCaptchaSliderConfigHandler(c *gin.Context) {
	data, err := captchaSliderService.GetCaptchaSliderConfig(c)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// PreviewCaptchaSliderHandler
// @Tags systemRbaccaptchaSliderApi
// @Summary PreviewCaptchaSliderHandler 获取滑块验证码预览-后台使用(调试用)
// @Description PreviewCaptchaSliderHandler 获取滑块验证码预览-后台使用(调试用)
// @Param data query req.PreviewCaptchaSliderReq false "请求参数"
// @Success 200 {object} vo.Result{data=_.PreviewCaptchaSliderRes}
// @Router /api/captcha/slider/preview [GET]
func (s *CaptchaSliderApi) PreviewCaptchaSliderHandler(c *gin.Context) {
	var req req.PreviewCaptchaSliderReq
	_ = c.ShouldBindQuery(&req)
	data, err := captchaSliderService.PreviewCaptchaSlider(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}
