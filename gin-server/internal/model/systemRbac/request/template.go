package request

type CreateTemplateReq struct {
	Code string `json:"code" form:"code"` // 模板编码
	Name string `json:"name" form:"name"` // 模板名称
	Tpl string `json:"tpl" form:"tpl"` // 模板内容(前端传入,后端自动解析变量填充Variables)
	Status string `json:"status" form:"status"` // 状态(active/inactive/draft)
	Description string `json:"description" form:"description"` // 模板描述
	Engine string `json:"engine" form:"engine"` // 模板引擎(默认html/template)
}

type DeleteTemplateReq struct {
	Id uint `json:"id" form:"id"`
}

type UpdateTemplateReq struct {
	Id uint `json:"id" form:"id"`
	Code string `json:"code" form:"code"` // 模板编码
	Name string `json:"name" form:"name"` // 模板名称
	Tpl string `json:"tpl" form:"tpl"` // 模板内容(前端传入,后端自动解析变量填充Variables)
	Status string `json:"status" form:"status"` // 状态(active/inactive/draft)
	Description string `json:"description" form:"description"` // 模板描述
	Engine string `json:"engine" form:"engine"` // 模板引擎
}

type GetTemplateDetailReq struct {
	Id uint `json:"id" form:"id"`
}

type GetTemplateListReq struct {
	Page int `json:"page" form:"page"`
	Size int `json:"size" form:"size"`
	Keyword string `json:"keyword" form:"keyword"`
	Status string `json:"status" form:"status"`
}
