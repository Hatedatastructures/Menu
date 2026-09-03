package systemRbac

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"shack/internal/global"
	"shack/internal/model/systemRbac"
	req "shack/internal/model/systemRbac/request"
	res "shack/internal/model/systemRbac/response"

	biz_err "shack/internal/error"

	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"shack/internal/utils"
	"gorm.io/gorm"
)

// 图片验证码驱动配置
var (
	imgCaptchaDriver = base64Captcha.NewDriverDigit(80, 240, 4, 0.7, 80)
	// 使用默认内存存储，生产环境可切换到redis
	captchaStore = base64Captcha.DefaultMemStore
)

type AuthService struct{}

// GetCaptcha 获取图片验证码-中台使用
func (s *AuthService) GetCaptcha(
	ctx context.Context,
) (rs res.GetCaptchaRes, err error) {
	cp := base64Captcha.NewCaptcha(imgCaptchaDriver, captchaStore)
	id, b64s, _, err := cp.Generate()
	if err != nil {
		return rs, biz_err.New(biz_err.CAPTCHA_ERROR, "生成验证码失败")
	}
	rs = res.GetCaptchaRes{
		Uuid: id,
		Img:  b64s,
	}
	return rs, nil
}

// AuthLogin 用户登录-后台使用
func (s *AuthService) AuthLogin(
	ctx context.Context,
	r req.AuthLoginReq,
) (rs res.AuthLoginRes, err error) {
	// 1. 获取登录配置
	loginConfig, err := s.getLoginConfig(ctx)
	if err != nil {
		return rs, err
	}

	// 2. 检查是否需要验证码
	if loginConfig.CaptchaEnabled {
		// 根据验证码类型进行验证
		switch r.CaptchaType {
		case "image":
			// 图片验证码
			if !captchaStore.Verify(r.CaptchaUuid, r.CaptchaCode, true) {
				return rs, biz_err.New(biz_err.CAPTCHA_ERROR, "验证码错误")
			}
		case "slider":
			// 滑块验证码
			if r.SliderToken == "" {
				return rs, biz_err.New(biz_err.CAPTCHA_ERROR, "请完成滑块验证")
			}
			sliderVerifyReq := req.VerifyCaptchaSliderReq{
				Token: r.SliderToken,
				X:     r.SliderX,
				Y:     r.SliderY,
			}
			captchaSliderService := &CaptchaSliderService{}
			c, ok := ctx.(*gin.Context)
			if !ok {
				return rs, biz_err.New(biz_err.PARAM_ERROR, "上下文错误")
			}
			sliderRes, err := captchaSliderService.VerifyCaptchaSlider(c, sliderVerifyReq)
			if err != nil {
				return rs, biz_err.New(biz_err.CAPTCHA_ERROR, "滑块验证失败")
			}
			if !sliderRes.Success {
				return rs, biz_err.New(biz_err.CAPTCHA_ERROR, "滑块验证失败")
			}
		case "click":
			// 点选验证码
			if r.ClickToken == "" {
				return rs, biz_err.New(biz_err.CAPTCHA_ERROR, "请完成点选验证")
			}
			clickVerifyReq := req.VerifyCaptchaClickReq{
				Token:  r.ClickToken,
				Points: r.ClickPoints,
			}
			captchaClickService := &CaptchaClickService{}
			c, ok := ctx.(*gin.Context)
			if !ok {
				return rs, biz_err.New(biz_err.PARAM_ERROR, "上下文错误")
			}
			clickRes, err := captchaClickService.VerifyCaptchaClick(c, clickVerifyReq)
			if err != nil {
				return rs, biz_err.New(biz_err.CAPTCHA_ERROR, "点选验证失败")
			}
			if !clickRes.Success {
				return rs, biz_err.New(biz_err.CAPTCHA_ERROR, "点选验证失败")
			}
		case "rotate":
			// 旋转验证码
			if r.RotateToken == "" {
				return rs, biz_err.New(biz_err.CAPTCHA_ERROR, "请完成旋转验证")
			}
			rotateVerifyReq := req.VerifyCaptchaRotateReq{
				Token: r.RotateToken,
				Angle: r.RotateAngle,
			}
			captchaRotateService := &CaptchaRotateService{}
			c, ok := ctx.(*gin.Context)
			if !ok {
				return rs, biz_err.New(biz_err.PARAM_ERROR, "上下文错误")
			}
			rotateRes, err := captchaRotateService.VerifyCaptchaRotate(c, rotateVerifyReq)
			if err != nil {
				return rs, biz_err.New(biz_err.CAPTCHA_ERROR, "旋转验证失败")
			}
			if !rotateRes.Success {
				return rs, biz_err.New(biz_err.CAPTCHA_ERROR, "旋转验证失败")
			}
		default:
			return rs, biz_err.New(biz_err.CAPTCHA_ERROR, "不支持的验证码类型")
		}
	}

	// 3. 查询用户
	var user systemRbac.User
	err = global.GVA_DB.Where("username = ?", r.Username).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return rs, biz_err.New(biz_err.INVALID_CREDENTIALS, "用户名或密码错误")
		}
		return rs, biz_err.New(biz_err.DB_ERROR, "查询用户失败")
	}

	// 4. 检查用户状态
	if user.Enable != 1 {
		return rs, biz_err.New(biz_err.AUTH_ERROR, "用户已被禁用")
	}

	// 5. 验证密码
	if !utils.BcryptCheck(r.Password, user.Password) {
		return rs, biz_err.New(biz_err.INVALID_CREDENTIALS, "用户名或密码错误")
	}

	// 6. 获取登录配置（检查单点登录）
	loginConfig, err = s.getLoginConfig(ctx)
	if err != nil {
		return rs, err
	}

	// 7. 单点登录处理：如果开启单点登录，将旧token加入黑名单
	if loginConfig.SingleLogin {
		// 从Redis获取用户旧的token
		oldToken, err := utils.GetRedisJWT(user.Username)
		if err == nil && oldToken != "" {
			// 将旧token加入黑名单
			global.BlackCache.Set(oldToken, 1, time.Hour*24*7)
			// 同时在Redis中记录黑名单（用于多实例部署）
			_ = utils.SetRedisBlacklist(oldToken, time.Hour*24*7)
		}
	}

	// 8. 生成 JWT token
	token, claims, err := utils.LoginToken(&user)
	if err != nil {
		return rs, biz_err.New(biz_err.AUTH_ERROR, "生成token失败")
	}

	// 9. 获取角色信息
	var authority systemRbac.SysAuthority
	if err := global.GVA_DB.Where("authority_id = ?", user.AuthorityId).First(&authority).Error; err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR, "获取角色信息失败")
	}

	// 10. 将token加入Redis（无论是否多点登录都需要记录，用于单点登录检查）
	c, ok := ctx.(*gin.Context)
	if ok {
		utils.SetToken(c, token, int(claims.RegisteredClaims.ExpiresAt.Unix()-time.Now().Unix()))
		// 将token存储到Redis
		_ = utils.SetRedisJWT(token, user.Username)
	}

	// 11. 构建返回数据（包含JWT token）
	rs = s.buildLoginResWithToken(user, authority, token, claims)
	return rs, nil
}

