package response

type CreateTemplateRes struct {
	ID uint `json:"ID"` // 模板ID
}

type GetTemplateDetailRes struct {
	ID uint `json:"ID"` // 模板ID
	Code string `json:"code"` // 模板编码
	Name string `json:"name"` // 模板名称
	Content string `json:"content"` // 模板内容
	Status string `json:"status"` // 状态(active/inactive/draft)
	Description string `json:"description"` // 模板描述
	Variables []string `json:"variables"` // 模板变量列表(后端自动解析)
	Engine string `json:"engine"` // 模板引擎
	CreatedAt string `json:"createdAt"` // 创建时间
	UpdatedAt string `json:"updatedAt"` // 更新时间
}

type GetTemplateListRes struct {

	Keyword string `json:"keyword"`
	Status string `json:"status"`
	List []GetTemplateListResList `json:"list"` // []
	Total int64 `json:"total"` // 总数
	Page int `json:"page"` // 当前页
	Size int `json:"size"` // 每页条数
}

type GetTemplateListResList struct {
	ID uint `json:"ID"` // 模板ID
	Code string `json:"code"` // 模板编码
	Name string `json:"name"` // 模板名称
	Content string `json:"content"` // 模板内容
	Status string `json:"status"` // 状态(active/inactive/draft)
	Description string `json:"description"` // 模板描述
	Engine string `json:"engine"` // 模板引擎
	Variables []string `json:"variables"` // 模板变量列表
	CreatedAt string `json:"createdAt"` // 创建时间
	UpdatedAt string `json:"updatedAt"` // 更新时间
}
