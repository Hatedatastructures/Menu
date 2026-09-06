package request

type CreateApiReq struct {
	Name string `json:"name" form:"name"` // 接口名称
	Path string `json:"path" form:"path"` // 请求路径
	Method string `json:"method" form:"method"` // 请求方法
	Status bool `json:"status" form:"status"` // 状态
}

type GetApiListReq struct {
	Keyword string `json:"keyword" form:"keyword"`
	Page int `json:"page" form:"page"`
	Size int `json:"size" form:"size"`
	Method string `json:"method" form:"method"`
}

type GetApiDetailReq struct {
	Id string `json:"id" form:"id"`
}

type UpdateApiReq struct {
	Id string `json:"id" form:"id"`
	Name string `json:"name" form:"name"` // 接口名称(可选)
	Path string `json:"path" form:"path"` // 请求路径(可选)
	Method string `json:"method" form:"method"` // 请求方法(可选)
	Status bool `json:"status" form:"status"` // 状态(可选)
}

type DeleteApiReq struct {
	Id string `json:"id" form:"id"`
}
