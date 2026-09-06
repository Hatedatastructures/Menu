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
	"github.com/golang/freetype/truetype"
	"github.com/wenlng/go-captcha-assets/bindata/chars"
	"github.com/wenlng/go-captcha-assets/resources/fonts/fzshengsksjw"
	"github.com/wenlng/go-captcha-assets/resources/imagesv2"
	gocaptcha "github.com/wenlng/go-captcha/v2/click"
	"github.com/wenlng/go-captcha/v2/base/option"
	"gorm.io/gorm"
)

const (
	// 配置编码
	configCode = "captcha_click"
	// 验证码过期时间
	captchaExpire = 3 * time.Minute
	// 验证码容差值
	captchaPadding = 5
)

// captchaClickInstance 点击验证码实例
var captchaClickInstance gocaptcha.Captcha

// initCaptchaClick 初始化点击验证码（使用默认配置）
func initCaptchaClick() error {
	if captchaClickInstance != nil {
		return nil
	}

	builder := gocaptcha.NewBuilder()

	// 加载字体
	fonts, err := fzshengsksjw.GetFont()
	if err != nil {
		return fmt.Errorf("加载字体失败: %w", err)
	}

	// 加载背景图
	imgs, err := imagesv2.GetImages()
	if err != nil {
		return fmt.Errorf("加载背景图失败: %w", err)
	}

	builder.SetResources(
		gocaptcha.WithChars(chars.GetChineseChars()),
		gocaptcha.WithFonts([]*truetype.Font{fonts}),
		gocaptcha.WithBackgrounds(imgs),
	)

	// 应用默认配置
	applyConfigToBuilder(builder, defaultClickConfig())

	captchaClickInstance = builder.Make()
	return nil
}

// initCaptchaClickWithConfig 使用指定配置初始化点击验证码
func initCaptchaClickWithConfig(cfg *ClickConfig) error {
	builder := gocaptcha.NewBuilder()

	// 加载字体
	fonts, err := fzshengsksjw.GetFont()
	if err != nil {
		return fmt.Errorf("加载字体失败: %w", err)
	}

	// 加载背景图
	imgs, err := imagesv2.GetImages()
	if err != nil {
		return fmt.Errorf("加载背景图失败: %w", err)
	}

	builder.SetResources(
		gocaptcha.WithChars(chars.GetChineseChars()),
		gocaptcha.WithFonts([]*truetype.Font{fonts}),
		gocaptcha.WithBackgrounds(imgs),
	)

	// 应用配置
	applyConfigToBuilder(builder, cfg)

	captchaClickInstance = builder.Make()
	return nil
}

// applyConfigToBuilder 将配置应用到 builder
func applyConfigToBuilder(builder gocaptcha.Builder, cfg *ClickConfig) {
	// 主图配置
	opts := []gocaptcha.Option{
		gocaptcha.WithImageSize(option.Size{
			Width:  cfg.MasterSizeWidth,
			Height: cfg.MasterSizeHeight,
		}),
		gocaptcha.WithRangeLen(option.RangeVal{
			Min: cfg.MasterRangeLenMin,
			Max: cfg.MasterRangeLenMax,
		}),
		gocaptcha.WithRangeAnglePos([]option.RangeVal{
			{Min: cfg.MasterRangeAnglePosMin, Max: cfg.MasterRangeAnglePosMax},
		}),
		gocaptcha.WithRangeSize(option.RangeVal{
			Min: cfg.MasterRangeSizeMin,
			Max: cfg.MasterRangeSizeMax,
		}),
		gocaptcha.WithImageAlpha(cfg.MasterImageAlpha),
	}

	// 缩略图配置
	opts = append(opts,
		gocaptcha.WithRangeThumbImageSize(option.Size{
			Width:  cfg.ThumbSizeWidth,
			Height: cfg.ThumbSizeHeight,
		}),
		gocaptcha.WithRangeVerifyLen(option.RangeVal{
			Min: cfg.ThumbRangeVerifyLenMin,
			Max: cfg.ThumbRangeVerifyLenMax,
		}),
		gocaptcha.WithRangeThumbBgDistort(cfg.ThumbBgDistort),
		gocaptcha.WithRangeThumbBgCirclesNum(cfg.ThumbBgCirclesNum),
		gocaptcha.WithRangeThumbBgSlimLineNum(cfg.ThumbBgSlimLineNum),
	)

	// 颜色配置
	if len(cfg.MasterRangeColors) > 0 {
		opts = append(opts, gocaptcha.WithRangeColors(cfg.MasterRangeColors))
	}
	if len(cfg.ThumbRangeThumbColors) > 0 {
		opts = append(opts, gocaptcha.WithRangeThumbColors(cfg.ThumbRangeThumbColors))
	}
	if len(cfg.ThumbRangeThumbBgColors) > 0 {
		opts = append(opts, gocaptcha.WithRangeThumbBgColors(cfg.ThumbRangeThumbBgColors))
	}

	// 阴影配置
	opts = append(opts, gocaptcha.WithDisplayShadow(cfg.MasterDisplayShadow))
	if cfg.MasterDisplayShadow {
		opts = append(opts, gocaptcha.WithShadowColor(cfg.MasterShadowColor))
	}

	// 其他配置
	opts = append(opts,
		gocaptcha.WithDisabledRangeVerifyLen(cfg.ThumbDisabledRangeVerifyLen),
		gocaptcha.WithIsThumbNonDeformAbility(cfg.ThumbIsThumbNonDeformAbility),
	)

	builder.SetOptions(opts...)
}

