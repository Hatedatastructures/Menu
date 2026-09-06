package request

type CreateConfigGroupReq struct {
	Code string `json:"code" form:"code"` // 分组编码
	Name string `json:"name" form:"name"` // 分组名称
	Config string `json:"config" form:"config"` // 配置JSON
	Sort int `json:"sort" form:"sort"` // 排序
	Status int `json:"status" form:"status"` // 状态
	Remark string `json:"remark" form:"remark"` // 备注
}

type GetConfigGroupListReq struct {
	Page int `json:"page" form:"page"`
	Size int `json:"size" form:"size"`
	Keyword string `json:"keyword" form:"keyword"`
}

type GetConfigGroupDetailReq struct {
	Id string `json:"id" form:"id"`
}

type UpdateConfigGroupReq struct {
	Id string `json:"id" form:"id"`
	Code string `json:"code" form:"code"` // 分组编码(可选)
	Name string `json:"name" form:"name"` // 分组名称(可选)
	Config string `json:"config" form:"config"` // 配置JSON(可选)
	Sort int `json:"sort" form:"sort"` // 排序(可选)
	Status int `json:"status" form:"status"` // 状态(可选)
	Remark string `json:"remark" form:"remark"` // 备注(可选)
}

type DeleteConfigGroupReq struct {
	Id string `json:"id" form:"id"`
}

type BatchDeleteConfigGroupsReq struct {
	Ids string `json:"ids" form:"ids"` // 分组ID数组
}

type GetConfigByCodeReq struct {
	Code string `json:"code" form:"code"`
}

type SaveConfigReq struct {
	Code string `json:"code" form:"code"`
	Config string `json:"config" form:"config"` // 配置JSON
}
