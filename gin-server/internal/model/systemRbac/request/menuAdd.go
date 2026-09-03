package request

type CreateMenuBaReq struct {
	ParentId uint   `json:"parentId" form:"parentId"` // 父级菜单ID(0为根)
	Name     string `json:"name" form:"name"`         // 菜单名称
	Path     string `json:"path" form:"path"`          // 路由路径
	Component string `json:"component" form:"component"` // 前端组件
	Icon     string `json:"icon" form:"icon"`         // 图标
	Sort     int    `json:"sort" form:"sort"`           // 排序
	Status   bool   `json:"status" form:"status"`       // 状态
}
