package menu

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	biz_err "shack/internal/error"
	"shack/internal/global"
	menuModel "shack/internal/model/menu"
	req "shack/internal/model/menu/request"
	res "shack/internal/model/menu/response"
	"shack/internal/utils"
)

type AuthService struct{}

// Register 注册新用户
func (s *AuthService) Register(ctx *gin.Context, r req.RegisterReq) (rs res.RegisterRes, err error) {
	email := strings.TrimSpace(strings.ToLower(r.Email))

	var count int64
	global.GVA_DB.Model(&menuModel.MenuUser{}).Where("Email = ?", email).Count(&count)
	if count > 0 {
		return rs, biz_err.New(biz_err.USER_ALREADY_EXISTS)
	}

	userId := "user." + uuid.New().String()[:16]
	user := menuModel.MenuUser{
		Id:           userId,
		Email:        email,
		PasswordHash: utils.BcryptHash(r.Password),
		DisplayName:  r.DisplayName,
		IsAdmin:      false,
	}
	if err := global.GVA_DB.Create(&user).Error; err != nil {
		return rs, biz_err.New(biz_err.REGISTER_FAILED_MENU)
	}

	tokenRes, tokenErr := s.issueTokens(user.Id, user.Email, user.DisplayName, user.IsAdmin)
	return tokenRes, tokenErr
}

// Login 用户登录
func (s *AuthService) Login(ctx *gin.Context, r req.LoginReq) (rs res.LoginRes, err error) {
	email := strings.TrimSpace(strings.ToLower(r.Email))

	var user menuModel.MenuUser
	if err := global.GVA_DB.Where("Email = ?", email).First(&user).Error; err != nil {
		return rs, biz_err.New(biz_err.INVALID_CREDENTIALS)
	}

	if !utils.BcryptCheck(r.Password, user.PasswordHash) {
		return rs, biz_err.New(biz_err.INVALID_CREDENTIALS)
	}

	tokenRes, tokenErr := s.issueTokens(user.Id, user.Email, user.DisplayName, user.IsAdmin)
	// RegisterRes 和 LoginRes 结构相同,手动映射
	return res.LoginRes{
		User: res.LoginResUser{
			Id:          tokenRes.User.Id,
			Email:       tokenRes.User.Email,
			DisplayName: tokenRes.User.DisplayName,
			IsAdmin:     tokenRes.User.IsAdmin,
		},
		AccessToken:               tokenRes.AccessToken,
		RefreshToken:              tokenRes.RefreshToken,
		AccessTokenExpiresInSeconds: tokenRes.AccessTokenExpiresInSeconds,
	}, tokenErr
}

// RefreshToken 刷新令牌
func (s *AuthService) RefreshToken(ctx *gin.Context, r req.RefreshTokenReq) (rs res.RefreshTokenRes, err error) {
	if r.RefreshToken == "" {
		return rs, biz_err.New(biz_err.TOKEN_INVALID)
	}

	tokenHash := utils.MD5V([]byte(r.RefreshToken))

	var refreshToken menuModel.MenuRefreshToken
	if err := global.GVA_DB.Where("TokenHash = ? AND RevokedAt IS NULL", tokenHash).First(&refreshToken).Error; err != nil {
		return rs, biz_err.New(biz_err.REFRESH_FAILED)
	}

	if refreshToken.ExpiresAt.Before(time.Now()) {
		return rs, biz_err.New(biz_err.TOKEN_EXPIRED)
	}

	var user menuModel.MenuUser
	if err := global.GVA_DB.Where("Id = ?", refreshToken.UserId).First(&user).Error; err != nil {
		return rs, biz_err.New(biz_err.USER_NOT_FOUND)
	}

	// 吊销旧刷新令牌
	now := time.Now()
	global.GVA_DB.Model(&menuModel.MenuRefreshToken{}).
		Where("Id = ?", refreshToken.Id).
		Update("RevokedAt", now)

	tokenRes, tokenErr := s.issueTokens(user.Id, user.Email, user.DisplayName, user.IsAdmin)
	return res.RefreshTokenRes{
		User: res.RefreshTokenResUser{
			Id:          tokenRes.User.Id,
			Email:       tokenRes.User.Email,
			DisplayName: tokenRes.User.DisplayName,
			IsAdmin:     tokenRes.User.IsAdmin,
		},
		AccessToken:               tokenRes.AccessToken,
		RefreshToken:              tokenRes.RefreshToken,
		AccessTokenExpiresInSeconds: tokenRes.AccessTokenExpiresInSeconds,
	}, tokenErr
}

// issueTokens 颁发 access + refresh 令牌对
func (s *AuthService) issueTokens(userId, email, displayName string, isAdmin bool) (res.RegisterRes, error) {
	accessToken := uuid.New().String()
	refreshToken := uuid.New().String()
	accessExpiresIn := 3600
	refreshExpiresAt := time.Now().AddDate(0, 0, 30)

	rt := menuModel.MenuRefreshToken{
		Id:        uuid.New().String(),
		UserId:    userId,
		TokenHash: utils.MD5V([]byte(refreshToken)),
		ExpiresAt: refreshExpiresAt,
	}
	if err := global.GVA_DB.Create(&rt).Error; err != nil {
		return res.RegisterRes{}, biz_err.New(biz_err.SYSTEM_ERROR)
	}

	return res.RegisterRes{
		User: res.RegisterResUser{
			Id:          userId,
			Email:       email,
			DisplayName: displayName,
			IsAdmin:     isAdmin,
		},
		AccessToken:               accessToken,
		RefreshToken:              refreshToken,
		AccessTokenExpiresInSeconds: accessExpiresIn,
	}, nil
}
