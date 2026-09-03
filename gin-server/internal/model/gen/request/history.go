package request

type GetHistoryListReq struct {
	Page int `json:"page" form:"page"`
	PageSize int `json:"pageSize" form:"pageSize"`
	Status string `json:"status" form:"status"`
}

type GetHistoryDetailReq struct {
	Id uint `json:"id" form:"id"`
}

type DeleteHistoryReq struct {
	Id uint `json:"id" form:"id"`
}
