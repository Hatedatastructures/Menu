package response

type GetHistoryListRes struct {
	Page int `json:"page"`
	PageSize int `json:"pageSize"`
	Status string `json:"status"`
	List []GetHistoryListResList `json:"list"` // []
	Total int64 `json:"total"` // 总数
}

type GetHistoryListResList struct {
	Id uint `json:"id"` // 历史记录ID
	FileName string `json:"fileName"` // 原始文件名
	FileId uint `json:"fileId"` // 文件ID
	QuestionCount int `json:"questionCount"` // 题目数量
	Status string `json:"status"` // 状态
	CreatedAt string `json:"createdAt"` // 生成时间
}

type GetHistoryDetailRes struct {
	Id uint `json:"id"` // 历史记录ID
	FileName string `json:"fileName"` // 原始文件名
	FileId uint `json:"fileId"` // 文件ID
	QuestionCount int `json:"questionCount"` // 题目数量
	Status string `json:"status"` // 状态
	Result map[string]interface{} `json:"result"` // 试卷JSON(结构见json/格式.md)
	CreatedAt string `json:"createdAt"` // 生成时间
}
