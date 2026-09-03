package response

type UpdateBaseMenuRes struct {
	Id int64 `json:"id"` // 菜单ID
}

type GetBaseMenuByIdRes struct {
	Id int64 `json:"id"` // 菜单ID
	ParentId uint `json:"parentId"` // 父菜单ID
	Path string `json:"path"` // 路由path
	Name string `json:"name"` // 路由name
	Hidden bool `json:"hidden"` // 是否隐藏
	Component string `json:"component"` // 组件路径
	Sort int `json:"sort"` // 排序
	Meta GetBaseMenuByIdResMeta `json:"meta"`
	Parameters []GetBaseMenuByIdResParameter `json:"parameters"` // []
	MenuBtn []GetBaseMenuByIdResMenubtn `json:"menuBtn"` // []
}

type GetBaseMenuByIdResMeta struct {
	ActiveName string `json:"activeName"` // 高亮菜单
	KeepAlive bool `json:"keepAlive"` // 是否缓存
	DefaultMenu bool `json:"defaultMenu"` // 是否基础路由
	Title string `json:"title"` // 菜单名
	Icon string `json:"icon"` // 菜单图标
	CloseTab bool `json:"closeTab"` // 自动关闭tab
	TransitionType string `json:"transitionType"` // 路由动画
}

type GetBaseMenuByIdResParameter struct {
	Id int64 `json:"id"` // 参数ID
	Type string `json:"type"` // 参数类型(params/query)
	Key string `json:"key"` // 参数key
	Value string `json:"value"` // 参数值
}

type GetBaseMenuByIdResMenubtn struct {
	Id int64 `json:"id"` // 按钮ID
	Name string `json:"name"` // 按钮key
	Desc string `json:"desc"` // 按钮描述
}
