package request

type SendEmailReq struct {
	To []string `json:"to" form:"to"` // 收件人邮箱列表
	Cc []string `json:"cc" form:"cc"` // 抄送邮箱列表
	Bcc []string `json:"bcc" form:"bcc"` // 密送邮箱列表
	Subject string `json:"subject" form:"subject"` // 邮件主题
	Content string `json:"content" form:"content"` // 邮件内容
	Type string `json:"type" form:"type"` // 邮件类型: html/text
	Attachments []SendEmailReqAttachment `json:"attachments" form:"attachments"` // []
}

type SendEmailReqAttachment struct {
	Filename string `json:"filename" form:"filename"` // 附件文件名
	Url string `json:"url" form:"url"` // 附件URL地址
}

type SendTemplateEmailReq struct {
	TemplateCode string                 `json:"templateCode" form:"templateCode"` // 模板编码
	To           []string               `json:"to" form:"to"`                     // 收件人邮箱列表
	Data         map[string]interface{} `json:"data" form:"data"`                 // 模板变量数据
}

type SendTemplateEmailReqRe struct {
	TaskId string `json:"taskId" form:"taskId"` // 邮件任务ID
	TemplateName string `json:"templateName" form:"templateName"` // 模板名称
}

type SendBatchEmailReq struct {
	TemplateCode string `json:"templateCode" form:"templateCode"` // 模板编码
	Recipients []SendBatchEmailReqRecipient `json:"recipients" form:"recipients"` // []
	SendMode string `json:"sendMode" form:"sendMode"` // 发送模式: instant/scheduled
	ScheduledTime string `json:"scheduledTime" form:"scheduledTime"` // 定时发送时间
}

type SendBatchEmailReqRecipient struct {
	Email string                 `json:"email" form:"email"` // 收件人邮箱
	Data  map[string]interface{} `json:"data" form:"data"`   // 该收件人的模板变量数据
}
