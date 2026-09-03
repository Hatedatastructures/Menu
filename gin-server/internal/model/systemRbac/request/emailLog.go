package request

type GetEmailLogListReq struct {
	Page int `json:"page" form:"page"`
	PageSize int `json:"pageSize" form:"pageSize"`
	ToEmail string `json:"toEmail" form:"toEmail"`
	Status string `json:"status" form:"status"`
	BizType string `json:"bizType" form:"bizType"`
	StartTime string `json:"startTime" form:"startTime"`
	EndTime string `json:"endTime" form:"endTime"`
}

type GetEmailLogReq struct {
	Id int `json:"id" form:"id"`
}

type ResendEmailReq struct {
	Id int `json:"id" form:"id"`
}
