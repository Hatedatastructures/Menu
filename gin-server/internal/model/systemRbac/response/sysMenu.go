package response

type GetMenuTreeRes struct {
	List []GetMenuTreeResList `json:"list"` // []
}

type GetMenuTreeResList struct {
	MenuId uint `json:"menuId"` // 菜单ID
	ParentId uint `json:"parentId"` // 父菜单ID
	Path string `json:"path"` // 路由path
	Name string `json:"name"` // 路由name
	Hidden bool `json:"hidden"` // 是否隐藏
	Component string `json:"component"` // 组件路径
	Sort int `json:"sort"` // 排序
	Meta GetMenuTreeResListMeta `json:"meta"`
	Children []GetMenuTreeResListChildren `json:"children"` // []
	Parameters []GetMenuTreeResListParameter `json:"parameters"` // []
	Btns map[string]uint `json:"btns"` // 按钮权限
}

type GetMenuTreeResListMeta struct {
	ActiveName string `json:"activeName"` // 高亮菜单
	KeepAlive bool `json:"keepAlive"` // 是否缓存
	DefaultMenu bool `json:"defaultMenu"` // 是否基础路由
	Title string `json:"title"` // 菜单名
	Icon string `json:"icon"` // 菜单图标
	CloseTab bool `json:"closeTab"` // 自动关闭tab
	TransitionType string `json:"transitionType"` // 路由动画
}

type GetMenuTreeResListChildren struct {
	MenuId uint `json:"menuId"` // 菜单ID
	ParentId uint `json:"parentId"` // 父菜单ID
	Path string `json:"path"` // 路由path
	Name string `json:"name"` // 路由name
	Hidden bool `json:"hidden"` // 是否隐藏
	Component string `json:"component"` // 组件路径
	Sort int `json:"sort"` // 排序
	Meta GetMenuTreeResListChildrenMeta `json:"meta"`
	Children []GetMenuTreeResListChildren `json:"children"` // 递归子菜单
}

type GetMenuTreeResListChildrenMeta struct {
	ActiveName string `json:"activeName"` // 高亮菜单
	KeepAlive bool `json:"keepAlive"` // 是否缓存
	DefaultMenu bool `json:"defaultMenu"` // 是否基础路由
	Title string `json:"title"` // 菜单名
	Icon string `json:"icon"` // 菜单图标
	CloseTab bool `json:"closeTab"` // 自动关闭tab
	TransitionType string `json:"transitionType"` // 路由动画
}

type GetMenuTreeResListParameter struct {
	Type string `json:"type"` // 参数类型(params/query)
	Key string `json:"key"` // 参数key
	Value string `json:"value"` // 参数值
}

type GetInfoListRes struct {
	List []GetInfoListResList `json:"list"` // []
	Total int64 `json:"total"` // 总数
}

type GetInfoListResList struct {
	Id int64 `json:"id"` // 菜单ID
	ParentId uint `json:"parentId"` // 父菜单ID
	Path string `json:"path"` // 路由path
	Name string `json:"name"` // 路由name
	Hidden bool `json:"hidden"` // 是否隐藏
	Component string `json:"component"` // 组件路径
	Sort int `json:"sort"` // 排序
	Meta GetInfoListResListMeta `json:"meta"`
	Children []GetInfoListResList `json:"children"` // 递归子菜单
	Parameters []GetInfoListResListParameter `json:"parameters"` // []
	MenuBtn []GetInfoListResListMenubtn `json:"menuBtn"` // []
}

type GetInfoListResListMeta struct {
	ActiveName string `json:"activeName"` // 高亮菜单
	KeepAlive bool `json:"keepAlive"` // 是否缓存
	DefaultMenu bool `json:"defaultMenu"` // 是否基础路由
	Title string `json:"title"` // 菜单名
	Icon string `json:"icon"` // 菜单图标
	CloseTab bool `json:"closeTab"` // 自动关闭tab
	TransitionType string `json:"transitionType"` // 路由动画
}

