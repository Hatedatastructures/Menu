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
	"github.com/wenlng/go-captcha-assets/resources/tiles"
	gocaptcha "github.com/wenlng/go-captcha/v2/slide"
	"github.com/wenlng/go-captcha/v2/base/option"
	"gorm.io/gorm"
)

const (
	// 配置编码
	sliderConfigCode = "captcha_slider"
	// 验证码过期时间
	sliderCaptchaExpire = 3 * time.Minute
	// 验证码容差值
	sliderCaptchaPadding = 5
)

// captchaSliderInstance 滑块验证码实例
var captchaSliderInstance gocaptcha.Captcha

// initCaptchaSlider 初始化滑块验证码（使用默认配置）
func initCaptchaSlider() error {
	if captchaSliderInstance != nil {
		return nil
	}

	builder := gocaptcha.NewBuilder()

	// 加载背景图
	imgs, err := imagesv2.GetImages()
	if err != nil {
		return fmt.Errorf("加载背景图失败: %w", err)
	}

	// 加载拼图块
	graphs, err := tiles.GetTiles()
	if err != nil {
		return fmt.Errorf("加载拼图块失败: %w", err)
	}

	var newGraphs = make([]*gocaptcha.GraphImage, 0, len(graphs))
	for i := 0; i < len(graphs); i++ {
		graph := graphs[i]
		newGraphs = append(newGraphs, &gocaptcha.GraphImage{
			OverlayImage: graph.OverlayImage,
			MaskImage:    graph.MaskImage,
			ShadowImage:  graph.ShadowImage,
		})
	}

	builder.SetResources(
		gocaptcha.WithGraphImages(newGraphs),
		gocaptcha.WithBackgrounds(imgs),
	)

	// 应用默认配置
	applySliderConfigToBuilder(builder, defaultSliderConfig())

	captchaSliderInstance = builder.Make()
	return nil
}

// initCaptchaSliderWithConfig 使用指定配置初始化滑块验证码
func initCaptchaSliderWithConfig(cfg *SliderConfig) error {
	builder := gocaptcha.NewBuilder()

	// 加载背景图
	imgs, err := imagesv2.GetImages()
	if err != nil {
		return fmt.Errorf("加载背景图失败: %w", err)
	}

	// 加载拼图块
	graphs, err := tiles.GetTiles()
	if err != nil {
		return fmt.Errorf("加载拼图块失败: %w", err)
	}

	var newGraphs = make([]*gocaptcha.GraphImage, 0, len(graphs))
	for i := 0; i < len(graphs); i++ {
		graph := graphs[i]
		newGraphs = append(newGraphs, &gocaptcha.GraphImage{
			OverlayImage: graph.OverlayImage,
			MaskImage:    graph.MaskImage,
			ShadowImage:  graph.ShadowImage,
		})
	}

	builder.SetResources(
		gocaptcha.WithGraphImages(newGraphs),
		gocaptcha.WithBackgrounds(imgs),
	)

	// 应用配置
	applySliderConfigToBuilder(builder, cfg)

	captchaSliderInstance = builder.Make()
	return nil
}

// applySliderConfigToBuilder 将配置应用到 builder
func applySliderConfigToBuilder(builder gocaptcha.Builder, cfg *SliderConfig) {
	opts := []gocaptcha.Option{
		gocaptcha.WithImageSize(option.Size{
			Width:  cfg.MasterSizeWidth,
			Height: cfg.MasterSizeHeight,
		}),
		gocaptcha.WithImageAlpha(cfg.MasterImageAlpha),
		gocaptcha.WithRangeGraphSize(option.RangeVal{
			Min: cfg.RangeGraphSizeMin,
			Max: cfg.RangeGraphSizeMax,
		}),
		gocaptcha.WithRangeGraphAnglePos([]option.RangeVal{
			{Min: cfg.RangeGraphAnglePosMin, Max: cfg.RangeGraphAnglePosMax},
		}),
		gocaptcha.WithGenGraphNumber(cfg.GenGraphNumber),
		gocaptcha.WithEnableGraphVerticalRandom(cfg.EnableGraphVerticalRandom),
	}

	// 盲区方向配置 - 转换为 slide.DeadZoneDirectionType
	if len(cfg.RangeDeadZoneDirections) > 0 {
		directions := convertToDeadZoneDirections(cfg.RangeDeadZoneDirections)
		opts = append(opts, gocaptcha.WithRangeDeadZoneDirections(directions))
	}

	builder.SetOptions(opts...)
}