// ClickConfig 点击验证码配置
type ClickConfig struct {
	MasterSizeWidth            int     `json:"masterSizeWidth"`
	MasterSizeHeight           int     `json:"masterSizeHeight"`
	MasterRangeLenMin         int     `json:"masterRangeLenMin"`
	MasterRangeLenMax         int     `json:"masterRangeLenMax"`
	MasterRangeAnglePosMin     int     `json:"masterRangeAnglePosMin"`
	MasterRangeAnglePosMax    int     `json:"masterRangeAnglePosMax"`
	MasterRangeSizeMin       int     `json:"masterRangeSizeMin"`
	MasterRangeSizeMax       int     `json:"masterRangeSizeMax"`
	MasterRangeColors       []string `json:"masterRangeColors"`
	MasterDisplayShadow     bool     `json:"masterDisplayShadow"`
	MasterShadowColor      string   `json:"masterShadowColor"`
	MasterShadowPointX    int     `json:"masterShadowPointX"`
	MasterShadowPointY    int     `json:"masterShadowPointY"`
	MasterImageAlpha       float32  `json:"masterImageAlpha"`
	MasterUseShapeOriginalColor bool `json:"masterUseShapeOriginalColor"`
	ThumbSizeWidth         int     `json:"thumbSizeWidth"`
	ThumbSizeHeight        int     `json:"thumbSizeHeight"`
	ThumbRangeVerifyLenMin int     `json:"thumbRangeVerifyLenMin"`
	ThumbRangeVerifyLenMax  int     `json:"thumbRangeVerifyLenMax"`
	ThumbDisabledRangeVerifyLen bool `json:"thumbDisabledRangeVerifyLen"`
	ThumbRangeThumbSizeMin int    `json:"thumbRangeThumbSizeMin"`
	ThumbRangeThumbSizeMax int     `json:"thumbRangeThumbSizeMax"`
	ThumbRangeThumbColors []string `json:"thumbRangeThumbColors"`
	ThumbRangeThumbBgColors []string `json:"thumbRangeThumbBgColors"`
	ThumbIsThumbNonDeformAbility bool `json:"thumbIsThumbNonDeformAbility"`
	ThumbBgDistort         int     `json:"thumbBgDistort"`
	ThumbBgCirclesNum      int     `json:"thumbBgCirclesNum"`
	ThumbBgSlimLineNum   int     `json:"thumbBgSlimLineNum"`
	Version              int     `json:"version"`
}

