package systemRbac

import (
	"shack/internal/global"

	"gorm.io/datatypes"
)

// Template 通用模板（支持邮件、页面、消息等）
type Template struct {
	global.GVA_MODEL
	Code        string   `json:"code" gorm:"type:varchar(50);uniqueIndex;not null;comment:模板编码"`
	Name        string   `json:"name" gorm:"type:varchar(100);not null;comment:模板名称"`
	Content     string   `json:"content" gorm:"type:longtext;not null;comment:模板内容"`
	Status      string   `json:"status" gorm:"type:varchar(20);not null;default:active;comment:状态(active/inactive/draft)"`
	Description string   `json:"description" gorm:"type:text;comment:模板描述"`
	Variables   datatypes.JSONType[[]string] `json:"variables" gorm:"type:json;comment:模板变量列表"`
	Engine      string   `json:"engine" gorm:"type:varchar(20);default:html/template;comment:模板引擎"`
}

// TableName 指定表名
func (Template) TableName() string {
	return "sys_templates"
}
