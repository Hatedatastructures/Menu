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

type ConfigEmailService struct{}

// GetEmailConfig 获取邮件配置-后台使用
func (s *ConfigEmailService) GetEmailConfig(
	ctx *gin.Context,
) (rs res.GetEmailConfigRes, err error) {
	var group systemRbac.ConfigGroup
	err = global.GVA_DB.Where("code = ?", "email").First(&group).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			rs = res.GetEmailConfigRes{
				Enabled:  false,
				Host:     "",
				Port:     465,
				Username: "",
				Password: "",
				FromName: "",
				Ssl:      true,
			}
			return rs, nil
		}
		return rs, biz_err.New(biz_err.DB_ERROR, "查询邮件配置失败")
	}

	var config map[string]interface{}
	if group.Config != nil {
		json.Unmarshal(group.Config, &config)
	}

	rs = res.GetEmailConfigRes{
		Enabled:  getBoolFromMap(config, "enabled", false),
		Host:     getStringFromMap(config, "host", ""),
		Port:     getIntFromMap(config, "port", 465),
		Username: getStringFromMap(config, "username", ""),
		Password: getStringFromMap(config, "password", ""),
		FromName: getStringFromMap(config, "fromName", ""),
		Ssl:      getBoolFromMap(config, "ssl", true),
	}
	return rs, nil
}

// SaveEmailConfig 保存邮件配置-后台使用
func (s *ConfigEmailService) SaveEmailConfig(
	ctx *gin.Context,
	r req.SaveEmailConfigReq,
) (rs res.SaveEmailConfigRes, err error) {
	configMap := map[string]interface{}{
		"enabled":  r.Enabled,
		"host":     r.Host,
		"port":     r.Port,
		"username": r.Username,
		"password": r.Password,
		"fromName": r.FromName,
		"ssl":      r.Ssl,
	}
	configJson, _ := json.Marshal(configMap)

	var group systemRbac.ConfigGroup
	err = global.GVA_DB.Where("code = ?", "email").First(&group).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			group = systemRbac.ConfigGroup{
				Code:   "email",
				Name:   "邮件配置",
				Config: datatypes.JSON(configJson),
				Status: 1,
			}
			err = global.GVA_DB.Create(&group).Error
			if err != nil {
				return rs, biz_err.New(biz_err.DB_ERROR, "创建邮件配置失败")
			}
			rs = res.SaveEmailConfigRes{Version: group.Version}
			return rs, nil
		}
		return rs, biz_err.New(biz_err.DB_ERROR, "查询邮件配置失败")
	}

	group.Config = datatypes.JSON(configJson)
	group.Version++
	err = global.GVA_DB.Save(&group).Error
	if err != nil {
		return rs, biz_err.New(biz_err.DB_ERROR, "保存邮件配置失败")
	}

	rs = res.SaveEmailConfigRes{Version: group.Version}
	return rs, nil
}

// TestEmail 测试邮件发送-后台使用
func (s *ConfigEmailService) TestEmail(
	ctx *gin.Context,
	r req.TestEmailReq,
) (rs res.TestEmailRes, err error) {
	if r.Address == "" {
		return rs, biz_err.New(biz_err.PARAM_MISSING, "收件人邮箱不能为空")
	}

	// 获取邮件配置
	config, err := GetEmailConfigFromDB()
	if err != nil {
		return rs, biz_err.New(biz_err.PARAM_ERROR, "邮件配置未设置")
	}

	if !config.Enabled {
		return rs, biz_err.New(biz_err.PARAM_ERROR, "邮件服务未启用")
	}

	// 创建邮件发送器
	mailer := NewMailer(*config)

	// 如果指定了模板ID，使用模板发送
	if r.TemplateId > 0 {
		err = mailer.SendTestEmailWithTemplate(r.Address, r.TemplateId)
	} else {
		// 否则使用默认测试邮件
		err = mailer.SendTestEmail(r.Address)
	}

	if err != nil {
		rs = res.TestEmailRes{
			Success: false,
			Message: "测试邮件发送失败: " + err.Error(),
		}
		return rs, nil
	}

	rs = res.TestEmailRes{
		Success: true,
		Message: "测试邮件已成功发送到 " + r.Address,
	}
	return rs, nil
}

