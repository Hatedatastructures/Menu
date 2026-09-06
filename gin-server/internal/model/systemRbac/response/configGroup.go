package response

type CreateConfigGroupRes struct {
	Id string `json:"id"` // 分组ID
}

type GetConfigGroupListRes struct {
	Page int `json:"page"`
	Size int `json:"size"`
	Keyword string `json:"keyword"`
	List []GetConfigGroupListResList `json:"list"` // []
	Total int64 `json:"total"` // 总数
}

type GetConfigGroupListResList struct {
	Id string `json:"id"` // ID
	Code string `json:"code"` // 分组编码
	Name string `json:"name"` // 分组名称
	Sort int `json:"sort"` // 排序
	Status int `json:"status"` // 状态
	Remark string `json:"remark"` // 备注
	CreatedAt string `json:"createdAt"` // 创建时间
	UpdatedAt string `json:"updatedAt"` // 更新时间
}

type GetConfigGroupDetailRes struct {
	Id string `json:"id"` // ID
	Code string `json:"code"` // 分组编码
	Name string `json:"name"` // 分组名称
	Config string `json:"config"` // 配置JSON
	Sort int `json:"sort"` // 排序
	Status int `json:"status"` // 状态
	Remark string `json:"remark"` // 备注
	Version int `json:"version"` // 版本号
	CreatedAt string `json:"createdAt"` // 创建时间
	UpdatedAt string `json:"updatedAt"` // 更新时间
}

type GetConfigByCodeRes struct {
	Id string `json:"id"` // ID
	Code string `json:"code"` // 分组编码
	Name string `json:"name"` // 分组名称
	Config string `json:"config"` // 配置JSON
	Version int `json:"version"` // 版本号
}

type SaveConfigRes struct {
	Version int `json:"version"` // 新版本号
}

type RefreshConfigCacheRes struct {
	Success bool `json:"success"` // 是否成功
}
