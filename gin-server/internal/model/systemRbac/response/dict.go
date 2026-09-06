package response


import "time"

type GetDictionaryListRes struct {
	List []GetDictionaryListResList `json:"list"` // []
}

type GetDictionaryListResList struct {
	Id int `json:"id"` // ID
	Name string `json:"name"` // 字典名(中)
	Type string `json:"type"` // 字典名(英)
	Status bool `json:"status"` // 状态
	Desc string `json:"desc"` // 描述
	CreatedAt time.Time `json:"createdAt"` // 创建时间
}

type CreateDictionaryRes struct {
	Id int `json:"id"` // ID
}

type GetDictionaryRes struct {
	Id int `json:"id"` // ID
	Name string `json:"name"` // 字典名(中)
	Type string `json:"type"` // 字典名(英)
	Status bool `json:"status"` // 状态
	Desc string `json:"desc"` // 描述
	SysDictionaryDetails []GetDictionaryResSysdictionarydetail `json:"sysDictionaryDetails"` // []
	CreatedAt time.Time `json:"createdAt"` // 创建时间
}

type GetDictionaryResSysdictionarydetail struct {
	Id int `json:"id"` // ID
	Label string `json:"label"` // 展示值
	Value int `json:"value"` // 字典值
	Extend string `json:"extend"` // 扩展值
	Status bool `json:"status"` // 启用状态
	Sort int `json:"sort"` // 排序标记
}

type GetDictsByTypesRes struct {
	Types string `json:"types"`
	Status []GetDictsByTypesResStatu `json:"status"` // []
	Gender []GetDictsByTypesResGender `json:"gender"` // []
}

type GetDictsByTypesResStatu struct {
	Label string `json:"label"` // 展示值
	Value int `json:"value"` // 字典值
}

type GetDictsByTypesResGender struct {
	Label string `json:"label"` // 展示值
	Value int `json:"value"` // 字典值
}

type GetAllDictsRes struct {
	Status []GetAllDictsResStatu `json:"status"` // []
	Gender []GetAllDictsResGender `json:"gender"` // []
}

type GetAllDictsResStatu struct {
	Label string `json:"label"` // 展示值
	Value int `json:"value"` // 字典值
}

type GetAllDictsResGender struct {
	Label string `json:"label"` // 展示值
	Value int `json:"value"` // 字典值
}

type RefreshDictsRes struct {
	Success bool `json:"success"` // 是否成功
	Message string `json:"message"` // 消息
}