// defaultClickConfig 获取默认配置
func defaultClickConfig() *ClickConfig {
	return &ClickConfig{
		MasterSizeWidth:             300,
		MasterSizeHeight:            220,
		MasterRangeLenMin:           3,  // 主图显示 3-5 个字符
		MasterRangeLenMax:           5,
		MasterRangeAnglePosMin:      -30,
		MasterRangeAnglePosMax:      30,
		MasterRangeSizeMin:          30,
		MasterRangeSizeMax:          45,
		MasterDisplayShadow:         true,
		MasterShadowColor:          "#666666",
		MasterImageAlpha:           1.0,
		ThumbSizeWidth:             150,
		ThumbSizeHeight:            50,
		ThumbRangeVerifyLenMin:      1,  // 缩略图显示 1-2 个字符（需要点击的）
		ThumbRangeVerifyLenMax:      2,
		ThumbDisabledRangeVerifyLen: false,
		ThumbBgCirclesNum:           15,
		ThumbBgSlimLineNum:          8,
		ThumbBgDistort:              0,
		ThumbIsThumbNonDeformAbility: false,
		Version:                    1,
	}
}

// getConfig 获取配置
func (s *CaptchaClickService) getConfig(ctx context.Context) (*ClickConfig, error) {
	// 先从Redis获取
	cacheKey := "captcha_click:config"
	cfgStr, err := global.GVA_REDIS.Get(ctx, cacheKey).Result()
	if err == nil && cfgStr != "" {
		var cfg ClickConfig
		if err := json.Unmarshal([]byte(cfgStr), &cfg); err == nil {
			return &cfg, nil
		}
	}

	// 从数据库获取
	var group systemRbac.ConfigGroup
	err = global.GVA_DB.Where("code = ?", configCode).First(&group).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return defaultClickConfig(), nil
		}
		return nil, err
	}

	// 解析配置
	var cfg ClickConfig
	if len(group.Config) > 0 {
		if err := json.Unmarshal([]byte(group.Config), &cfg); err != nil {
			return defaultClickConfig(), nil
		}
	}

	// 缓存到Redis
	cfgBytes, _ := json.Marshal(cfg)
	global.GVA_REDIS.Set(ctx, cacheKey, string(cfgBytes), time.Hour*24)

	return &cfg, nil
}

