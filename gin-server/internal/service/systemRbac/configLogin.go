package systemRbac

import (
	"encoding/json"

	"shack/internal/global"
	"shack/internal/model/systemRbac"
	req "shack/internal/model/systemRbac/request"
	res "shack/internal/model/systemRbac/response"

	biz_err "shack/internal/error"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ConfigLoginService struct{}

// GetLoginConfig 获取登录配置-后台使用
func (s *ConfigLoginService) GetLoginConfig(
	ctx *gin.Context,
) (rs res.GetLoginConfigRes, err error) {
	var group systemRbac.ConfigGroup
	err = global.GVA_DB.Where("code = ?", "login").First(&group).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			rs = res.GetLoginConfigRes{
				CaptchaEnabled:        false,
				CaptchaType:           "image",
				MaxRetryCount:         5,
				LockTime:              30,
				RememberMe:            true,
				SingleLogin:           false,
				SliderCaptchaWidth:    320,
				SliderCaptchaHeight:   160,
				SliderThumbWidth:      50,
				SliderThumbHeight:     50,
				SliderVerticalPadding: 10,
				SliderHorizontalPadding: 10,
				SliderShowTheme:       true,
				SliderTitle:           "请向右滑动完成验证",
				SliderButtonText:      "向右滑动",
				SliderIconSize:        24,
				SliderDotSize:         6,
			}
			return rs, nil
		}
		return rs, biz_err.New(biz_err.DB_ERROR, "查询登录配置失败")
	}

	var config map[string]interface{}
	if group.Config != nil {
		json.Unmarshal(group.Config, &config)
	}

	rs = res.GetLoginConfigRes{
		CaptchaEnabled:        getBoolFromMap(config, "captchaEnabled", false),
		CaptchaType:           getStringFromMap(config, "captchaType", "image"),
		MaxRetryCount:        getIntFromMap(config, "maxRetryCount", 5),
		LockTime:             getIntFromMap(config, "lockTime", 30),
		RememberMe:           getBoolFromMap(config, "rememberMe", true),
		SingleLogin:          getBoolFromMap(config, "singleLogin", false),
		SliderCaptchaWidth:   getIntFromMap(config, "sliderCaptchaWidth", 320),
		SliderCaptchaHeight:  getIntFromMap(config, "sliderCaptchaHeight", 160),
		SliderThumbWidth:     getIntFromMap(config, "sliderThumbWidth", 50),
		SliderThumbHeight:    getIntFromMap(config, "sliderThumbHeight", 50),
		SliderVerticalPadding: getIntFromMap(config, "sliderVerticalPadding", 10),
		SliderHorizontalPadding: getIntFromMap(config, "sliderHorizontalPadding", 10),
		SliderShowTheme:      getBoolFromMap(config, "sliderShowTheme", true),
		SliderTitle:          getStringFromMap(config, "sliderTitle", "请向右滑动完成验证"),
		SliderButtonText:     getStringFromMap(config, "sliderButtonText", "向右滑动"),
		SliderIconSize:       getIntFromMap(config, "sliderIconSize", 24),
		SliderDotSize:        getIntFromMap(config, "sliderDotSize", 6),
	}
	return rs, nil
}

// SaveLoginConfig 保存登录配置-后台使用
func (s *ConfigLoginService) SaveLoginConfig(
	ctx *gin.Context,
	r req.SaveLoginConfigReq,
) (rs res.SaveLoginConfigRes, err error) {
	configMap := map[string]interface{}{
		"captchaEnabled":           r.CaptchaEnabled,
		"captchaType":              r.CaptchaType,
		"maxRetryCount":            r.MaxRetryCount,
		"lockTime":                 r.LockTime,
		"rememberMe":               r.RememberMe,
		"singleLogin":              r.SingleLogin,
		"sliderCaptchaWidth":       r.SliderCaptchaWidth,
		"sliderCaptchaHeight":      r.SliderCaptchaHeight,
		"sliderThumbWidth":         r.SliderThumbWidth,
		"sliderThumbHeight":        r.SliderThumbHeight,
		"sliderVerticalPadding":    r.SliderVerticalPadding,
		"sliderHorizontalPadding":  r.SliderHorizontalPadding,
		"sliderShowTheme":          r.SliderShowTheme,
		"sliderTitle":              r.SliderTitle,
		"sliderButtonText":         r.SliderButtonText,
		"sliderIconSize":           r.SliderIconSize,
		"sliderDotSize":            r.SliderDotSize,
	}
	configJson, _ := json.Marshal(configMap)

	var group systemRbac.ConfigGroup
	err = global.GVA_DB.Where("code = ?", "login").First(&group).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			group = systemRbac.ConfigGroup{
				Code:   "login",
				Name:   "登录配置",
				Config: datatypes.JSON(configJson),
				Status: 1,
			}
			err = global.GVA_DB.Create(&group).Error
			if err != nil {
				return rs, biz_err.New(biz_err.DB_ERROR, "创建登录配置失败")
			}
			rs = res.SaveLoginConfigRes{Version: group.Version}
			return rs, nil
		}
		return rs, biz_err.New(biz_err.DB_ERROR, "查询登录配置失败")
	}

	group.Config = datatypes.JSON(configJson)
	group.Version++
	err = global.GVA_DB.Save(&group).Error
	if err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR, "保存登录配置失败")
	}

	rs = res.SaveLoginConfigRes{Version: group.Version}
	return rs, nil
}

// 辅助函数：从map中获取int
func getIntFromMap(m map[string]interface{}, key string, defaultVal int) int {
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

