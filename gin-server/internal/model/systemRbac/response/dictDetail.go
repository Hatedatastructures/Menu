package response


import "time"

type GetDictionaryDetailListRes struct {
	Page int `json:"page"`
	Size int `json:"size"`
	Label string `json:"label"`
	Value string `json:"value"`
	Status bool `json:"status"`
	SysDictionaryID int `json:"sysDictionaryID"`
	List []GetDictionaryDetailListResList `json:"list"` // []
	Total int64 `json:"total"` // 总数
}

type GetDictionaryDetailListResList struct {
	Id int `json:"id"` // ID
	Label string `json:"label"` // 展示值
	Value int `json:"value"` // 字典值
	Extend string `json:"extend"` // 扩展值
	Status bool `json:"status"` // 启用状态
	Sort int `json:"sort"` // 排序标记
	SysDictionaryID int `json:"sysDictionaryID"` // 关联字典ID
	CreatedAt time.Time `json:"createdAt"` // 创建时间
}

type CreateDictionaryDetailRes struct {
	Id int `json:"id"` // ID
}

type GetDictionaryDetailRes struct {
	Id int `json:"id"` // ID
	Label string `json:"label"` // 展示值
	Value int `json:"value"` // 字典值
	Extend string `json:"extend"` // 扩展值
	Status bool `json:"status"` // 启用状态
	Sort int `json:"sort"` // 排序标记
	SysDictionaryID int `json:"sysDictionaryID"` // 关联字典ID
	CreatedAt time.Time `json:"createdAt"` // 创建时间
}

type GetDictionaryListByIdRes struct {
	List []GetDictionaryListByIdResList `json:"list"` // []
}

type GetDictionaryListByIdResList struct {
	Id int `json:"id"` // ID
	Label string `json:"label"` // 展示值
	Value int `json:"value"` // 字典值
	Extend string `json:"extend"` // 扩展值
	Status bool `json:"status"` // 启用状态
	Sort int `json:"sort"` // 排序标记
}

type GetDictionaryListByTypeRes struct {
	List []GetDictionaryListByTypeResList `json:"list"` // []
}

type GetDictionaryListByTypeResList struct {
	Id int `json:"id"` // ID
	Label string `json:"label"` // 展示值
	Value int `json:"value"` // 字典值
	Extend string `json:"extend"` // 扩展值
	Status bool `json:"status"` // 启用状态
	Sort int `json:"sort"` // 排序标记
}

type GetDictionaryInfoByValueRes struct {
	Id int `json:"id"` // ID
	Label string `json:"label"` // 展示值
	Value int `json:"value"` // 字典值
	Extend string `json:"extend"` // 扩展值
	Status bool `json:"status"` // 启用状态
	Sort int `json:"sort"` // 排序标记
}

type GetDictionaryInfoByTypeValueRes struct {
	Id int `json:"id"` // ID
	Label string `json:"label"` // 展示值
	Value int `json:"value"` // 字典值
	Extend string `json:"extend"` // 扩展值
	Status bool `json:"status"` // 启用状态
	Sort int `json:"sort"` // 排序标记
}
