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

type ConfigRegisterService struct{}

// GetRegisterConfig 获取注册配置-后台使用
func (s *ConfigRegisterService) GetRegisterConfig(
	ctx *gin.Context,
) (rs res.GetRegisterConfigRes, err error) {
	var group systemRbac.ConfigGroup
	err = global.GVA_DB.Where("code = ?", "register").First(&group).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 返回默认值
			rs = res.GetRegisterConfigRes{
				Enabled:      true,
				VerifyEmail:  false,
				VerifyPhone:  false,
				DefaultRole:  "user",
				NeedAudit:    false,
				DefaultImage: "",
			}
			return rs, nil
		}
		return rs, biz_err.New(biz_err.DB_ERROR, "查询注册配置失败")
	}

	var config map[string]interface{}
	if group.Config != nil {
		json.Unmarshal(group.Config, &config)
	}

	rs = res.GetRegisterConfigRes{
		Enabled:      getBoolFromMap(config, "enabled", true),
		VerifyEmail:  getBoolFromMap(config, "verifyEmail", false),
		VerifyPhone:  getBoolFromMap(config, "verifyPhone", false),
		DefaultRole:  getStringFromMap(config, "defaultRole", "user"),
		NeedAudit:    getBoolFromMap(config, "needAudit", false),
		DefaultImage: getStringFromMap(config, "defaultImage", ""),
	}
	return rs, nil
}

// SaveRegisterConfig 保存注册配置-后台使用
func (s *ConfigRegisterService) SaveRegisterConfig(
	ctx *gin.Context,
	r req.SaveRegisterConfigReq,
) (rs res.SaveRegisterConfigRes, err error) {
	configMap := map[string]interface{}{
		"enabled":      r.Enabled,
		"verifyEmail":  r.VerifyEmail,
		"verifyPhone":  r.VerifyPhone,
		"defaultRole":  r.DefaultRole,
		"needAudit":    r.NeedAudit,
		"defaultImage": r.DefaultImage,
	}
	configJson, _ := json.Marshal(configMap)

	var group systemRbac.ConfigGroup
	err = global.GVA_DB.Where("code = ?", "register").First(&group).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			group = systemRbac.ConfigGroup{
				Code:   "register",
				Name:   "注册配置",
				Config: datatypes.JSON(configJson),
				Status: 1,
			}
			err = global.GVA_DB.Create(&group).Error
			if err != nil {
				return rs, biz_err.New(biz_err.DB_ERROR, "创建注册配置失败")
			}
			rs = res.SaveRegisterConfigRes{Version: group.Version}
			return rs, nil
		}
		return rs, biz_err.New(biz_err.DB_ERROR, "查询注册配置失败")
	}

	group.Config = datatypes.JSON(configJson)
	group.Version++
	err = global.GVA_DB.Save(&group).Error
	if err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR, "保存注册配置失败")
	}

	rs = res.SaveRegisterConfigRes{Version: group.Version}
	return rs, nil
}