// Logout 用户退出登录-后台使用
func (s *AuthService) Logout(
	ctx context.Context,
) (err error) {
	// 从gin.Context获取token并加入黑名单
	c, ok := ctx.(*gin.Context)
	if !ok {
		return nil
	}
	// 获取token
	token := c.GetHeader("x-token")
	if token != "" {
		global.BlackCache.Set(token, 1, time.Hour*24*7)
	}
	return nil
}

// UpdatePassword 修改密码-后台使用
func (s *AuthService) UpdatePassword(
	ctx context.Context,
	r req.UpdatePasswordReq,
) (err error) {
	// 1. 获取密码配置
	passwordConfig, err := s.getPasswordConfig(ctx)
	if err != nil {
		return err
	}

	// 2. 验证新密码格式
	if err := validatePassword(r.NewPassword, passwordConfig); err != nil {
		return err
	}

	// 3. 获取当前用户
	c, ok := ctx.(*gin.Context)
	if !ok {
		return biz_err.New(biz_err.PARAM_ERROR, "上下文错误")
	}
	claimsVal, exists := c.Get("claims")
	if !exists {
		return biz_err.New(biz_err.UNAUTHORIZED, "未登录")
	}

	// 获取用户ID - 根据CustomClaims的实际结构
	claims, ok := claimsVal.(req.CustomClaims)
	if !ok {
		return biz_err.New(biz_err.TOKEN_INVALID, "令牌解析失败")
	}

	// 4. 查询用户 - 使用claims中的ID
	// BaseClaims embedded in CustomClaims
	var baseClaims req.BaseClaims = claims.BaseClaims
	userID := baseClaims.ID
	var user systemRbac.User
	if err := global.GVA_DB.Where("id = ?", userID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return biz_err.New(biz_err.USER_NOT_FOUND, "用户不存在")
		}
		return biz_err.New(biz_err.DB_ERROR, "查询用户失败")
	}

	// 5. 验证旧密码
	if !utils.BcryptCheck(r.OldPassword, user.Password) {
		return biz_err.New(biz_err.INVALID_CREDENTIALS, "原密码错误")
	}

	// 6. 加密新密码并保存
	user.Password = utils.BcryptHash(r.NewPassword)
	if err := global.GVA_DB.Save(&user).Error; err != nil {
		return biz_err.New(biz_err.USER_UPDATE_FAILED, "保存密码失败")
	}

	return nil
}

