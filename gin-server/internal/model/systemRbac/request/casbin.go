package request

type UpdateCasbinReq struct {
	AuthorityId uint `json:"authorityId" form:"authorityId"` // 权限ID
	CasbinInfos []UpdateCasbinReqCasbininfo `json:"casbinInfos" form:"casbinInfos"` // []
}

type UpdateCasbinReqCasbininfo struct {
	Path string `json:"path" form:"path"` // API路径
	Method string `json:"method" form:"method"` // 请求方法
}

type GetPolicyPathByAuthorityIdReq struct {
	AuthorityId uint `json:"authorityId" form:"authorityId"`
}

type ClearCasbinReq struct {
	AuthorityId uint `json:"authorityId" form:"authorityId"` // 权限ID
}