type GetInfoListResListParameter struct {
	Type string `json:"type"` // 参数类型(params/query)
	Key string `json:"key"` // 参数key
	Value string `json:"value"` // 参数值
}

type GetInfoListResListMenubtn struct {
	Id int64 `json:"id"` // 按钮ID
	Name string `json:"name"` // 按钮key
	Desc string `json:"desc"` // 按钮描述
}

type AddBaseMenuRes struct {
	Id int64 `json:"id"` // 菜单ID
	ParentId uint `json:"parentId"` // 父菜单ID
	Path string `json:"path"` // 路由path
	Name string `json:"name"` // 路由name
	Hidden bool `json:"hidden"` // 是否隐藏
	Component string `json:"component"` // 组件路径
	Sort int `json:"sort"` // 排序
}

type GetBaseMenuTreeRes struct {
	List []GetBaseMenuTreeResList `json:"list"` // []
}

type GetBaseMenuTreeResList struct {
	Id int64 `json:"id"` // 菜单ID
	ParentId uint `json:"parentId"` // 父菜单ID
	Path string `json:"path"` // 路由path
	Name string `json:"name"` // 路由name
	Hidden bool `json:"hidden"` // 是否隐藏
	Component string `json:"component"` // 组件路径
	Sort int `json:"sort"` // 排序
	Meta GetBaseMenuTreeResListMeta `json:"meta"`
	Children []GetBaseMenuTreeResList `json:"children"` // 递归子菜单
	Parameters []GetBaseMenuTreeResListParameter `json:"parameters"` // []
	MenuBtn []GetBaseMenuTreeResListMenubtn `json:"menuBtn"` // []
}

type GetBaseMenuTreeResListMeta struct {
	ActiveName string `json:"activeName"` // 高亮菜单
	KeepAlive bool `json:"keepAlive"` // 是否缓存
	DefaultMenu bool `json:"defaultMenu"` // 是否基础路由
	Title string `json:"title"` // 菜单名
	Icon string `json:"icon"` // 菜单图标
	CloseTab bool `json:"closeTab"` // 自动关闭tab
	TransitionType string `json:"transitionType"` // 路由动画
}

type GetBaseMenuTreeResListParameter struct {
	Type string `json:"type"` // 参数类型(params/query)
	Key string `json:"key"` // 参数key
	Value string `json:"value"` // 参数值
}

type GetBaseMenuTreeResListMenubtn struct {
	Id int64 `json:"id"` // 按钮ID
	Name string `json:"name"` // 按钮key
	Desc string `json:"desc"` // 按钮描述
}

type GetMenuAuthorityRes struct {
	List []GetMenuAuthorityResList `json:"list"` // []
}

type GetMenuAuthorityResList struct {
	MenuId uint `json:"menuId"` // 菜单ID
	ParentId uint `json:"parentId"` // 父菜单ID
	Path string `json:"path"` // 路由path
	Name string `json:"name"` // 路由name
	Hidden bool `json:"hidden"` // 是否隐藏
	Component string `json:"component"` // 组件路径
	Sort int `json:"sort"` // 排序
	Meta GetMenuAuthorityResListMeta `json:"meta"`
	Children []GetMenuAuthorityResList `json:"children"` // 递归子菜单
	Parameters []GetMenuAuthorityResListParameter `json:"parameters"` // []
}

type GetMenuAuthorityResListMeta struct {
	ActiveName string `json:"activeName"` // 高亮菜单
	KeepAlive bool `json:"keepAlive"` // 是否缓存
	DefaultMenu bool `json:"defaultMenu"` // 是否基础路由
	Title string `json:"title"` // 菜单名
	Icon string `json:"icon"` // 菜单图标
	CloseTab bool `json:"closeTab"` // 自动关闭tab
	TransitionType string `json:"transitionType"` // 路由动画
}

type GetMenuAuthorityResListParameter struct {
	Type string `json:"type"` // 参数类型(params/query)
	Key string `json:"key"` // 参数key
	Value string `json:"value"` // 参数值
}
