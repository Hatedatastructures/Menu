package response

type GetEmailLogListRes struct {
	Page int `json:"page"`
	PageSize int `json:"pageSize"`
	ToEmail string `json:"toEmail"`
	Status string `json:"status"`
	BizType string `json:"bizType"`
	StartTime string `json:"startTime"`
	EndTime string `json:"endTime"`
	Total int64 `json:"total"` // 总数
	List []GetEmailLogListResList `json:"list"` // []
}

type GetEmailLogListResList struct {
	Id int64 `json:"id"` // 日志ID
	ToEmail string `json:"toEmail"` // 收件人邮箱
	Subject string `json:"subject"` // 邮件主题
	TemplateCode string `json:"templateCode"` // 模板编码
	Status string `json:"status"` // 发送状态
	BizType string `json:"bizType"` // 业务类型
	SentAt string `json:"sentAt"` // 发送时间
	DeliveredAt string `json:"deliveredAt"` // 送达时间
	ErrorMsg string `json:"errorMsg"` // 错误信息
	RetryCount int `json:"retryCount"` // 重试次数
}

type GetEmailLogRes struct {
	Id           int64                  `json:"id"`                      // 日志ID
	ToEmail      string                 `json:"toEmail"`                 // 收件人邮箱
	Subject      string                 `json:"subject"`                 // 邮件主题
	Content      string                 `json:"content"`                 // 邮件内容
	TemplateCode string                 `json:"templateCode"`             // 模板编码
	TemplateData map[string]interface{} `json:"templateData"`            // 模板数据
	Status       string                 `json:"status"`                  // 发送状态
	BizType      string                 `json:"bizType"`                 // 业务类型
	SentAt       string                 `json:"sentAt"`                  // 发送时间
	DeliveredAt  string                 `json:"deliveredAt"`             // 送达时间
	OpenedAt     string                 `json:"openedAt"`                // 打开时间
	ClickedAt    string                 `json:"clickedAt"`               // 点击时间
	ErrorMsg     string                 `json:"errorMsg"`                // 错误信息
	RetryCount   int                    `json:"retryCount"`              // 重试次数
	Provider     string                 `json:"provider"`                // 邮件服务商
}

type ResendEmailRes struct {
	NewTaskId string `json:"newTaskId"` // 新任务ID
}
