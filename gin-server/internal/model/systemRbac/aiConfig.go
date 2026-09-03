package systemRbac

import (
	"time"

	"gorm.io/gorm"
)

// AiConfig AI配置模型
type AiConfig struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Name      string         `gorm:"size:100;not null"  json:"name"`                // 配置名称
	Provider string          	`gorm:"size:100;not null"  json:"provider"`
	BaseURL   string         `gorm:"size:500;not null" json:"baseUrl"`             // 接口地址
	ApiKey    string         `gorm:"type:text;not null" json:"apiKey"`              // 密钥（使用TEXT类型支持更长的key）
	Model     string         `gorm:"size:100;not null" json:"model"`               // 模型名称
	Type      string         `gorm:"size:50;not null;default:openai" json:"type"`  // 类型(openai/兼容openai)
	Enabled   bool           `gorm:"default:true" json:"enabled"`                 // 是否启用
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (AiConfig) TableName() string {
	return "sys_ai_configs"
}
