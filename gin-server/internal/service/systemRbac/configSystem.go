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

type ConfigSystemService struct{}

// GetSystemConfig 获取系统配置-后台使用
func (s *ConfigSystemService) GetSystemConfig(
	ctx *gin.Context,
) (rs res.GetSystemConfigRes, err error) {
	var group systemRbac.ConfigGroup
	err = global.GVA_DB.Where("code = ?", "system").First(&group).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 返回默认值
			rs = res.GetSystemConfigRes{
				SiteName:           "Mars Admin",
				SiteDescription:    "",
				SiteLogo:           "",
				Copyright:          "",
				Icp:                "",
				WatermarkEnabled:   true,
				WatermarkType:      "username",
				WatermarkCustomText: "",
				WatermarkOpacity:   0.1,
			}
			return rs, nil
		}
		return rs, biz_err.New(biz_err.DB_ERROR, "查询系统配置失败")
	}

	// 解析JSON配置
	var config map[string]interface{}
	if group.Config != nil {
		err = json.Unmarshal(group.Config, &config)
		if err != nil {
			return rs, biz_err.New(biz_err.JSON_PARSE, "解析系统配置失败")
		}
	}

	// 填充返回数据
	rs = res.GetSystemConfigRes{
		SiteName:           getStringFromMap(config, "siteName", "Mars Admin"),
		SiteDescription:    getStringFromMap(config, "siteDescription", ""),
		SiteLogo:           getStringFromMap(config, "siteLogo", ""),
		Copyright:          getStringFromMap(config, "copyright", ""),
		Icp:                getStringFromMap(config, "icp", ""),
		WatermarkEnabled:   getBoolFromMap(config, "watermarkEnabled", true),
		WatermarkType:      getStringFromMap(config, "watermarkType", "username"),
		WatermarkCustomText: getStringFromMap(config, "watermarkCustomText", ""),
		WatermarkOpacity:   getFloat32FromMap(config, "watermarkOpacity", 0.1),
	}
	return rs, nil
}

// SaveSystemConfig 保存系统配置-后台使用
func (s *ConfigSystemService) SaveSystemConfig(
	ctx *gin.Context,
	r req.SaveSystemConfigReq,
) (rs res.SaveSystemConfigRes, err error) {
	// 构建配置JSON
	configMap := map[string]interface{}{
		"siteName":             r.SiteName,
		"siteDescription":      r.SiteDescription,
		"siteLogo":             r.SiteLogo,
		"copyright":            r.Copyright,
		"icp":                 r.Icp,
		"watermarkEnabled":    r.WatermarkEnabled,
		"watermarkType":       r.WatermarkType,
		"watermarkCustomText": r.WatermarkCustomText,
		"watermarkOpacity":    r.WatermarkOpacity,
	}
	configJson, err := json.Marshal(configMap)
	if err != nil {
		return rs, biz_err.New(biz_err.JSON_PARSE, "序列化配置失败")
	}

		var group systemRbac.ConfigGroup
	err = global.GVA_DB.Where("code = ?", "system").First(&group).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 不存在则创建
			group = systemRbac.ConfigGroup{
				Code:   "system",
				Name:   "系统配置",
				Config: datatypes.JSON(configJson),
				Status: 1,
			}
			err = global.GVA_DB.Create(&group).Error
			if err != nil {
				return rs, biz_err.New(biz_err.DB_ERROR, "创建系统配置失败")
			}
			rs = res.SaveSystemConfigRes{
				Version: group.Version,
			}
			return rs, nil
		}
		return rs, biz_err.New(biz_err.DB_ERROR, "查询系统配置失败")
	}

	// 更新配置
	group.Config = datatypes.JSON(configJson)
	group.Version++

	err = global.GVA_DB.Save(&group).Error
	if err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR, "保存系统配置失败")
	}

	rs = res.SaveSystemConfigRes{
		Version: group.Version,
	}
	return rs, nil
}

// 辅助函数：从map中获取string
func getStringFromMap(m map[string]interface{}, key string, defaultVal string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return defaultVal
}

// 辅助函数：从map中获取bool
func getBoolFromMap(m map[string]interface{}, key string, defaultVal bool) bool {
	if v, ok := m[key]; ok {
		if b, ok := v.(bool); ok {
			return b
		}
	}
	return defaultVal
}

// 辅助函数：从map中获取float32
func getFloat32FromMap(m map[string]interface{}, key string, defaultVal float32) float32 {
	if v, ok := m[key]; ok {
		switch val := v.(type) {
		case float64:
			return float32(val)
		case float32:
			return val
		}
	}
	return defaultVal
}