// AuthRegister 用户注册-前台使用
func (s *AuthService) AuthRegister(
	ctx context.Context,
	r req.AuthRegisterReq,
) (err error) {
	// 1. 获取注册配置
	registerConfig, err := s.getRegisterConfig(ctx)
	if err != nil {
		return err
	}

	// 检查是否开放注册
	if !registerConfig.Enabled {
		return biz_err.New(biz_err.FORBIDDEN, "暂不支持注册")
	}

	// 2. 获取密码配置
	passwordConfig, err := s.getPasswordConfig(ctx)
	if err != nil {
		return err
	}

	// 3. 验证密码格式
	if err := validatePassword(r.Password, passwordConfig); err != nil {
		return err
	}

	// 4. 检查用户名格式
	if !isValidUsername(r.Username) {
		return biz_err.New(biz_err.PARAM_FORMAT, "用户名只能包含字母、数字、下划线，长度4-20位")
	}

	// 5. 检查用户名是否已存在
	var existUser systemRbac.User
	if err := global.GVA_DB.Where("username = ?", r.Username).First(&existUser).Error; err == nil {
		return biz_err.New(biz_err.USER_ALREADY_EXISTS, "用户名已存在")
	}

	// 6. 邮箱验证
	if registerConfig.VerifyEmail {
		if r.Email == "" {
			return biz_err.New(biz_err.PARAM_ERROR, "请填写邮箱地址")
		}
		if r.Code == "" {
			return biz_err.New(biz_err.PARAM_ERROR, "请填写邮箱验证码")
		}

		// 验证邮箱验证码
		c, ok := ctx.(*gin.Context)
		if !ok {
			return biz_err.New(biz_err.PARAM_ERROR, "上下文错误")
		}

		verifyCodeReq := req.VerifyCodeReq{
			Email: r.Email,
			Code:  r.Code,
			Type:  "register",
		}

		verifyCodeService := &VerifyCodeService{}
		_, err := verifyCodeService.VerifyCode(c, verifyCodeReq)
		if err != nil {
			return biz_err.New(biz_err.VERIFY_CODE_ERROR, "邮箱验证码错误或已过期")
		}

		// 检查邮箱是否已被使用
		var existEmailUser systemRbac.User
		if err := global.GVA_DB.Where("email = ?", r.Email).First(&existEmailUser).Error; err == nil {
			return biz_err.New(biz_err.USER_ALREADY_EXISTS, "该邮箱已被注册")
		}
	}

	// 7. 加密密码
	hashedPassword := utils.BcryptHash(r.Password)

	// 8. 获取默认角色
	roleId := uint(888) // 默认角色ID
	if registerConfig.DefaultRole != "" {
		// 前端保存的是角色ID，直接使用
		// 尝试将字符串转换为uint
		roleUint, err := strconv.ParseUint(registerConfig.DefaultRole, 10, 32)
		if err == nil {
			// 验证角色是否存在
			var count int64
			if err := global.GVA_DB.Model(&systemRbac.SysAuthority{}).Where("authority_id = ?", roleUint).Count(&count).Error; err == nil && count > 0 {
				roleId = uint(roleUint)
			}
		} else {
			// 如果转换失败，尝试按角色名称查询（兼容旧数据）
			var role systemRbac.SysAuthority
			if err := global.GVA_DB.Where("authority_name = ?", registerConfig.DefaultRole).First(&role).Error; err == nil {
				roleId = role.AuthorityId
			}
		}
	}

	// 9. 获取默认头像
	defaultHeaderImg := registerConfig.DefaultImage
	if defaultHeaderImg == "" {
		defaultHeaderImg = "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif"
	}

	// 10. 创建用户
	user := systemRbac.User{
		Username:    r.Username,
		Password:    hashedPassword,
		NickName:    r.Nickname,
		AuthorityId: roleId,
		Email:       r.Email,
		Phone:       r.Phone,
		HeaderImg:   defaultHeaderImg,
		Enable:      1,
	}
	// 注册需要审核时，设置为待审核状态
	if registerConfig.NeedAudit {
		user.Enable = 2
	}

	if err := global.GVA_DB.Create(&user).Error; err != nil {
		return biz_err.New(biz_err.REGISTRATION_FAILED, "注册失败")
	}

	return nil
}

