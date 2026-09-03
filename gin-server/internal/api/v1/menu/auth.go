package menu

import (

	"github.com/gin-gonic/gin"
	req "shack/internal/model/menu/request"
	_ "shack/internal/model/menu/response"
	"shack/internal/vo"
	"shack/internal/utils/validator"

)

type AuthApi struct{}

// RegisterHandler
// @Tags menuauthApi
// @Summary RegisterHandler 用户认证模块 对应 C++ server AuthRoutes (register/login/refresh) 所有接口均为 POST,无需鉴权 注册新用户
// @Description RegisterHandler 用户认证模块 对应 C++ server AuthRoutes (register/login/refresh) 所有接口均为 POST,无需鉴权 注册新用户
// @Param data body req.RegisterReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.RegisterRes}
// @Router /api/v1/auth/register [POST]
func (s *AuthApi) RegisterHandler(c *gin.Context) {
	var req req.RegisterReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := authService.Register(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// LoginHandler
// @Tags menuauthApi
// @Summary LoginHandler 用户登录
// @Description LoginHandler 用户登录
// @Param data body req.LoginReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.LoginRes}
// @Router /api/v1/auth/login [POST]
func (s *AuthApi) LoginHandler(c *gin.Context) {
	var req req.LoginReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := authService.Login(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// RefreshTokenHandler
// @Tags menuauthApi
// @Summary RefreshTokenHandler 刷新令牌
// @Description RefreshTokenHandler 刷新令牌
// @Param data body req.RefreshTokenReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.RefreshTokenRes}
// @Router /api/v1/auth/refresh [POST]
func (s *AuthApi) RefreshTokenHandler(c *gin.Context) {
	var req req.RefreshTokenReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := authService.RefreshToken(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}