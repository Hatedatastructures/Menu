package systemRbac

import (
	"shack/internal/global"
)

// ApiTestLog API测试记录模型
type ApiTestLog struct {
	global.GVA_MODEL
	Method      string         `gorm:"size:10;not null" json:"method"`           // 请求方法
	URL         string         `gorm:"size:500;not null" json:"url"`             // 请求URL
	StatusCode  int            `gorm:"default:0" json:"statusCode"`              // 响应状态码
	Duration    int            `gorm:"default:0" json:"duration"`                // 耗时(毫秒)
	ReqHeaders  string         `gorm:"type:text" json:"reqHeaders"`              // 请求头JSON
	ReqBody     string         `gorm:"type:text" json:"reqBody"`                 // 请求体
	ResHeaders  string         `gorm:"type:text" json:"resHeaders"`              // 响应头JSON
	ResBody     string         `gorm:"type:longtext" json:"resBody"`             // 响应体
	Params      string         `gorm:"type:text" json:"params"`                  // 请求参数JSON
	Description string         `gorm:"size:200" json:"description"`              // 描述

}

// TableName 指定表名
func (ApiTestLog) TableName() string {
	return "sys_api_test_logs"
}
