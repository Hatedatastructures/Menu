package systemRbac

import (
	req "shack/internal/model/systemRbac/request"
	_ "shack/internal/model/systemRbac/response"
	"shack/internal/utils/validator"
	"shack/internal/vo"

	"github.com/gin-gonic/gin"
)

type AuthApi struct{}

// GetCaptchaHandler
// @Tags systemRbacauthApi
// @Summary GetCaptchaHandler 获取验证码-后台使用
// @Description GetCaptchaHandler 获取验证码-后台使用
// @Success 200 {object} vo.Result{data=_.GetCaptchaRes}
// @Router /auth/captcha [GET]
func (s *AuthApi) GetCaptchaHandler(c *gin.Context) {
	data, err := authService.GetCaptcha(c)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// AuthLoginHandler
// @Tags systemRbacauthApi
// @Summary AuthLoginHandler 用户登录-后台使用
// @Description AuthLoginHandler 用户登录-后台使用
// @Param data body req.AuthLoginReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.AuthLoginRes}
// @Router /auth/login [POST]
func (s *AuthApi) AuthLoginHandler(c *gin.Context) {
	var req req.AuthLoginReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := authService.AuthLogin(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// LogoutHandler
// @Tags systemRbacauthApi
// @Summary LogoutHandler 用户退出登录-后台使用
// @Description LogoutHandler 用户退出登录-后台使用
// @Success 200 {object} vo.Result{}
// @Router /auth/logout [POST]
func (s *AuthApi) LogoutHandler(c *gin.Context) {
	err := authService.Logout(c)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, "操作成功"))
}

// UpdatePasswordHandler
// @Tags systemRbacauthApi
// @Summary UpdatePasswordHandler 修改密码-后台使用
// @Description UpdatePasswordHandler 修改密码-后台使用
// @Param data body req.UpdatePasswordReq true "请求参数"
// @Success 200 {object} vo.Result{}
// @Router /auth/password [POST]
func (s *AuthApi) UpdatePasswordHandler(c *gin.Context) {
	var req req.UpdatePasswordReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	err := authService.UpdatePassword(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, "操作成功"))
}

// AuthRegisterHandler
// @Tags systemRbacauthApi
// @Summary AuthRegisterHandler 用户注册-前台使用
// @Description AuthRegisterHandler 用户注册-前台使用
// @Param data body req.AuthRegisterReq true "请求参数"
// @Success 200 {object} vo.Result{}
// @Router /auth/register [POST]
func (s *AuthApi) AuthRegisterHandler(c *gin.Context) {
	var req req.AuthRegisterReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	err := authService.AuthRegister(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, "操作成功"))
}

// GetPublicConfigHandler
// @Tags systemRbacauthApi
// @Summary GetPublicConfigHandler 获取公开配置(含注册配置)-中台使用
// @Description GetPublicConfigHandler 获取公开配置(含注册配置)-中台使用
// @Success 200 {object} vo.Result{data=_.GetPublicConfigRes}
// @Router /sys/config-group/public [GET]
func (s *AuthApi) GetPublicConfigHandler(c *gin.Context) {
	data, err := authService.GetPublicConfig(c)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}
