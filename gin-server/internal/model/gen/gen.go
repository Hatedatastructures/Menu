package gen

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"shack/internal/global"
)

// JSON json字段类型,支持数据库存取
type JSON json.RawMessage

func (j JSON) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return string(j), nil
}

func (j *JSON) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	switch v := value.(type) {
	case []byte:
		*j = append((*j)[0:0], v...)
	case string:
		*j = append((*j)[0:0], []byte(v)...)
	default:
		return errors.New("gen.JSON.Scan: invalid value type")
	}
	return nil
}

func (j JSON) MarshalJSON() ([]byte, error) {
	if j == nil {
		return []byte("null"), nil
	}
	return j, nil
}

func (j *JSON) UnmarshalJSON(data []byte) error {
	if j == nil {
		return errors.New("gen.JSON.UnmarshalJSON: nil pointer")
	}
	*j = append((*j)[0:0], data...)
	return nil
}

// GenApiKey 用户的 DeepSeek API Key 存储
type GenApiKey struct {
	global.GVA_MODEL
	UserID  uint   `gorm:"index;not null;comment:用户ID" json:"userId"`
	ApiKey  string `gorm:"type:text;not null;comment:DeepSeek API Key" json:"apiKey"`
	BaseURL string `gorm:"size:500;default:https://api.deepseek.com;comment:API地址" json:"baseUrl"`
	Model   string `gorm:"size:100;default:deepseek-chat;comment:模型名称" json:"model"`
}

func (GenApiKey) TableName() string {
	return "gen_api_keys"
}

// GenHistory 生成历史记录
type GenHistory struct {
	global.GVA_MODEL
	UserID        uint `gorm:"index;not null;comment:用户ID" json:"userId"`
	FileName      string `gorm:"size:255;comment:原始文件名" json:"fileName"`
	FileID        uint   `gorm:"index;comment:关联的文件ID(sys_rbac_files)" json:"fileId"`
	QuestionCount int    `gorm:"default:0;comment:题目数量" json:"questionCount"`
	Status        string `gorm:"size:20;default:pending;comment:状态: pending/generating/done/failed" json:"status"`
	ErrMsg        string `gorm:"type:text;comment:错误信息" json:"errMsg"`
	RawContent    string `gorm:"type:text;comment:PDF提取的原始文本" json:"rawContent,omitempty"`
	Result        JSON   `gorm:"type:text;comment:AI生成的答案及解析(JSON)" json:"result,omitempty"`
}

func (GenHistory) TableName() string {
	return "gen_histories"
}
