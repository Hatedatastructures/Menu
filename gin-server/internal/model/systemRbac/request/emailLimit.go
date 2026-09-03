package request

type GetEmailLimitLogReq struct {
	Page int `json:"page" form:"page"`
	PageSize int `json:"pageSize" form:"pageSize"`
	Email string `json:"email" form:"email"`
	Ip string `json:"ip" form:"ip"`
	StartTime string `json:"startTime" form:"startTime"`
	EndTime string `json:"endTime" form:"endTime"`
}

type AddEmailBlacklistReq struct {
	Email string `json:"email" form:"email"` // 被拉黑的邮箱
	Reason string `json:"reason" form:"reason"` // 拉黑原因
	Type string `json:"type" form:"type"` // 拉黑类型: manual/auto
}

type DeleteEmailBlacklistReq struct {
	Id int `json:"id" form:"id"`
}

type GetEmailBlacklistListReq struct {
	Page int `json:"page" form:"page"`
	PageSize int `json:"pageSize" form:"pageSize"`
}
