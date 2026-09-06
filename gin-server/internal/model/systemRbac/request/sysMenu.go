package request

type GetMenuTreeReq struct {
	AuthorityId uint `json:"authorityId" form:"authorityId"`
}

type GetInfoListReq struct {
	AuthorityId uint `json:"authorityId" form:"authorityId"`
}

type AddBaseMenuReq struct {
	ParentId uint `json:"parentId" form:"parentId"` // 父菜单ID
	Path string `json:"path" form:"path"` // 路由path
	Name string `json:"name" form:"name"` // 路由name
	Hidden bool `json:"hidden" form:"hidden"` // 是否隐藏
	Component string `json:"component" form:"component"` // 组件路径
	Sort int `json:"sort" form:"sort"` // 排序
	Meta AddBaseMenuReqMeta `json:"meta" form:"meta"`
}

type AddBaseMenuReqMeta struct {
	ActiveName string `json:"activeName" form:"activeName"` // 高亮菜单
	KeepAlive bool `json:"keepAlive" form:"keepAlive"` // 是否缓存
	DefaultMenu bool `json:"defaultMenu" form:"defaultMenu"` // 是否基础路由
	Title string `json:"title" form:"title"` // 菜单名
	Icon string `json:"icon" form:"icon"` // 菜单图标
	CloseTab bool `json:"closeTab" form:"closeTab"` // 自动关闭tab
	TransitionType string `json:"transitionType" form:"transitionType"` // 路由动画
}

type GetBaseMenuTreeReq struct {
	AuthorityId uint `json:"authorityId" form:"authorityId"`
}

type AddMenuAuthorityReq struct {
	Menus []AddMenuAuthorityReqMenu `json:"menus" form:"menus"` // []
	AdminAuthorityId uint `json:"adminAuthorityId" form:"adminAuthorityId"` // 管理员角色ID
	AuthorityId uint `json:"authorityId" form:"authorityId"` // 目标角色ID
}

type AddMenuAuthorityReqMenu struct {
	Id int64 `json:"id" form:"id"` // 菜单ID
	ParentId uint `json:"parentId" form:"parentId"` // 父菜单ID
	Path string `json:"path" form:"path"` // 路由path
	Name string `json:"name" form:"name"` // 路由name
	Hidden bool `json:"hidden" form:"hidden"` // 是否隐藏
	Component string `json:"component" form:"component"` // 组件路径
	Sort int `json:"sort" form:"sort"` // 排序
	Meta AddMenuAuthorityReqMenuMeta `json:"meta" form:"meta"`
}

type AddMenuAuthorityReqMenuMeta struct {
	ActiveName string `json:"activeName" form:"activeName"` // 高亮菜单
	KeepAlive bool `json:"keepAlive" form:"keepAlive"` // 是否缓存
	DefaultMenu bool `json:"defaultMenu" form:"defaultMenu"` // 是否基础路由
	Title string `json:"title" form:"title"` // 菜单名
	Icon string `json:"icon" form:"icon"` // 菜单图标
	CloseTab bool `json:"closeTab" form:"closeTab"` // 自动关闭tab
	TransitionType string `json:"transitionType" form:"transitionType"` // 路由动画
}

type GetMenuAuthorityReq struct {
	AuthorityId uint `json:"authorityId" form:"authorityId"`
}
