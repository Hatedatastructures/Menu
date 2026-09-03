package menu

import (
	"github.com/gin-gonic/gin"
	req "shack/internal/model/menu/request"
	res "shack/internal/model/menu/response"
)

type AuthService struct{}

// 用户认证模块 对应 C++ server AuthRoutes (register/login/refresh) 所有接口均为 POST,无需鉴权 注册新用户
// Auth: shack 
// Github: https://github.com/hishack
// Time : 2026年09月02日 23:41:50
func (s *AuthService) Register(
	ctx *gin.Context,
	r req.RegisterReq,
) (rs res.RegisterRes, err error) {
	return rs, nil
}

// 用户登录
// Auth: shack 
// Github: https://github.com/hishack
// Time : 2026年09月02日 23:41:50
func (s *AuthService) Login(
	ctx *gin.Context,
	r req.LoginReq,
) (rs res.LoginRes, err error) {
	return rs, nil
}

// 刷新令牌
// Auth: shack 
// Github: https://github.com/hishack
// Time : 2026年09月02日 23:41:50
func (s *AuthService) RefreshToken(
	ctx *gin.Context,
	r req.RefreshTokenReq,
) (rs res.RefreshTokenRes, err error) {
	return rs, nil
}

