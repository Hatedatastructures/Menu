package request

type DeleteBaseMenuReq struct {
	Id int `json:"id" form:"id"`
}

type UpdateBaseMenuReq struct {
	Id int64 `json:"id" form:"id"` // 菜单ID
	ParentId uint `json:"parentId" form:"parentId"` // 父菜单ID
	Path string `json:"path" form:"path"` // 路由path
	Name string `json:"name" form:"name"` // 路由name
	Hidden bool `json:"hidden" form:"hidden"` // 是否隐藏
	Component string `json:"component" form:"component"` // 组件路径
	Sort int `json:"sort" form:"sort"` // 排序
	Meta UpdateBaseMenuReqMeta `json:"meta" form:"meta"`
	Parameters []UpdateBaseMenuReqParameter `json:"parameters" form:"parameters"` // []
	MenuBtn []UpdateBaseMenuReqMenubtn `json:"menuBtn" form:"menuBtn"` // []
}

type UpdateBaseMenuReqMeta struct {
	ActiveName string `json:"activeName" form:"activeName"` // 高亮菜单
	KeepAlive bool `json:"keepAlive" form:"keepAlive"` // 是否缓存
	DefaultMenu bool `json:"defaultMenu" form:"defaultMenu"` // 是否基础路由
	Title string `json:"title" form:"title"` // 菜单名
	Icon string `json:"icon" form:"icon"` // 菜单图标
	CloseTab bool `json:"closeTab" form:"closeTab"` // 自动关闭tab
	TransitionType string `json:"transitionType" form:"transitionType"` // 路由动画
}

type UpdateBaseMenuReqParameter struct {
	Type string `json:"type" form:"type"` // 参数类型(params/query)
	Key string `json:"key" form:"key"` // 参数key
	Value string `json:"value" form:"value"` // 参数值
}

type UpdateBaseMenuReqMenubtn struct {
	Name string `json:"name" form:"name"` // 按钮key
	Desc string `json:"desc" form:"desc"` // 按钮描述
}

type GetBaseMenuByIdReq struct {
	Id int `json:"id" form:"id"`
}
