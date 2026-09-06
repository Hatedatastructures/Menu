package systemRbac

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ConfigGroup struct {
	ID     uint   `gorm:"primarykey" json:"id"`
	Code   string `gorm:"size:50;not null;uniqueIndex" json:"code"` // system/login
	Name   string `gorm:"size:100;not null" json:"name"`

	Config datatypes.JSON `gorm:"type:json;not null" json:"config"` // 🔥关键优化

	Sort   int   `gorm:"default:0" json:"sort"`
	Status int  `gorm:"default:1" json:"status"` // 1=启用 0=禁用

	Remark string `gorm:"size:255" json:"remark"`

	Version int `gorm:"default:1" json:"version"` // 🔥版本控制

	CreatedAt time.Time     `json:"createdAt"`
	UpdatedAt time.Time     `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (ConfigGroup) TableName() string {
	return "sys_config_groups"
}