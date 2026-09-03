package request

type GetDictionaryDetailListReq struct {
	Page int `json:"page" form:"page"`
	Size int `json:"size" form:"size"`
	Label string `json:"label" form:"label"`
	Value int `json:"value" form:"value"`
	Status bool `json:"status" form:"status"`
	SysDictionaryID int `json:"sysDictionaryID" form:"sysDictionaryID"`
}

type CreateDictionaryDetailReq struct {
	Label string `json:"label" form:"label"` // 展示值
	Value int `json:"value" form:"value"` // 字典值
	Extend string `json:"extend" form:"extend"` // 扩展值
	Status bool `json:"status" form:"status"` // 启用状态
	Sort int `json:"sort" form:"sort"` // 排序标记
	SysDictionaryID int `json:"sysDictionaryID" form:"sysDictionaryID"` // 关联字典ID
}

type GetDictionaryDetailReq struct {
	Id int `json:"id" form:"id"`
}

type UpdateDictionaryDetailReq struct {
	Id int `json:"id" form:"id"`
	Label string `json:"label" form:"label"` // 展示值(可选)
	Value int `json:"value" form:"value"` // 字典值(可选)
	Extend string `json:"extend" form:"extend"` // 扩展值(可选)
	Status bool `json:"status" form:"status"` // 启用状态(可选)
	Sort int `json:"sort" form:"sort"` // 排序标记(可选)
	SysDictionaryID int `json:"sysDictionaryID" form:"sysDictionaryID"` // 关联字典ID(可选)
}

type DeleteDictionaryDetailReq struct {
	Id int `json:"id" form:"id"`
}

type GetDictionaryListByIdReq struct {
	DictionaryId int `json:"dictionaryId" form:"dictionaryId"`
}

type GetDictionaryListByTypeReq struct {
	Type string `json:"type" form:"type"`
}

type GetDictionaryInfoByValueReq struct {
	DictionaryId int `json:"dictionaryId" form:"dictionaryId"`
	Value int `json:"value" form:"value"`
}

type GetDictionaryInfoByTypeValueReq struct {
	Type string `json:"type" form:"type"`
	Value int `json:"value" form:"value"`
}
