package response

type GenerateRes struct {
	HistoryId uint `json:"historyId"` // 历史记录ID
	Status string `json:"status"` // 状态
}

type GetGenerateResultRes struct {
	Id uint `json:"id"` // 历史记录ID
	Status string `json:"status"` // 状态(pending/generating/done/failed)
	ErrMsg string `json:"errMsg"` // 错误信息
	QuestionCount int `json:"questionCount"` // 题目数量
	Result map[string]interface{} `json:"result"` // 试卷JSON(结构见json/格式.md)
}