// convertToDeadZoneDirections 将字符串数组转换为 DeadZoneDirectionType 数组
func convertToDeadZoneDirections(directions []string) []gocaptcha.DeadZoneDirectionType {
	var result []gocaptcha.DeadZoneDirectionType
	for _, d := range directions {
		switch d {
		case "left":
			result = append(result, gocaptcha.DeadZoneDirectionTypeLeft)
		case "right":
			result = append(result, gocaptcha.DeadZoneDirectionTypeRight)
		case "top":
			result = append(result, gocaptcha.DeadZoneDirectionTypeTop)
		case "bottom":
			result = append(result, gocaptcha.DeadZoneDirectionTypeBottom)
		case "top-left":
			result = append(result, gocaptcha.DeadZoneDirectionTypeTop, gocaptcha.DeadZoneDirectionTypeLeft)
		case "top-right":
			result = append(result, gocaptcha.DeadZoneDirectionTypeTop, gocaptcha.DeadZoneDirectionTypeRight)
		case "bottom-left":
			result = append(result, gocaptcha.DeadZoneDirectionTypeBottom, gocaptcha.DeadZoneDirectionTypeLeft)
		case "bottom-right":
			result = append(result, gocaptcha.DeadZoneDirectionTypeBottom, gocaptcha.DeadZoneDirectionTypeRight)
		}
	}
	return result
}

// SliderConfig 滑块验证码配置
type SliderConfig struct {
	MasterSizeWidth           int      `json:"masterSizeWidth"`
	MasterSizeHeight          int      `json:"masterSizeHeight"`
	MasterImageAlpha         float32  `json:"masterImageAlpha"`
	RangeGraphSizeMin         int      `json:"rangeGraphSizeMin"`
	RangeGraphSizeMax         int      `json:"rangeGraphSizeMax"`
	RangeGraphAnglePosMin     int      `json:"rangeGraphAnglePosMin"`
	RangeGraphAnglePosMax     int      `json:"rangeGraphAnglePosMax"`
	GenGraphNumber            int      `json:"genGraphNumber"`
	EnableGraphVerticalRandom bool     `json:"enableGraphVerticalRandom"`
	RangeDeadZoneDirections   []string `json:"rangeDeadZoneDirections"`
	Version                   int      `json:"version"`
}

// defaultSliderConfig 获取默认配置
func defaultSliderConfig() *SliderConfig {
	return &SliderConfig{
		MasterSizeWidth:           300,
		MasterSizeHeight:          220,
		MasterImageAlpha:         1.0,
		RangeGraphSizeMin:         40,
		RangeGraphSizeMax:         50,
		RangeGraphAnglePosMin:     -45,
		RangeGraphAnglePosMax:     45,
		GenGraphNumber:            1,
		EnableGraphVerticalRandom: false,
		Version:                   1,
	}
}

// getSliderConfig 获取配置
func (s *CaptchaSliderService) getSliderConfig(ctx context.Context) (*SliderConfig, error) {
	// 先从Redis获取
	cacheKey := "captcha_slider:config"
	cfgStr, err := global.GVA_REDIS.Get(ctx, cacheKey).Result()
	if err == nil && cfgStr != "" {
		var cfg SliderConfig
		if err := json.Unmarshal([]byte(cfgStr), &cfg); err == nil {
			return &cfg, nil
		}
	}

	// 从数据库获取
	var group systemRbac.ConfigGroup
	err = global.GVA_DB.Where("code = ?", sliderConfigCode).First(&group).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return defaultSliderConfig(), nil
		}
		return nil, err
	}

	// 解析配置
	var cfg SliderConfig
	if len(group.Config) > 0 {
		if err := json.Unmarshal([]byte(group.Config), &cfg); err != nil {
			return defaultSliderConfig(), nil
		}
	}

	// 缓存到Redis
	cfgBytes, _ := json.Marshal(cfg)
	global.GVA_REDIS.Set(ctx, cacheKey, string(cfgBytes), time.Hour*24)

	return &cfg, nil
}

