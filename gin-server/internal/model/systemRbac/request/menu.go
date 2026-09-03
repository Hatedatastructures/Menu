package request

type CreateMenuReq struct {
	ParentId uint `json:"parentId" form:"parentId"` // 父级菜单ID(0为根)
	Name string `json:"name" form:"name"` // 菜单名称
	Path string `json:"path" form:"path"` // 路由路径
	Component string `json:"component" form:"component"` // 前端组件
	Icon string `json:"icon" form:"icon"` // 图标
	Sort int `json:"sort" form:"sort"` // 排序
	Status bool `json:"status" form:"status"` // 状态
}

type GetMenuDetailReq struct {
	Id string `json:"id" form:"id"`
}

type UpdateMenuReq struct {
	Id string `json:"id" form:"id"`
	ParentId uint `json:"parentId" form:"parentId"` // 父级菜单ID(可选)
	Name string `json:"name" form:"name"` // 菜单名称(可选)
	Path string `json:"path" form:"path"` // 路由路径(可选)
	Component string `json:"component" form:"component"` // 前端组件(可选)
	Icon string `json:"icon" form:"icon"` // 图标(可选)
	Sort int `json:"sort" form:"sort"` // 排序(可选)
	Status bool `json:"status" form:"status"` // 状态(可选)
}

type DeleteMenuReq struct {
	Id string `json:"id" form:"id"`
}
