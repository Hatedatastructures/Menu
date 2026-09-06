package systemRbac

import (
	"shack/internal/global"
)

// EmailTemplate 邮件模板
type EmailTemplate struct {
	global.GVA_MODEL
	Code        string `json:"code" gorm:"type:varchar(50);uniqueIndex;not null;comment:模板编码"`
	Name        string `json:"name" gorm:"type:varchar(100);not null;comment:模板名称"`
	Subject     string `json:"subject" gorm:"type:varchar(200);not null;comment:邮件主题"`
	Content     string `json:"content" gorm:"type:longtext;not null;comment:邮件内容HTML"`
	Type        string `json:"type" gorm:"type:varchar(20);not null;default:system;comment:模板类型"`
	Status      string `json:"status" gorm:"type:varchar(20);not null;default:active;comment:状态"`
	Description string `json:"description" gorm:"type:text;comment:模板描述"`
	Variables   string `json:"variables" gorm:"type:json;comment:模板变量定义JSON"`
	Tags        string `json:"tags" gorm:"type:varchar(200);comment:标签"`
	CreatedBy   uint   `json:"createdBy" gorm:"column:created_by;comment:创建人ID"`
	UpdatedBy   uint   `json:"updatedBy" gorm:"column:updated_by;comment:更新人ID"`
}

// TableName 指定表名
func (EmailTemplate) TableName() string {
	return "sys_email_templates"
}
