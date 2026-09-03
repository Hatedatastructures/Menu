package request

type CreatePermissionReq struct {
	Name string `json:"name" form:"name"` // 权限名称
	Code string `json:"code" form:"code"` // 权限编码
	Type string `json:"type" form:"type"` // 权限类型(menu/api/button)
	Status bool `json:"status" form:"status"` // 状态
}

type GetPermissionListReq struct {
	Keyword string `json:"keyword" form:"keyword"`
	Page int `json:"page" form:"page"`
	Size int `json:"size" form:"size"`
	Type string `json:"type" form:"type"`
}

type GetPermissionDetailReq struct {
	Id string `json:"id" form:"id"`
}

type UpdatePermissionReq struct {
	Id string `json:"id" form:"id"`
	Name string `json:"name" form:"name"` // 权限名称(可选)
	Code string `json:"code" form:"code"` // 权限编码(可选)
	Type string `json:"type" form:"type"` // 权限类型(可选)
	Status bool `json:"status" form:"status"` // 状态(可选)
}

type DeletePermissionReq struct {
	Id string `json:"id" form:"id"`
}
