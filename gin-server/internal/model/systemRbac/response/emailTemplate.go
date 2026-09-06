package response

type CreateEmailTemplateRes struct {
	Id int64 `json:"id"` // 模板ID
	Code string `json:"code"` // 模板编码
	CreatedAt string `json:"createdAt"` // 创建时间
}

type GetEmailTemplateListRes struct {
	Page int `json:"page"`
	PageSize int `json:"pageSize"`
	Type string `json:"type"`
	Keyword string `json:"keyword"`
	Status string `json:"status"`
	Total int64 `json:"total"` // 总数
	List []GetEmailTemplateListResList `json:"list"` // []
}

type GetEmailTemplateListResList struct {
	Id int64 `json:"id"` // 模板ID
	Code string `json:"code"` // 模板编码
	Name string `json:"name"` // 模板名称
	Type string `json:"type"` // 模板类型
	Status string `json:"status"` // 状态: active/inactive
	CreatedAt string `json:"createdAt"` // 创建时间
	UpdatedAt string `json:"updatedAt"` // 更新时间
}

type GetEmailTemplateRes struct {
	Id int64 `json:"id"` // 模板ID
	Code string `json:"code"` // 模板编码
	Name string `json:"name"` // 模板名称
	Subject string `json:"subject"` // 邮件主题
	Content string `json:"content"` // 邮件内容
	Type string `json:"type"` // 模板类型
	Status string `json:"status"` // 状态
	Description string `json:"description"` // 模板描述
	Variables []GetEmailTemplateResVariable `json:"variables"` // []
	CreatedAt string `json:"createdAt"` // 创建时间
	UpdatedAt string `json:"updatedAt"` // 更新时间
}

type GetEmailTemplateResVariable struct {
	Name string `json:"name"` // 变量名
	Description string `json:"description"` // 变量描述
	Required bool `json:"required"` // 是否必填
}

type AiGenerateEmailTemplateRes struct {
	Subject string `json:"subject"` // 生成的主题
	Html string `json:"html"` // 生成的HTML内容
	Text string `json:"text"` // 纯文本内容
	Variables []string `json:"variables"` // 识别的变量列表
	PreviewUrl string `json:"previewUrl"` // 预览URL
}
