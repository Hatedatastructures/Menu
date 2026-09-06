package systemRbac

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"shack/internal/global"
	"shack/internal/model/systemRbac"
	req "shack/internal/model/systemRbac/request"
	res "shack/internal/model/systemRbac/response"

	biz_err "shack/internal/error"

	"github.com/gin-gonic/gin"
	"github.com/wenlng/go-captcha-assets/resources/imagesv2"
	gocaptcha "github.com/wenlng/go-captcha/v2/rotate"
	"github.com/wenlng/go-captcha/v2/base/option"
	"gorm.io/gorm"
)

const (
	// 配置编码
	rotateConfigCode = "captcha_rotate"
	// 验证码过期时间
	rotateCaptchaExpire = 3 * time.Minute
	// 验证码容差值
	rotateCaptchaPadding = 5
)

// captchaRotateInstance 旋转验证码实例
var captchaRotateInstance gocaptcha.Captcha

// initCaptchaRotate 初始化旋转验证码（使用默认配置）
func initCaptchaRotate() error {
	if captchaRotateInstance != nil {
		return nil
	}

	builder := gocaptcha.NewBuilder()

	// 加载背景图
	imgs, err := imagesv2.GetImages()
	if err != nil {
		return fmt.Errorf("加载背景图失败: %w", err)
	}

	builder.SetResources(
		gocaptcha.WithImages(imgs),
	)

	// 应用默认配置
	applyRotateConfigToBuilder(builder, defaultRotateConfig())

	captchaRotateInstance = builder.Make()
	return nil
}

// initCaptchaRotateWithConfig 使用指定配置初始化旋转验证码
func initCaptchaRotateWithConfig(cfg *RotateConfig) error {
	builder := gocaptcha.NewBuilder()

	// 加载背景图
	imgs, err := imagesv2.GetImages()
	if err != nil {
		return fmt.Errorf("加载背景图失败: %w", err)
	}

	builder.SetResources(
		gocaptcha.WithImages(imgs),
	)

	// 应用配置
	applyRotateConfigToBuilder(builder, cfg)

	captchaRotateInstance = builder.Make()
	return nil
}

// applyRotateConfigToBuilder 将配置应用到 builder
func applyRotateConfigToBuilder(builder gocaptcha.Builder, cfg *RotateConfig) {
	opts := []gocaptcha.Option{
		gocaptcha.WithImageSquareSize(cfg.ImageSquareSize),
		gocaptcha.WithRangeAnglePos([]option.RangeVal{
			{Min: cfg.RangeAnglePosMin, Max: cfg.RangeAnglePosMax},
		}),
		gocaptcha.WithThumbImageAlpha(cfg.ThumbImageAlpha),
	}

	// 缩略图大小
	if len(cfg.ThumbImageSquareSize) >= 2 {
		opts = append(opts, gocaptcha.WithRangeThumbImageSquareSize([]int{
			cfg.ThumbImageSquareSize[0],
			cfg.ThumbImageSquareSize[1],
		}))
	}

	builder.SetOptions(opts...)
}

// RotateConfig 旋转验证码配置
type RotateConfig struct {
	ImageSquareSize        int     `json:"imageSquareSize"`
	ThumbImageSquareSize   []int   `json:"thumbImageSquareSize"`
	RangeAnglePosMin       int     `json:"rangeAnglePosMin"`
	RangeAnglePosMax       int     `json:"rangeAnglePosMax"`
	ThumbImageAlpha        float32 `json:"thumbImageAlpha"`
	Version                int     `json:"version"`
}

// defaultRotateConfig 获取默认配置
func defaultRotateConfig() *RotateConfig {
	return &RotateConfig{
		ImageSquareSize:      220,
		ThumbImageSquareSize: []int{80, 80},
		RangeAnglePosMin:     20,
		RangeAnglePosMax:     330,
		ThumbImageAlpha:      1.0,
		Version:              1,
	}
}

// getRotateConfig 获取配置
func (s *CaptchaRotateService) getRotateConfig(ctx context.Context) (*RotateConfig, error) {
	// 先从Redis获取
	cacheKey := "captcha_rotate:config"
	cfgStr, err := global.GVA_REDIS.Get(ctx, cacheKey).Result()
	if err == nil && cfgStr != "" {
		var cfg RotateConfig
		if err := json.Unmarshal([]byte(cfgStr), &cfg); err == nil {
			return &cfg, nil
		}
	}

	// 从数据库获取
	var group systemRbac.ConfigGroup
	err = global.GVA_DB.Where("code = ?", rotateConfigCode).First(&group).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return defaultRotateConfig(), nil
		}
		return nil, err
	}

	// 解析配置
	var cfg RotateConfig
	if len(group.Config) > 0 {
		if err := json.Unmarshal([]byte(group.Config), &cfg); err != nil {
			return defaultRotateConfig(), nil
		}
	}

	// 缓存到Redis
	cfgBytes, _ := json.Marshal(cfg)
	global.GVA_REDIS.Set(ctx, cacheKey, string(cfgBytes), time.Hour*24)

	return &cfg, nil
}