// GetPublicConfig 获取公开配置(含注册配置)-中台使用
func (s *AuthService) GetPublicConfig(
	ctx context.Context,
) (rs res.GetPublicConfigRes, err error) {
	// 获取系统配置
	systemConfig := s.getSystemConfig()

	// 获取登录配置
	loginConfig, _ := s.getLoginConfig(ctx)

	// 获取注册配置
	registerConfig, _ := s.getRegisterConfig(ctx)

	// 获取密码配置
	passwordConfig, _ := s.getPasswordConfig(ctx)

	// 构建返回
	rs = res.GetPublicConfigRes{
		System: systemConfig,
		Login: res.GetPublicConfigResLogin{
			CaptchaEnabled: loginConfig.CaptchaEnabled,
			CaptchaType:    loginConfig.CaptchaType,
			MaxRetryCount:  loginConfig.MaxRetryCount,
			RememberMe:     loginConfig.RememberMe,
		},
		Register: res.GetPublicConfigResRegister{
			Enabled:      registerConfig.Enabled,
			VerifyEmail:  registerConfig.VerifyEmail,
			VerifyPhone:  registerConfig.VerifyPhone,
			NeedAudit:    registerConfig.NeedAudit,
		},
		Password: res.GetPublicConfigResPassword{
			MinLength:        passwordConfig.MinLength,
			MaxLength:        passwordConfig.MaxLength,
			RequireUppercase: passwordConfig.RequireUppercase,
			RequireLowercase: passwordConfig.RequireLowercase,
			RequireNumber:    passwordConfig.RequireNumber,
			RequireSpecial:   passwordConfig.RequireSpecial,
		},
	}
	return rs, nil
}

// ========== 私有方法 ==========

// getLoginConfig 获取登录配置
func (s *AuthService) getLoginConfig(ctx context.Context) (res.GetLoginConfigRes, error) {
	var group systemRbac.ConfigGroup
	err := global.GVA_DB.Where("code = ?", "login").First(&group).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return res.GetLoginConfigRes{
				CaptchaEnabled: false,
				CaptchaType:    "image",
				MaxRetryCount:  5,
			}, nil
		}
		return res.GetLoginConfigRes{}, biz_err.New(biz_err.DB_ERROR, "获取登录配置失败")
	}

	var config map[string]interface{}
	if group.Config != nil {
		json.Unmarshal(group.Config, &config)
	}

	return res.GetLoginConfigRes{
		CaptchaEnabled: authGetBoolFromMap(config, "captchaEnabled", false),
		CaptchaType:    authGetStringFromMap(config, "captchaType", "image"),
		MaxRetryCount:  authGetIntFromMap(config, "maxRetryCount", 5),
		LockTime:       authGetIntFromMap(config, "lockTime", 30),
		RememberMe:     authGetBoolFromMap(config, "rememberMe", true),
		SingleLogin:    authGetBoolFromMap(config, "singleLogin", false),
	}, nil
}

