package systemRbac

import (
	"shack/internal/global"
)

// EmailLimitLog 限流日志
type EmailLimitLog struct {
	global.GVA_MODEL
	Email       string          `json:"email" gorm:"type:varchar(200);column:email;index:idx_email;comment:触犯限流的邮箱"`
	IP          string          `json:"ip" gorm:"type:varchar(50);column:ip;index:idx_ip;comment:触犯限流的IP"`
	LimitType   string          `json:"limitType" gorm:"type:varchar(20);not null;column:limit_type;index:idx_limit_type;comment:限流类型"`
	Action      string          `json:"action" gorm:"type:varchar(50);not null;comment:触犯限流的操作"`
	Blocked     bool            `json:"blocked" gorm:"default:true;comment:是否被阻止"`
	Reason      string          `json:"reason" gorm:"type:text;comment:限流原因"`
	RequestData string          `json:"requestData" gorm:"type:json;column:request_data;comment:请求数据JSON"`
}

// TableName 指定表名
func (EmailLimitLog) TableName() string {
	return "sys_email_limit_logs"
}
