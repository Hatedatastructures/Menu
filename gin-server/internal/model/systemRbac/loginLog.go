package systemRbac

import (
	"time"
	"gorm.io/gorm"
)

// LoginLog 用户登录日志
type LoginLog struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"size:50;not null" json:"username"`
	IP        string         `gorm:"size:50;not null" json:"ip"`
	Location  string         `gorm:"size:100" json:"location"`
	Browser   string         `gorm:"size:50" json:"browser"`
	OS        string         `gorm:"size:100" json:"os"`
	Status    string         `gorm:"size:20" json:"status"`
	LoginTime time.Time      `json:"loginTime"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
func (LoginLog) TableName() string {
	return "sys_rbac_LoginLog"
}