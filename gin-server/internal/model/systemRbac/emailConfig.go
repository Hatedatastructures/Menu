package systemRbac

import (
	"shack/internal/global"
)

// EmailConfig 邮件服务商配置
type EmailConfig struct {
	global.GVA_MODEL
	Name     string `json:"name" gorm:"type:varchar(100);not null;comment:配置名称"`
	Provider string `json:"provider" gorm:"type:varchar(50);not null;index:idx_provider;comment:服务商"`
	Config   string `json:"config" gorm:"type:json;not null;comment:配置信息JSON"`
	IsDefault bool  `json:"isDefault" gorm:"column:is_default;default:false;index:idx_is_default;comment:是否默认配置"`
	Priority int    `json:"priority" gorm:"default:0;index:idx_priority;comment:优先级"`
	Status   string `json:"status" gorm:"type:varchar(20);not null;default:active;comment:状态"`
	CreatedBy uint  `json:"createdBy" gorm:"column:created_by;comment:创建人ID"`
	UpdatedBy uint  `json:"updatedBy" gorm:"column:updated_by;comment:更新人ID"`
}

// TableName 指定表名
func (EmailConfig) TableName() string {
	return "sys_email_configs"
}