// saveConfig 保存配置到数据库
func (s *CaptchaClickService) saveConfig(ctx context.Context, cfg *ClickConfig) error {
	cfgJson, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	var group systemRbac.ConfigGroup
	err = global.GVA_DB.Where("code = ?", configCode).First(&group).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			group = systemRbac.ConfigGroup{
				Code:   configCode,
				Name:  "点击验证码配置",
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

type CaptchaClickService struct{}

// GenerateCaptchaClick 生成点击验证码-前台使用
func (s *CaptchaClickService) GenerateCaptchaClick(
	ctx *gin.Context,
	r req.GenerateCaptchaClickReq,
) (rs res.GenerateCaptchaClickRes, err error) {
	// 获取配置
	cfg, err := s.getConfig(ctx)
	if err != nil {
		return rs, biz_err.New(biz_err.SYSTEM_ERROR, "获取配置失败")
	}

	// 使用配置初始化验证码
	if err := initCaptchaClickWithConfig(cfg); err != nil {
		return rs, biz_err.New(biz_err.SYSTEM_ERROR, "验证码初始化失败")
	}

	// 生成验证码
	captData, err := captchaClickInstance.Generate()
	if err != nil {
		return rs, biz_err.New(biz_err.SYSTEM_ERROR, "生成验证码失败")
	}

	dotData := captData.GetData()
	if dotData == nil {
		return rs, biz_err.New(biz_err.SYSTEM_ERROR, "获取验证码数据失败")
	}

	// 生成token
	tokenStr := fmt.Sprintf("click_%d", time.Now().UnixNano())

	// 存储验证码数据到Redis
	dotBytes, _ := json.Marshal(dotData)
	storeKey := fmt.Sprintf("captcha_click:%s", tokenStr)
	if err := global.GVA_REDIS.Set(ctx, storeKey, string(dotBytes), captchaExpire).Err(); err != nil {
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

	rs = res.GenerateCaptchaClickRes{
		Token:       tokenStr,
		MasterImage: masterImage,
		ThumbImage:  thumbImage,
		Width:       cfg.MasterSizeWidth,
		Height:      cfg.MasterSizeHeight,
	}
	return rs, nil
}

// VerifyCaptchaClick 验证点击验证码-前台使用
func (s *CaptchaClickService) VerifyCaptchaClick(
	ctx *gin.Context,
	r req.VerifyCaptchaClickReq,
) (rs res.VerifyCaptchaClickRes, err error) {
	// 参数校验
	if r.Token == "" {
		return rs, biz_err.New(biz_err.PARAM_MISSING, "验证码令牌不能为空")
	}
	if len(r.Points) == 0 {
		return rs, biz_err.New(biz_err.PARAM_MISSING, "点击坐标不能为空")
	}

	// 获取存储的验证码数据
	storeKey := fmt.Sprintf("captcha_click:%s", r.Token)
	dotStr, err := global.GVA_REDIS.Get(ctx, storeKey).Result()
	if err != nil {
		return rs, biz_err.New(biz_err.CAPTCHA_ERROR, "验证码已过期")
	}

	// 解析验证码数据
	var dotData map[int]*gocaptcha.Dot
	if err := json.Unmarshal([]byte(dotStr), &dotData); err != nil {
		return rs, biz_err.New(biz_err.CAPTCHA_ERROR, "验证码数据错误")
	}

	// 清理Redis
	global.GVA_REDIS.Del(ctx, storeKey)

	// 验证点击 - 按顺序验证所有点
	if len(r.Points) != len(dotData) {
		rs = res.VerifyCaptchaClickRes{
			Success: false,
			Message: fmt.Sprintf("验证失败，需要点击%d个字符", len(dotData)),
		}
		return rs, nil
	}

	success := true
	for i, userPoint := range r.Points {
		// 按顺序验证每个点
		dot, exists := dotData[i]
		if !exists {
			success = false
			break
		}

		// 验证坐标是否在容差范围内
		ok := gocaptcha.Validate(
			userPoint.X,
			userPoint.Y,
			dot.X,
			dot.Y,
			0, 0, // width, height
			captchaPadding,
		)
		if !ok {
			success = false
			break
		}
	}

	if !success {
		rs = res.VerifyCaptchaClickRes{
			Success: false,
			Message: "验证失败，请重新尝试",
		}
		return rs, nil
	}

	rs = res.VerifyCaptchaClickRes{
		Success: true,
		Message: "验证成功",
	}
	return rs, nil
}

// SaveCaptchaClickConfig 保存点击验证码配置-后台使用
func (s *CaptchaClickService) SaveCaptchaClickConfig(
	ctx *gin.Context,
	r req.SaveCaptchaClickConfigReq,
) (rs res.SaveCaptchaClickConfigRes, err error) {
	// 构建配置
	cfg := ClickConfig{
		MasterSizeWidth:             r.MasterSizeWidth,
		MasterSizeHeight:            r.MasterSizeHeight,
		MasterRangeLenMin:           r.MasterRangeLenMin,
		MasterRangeLenMax:           r.MasterRangeLenMax,
		MasterRangeAnglePosMin:       r.MasterRangeAnglePosMin,
		MasterRangeAnglePosMax:      r.MasterRangeAnglePosMax,
		MasterRangeSizeMin:           r.MasterRangeSizeMin,
		MasterRangeSizeMax:           r.MasterRangeSizeMax,
		MasterRangeColors:           r.MasterRangeColors,
		MasterDisplayShadow:         r.MasterDisplayShadow,
		MasterShadowColor:          r.MasterShadowColor,
		MasterShadowPointX:         r.MasterShadowPointX,
		MasterShadowPointY:         r.MasterShadowPointY,
		MasterImageAlpha:          r.MasterImageAlpha,
		MasterUseShapeOriginalColor: r.MasterUseShapeOriginalColor,
		ThumbSizeWidth:           r.ThumbSizeWidth,
		ThumbSizeHeight:         r.ThumbSizeHeight,
		ThumbRangeVerifyLenMin:    r.ThumbRangeVerifyLenMin,
		ThumbRangeVerifyLenMax:   r.ThumbRangeVerifyLenMax,
		ThumbDisabledRangeVerifyLen: r.ThumbDisabledRangeVerifyLen,
		ThumbRangeThumbSizeMin:    r.ThumbRangeThumbSizeMin,
		ThumbRangeThumbSizeMax:    r.ThumbRangeThumbSizeMax,
		ThumbRangeThumbColors:      r.ThumbRangeThumbColors,
		ThumbRangeThumbBgColors:  r.ThumbRangeThumbBgColors,
		ThumbIsThumbNonDeformAbility: r.ThumbIsThumbNonDeformAbility,
		ThumbBgDistort:           r.ThumbBgDistort,
		ThumbBgCirclesNum:        r.ThumbBgCirclesNum,
		ThumbBgSlimLineNum:       r.ThumbBgSlimLineNum,
	}

	// 设置默认值
	defaultCfg := defaultClickConfig()
	if cfg.MasterSizeWidth == 0 {
		cfg.MasterSizeWidth = defaultCfg.MasterSizeWidth
	}
	if cfg.MasterSizeHeight == 0 {
		cfg.MasterSizeHeight = defaultCfg.MasterSizeHeight
	}
	if cfg.MasterRangeLenMin == 0 {
		cfg.MasterRangeLenMin = defaultCfg.MasterRangeLenMin
	}
	if cfg.MasterRangeLenMax == 0 {
		cfg.MasterRangeLenMax = defaultCfg.MasterRangeLenMax
	}
	if cfg.MasterRangeAnglePosMin == 0 {
		cfg.MasterRangeAnglePosMin = defaultCfg.MasterRangeAnglePosMin
	}
	if cfg.MasterRangeAnglePosMax == 0 {
		cfg.MasterRangeAnglePosMax = defaultCfg.MasterRangeAnglePosMax
	}
	if cfg.MasterRangeSizeMin == 0 {
		cfg.MasterRangeSizeMin = defaultCfg.MasterRangeSizeMin
	}
	if cfg.MasterRangeSizeMax == 0 {
		cfg.MasterRangeSizeMax = defaultCfg.MasterRangeSizeMax
	}
	if cfg.MasterImageAlpha == 0 {
		cfg.MasterImageAlpha = defaultCfg.MasterImageAlpha
	}
	if !cfg.MasterDisplayShadow {
		cfg.MasterDisplayShadow = defaultCfg.MasterDisplayShadow
	}
	if cfg.MasterShadowColor == "" {
		cfg.MasterShadowColor = defaultCfg.MasterShadowColor
	}
	if cfg.ThumbSizeWidth == 0 {
		cfg.ThumbSizeWidth = defaultCfg.ThumbSizeWidth
	}
	if cfg.ThumbSizeHeight == 0 {
		cfg.ThumbSizeHeight = defaultCfg.ThumbSizeHeight
	}
	if cfg.ThumbRangeVerifyLenMin == 0 {
		cfg.ThumbRangeVerifyLenMin = defaultCfg.ThumbRangeVerifyLenMin
	}
	if cfg.ThumbRangeVerifyLenMax == 0 {
		cfg.ThumbRangeVerifyLenMax = defaultCfg.ThumbRangeVerifyLenMax
	}
	if cfg.ThumbBgCirclesNum == 0 {
		cfg.ThumbBgCirclesNum = defaultCfg.ThumbBgCirclesNum
	}
	if cfg.ThumbBgSlimLineNum == 0 {
		cfg.ThumbBgSlimLineNum = defaultCfg.ThumbBgSlimLineNum
	}

	// 保存到数据库
	err = s.saveConfig(ctx, &cfg)
	if err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR, "保存配置失败")
	}

	// 清除Redis缓存
	global.GVA_REDIS.Del(ctx, "captcha_click:config")

	rs = res.SaveCaptchaClickConfigRes{
		Version: cfg.Version,
	}
	return rs, nil
}

// GetCaptchaClickConfig 获取点击验证码配置-后台使用
func (s *CaptchaClickService) GetCaptchaClickConfig(
	ctx *gin.Context,
) (rs res.GetCaptchaClickConfigRes, err error) {
	cfg, err := s.getConfig(ctx)
	if err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR, "获取配置失败")
	}

	rs = res.GetCaptchaClickConfigRes{
		MasterSizeWidth:            cfg.MasterSizeWidth,
		MasterSizeHeight:           cfg.MasterSizeHeight,
		MasterRangeLenMin:         cfg.MasterRangeLenMin,
		MasterRangeLenMax:         cfg.MasterRangeLenMax,
		MasterRangeAnglePosMin:     cfg.MasterRangeAnglePosMin,
		MasterRangeAnglePosMax:    cfg.MasterRangeAnglePosMax,
		MasterRangeSizeMin:       cfg.MasterRangeSizeMin,
		MasterRangeSizeMax:        cfg.MasterRangeSizeMax,
		MasterRangeColors:       cfg.MasterRangeColors,
		MasterDisplayShadow:    cfg.MasterDisplayShadow,
		MasterShadowColor:      cfg.MasterShadowColor,
		MasterShadowPointX:    cfg.MasterShadowPointX,
		MasterShadowPointY:    cfg.MasterShadowPointY,
		MasterImageAlpha:      cfg.MasterImageAlpha,
		MasterUseShapeOriginalColor: cfg.MasterUseShapeOriginalColor,
		ThumbSizeWidth:       cfg.ThumbSizeWidth,
		ThumbSizeHeight:      cfg.ThumbSizeHeight,
		ThumbRangeVerifyLenMin: cfg.ThumbRangeVerifyLenMin,
		ThumbRangeVerifyLenMax: cfg.ThumbRangeVerifyLenMax,
		ThumbDisabledRangeVerifyLen: cfg.ThumbDisabledRangeVerifyLen,
		ThumbRangeThumbSizeMin:  cfg.ThumbRangeThumbSizeMin,
		ThumbRangeThumbSizeMax:  cfg.ThumbRangeThumbSizeMax,
		ThumbRangeThumbColors:  cfg.ThumbRangeThumbColors,
		ThumbRangeThumbBgColors: cfg.ThumbRangeThumbBgColors,
		ThumbIsThumbNonDeformAbility: cfg.ThumbIsThumbNonDeformAbility,
		ThumbBgDistort:     cfg.ThumbBgDistort,
		ThumbBgCirclesNum:  cfg.ThumbBgCirclesNum,
		ThumbBgSlimLineNum: cfg.ThumbBgSlimLineNum,
		Version:         cfg.Version,
	}
	return rs, nil
}

