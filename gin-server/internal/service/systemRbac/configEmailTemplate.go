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

type ConfigEmailTemplateService struct{}

// GetEmailTemplateConfig 获取邮件模板配置-后台使用
func (s *ConfigEmailTemplateService) GetEmailTemplateConfig(
	ctx *gin.Context,
) (rs res.GetEmailTemplateConfigRes, err error) {
	var group systemRbac.ConfigGroup
	err = global.GVA_DB.Where("code = ?", "emailTemplate").First(&group).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			rs = res.GetEmailTemplateConfigRes{
				VerifyCode:   "",
				ResetPassword: "",
				Welcome:      "",
			}
			return rs, nil
		}
		return rs, biz_err.New(biz_err.DB_ERROR, "查询邮件模板配置失败")
	}

	var config map[string]interface{}
	if group.Config != nil {
		json.Unmarshal(group.Config, &config)
	}

	rs = res.GetEmailTemplateConfigRes{
		VerifyCode:   getStringFromMap(config, "verifyCode", ""),
		ResetPassword: getStringFromMap(config, "resetPassword", ""),
		Welcome:      getStringFromMap(config, "welcome", ""),
	}
	return rs, nil
}

// SaveEmailTemplateConfig 保存邮件模板配置-后台使用
func (s *ConfigEmailTemplateService) SaveEmailTemplateConfig(
	ctx *gin.Context,
	r req.SaveEmailTemplateConfigReq,
) (rs res.SaveEmailTemplateConfigRes, err error) {
	configMap := map[string]interface{}{
		"verifyCode":    r.VerifyCode,
		"resetPassword": r.ResetPassword,
		"welcome":       r.Welcome,
	}
	configJson, _ := json.Marshal(configMap)

	var group systemRbac.ConfigGroup
	err = global.GVA_DB.Where("code = ?", "emailTemplate").First(&group).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			group = systemRbac.ConfigGroup{
				Code:   "emailTemplate",
				Name:   "邮件模板配置",
				Config: datatypes.JSON(configJson),
				Status: 1,
			}
			err = global.GVA_DB.Create(&group).Error
			if err != nil {
				return rs, biz_err.New(biz_err.DB_ERROR, "创建邮件模板配置失败")
			}
			rs = res.SaveEmailTemplateConfigRes{Version: group.Version}
			return rs, nil
		}
		return rs, biz_err.New(biz_err.DB_ERROR, "查询邮件模板配置失败")
	}

	group.Config = datatypes.JSON(configJson)
	group.Version++
	err = global.GVA_DB.Save(&group).Error
	if err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR, "保存邮件模板配置失败")
	}

	rs = res.SaveEmailTemplateConfigRes{Version: group.Version}
	return rs, nil
}