// saveRotateConfig 保存配置到数据库
func (s *CaptchaRotateService) saveRotateConfig(ctx context.Context, cfg *RotateConfig) error {
	cfgJson, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	var group systemRbac.ConfigGroup
	err = global.GVA_DB.Where("code = ?", rotateConfigCode).First(&group).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			group = systemRbac.ConfigGroup{
				Code:   rotateConfigCode,
				Name:  "旋转验证码配置",
				Config: cfgJson,
				Status: 1,
			}
			return global.GVA_DB.Create(&group).Error
		}
		return err
	}

	group.Config = cfgJson
	group.Version++
	cfg.Version = group.Version
	return global.GVA_DB.Save(&group).Error
}

type CaptchaRotateService struct{}

// GenerateCaptchaRotate 生成旋转验证码-前台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年04月05日 16:03:16
func (s *CaptchaRotateService) GenerateCaptchaRotate(
	ctx *gin.Context,
	r req.GenerateCaptchaRotateReq,
) (rs res.GenerateCaptchaRotateRes, err error) {
	// 获取配置
	cfg, err := s.getRotateConfig(ctx)
	if err != nil {
		return rs, biz_err.New(biz_err.SYSTEM_ERROR, "获取配置失败")
	}

	// 如果有传入参数，临时覆盖配置
	if r.ImageSquareSize > 0 {
		cfg.ImageSquareSize = r.ImageSquareSize
	}
	if len(r.ThumbImageSquareSize) > 0 {
		cfg.ThumbImageSquareSize = r.ThumbImageSquareSize
	}
	if r.RangeAnglePosMin != 0 {
		cfg.RangeAnglePosMin = r.RangeAnglePosMin
	}
	if r.RangeAnglePosMax != 0 {
		cfg.RangeAnglePosMax = r.RangeAnglePosMax
	}
	if r.ThumbImageAlpha > 0 {
		cfg.ThumbImageAlpha = r.ThumbImageAlpha
	}

	// 使用配置初始化验证码
	if err := initCaptchaRotateWithConfig(cfg); err != nil {
		return rs, biz_err.New(biz_err.SYSTEM_ERROR, "验证码初始化失败")
	}

	// 生成验证码
	captData, err := captchaRotateInstance.Generate()
	if err != nil {
		return rs, biz_err.New(biz_err.SYSTEM_ERROR, "生成验证码失败")
	}

	blockData := captData.GetData()
	if blockData == nil {
		return rs, biz_err.New(biz_err.SYSTEM_ERROR, "获取验证码数据失败")
	}

	// 生成token
	tokenStr := fmt.Sprintf("rotate_%d", time.Now().UnixNano())

	// 存储验证码数据到Redis
	blockBytes, _ := json.Marshal(blockData)
	storeKey := fmt.Sprintf("captcha_rotate:%s", tokenStr)
	if err := global.GVA_REDIS.Set(ctx, storeKey, string(blockBytes), rotateCaptchaExpire).Err(); err != nil {
		return rs, biz_err.New(biz_err.CACHE_ERROR, "存储验证码失败")
	}

	// 获取图片
	masterImage, err := captData.GetMasterImage().ToBase64()
	if err != nil {
		return rs, biz_err.New(biz_err.SYSTEM_ERROR, "获取主图失败")
	}

	thumbImage, err := captData.GetThumbImage().ToBase64()
	if err != nil {
		return rs, biz_err.New(biz_err.SYSTEM_ERROR, "获取缩略图失败")
	}

	rs = res.GenerateCaptchaRotateRes{
		Token:       tokenStr,
		MasterImage: masterImage,
		ThumbImage:  thumbImage,
		Width:       cfg.ImageSquareSize,
		Height:      cfg.ImageSquareSize,
	}
	return rs, nil
}

// VerifyCaptchaRotate 验证旋转验证码-前台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年04月05日 16:03:16
func (s *CaptchaRotateService) VerifyCaptchaRotate(
	ctx *gin.Context,
	r req.VerifyCaptchaRotateReq,
) (rs res.VerifyCaptchaRotateRes, err error) {
	// 参数校验
	if r.Token == "" {
		return rs, biz_err.New(biz_err.PARAM_MISSING, "验证码令牌不能为空")
	}

	// 获取存储的验证码数据
	storeKey := fmt.Sprintf("captcha_rotate:%s", r.Token)
	blockStr, err := global.GVA_REDIS.Get(ctx, storeKey).Result()
	if err != nil {
		return rs, biz_err.New(biz_err.CAPTCHA_ERROR, "验证码已过期")
	}

	// 解析验证码数据
	var blockData gocaptcha.Block
	if err := json.Unmarshal([]byte(blockStr), &blockData); err != nil {
		return rs, biz_err.New(biz_err.CAPTCHA_ERROR, "验证码数据错误")
	}

	// 清理Redis
	global.GVA_REDIS.Del(ctx, storeKey)

	// 验证旋转角度
	ok := gocaptcha.Validate(
		int(r.Angle),
		blockData.Angle,
		rotateCaptchaPadding,
	)

	if !ok {
		rs = res.VerifyCaptchaRotateRes{
			Success: false,
			Message: "验证失败，请重新尝试",
		}
		return rs, nil
	}

	rs = res.VerifyCaptchaRotateRes{
		Success: true,
		Message: "验证成功",
	}
	return rs, nil
}