// saveSliderConfig 保存配置到数据库
func (s *CaptchaSliderService) saveSliderConfig(ctx context.Context, cfg *SliderConfig) error {
	cfgJson, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	var group systemRbac.ConfigGroup
	err = global.GVA_DB.Where("code = ?", sliderConfigCode).First(&group).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			group = systemRbac.ConfigGroup{
				Code:   sliderConfigCode,
				Name:  "滑块验证码配置",
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

type CaptchaSliderService struct{}

// GenerateCaptchaSlider 生成滑块验证码-前台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年04月03日 22:40:06
func (s *CaptchaSliderService) GenerateCaptchaSlider(
	ctx *gin.Context,
	r req.GenerateCaptchaSliderReq,
) (rs res.GenerateCaptchaSliderRes, err error) {
	// 获取配置
	cfg, err := s.getSliderConfig(ctx)
	if err != nil {
		return rs, biz_err.New(biz_err.SYSTEM_ERROR, "获取配置失败")
	}

	// 使用配置初始化验证码
	if err := initCaptchaSliderWithConfig(cfg); err != nil {
		return rs, biz_err.New(biz_err.SYSTEM_ERROR, "验证码初始化失败")
	}

	// 生成验证码
	captData, err := captchaSliderInstance.Generate()
	if err != nil {
		return rs, biz_err.New(biz_err.SYSTEM_ERROR, "生成验证码失败")
	}

	blockData := captData.GetData()
	if blockData == nil {
		return rs, biz_err.New(biz_err.SYSTEM_ERROR, "获取验证码数据失败")
	}

	// 生成token
	tokenStr := fmt.Sprintf("slider_%d", time.Now().UnixNano())

	// 存储验证码数据到Redis
	blockBytes, _ := json.Marshal(blockData)
	storeKey := fmt.Sprintf("captcha_slider:%s", tokenStr)
	if err := global.GVA_REDIS.Set(ctx, storeKey, string(blockBytes), sliderCaptchaExpire).Err(); err != nil {
		return rs, biz_err.New(biz_err.CACHE_ERROR, "存储验证码失败")
	}

	// 获取图片
	masterImage, err := captData.GetMasterImage().ToBase64()
	if err != nil {
		return rs, biz_err.New(biz_err.SYSTEM_ERROR, "获取主图失败")
	}

	tileImage, err := captData.GetTileImage().ToBase64()
	if err != nil {
		return rs, biz_err.New(biz_err.SYSTEM_ERROR, "获取拼图块失败")
	}

	rs = res.GenerateCaptchaSliderRes{
		Token:       tokenStr,
		MasterImage: masterImage,
		TileImage:   tileImage,
		X:           blockData.X,
		Y:           blockData.Y,
		TileWidth:   blockData.Width,
		TileHeight:  blockData.Height,
		Width:       cfg.MasterSizeWidth,
		Height:      cfg.MasterSizeHeight,
	}
	return rs, nil
}

// VerifyCaptchaSlider 验证滑块验证码-前台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年04月03日 22:40:06
func (s *CaptchaSliderService) VerifyCaptchaSlider(
	ctx *gin.Context,
	r req.VerifyCaptchaSliderReq,
) (rs res.VerifyCaptchaSliderRes, err error) {
	// 参数校验
	if r.Token == "" {
		return rs, biz_err.New(biz_err.PARAM_MISSING, "验证码令牌不能为空")
	}

	// 获取存储的验证码数据
	storeKey := fmt.Sprintf("captcha_slider:%s", r.Token)
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

	// 验证滑动位置
	ok := gocaptcha.Validate(
		r.X,
		r.Y,
		blockData.X,
		blockData.Y,
		sliderCaptchaPadding,
	)

	if !ok {
		rs = res.VerifyCaptchaSliderRes{
			Success: false,
			Message: "验证失败，请重新尝试",
		}
		return rs, nil
	}

	rs = res.VerifyCaptchaSliderRes{
		Success: true,
		Message: "验证成功",
	}
	return rs, nil
}

// SaveCaptchaSliderConfig 保存滑块验证码配置-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年04月03日 22:40:06
func (s *CaptchaSliderService) SaveCaptchaSliderConfig(
	ctx *gin.Context,
	r req.SaveCaptchaSliderConfigReq,
) (rs res.SaveCaptchaSliderConfigRes, err error) {
	// 构建配置
	cfg := SliderConfig{
		MasterSizeWidth:           r.MasterSizeWidth,
		MasterSizeHeight:          r.MasterSizeHeight,
		MasterImageAlpha:         r.MasterImageAlpha,
		RangeGraphSizeMin:         r.RangeGraphSizeMin,
		RangeGraphSizeMax:         r.RangeGraphSizeMax,
		RangeGraphAnglePosMin:     r.RangeGraphAnglePosMin,
		RangeGraphAnglePosMax:     r.RangeGraphAnglePosMax,
		GenGraphNumber:            r.GenGraphNumber,
		EnableGraphVerticalRandom: r.EnableGraphVerticalRandom,
		RangeDeadZoneDirections:   r.RangeDeadZoneDirections,
	}

	// 设置默认值
	defaultCfg := defaultSliderConfig()
	if cfg.MasterSizeWidth == 0 {
		cfg.MasterSizeWidth = defaultCfg.MasterSizeWidth
	}
	if cfg.MasterSizeHeight == 0 {
		cfg.MasterSizeHeight = defaultCfg.MasterSizeHeight
	}
	if cfg.MasterImageAlpha == 0 {
		cfg.MasterImageAlpha = defaultCfg.MasterImageAlpha
	}
	if cfg.RangeGraphSizeMin == 0 {
		cfg.RangeGraphSizeMin = defaultCfg.RangeGraphSizeMin
	}
	if cfg.RangeGraphSizeMax == 0 {
		cfg.RangeGraphSizeMax = defaultCfg.RangeGraphSizeMax
	}
	if cfg.RangeGraphAnglePosMin == 0 {
		cfg.RangeGraphAnglePosMin = defaultCfg.RangeGraphAnglePosMin
	}
	if cfg.RangeGraphAnglePosMax == 0 {
		cfg.RangeGraphAnglePosMax = defaultCfg.RangeGraphAnglePosMax
	}
	if cfg.GenGraphNumber == 0 {
		cfg.GenGraphNumber = defaultCfg.GenGraphNumber
	}

	// 保存到数据库
	err = s.saveSliderConfig(ctx, &cfg)
	if err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR, "保存配置失败")
	}

	// 清除Redis缓存
	global.GVA_REDIS.Del(ctx, "captcha_slider:config")

	rs = res.SaveCaptchaSliderConfigRes{
		Version: cfg.Version,
	}
	return rs, nil
}

