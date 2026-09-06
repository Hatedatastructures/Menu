package systemRbac

import (
	"shack/internal/global"
)

// EmailBlacklist 邮件黑名单
type EmailBlacklist struct {
	global.GVA_MODEL
	Email     string `json:"email" gorm:"type:varchar(200);uniqueIndex;not null;comment:被拉黑的邮箱"`
	Reason    string `json:"reason" gorm:"type:text;comment:拉黑原因"`
	Type      string `json:"type" gorm:"type:varchar(20);not null;default:manual;comment:拉黑类型"`
	Status    string `json:"status" gorm:"type:varchar(20);not null;default:active;comment:状态"`
	CreatedBy uint   `json:"createdBy" gorm:"column:created_by;comment:创建人ID"`
}

// TableName 指定表名
func (EmailBlacklist) TableName() string {
	return "sys_email_blacklist"
}


