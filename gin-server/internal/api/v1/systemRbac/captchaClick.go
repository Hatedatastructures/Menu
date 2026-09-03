package systemRbac

import (
	biz_err "shack/internal/error"
	req "shack/internal/model/systemRbac/request"
	_ "shack/internal/model/systemRbac/response"
	"shack/internal/utils/validator"
	"shack/internal/vo"

	"github.com/gin-gonic/gin"
)

type CaptchaClickApi struct{}

// GenerateCaptchaClickHandler
// @Tags systemRbaccaptchaClickApi
// @Summary GenerateCaptchaClickHandler 生成点击验证码-前台使用
// @Description GenerateCaptchaClickHandler 生成点击验证码-前台使用
// @Param data body req.GenerateCaptchaClickReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.GenerateCaptchaClickRes}
// @Router /api/captcha/click/generate [POST]
func (s *CaptchaClickApi) GenerateCaptchaClickHandler(c *gin.Context) {
	var req req.GenerateCaptchaClickReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := captchaClickService.GenerateCaptchaClick(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// VerifyCaptchaClickHandler
// @Tags systemRbaccaptchaClickApi
// @Summary VerifyCaptchaClickHandler 验证点击验证码-前台使用
// @Description VerifyCaptchaClickHandler 验证点击验证码-前台使用
// @Param data body req.VerifyCaptchaClickReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.VerifyCaptchaClickRes}
// @Router /api/captcha/click/verify [POST]
func (s *CaptchaClickApi) VerifyCaptchaClickHandler(c *gin.Context) {
	var req req.VerifyCaptchaClickReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := captchaClickService.VerifyCaptchaClick(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// SaveCaptchaClickConfigHandler
// @Tags systemRbaccaptchaClickApi
// @Summary SaveCaptchaClickConfigHandler 保存点击验证码配置-后台使用
// @Description SaveCaptchaClickConfigHandler 保存点击验证码配置-后台使用
// @Param data body req.SaveCaptchaClickConfigReq true "请求参数"
// @Success 200 {object} vo.Result{}
// @Router /api/captcha/click/config [POST]
func (s *CaptchaClickApi) SaveCaptchaClickConfigHandler(c *gin.Context) {
	var req req.SaveCaptchaClickConfigReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", biz_err.New(200, "参数错误")))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := captchaClickService.SaveCaptchaClickConfig(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// GetCaptchaClickConfigHandler
// @Tags systemRbaccaptchaClickApi
// @Summary GetCaptchaClickConfigHandler 获取点击验证码配置-后台使用
// @Description GetCaptchaClickConfigHandler 获取点击验证码配置-后台使用
// @Success 200 {object} vo.Result{data=_.GetCaptchaClickConfigRes}
// @Router /api/captcha/click/config [GET]
func (s *CaptchaClickApi) GetCaptchaClickConfigHandler(c *gin.Context) {
	data, err := captchaClickService.GetCaptchaClickConfig(c)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// PreviewCaptchaClickHandler
// @Tags systemRbaccaptchaClickApi
// @Summary PreviewCaptchaClickHandler 获取点击验证码预览-后台使用(调试用)
// @Description PreviewCaptchaClickHandler 获取点击验证码预览-后台使用(调试用)
// @Param data body req.PreviewCaptchaClickReq false "请求参数"
// @Success 200 {object} vo.Result{data=_.PreviewCaptchaClickRes}
// @Router /api/captcha/click/preview [GET]
func (s *CaptchaClickApi) PreviewCaptchaClickHandler(c *gin.Context) {
	var req req.PreviewCaptchaClickReq
	_ = c.ShouldBindQuery(&req)
	data, err := captchaClickService.PreviewCaptchaClick(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}