// GetCaptchaSliderConfig 获取滑块验证码配置-后台使用
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年04月03日 22:40:06
func (s *CaptchaSliderService) GetCaptchaSliderConfig(
	ctx *gin.Context,
) (rs res.GetCaptchaSliderConfigRes, err error) {
	cfg, err := s.getSliderConfig(ctx)
	if err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR, "获取配置失败")
	}

	rs = res.GetCaptchaSliderConfigRes{
		MasterSizeWidth:           cfg.MasterSizeWidth,
		MasterSizeHeight:          cfg.MasterSizeHeight,
		MasterImageAlpha:         cfg.MasterImageAlpha,
		RangeGraphSizeMin:         cfg.RangeGraphSizeMin,
		RangeGraphSizeMax:         cfg.RangeGraphSizeMax,
		RangeGraphAnglePosMin:     cfg.RangeGraphAnglePosMin,
		RangeGraphAnglePosMax:     cfg.RangeGraphAnglePosMax,
		GenGraphNumber:            cfg.GenGraphNumber,
		EnableGraphVerticalRandom: cfg.EnableGraphVerticalRandom,
		RangeDeadZoneDirections:   cfg.RangeDeadZoneDirections,
		Version:                   cfg.Version,
	}
	return rs, nil
}

// PreviewCaptchaSlider 获取滑块验证码预览-后台使用(调试用)
// Auth: shack
// Github: https://github.com/hishack
// Time : 2026年04月03日 22:40:06
func (s *CaptchaSliderService) PreviewCaptchaSlider(
	ctx *gin.Context,
	r req.PreviewCaptchaSliderReq,
) (rs res.PreviewCaptchaSliderRes, err error) {
	// 获取配置
	cfg, err := s.getSliderConfig(ctx)
	if err != nil {
		return rs, biz_err.New(biz_err.SYSTEM_ERROR, "获取配置失败")
	}

	// 使用配置初始化验证码
	if err := initCaptchaSliderWithConfig(cfg); err != nil {
		return rs, biz_err.New(biz_err.SYSTEM_ERROR, "验证码初始化失败")
	}

	// 生成验证码
	captData, err := captchaSliderInstance.Generate()
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

	tileImage, err := captData.GetTileImage().ToBase64()
	if err != nil {
		return rs, biz_err.New(biz_err.SYSTEM_ERROR, "获取拼图块失败")
	}

	rs = res.PreviewCaptchaSliderRes{
		MasterImage: masterImage,
		TileImage:   tileImage,
		X:           blockData.X,
		Y:           blockData.Y,
		TileWidth:   blockData.Width,
		TileHeight:  blockData.Height,
	}
	return rs, nil
}

