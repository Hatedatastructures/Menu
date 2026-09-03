package response

type CreateApiRes struct {
	Id string `json:"id"` // API,ID
}

type GetApiListRes struct {
	Keyword string `json:"keyword"`
	Page int `json:"page"`
	Size int `json:"size"`
	Method string `json:"method"`
	List []GetApiListResList `json:"list"` // []
	Total int64 `json:"total"` // 总数
}

type GetApiListResList struct {
	Id string `json:"id"` // ID
	Name string `json:"name"` // 接口名称
	Path string `json:"path"` // 请求路径
	Method string `json:"method"` // 请求方法
	Status bool `json:"status"` // 状态
	CreatedAt string `json:"createdAt"` // 创建时间
}

type GetApiDetailRes struct {
	Id string `json:"id"` // ID
	Name string `json:"name"` // 接口名称
	Path string `json:"path"` // 请求路径
	Method string `json:"method"` // 请求方法
	Status bool `json:"status"` // 状态
	CreatedAt string `json:"createdAt"` // 创建时间
}
