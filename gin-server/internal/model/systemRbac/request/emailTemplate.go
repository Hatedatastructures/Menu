package request

type CreateEmailTemplateReq struct {
	Code string `json:"code" form:"code"` // 模板编码(唯一)
	Name string `json:"name" form:"name"` // 模板名称
	Subject string `json:"subject" form:"subject"` // 邮件主题
	Content string `json:"content" form:"content"` // 邮件内容(HTML)
	Type string `json:"type" form:"type"` // 模板类型:system/marketing/notification
	Description string `json:"description" form:"description"` // 模板描述
	Variables []CreateEmailTemplateReqVariable `json:"variables" form:"variables"` // []
}

type CreateEmailTemplateReqVariable struct {
	Name string `json:"name" form:"name"` // 变量名
	Description string `json:"description" form:"description"` // 变量描述
	Required bool `json:"required" form:"required"` // 是否必填
}

type UpdateEmailTemplateReq struct {
	Id int `json:"id" form:"id"`
	Code string `json:"code" form:"code"` // 模板编码(唯一)
	Name string `json:"name" form:"name"` // 模板名称
	Subject string `json:"subject" form:"subject"` // 邮件主题
	Content string `json:"content" form:"content"` // 邮件内容(HTML)
	Type string `json:"type" form:"type"` // 模板类型:system/marketing/notification
	Description string `json:"description" form:"description"` // 模板描述
	Variables []UpdateEmailTemplateReqVariable `json:"variables" form:"variables"` // []
}

type UpdateEmailTemplateReqVariable struct {
	Name string `json:"name" form:"name"` // 变量名
	Description string `json:"description" form:"description"` // 变量描述
	Required bool `json:"required" form:"required"` // 是否必填
}

type DeleteEmailTemplateReq struct {
	Id int `json:"id" form:"id"`
}

type GetEmailTemplateListReq struct {
	Page int `json:"page" form:"page"`
	PageSize int `json:"pageSize" form:"pageSize"`
	Type string `json:"type" form:"type"`
	Keyword string `json:"keyword" form:"keyword"`
	Status string `json:"status" form:"status"`
}

type GetEmailTemplateReq struct {
	Id int `json:"id" form:"id"`
}

type PreviewEmailTemplateReq struct {
	TemplateId int64                  `json:"templateId" form:"templateId"` // 模板ID
	Data       map[string]interface{} `json:"data" form:"data"`              // 模板变量数据
}

type PreviewEmailTemplateReqRe struct {
	Subject string `json:"subject" form:"subject"` // 渲染后的主题
	Html string `json:"html" form:"html"` // 渲染后的HTML内容
	Text string `json:"text" form:"text"` // 纯文本内容
}

type AiGenerateEmailTemplateReq struct {
	Prompt string `json:"prompt" form:"prompt"` // 模板描述
	Style string `json:"style" form:"style"` // 风格: modern/classic/minimal
	Language string `json:"language" form:"language"` // 语言: zh-CN/en-US
}
