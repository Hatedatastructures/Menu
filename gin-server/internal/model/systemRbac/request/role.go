package request

type CreateRoleReq struct {
	Name string `json:"name" form:"name"` // 角色名称
	Code string `json:"code" form:"code"` // 角色编码
	Description string `json:"description" form:"description"` // 描述
	Status bool `json:"status" form:"status"` // 状态
}

type GetRoleListReq struct {
	Keyword string `json:"keyword" form:"keyword"`
	Page int `json:"page" form:"page"`
	Size int `json:"size" form:"size"`
}

type GetRoleDetailReq struct {
	Id string `json:"id" form:"id"`
}

type UpdateRoleReq struct {
	Id string `json:"id" form:"id"`
	Name string `json:"name" form:"name"` // 角色名称(可选)
	Code string `json:"code" form:"code"` // 角色编码(可选)
	Description string `json:"description" form:"description"` // 描述(可选)
	Status bool `json:"status" form:"status"` // 状态(可选)
}

type DeleteRoleReq struct {
	Id string `json:"id" form:"id"`
}

type AssignRolePermissionsReq struct {
	Id string `json:"id" form:"id"`
	PermissionIds string `json:"permissionIds" form:"permissionIds"` // 权限ID数组
}

type GetRolePermissionsReq struct {
	Id string `json:"id" form:"id"`
}