// SaveCaptchaRotateConfig 保存旋转验证码配置-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年04月05日 16:03:16
func (s *CaptchaRotateService) SaveCaptchaRotateConfig(
	ctx *gin.Context,
	r req.SaveCaptchaRotateConfigReq,
) (rs res.SaveCaptchaRotateConfigRes, err error) {
	// 构建配置
	cfg := RotateConfig{
		ImageSquareSize:      r.ImageSquareSize,
		ThumbImageSquareSize: r.ThumbImageSquareSize,
		RangeAnglePosMin:     r.RangeAnglePosMin,
		RangeAnglePosMax:     r.RangeAnglePosMax,
		ThumbImageAlpha:      r.ThumbImageAlpha,
	}

	// 设置默认值
	defaultCfg := defaultRotateConfig()
	if cfg.ImageSquareSize == 0 {
		cfg.ImageSquareSize = defaultCfg.ImageSquareSize
	}
	if len(cfg.ThumbImageSquareSize) == 0 {
		cfg.ThumbImageSquareSize = defaultCfg.ThumbImageSquareSize
	}
	if cfg.RangeAnglePosMin == 0 {
		cfg.RangeAnglePosMin = defaultCfg.RangeAnglePosMin
	}
	if cfg.RangeAnglePosMax == 0 {
		cfg.RangeAnglePosMax = defaultCfg.RangeAnglePosMax
	}
	if cfg.ThumbImageAlpha == 0 {
		cfg.ThumbImageAlpha = defaultCfg.ThumbImageAlpha
	}

	// 保存到数据库
	err = s.saveRotateConfig(ctx, &cfg)
	if err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR, "保存配置失败")
	}

	// 清除Redis缓存
	global.GVA_REDIS.Del(ctx, "captcha_rotate:config")

	rs = res.SaveCaptchaRotateConfigRes{
		Version: cfg.Version,
	}
	return rs, nil
}

// GetCaptchaRotateConfig 获取旋转验证码配置-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年04月05日 16:03:16
func (s *CaptchaRotateService) GetCaptchaRotateConfig(
	ctx *gin.Context,
) (rs res.GetCaptchaRotateConfigRes, err error) {
	cfg, err := s.getRotateConfig(ctx)
	if err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR, "获取配置失败")
	}

	rs = res.GetCaptchaRotateConfigRes{
		ImageSquareSize:      cfg.ImageSquareSize,
		ThumbImageSquareSize: cfg.ThumbImageSquareSize,
		RangeAnglePosMin:     cfg.RangeAnglePosMin,
		RangeAnglePosMax:     cfg.RangeAnglePosMax,
		ThumbImageAlpha:      cfg.ThumbImageAlpha,
		Version:              cfg.Version,
	}
	return rs, nil
}

// PreviewCaptchaRotate 获取旋转验证码预览-后台使用(调试用)
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年04月05日 16:03:16
func (s *CaptchaRotateService) PreviewCaptchaRotate(
	ctx *gin.Context,
) (rs res.PreviewCaptchaRotateRes, err error) {
	// 获取配置
	cfg, err := s.getRotateConfig(ctx)
	if err != nil {
		return rs, biz_err.New(biz_err.SYSTEM_ERROR, "获取配置失败")
	}

	// 使用配置初始化验证码
	if err := initCaptchaRotateWithConfig(cfg); err != nil {
		return rs, biz_err.New(biz_err.SYSTEM_ERROR, "验证码初始化失败")
	}

	// 生成验证码
	captData, err := captchaRotateInstance.Generate()
	if err != nil {
		return rs, biz_err.New(biz_err.SYSTEM_ERROR, "生成验证码失败")
	}

	blockData := captData.GetData()
	if blockData == nil {
		return rs, biz_err.New(biz_err.SYSTEM_ERROR, "获取验证码数据失败")
	}

	// 获取图片
	masterImage, err := captData.GetMasterImage().ToBase64()
	if err != nil {
		return rs, biz_err.New(biz_err.SYSTEM_ERROR, "获取主图失败")
	}

	thumbImage, err := captData.GetThumbImage().ToBase64()
	if err != nil {
		return rs, biz_err.New(biz_err.SYSTEM_ERROR, "获取缩略图失败")
	}

	rs = res.PreviewCaptchaRotateRes{
		MasterImage: masterImage,
		ThumbImage:  thumbImage,
		Angle:       blockData.Angle,
	}
	return rs, nil
}