// getRegisterConfig 获取注册配置
func (s *AuthService) getRegisterConfig(ctx context.Context) (res.GetRegisterConfigRes, error) {
	var group systemRbac.ConfigGroup
	err := global.GVA_DB.Where("code = ?", "register").First(&group).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return res.GetRegisterConfigRes{
				Enabled:      true,
				DefaultRole:  "user",
				DefaultImage: "",
			}, nil
		}
		return res.GetRegisterConfigRes{}, biz_err.New(biz_err.DB_ERROR, "获取注册配置失败")
	}

	var config map[string]interface{}
	if group.Config != nil {
		json.Unmarshal(group.Config, &config)
	}

	return res.GetRegisterConfigRes{
		Enabled:      authGetBoolFromMap(config, "enabled", true),
		VerifyEmail:  authGetBoolFromMap(config, "verifyEmail", false),
		VerifyPhone:  authGetBoolFromMap(config, "verifyPhone", false),
		DefaultRole:  authGetStringFromMap(config, "defaultRole", "user"),
		NeedAudit:    authGetBoolFromMap(config, "needAudit", false),
		DefaultImage: authGetStringFromMap(config, "defaultImage", ""),
	}, nil
}

// getPasswordConfig 获取密码配置
func (s *AuthService) getPasswordConfig(ctx context.Context) (res.GetPasswordConfigRes, error) {
	var group systemRbac.ConfigGroup
	err := global.GVA_DB.Where("code = ?", "password").First(&group).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return res.GetPasswordConfigRes{
				MinLength: 6,
				MaxLength: 20,
			}, nil
		}
		return res.GetPasswordConfigRes{}, biz_err.New(biz_err.DB_ERROR, "获取密码配置失败")
	}

	var config map[string]interface{}
	if group.Config != nil {
		json.Unmarshal(group.Config, &config)
	}

	return res.GetPasswordConfigRes{
		MinLength:        authGetIntFromMap(config, "minLength", 6),
		MaxLength:        authGetIntFromMap(config, "maxLength", 20),
		RequireUppercase: authGetBoolFromMap(config, "requireUppercase", false),
		RequireLowercase: authGetBoolFromMap(config, "requireLowercase", false),
		RequireNumber:    authGetBoolFromMap(config, "requireNumber", false),
		RequireSpecial:   authGetBoolFromMap(config, "requireSpecial", false),
		ExpireDays:       authGetIntFromMap(config, "expireDays", 0),
	}, nil
}

// getSystemConfig 获取系统配置
func (s *AuthService) getSystemConfig() res.GetPublicConfigResSystem {
	var group systemRbac.ConfigGroup
	err := global.GVA_DB.Where("code = ?", "system").First(&group).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return res.GetPublicConfigResSystem{
				SiteName: "BQ-ADMIN",
			}
		}
		return res.GetPublicConfigResSystem{}
	}

	var config map[string]interface{}
	if group.Config != nil {
		json.Unmarshal(group.Config, &config)
	}

	return res.GetPublicConfigResSystem{
		SiteName:        authGetStringFromMap(config, "siteName", "BQ-ADMIN"),
		SiteDescription: authGetStringFromMap(config, "siteDescription", ""),
		SiteLogo:        authGetStringFromMap(config, "siteLogo", ""),
		Copyright:       authGetStringFromMap(config, "copyright", ""),
	}
}

// buildLoginRes 构建登录响应
func (s *AuthService) buildLoginRes(user systemRbac.User, authority systemRbac.SysAuthority) res.AuthLoginRes {
	// 处理ParentId为空的情况
	var parentId uint
	if authority.ParentId != nil {
		parentId = *authority.ParentId
	}

	return res.AuthLoginRes{
		ID:          user.ID,
		Uuid:        user.UUID.String(),
		UserName:    user.Username,
		NickName:    user.NickName,
		HeaderImg:   user.HeaderImg,
		AuthorityId: user.AuthorityId,
		Authority: res.AuthLoginResAuthority{
			AuthorityId:    authority.AuthorityId,
			AuthorityName:  authority.AuthorityName,
			ParentId:       parentId,
			DefaultRouter:  authority.DefaultRouter,
		},
		Phone:  user.Phone,
		Email:  user.Email,
		Enable: user.Enable,
	}
}

