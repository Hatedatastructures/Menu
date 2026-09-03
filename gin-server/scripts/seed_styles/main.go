// scripts/seed_styles/main.go
// Usage: go run scripts/seed_styles/main.go
package main

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// StringArray 用于存储字符串数组
type StringArray []string

func (sa StringArray) Value() (driver.Value, error) {
	if sa == nil {
		return nil, nil
	}
	return json.Marshal(sa)
}

func (sa *StringArray) Scan(value interface{}) error {
	if value == nil {
		*sa = make([]string, 0)
		return nil
	}
	var err error
	switch value.(type) {
	case []byte:
		err = json.Unmarshal(value.([]byte), sa)
	case string:
		err = json.Unmarshal([]byte(value.(string)), sa)
	default:
		err = errors.New("StringArray.Scan: invalid value type")
	}
	if err != nil {
		return err
	}
	return nil
}

// Style 风格表
type Style struct {
	ID             uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	Name           string         `gorm:"type:varchar(100);not null;comment:风格名称" json:"name"`
	Description    string         `gorm:"type:text;comment:风格描述" json:"description"`
	MarkdownConfig string         `gorm:"type:text;comment:Markdown编辑器配置代码" json:"markdownConfig"`
	Tags           StringArray    `gorm:"type:json;comment:标签数组" json:"tags"`
	PreviewImage   string         `gorm:"type:varchar(500);comment:预览图URL" json:"previewImage"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Style) TableName() string {
	return "ppt_styles"
}

type StyleJSON struct {
	Style       string   `json:"style"`
	StyleName   string   `json:"styleName"`
	Aliases     []string `json:"aliases"`
	Description string   `json:"description"`
	Category    string   `json:"category"`
	Source      string   `json:"source"`
	StyleSkill  string   `json:"styleSkill"`
}

func main() {
	// 读取 styles.json
	data, err := os.ReadFile("d:/coder/reactGin/html-ppt/html-my-ppt/docs/styles.json")
	if err != nil {
		log.Fatalf("读取 styles.json 失败: %v", err)
	}

	var styles []StyleJSON
	if err := json.Unmarshal(data, &styles); err != nil {
		log.Fatalf("解析 JSON 失败: %v", err)
	}

	fmt.Printf("读取到 %d 个风格\n", len(styles))

	// 连接数据库
	// 修改这里的DSN为你的数据库配置
	dsn := "root:root@tcp(127.0.0.1:3306)/html-ppt?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}

	// 迁移表结构
	if err := db.AutoMigrate(&Style{}); err != nil {
		log.Fatalf("迁移表结构失败: %v", err)
	}

	// 插入数据
	inserted := 0
	skipped := 0
	for _, s := range styles {
		// 检查是否已存在
		var count int64
		db.Model(&Style{}).Where("name = ?", s.StyleName).Count(&count)
		if count > 0 {
			fmt.Printf("跳过已存在: %s\n", s.StyleName)
			skipped++
			continue
		}

		// 合并 aliases 到 tags
		tags := s.Aliases
		if s.Category != "" {
			tags = append(tags, s.Category)
		}

		style := Style{
			ID:             uuid.New(),
			Name:           s.StyleName,
			Description:    s.Description,
			MarkdownConfig: s.StyleSkill,
			Tags:           StringArray(tags),
			PreviewImage:   "", // styles.json 中没有预览图
		}

		if err := db.Create(&style).Error; err != nil {
			log.Printf("插入失败 %s: %v", s.StyleName, err)
			continue
		}
		fmt.Printf("插入成功: %s\n", s.StyleName)
		inserted++
	}

	fmt.Printf("\n完成! 新增: %d, 跳过: %d\n", inserted, skipped)
}
