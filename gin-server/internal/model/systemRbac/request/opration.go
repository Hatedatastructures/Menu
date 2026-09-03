package request

type GetOperationRecordsByUserIdReq struct {
	UserId int `json:"userId" form:"userId"`
	Page int `json:"page" form:"page"`
	Size int `json:"size" form:"size"`
}

type GetOperationRecordsByTimeRangeReq struct {
	StartTime string `json:"startTime" form:"startTime"`
	EndTime string `json:"endTime" form:"endTime"`
	Page int `json:"page" form:"page"`
	Size int `json:"size" form:"size"`
}

type GetRecentOperationRecordsReq struct {
	Limit int `json:"limit" form:"limit"`
}

type GetErrorRecordsReq struct {
	Page int `json:"page" form:"page"`
	Size int `json:"size" form:"size"`
}

type DeleteExpiredRecordsReq struct {
	Days int `json:"days" form:"days"` // 保留天数(删除该天数之前的记录)
}

type GetOperationRecordsByMethodReq struct {
	Method string `json:"method" form:"method"`
	Page int `json:"page" form:"page"`
	Size int `json:"size" form:"size"`
}

type GetOperationRecordsByPathReq struct {
	Path string `json:"path" form:"path"`
	Page int `json:"page" form:"page"`
	Size int `json:"size" form:"size"`
}