// buildLoginResWithToken 构建登录响应（包含JWT token）
func (s *AuthService) buildLoginResWithToken(user systemRbac.User, authority systemRbac.SysAuthority, token string, claims req.CustomClaims) res.AuthLoginRes {
	// 处理ParentId为空的情况
	var parentId uint
	if authority.ParentId != nil {
		parentId = *authority.ParentId
	}

	return res.AuthLoginRes{
		ID:          user.ID,
		Uuid:        user.UUID.String(),
		UserName:    user.Username,
		NickName:    user.NickName,
		HeaderImg:   user.HeaderImg,
		AuthorityId: user.AuthorityId,
		Authority: res.AuthLoginResAuthority{
			AuthorityId:    authority.AuthorityId,
			AuthorityName:  authority.AuthorityName,
			ParentId:       parentId,
			DefaultRouter:  authority.DefaultRouter,
		},
		Phone:     user.Phone,
		Email:     user.Email,
		Enable:    user.Enable,
		Token:     token, // 添加JWT token
		ExpiresAt: claims.RegisteredClaims.ExpiresAt.Unix() * 1000, // 转换为毫秒时间戳
	}
}

// validatePassword 验证密码格式
func validatePassword(password string, config res.GetPasswordConfigRes) error {
	length := len(password)
	if length < config.MinLength || length > config.MaxLength {
		return biz_err.New(biz_err.PARAM_FORMAT, fmt.Sprintf("密码长度必须为%d-%d位", config.MinLength, config.MaxLength))
	}

	if config.RequireUppercase {
		if !regexp.MustCompile(`[A-Z]`).MatchString(password) {
			return biz_err.New(biz_err.PARAM_FORMAT, "密码必须包含大写字母")
		}
	}
	if config.RequireLowercase {
		if !regexp.MustCompile(`[a-z]`).MatchString(password) {
			return biz_err.New(biz_err.PARAM_FORMAT, "密码必须包含小写字母")
		}
	}
	if config.RequireNumber {
		if !regexp.MustCompile(`[0-9]`).MatchString(password) {
			return biz_err.New(biz_err.PARAM_FORMAT, "密码必须包含数字")
		}
	}
	if config.RequireSpecial {
		if !regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(password) {
			return biz_err.New(biz_err.PARAM_FORMAT, "密码必须包含特殊字符")
		}
	}
	return nil
}

// isValidUsername 验证用户名格式
func isValidUsername(username string) bool {
	if len(username) < 4 || len(username) > 20 {
		return false
	}
	return regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]*$`).MatchString(username)
}

// authGetStringFromMapInterface 从map[string]interface获取string
func authGetStringFromMapInterface(m map[string]interface{}, key string, defaultVal string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return defaultVal
}

// authGetIntFromMapInterface 从map[string]interface获取int
func authGetIntFromMapInterface(m map[string]interface{}, key string, defaultVal int) int {
	if v, ok := m[key]; ok {
		switch val := v.(type) {
		case float64:
			return int(val)
		case int:
			return val
		}
	}
	return defaultVal
}

// 辅助函数：从map中获取bool (auth专用)
func authGetBoolFromMap(m map[string]interface{}, key string, defaultVal bool) bool {
	if v, ok := m[key]; ok {
		if b, ok := v.(bool); ok {
			return b
		}
	}
	return defaultVal
}

// 辅助函数：从map中获取string (auth专用)
func authGetStringFromMap(m map[string]interface{}, key string, defaultVal string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return defaultVal
}

// 辅助函数：从map中获取int (auth专用)
func authGetIntFromMap(m map[string]interface{}, key string, defaultVal int) int {
	if v, ok := m[key]; ok {
		switch val := v.(type) {
		case float64:
			return int(val)
		case int:
			return val
		}
	}
	return defaultVal
}
