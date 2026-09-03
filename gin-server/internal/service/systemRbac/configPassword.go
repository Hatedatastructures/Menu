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

type ConfigPasswordService struct{}

// GetPasswordConfig 获取密码配置-后台使用
func (s *ConfigPasswordService) GetPasswordConfig(
	ctx *gin.Context,
) (rs res.GetPasswordConfigRes, err error) {
	var group systemRbac.ConfigGroup
	err = global.GVA_DB.Where("code = ?", "password").First(&group).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			rs = res.GetPasswordConfigRes{
				MinLength:        6,
				MaxLength:        20,
				RequireUppercase: false,
				RequireLowercase: false,
				RequireNumber:    false,
				RequireSpecial:   false,
				ExpireDays:       0,
			}
			return rs, nil
		}
		return rs, biz_err.New(biz_err.DB_ERROR, "查询密码配置失败")
	}

	var config map[string]interface{}
	if group.Config != nil {
		json.Unmarshal(group.Config, &config)
	}

	rs = res.GetPasswordConfigRes{
		MinLength:        getIntFromMap(config, "minLength", 6),
		MaxLength:        getIntFromMap(config, "maxLength", 20),
		RequireUppercase: getBoolFromMap(config, "requireUppercase", false),
		RequireLowercase: getBoolFromMap(config, "requireLowercase", false),
		RequireNumber:    getBoolFromMap(config, "requireNumber", false),
		RequireSpecial:   getBoolFromMap(config, "requireSpecial", false),
		ExpireDays:       getIntFromMap(config, "expireDays", 0),
	}
	return rs, nil
}

// SavePasswordConfig 保存密码配置-后台使用
func (s *ConfigPasswordService) SavePasswordConfig(
	ctx *gin.Context,
	r req.SavePasswordConfigReq,
) (rs res.SavePasswordConfigRes, err error) {
	configMap := map[string]interface{}{
		"minLength":         r.MinLength,
		"maxLength":         r.MaxLength,
		"requireUppercase": r.RequireUppercase,
		"requireLowercase": r.RequireLowercase,
		"requireNumber":    r.RequireNumber,
		"requireSpecial":   r.RequireSpecial,
		"expireDays":       r.ExpireDays,
	}
	configJson, _ := json.Marshal(configMap)

	var group systemRbac.ConfigGroup
	err = global.GVA_DB.Where("code = ?", "password").First(&group).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			group = systemRbac.ConfigGroup{
				Code:   "password",
				Name:   "密码配置",
				Config: datatypes.JSON(configJson),
				Status: 1,
			}
			err = global.GVA_DB.Create(&group).Error
			if err != nil {
				return rs, biz_err.New(biz_err.DB_ERROR, "创建密码配置失败")
			}
			rs = res.SavePasswordConfigRes{Version: group.Version}
			return rs, nil
		}
		return rs, biz_err.New(biz_err.DB_ERROR, "查询密码配置失败")
	}

	group.Config = datatypes.JSON(configJson)
	group.Version++
	err = global.GVA_DB.Save(&group).Error
	if err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR, "保存密码配置失败")
	}

	rs = res.SavePasswordConfigRes{Version: group.Version}
	return rs, nil
}