// PreviewCaptchaClick 获取点击验证码预览-后台使用(调试用)
func (s *CaptchaClickService) PreviewCaptchaClick(
	ctx *gin.Context,
	r req.PreviewCaptchaClickReq,
) (rs res.PreviewCaptchaClickRes, err error) {
	// 获取配置
	cfg, err := s.getConfig(ctx)
	if err != nil {
		return rs, biz_err.New(biz_err.SYSTEM_ERROR, "获取配置失败")
	}

	// 使用配置初始化验证码
	if err := initCaptchaClickWithConfig(cfg); err != nil {
		return rs, biz_err.New(biz_err.SYSTEM_ERROR, "验证码初始化失败")
	}

	// 生成验证码
	captData, err := captchaClickInstance.Generate()
	if err != nil {
		return rs, biz_err.New(biz_err.SYSTEM_ERROR, "生成验证码失败")
	}

	dotData := captData.GetData()
	if dotData == nil {
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

	// 转换点数据
	dots := make([]res.PreviewCaptchaClickResDot, 0, len(dotData))
	for i, dot := range dotData {
		dots = append(dots, res.PreviewCaptchaClickResDot{
			X:       dot.X,
			Y:       dot.Y,
			W:       0,
			H:       0,
			Content: fmt.Sprintf("点%d", i+1),
		})
	}

	rs = res.PreviewCaptchaClickRes{
		MasterImage: masterImage,
		ThumbImage:  thumbImage,
		Dots:       dots,
	}
	return rs, nil
}