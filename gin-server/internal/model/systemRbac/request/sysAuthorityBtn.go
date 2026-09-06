package request

type GetAuthorityBtnReq struct {
	AuthorityId uint `json:"authorityId" form:"authorityId"`
	MenuID uint `json:"menuID" form:"menuID"`
}

type SetAuthorityBtnReq struct {
	AuthorityId uint `json:"authorityId" form:"authorityId"` // 角色ID
	MenuID uint `json:"menuID" form:"menuID"` // 菜单ID
	Selected []uint `json:"selected" form:"selected"` // 按钮ID数组
}

type CanRemoveAuthorityBtnReq struct {
	Id string `json:"id" form:"id"`
}
