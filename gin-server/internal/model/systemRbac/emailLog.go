package systemRbac

import (
	"shack/internal/global"
	"time"
)

// EmailLog 邮件发送日志
type EmailLog struct {
	global.GVA_MODEL
	TaskID       *uint            `json:"taskId" gorm:"column:task_id;comment:关联的定时任务ID"`
	BatchID      string           `json:"batchId" gorm:"type:varchar(50);column:batch_id;comment:批次ID"`
	ToEmail      string           `json:"toEmail" gorm:"type:varchar(200);not null;column:to_email;index:idx_to_email;comment:收件人邮箱"`
	CcEmail      string           `json:"ccEmail" gorm:"type:text;column:cc_email;comment:抄送邮箱"`
	BccEmail     string           `json:"bccEmail" gorm:"type:text;column:bcc_email;comment:密送邮箱"`
	Subject      string           `json:"subject" gorm:"type:varchar(500);not null;comment:邮件主题"`
	Content      string           `json:"content" gorm:"type:longtext;comment:邮件内容"`
	TemplateCode string           `json:"templateCode" gorm:"type:varchar(50);column:template_code;index:idx_template_code;comment:模板编码"`
	TemplateData string           `json:"templateData" gorm:"type:json;column:template_data;comment:模板数据JSON"`
	Status       string           `json:"status" gorm:"type:varchar(20);not null;default:pending;index:idx_status;comment:发送状态"`
	BizType      string           `json:"bizType" gorm:"type:varchar(50);column:biz_type;index:idx_biz_type;comment:业务类型"`
	Provider     string           `json:"provider" gorm:"type:varchar(50);comment:邮件服务商"`
	ErrorMsg     string           `json:"errorMsg" gorm:"type:text;column:error_msg;comment:错误信息"`
	RetryCount   int              `json:"retryCount" gorm:"column:retry_count;default:0;comment:重试次数"`
	SentAt       *time.Time       `json:"sentAt" gorm:"column:sent_at;index:idx_sent_at;comment:发送时间"`
	DeliveredAt  *time.Time       `json:"deliveredAt" gorm:"column:delivered_at;comment:送达时间"`
	OpenedAt     *time.Time       `json:"openedAt" gorm:"column:opened_at;comment:打开时间"`
	ClickedAt    *time.Time       `json:"clickedAt" gorm:"column:clicked_at;comment:点击时间"`
}

// TableName 指定表名
func (EmailLog) TableName() string {
	return "sys_email_logs"
}
