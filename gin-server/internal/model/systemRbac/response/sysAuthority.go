package response

type CreateAuthorityRes struct {
	AuthorityId uint `json:"authorityId"` // 角色ID
	AuthorityName string `json:"authorityName"` // 角色名
	ParentId uint `json:"parentId"` // 父角色ID
	DefaultRouter string `json:"defaultRouter"` // 默认菜单
	Children []CreateAuthorityResChildren `json:"children"` // []
}

type CreateAuthorityResChildren struct {
	AuthorityId uint `json:"authorityId"` // 角色ID
	AuthorityName string `json:"authorityName"` // 角色名
}

type CopyAuthorityRes struct {
	AuthorityId uint `json:"authorityId"` // 角色ID
	AuthorityName string `json:"authorityName"` // 角色名
	ParentId uint `json:"parentId"` // 父角色ID
	DefaultRouter string `json:"defaultRouter"` // 默认菜单
}

type UpdateAuthorityRes struct {
	AuthorityId uint `json:"authorityId"` // 角色ID
	AuthorityName string `json:"authorityName"` // 角色名
	ParentId uint `json:"parentId"` // 父角色ID
	DefaultRouter string `json:"defaultRouter"` // 默认菜单
}

type GetAuthorityInfoListRes struct {
	AuthorityId uint `json:"authorityId"`
	List []GetAuthorityInfoListResList `json:"list"` // []
	Total int64 `json:"total"` // 总数
}

type GetAuthorityInfoListResList struct {
	AuthorityId uint `json:"authorityId"` // 角色ID
	AuthorityName string `json:"authorityName"` // 角色名
	ParentId uint `json:"parentId"` // 父角色ID
	DefaultRouter string `json:"defaultRouter"` // 默认菜单
	DataAuthorityId []GetAuthorityInfoListResListDataauthorityid `json:"dataAuthorityId"` // []
	Children []GetAuthorityInfoListResListChildren `json:"children"` // []
}

type GetAuthorityInfoListResListDataauthorityid struct {
	AuthorityId uint `json:"authorityId"` // 数据权限角色ID
	AuthorityName string `json:"authorityName"` // 角色名
}

type GetAuthorityInfoListResListChildren struct {
	AuthorityId uint `json:"authorityId"` // 子角色ID
	AuthorityName string `json:"authorityName"` // 子角色名
	ParentId uint `json:"parentId"` // 父角色ID
	DefaultRouter string `json:"defaultRouter"` // 默认菜单
}

type GetStructAuthorityListRes struct {
	List []uint `json:"list"` // 角色ID数组
}

type GetAuthorityInfoRes struct {
	AuthorityId uint `json:"authorityId"` // 角色ID
	AuthorityName string `json:"authorityName"` // 角色名
	ParentId uint `json:"parentId"` // 父角色ID
	DefaultRouter string `json:"defaultRouter"` // 默认菜单
	DataAuthorityId []GetAuthorityInfoResDataauthorityid `json:"dataAuthorityId"` // []
	SysBaseMenus []GetAuthorityInfoResSysbasemenu `json:"sysBaseMenus"` // []
}

type GetAuthorityInfoResDataauthorityid struct {
	AuthorityId uint `json:"authorityId"` // 数据权限角色ID
	AuthorityName string `json:"authorityName"` // 角色名
}

type GetAuthorityInfoResSysbasemenu struct {
	Id int64 `json:"id"` // 菜单ID
	ParentId uint `json:"parentId"` // 父菜单ID
	Path string `json:"path"` // 路由path
	Name string `json:"name"` // 路由name
	Hidden bool `json:"hidden"` // 是否隐藏
	Component string `json:"component"` // 组件路径
	Sort int `json:"sort"` // 排序
	Meta GetAuthorityInfoResSysbasemenuMeta `json:"meta"`
}

type GetAuthorityInfoResSysbasemenuMeta struct {
	ActiveName string `json:"activeName"` // 高亮菜单
	KeepAlive bool `json:"keepAlive"` // 是否缓存
	DefaultMenu bool `json:"defaultMenu"` // 是否基础路由
	Title string `json:"title"` // 菜单名
	Icon string `json:"icon"` // 菜单图标
	CloseTab bool `json:"closeTab"` // 自动关闭tab
	TransitionType string `json:"transitionType"` // 路由动画
}

type GetParentAuthorityIDRes struct {
	ParentId uint `json:"parentId"` // 父角色ID
}
