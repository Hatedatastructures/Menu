package request

type CreateAuthorityReq struct {
	AuthorityId uint `json:"authorityId" form:"authorityId"` // 角色ID
	AuthorityName string `json:"authorityName" form:"authorityName"` // 角色名
	ParentId uint `json:"parentId" form:"parentId"` // 父角色ID
	DefaultRouter string `json:"defaultRouter" form:"defaultRouter"` // 默认菜单
}

type CopyAuthorityReq struct {
	Authority CopyAuthorityReqAuthority `json:"authority" form:"authority"`
	OldAuthorityId uint `json:"oldAuthorityId" form:"oldAuthorityId"` // 旧角色ID
}

type CopyAuthorityReqAuthority struct {
	AuthorityId uint `json:"authorityId" form:"authorityId"` // 新角色ID
	AuthorityName string `json:"authorityName" form:"authorityName"` // 新角色名
	ParentId uint `json:"parentId" form:"parentId"` // 父角色ID
	DefaultRouter string `json:"defaultRouter" form:"defaultRouter"` // 默认菜单
}

type UpdateAuthorityReq struct {
	AuthorityId uint `json:"authorityId" form:"authorityId"` // 角色ID
	AuthorityName string `json:"authorityName" form:"authorityName"` // 角色名
	ParentId uint `json:"parentId" form:"parentId"` // 父角色ID
	DefaultRouter string `json:"defaultRouter" form:"defaultRouter"` // 默认菜单
}

type DeleteAuthorityReq struct {
	AuthorityId uint `json:"authorityId" form:"authorityId"`
}

type GetAuthorityInfoListReq struct {
	AuthorityId uint `json:"authorityId" form:"authorityId"`
}

type GetStructAuthorityListReq struct {
	AuthorityId uint `json:"authorityId" form:"authorityId"`
}

type GetAuthorityInfoReq struct {
	AuthorityId uint `json:"authorityId" form:"authorityId"`
}

type SetDataAuthorityReq struct {
	AuthorityId uint `json:"authorityId" form:"authorityId"` // 角色ID
	DataAuthorityId []uint `json:"dataAuthorityId" form:"dataAuthorityId"` // 数据权限角色ID数组
}

type SetMenuAuthorityReq struct {
	AuthorityId uint `json:"authorityId" form:"authorityId"` // 角色ID
	SysBaseMenus []SetMenuAuthorityReqSysbasemenu `json:"sysBaseMenus" form:"sysBaseMenus"` // []
}

type SetMenuAuthorityReqSysbasemenu struct {
	Id int64 `json:"id" form:"id"` // 菜单ID
	ParentId uint `json:"parentId" form:"parentId"` // 父菜单ID
	Path string `json:"path" form:"path"` // 路由path
	Name string `json:"name" form:"name"` // 路由name
	Hidden bool `json:"hidden" form:"hidden"` // 是否隐藏
	Component string `json:"component" form:"component"` // 组件路径
	Sort int `json:"sort" form:"sort"` // 排序
	Meta SetMenuAuthorityReqSysbasemenuMeta `json:"meta" form:"meta"`
}

type SetMenuAuthorityReqSysbasemenuMeta struct {
	ActiveName string `json:"activeName" form:"activeName"` // 高亮菜单
	KeepAlive bool `json:"keepAlive" form:"keepAlive"` // 是否缓存
	DefaultMenu bool `json:"defaultMenu" form:"defaultMenu"` // 是否基础路由
	Title string `json:"title" form:"title"` // 菜单名
	Icon string `json:"icon" form:"icon"` // 菜单图标
	CloseTab bool `json:"closeTab" form:"closeTab"` // 自动关闭tab
	TransitionType string `json:"transitionType" form:"transitionType"` // 路由动画
}

type GetParentAuthorityIDReq struct {
	AuthorityId uint `json:"authorityId" form:"authorityId"`
}
