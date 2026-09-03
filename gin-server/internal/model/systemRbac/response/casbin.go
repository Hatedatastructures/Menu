package response

type GetPolicyPathByAuthorityIdRes struct {
	List []GetPolicyPathByAuthorityIdResList `json:"list"` // []
}

type GetPolicyPathByAuthorityIdResList struct {
	Path string `json:"path"` // API路径
	Method string `json:"method"` // 请求方法
}
