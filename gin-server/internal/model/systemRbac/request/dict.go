package request

type CreateDictionaryReq struct {
	Name string `json:"name" form:"name"` // 字典名(中)
	Type string `json:"type" form:"type"` // 字典名(英)
	Status bool `json:"status" form:"status"` // 状态
	Desc string `json:"desc" form:"desc"` // 描述
}

type GetDictionaryReq struct {
	Id int `json:"id" form:"id"`
}

type UpdateDictionaryReq struct {
	Id int `json:"id" form:"id"`
	Name string `json:"name" form:"name"` // 字典名(中)(可选)
	Type string `json:"type" form:"type"` // 字典名(英)(可选)
	Status bool `json:"status" form:"status"` // 状态(可选)
	Desc string `json:"desc" form:"desc"` // 描述(可选)
}

type DeleteDictionaryReq struct {
	Id int `json:"id" form:"id"`
}

type GetDictsByTypesReq struct {
	Types string `json:"types" form:"types"`
}
