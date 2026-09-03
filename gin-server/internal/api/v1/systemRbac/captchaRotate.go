package systemRbac

import (
	req "shack/internal/model/systemRbac/request"
	_ "shack/internal/model/systemRbac/response"
	"shack/internal/utils/validator"
	"shack/internal/vo"

	"github.com/gin-gonic/gin"
)

type CaptchaRotateApi struct{}

// GenerateCaptchaRotateHandler
// @Tags systemRbaccaptchaRotateApi
// @Summary GenerateCaptchaRotateHandler 生成旋转验证码-前台使用
// @Description GenerateCaptchaRotateHandler 生成旋转验证码-前台使用
// @Param data body req.GenerateCaptchaRotateReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.GenerateCaptchaRotateRes}
// @Router /api/captcha/rotate/generate [POST]
func (s *CaptchaRotateApi) GenerateCaptchaRotateHandler(c *gin.Context) {
	var req req.GenerateCaptchaRotateReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := captchaRotateService.GenerateCaptchaRotate(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// VerifyCaptchaRotateHandler
// @Tags systemRbaccaptchaRotateApi
// @Summary VerifyCaptchaRotateHandler 验证旋转验证码-前台使用
// @Description VerifyCaptchaRotateHandler 验证旋转验证码-前台使用
// @Param data body req.VerifyCaptchaRotateReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.VerifyCaptchaRotateRes}
// @Router /api/captcha/rotate/verify [POST]
func (s *CaptchaRotateApi) VerifyCaptchaRotateHandler(c *gin.Context) {
	var req req.VerifyCaptchaRotateReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := captchaRotateService.VerifyCaptchaRotate(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// SaveCaptchaRotateConfigHandler
// @Tags systemRbaccaptchaRotateApi
// @Summary SaveCaptchaRotateConfigHandler 保存旋转验证码配置-后台使用
// @Description SaveCaptchaRotateConfigHandler 保存旋转验证码配置-后台使用
// @Param data body req.SaveCaptchaRotateConfigReq true "请求参数"
// @Success 200 {object} vo.Result{}
// @Router /api/captcha/rotate/config [POST]
func (s *CaptchaRotateApi) SaveCaptchaRotateConfigHandler(c *gin.Context) {
	var req req.SaveCaptchaRotateConfigReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	_, err := captchaRotateService.SaveCaptchaRotateConfig(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, "操作成功"))
}

// GetCaptchaRotateConfigHandler
// @Tags systemRbaccaptchaRotateApi
// @Summary GetCaptchaRotateConfigHandler 获取旋转验证码配置-后台使用
// @Description GetCaptchaRotateConfigHandler 获取旋转验证码配置-后台使用
// @Success 200 {object} vo.Result{data=_.GetCaptchaRotateConfigRes}
// @Router /api/captcha/rotate/config [GET]
func (s *CaptchaRotateApi) GetCaptchaRotateConfigHandler(c *gin.Context) {
	data, err := captchaRotateService.GetCaptchaRotateConfig(c)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// PreviewCaptchaRotateHandler
// @Tags systemRbaccaptchaRotateApi
// @Summary PreviewCaptchaRotateHandler 获取旋转验证码预览-后台使用(调试用)
// @Description PreviewCaptchaRotateHandler 获取旋转验证码预览-后台使用(调试用)
// @Success 200 {object} vo.Result{data=_.PreviewCaptchaRotateRes}
// @Router /api/captcha/rotate/preview [GET]
func (s *CaptchaRotateApi) PreviewCaptchaRotateHandler(c *gin.Context) {
	data, err := captchaRotateService.PreviewCaptchaRotate(c)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}
