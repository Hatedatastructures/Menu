package request

type GetLoginLogListReq struct {
	Page int `json:"page" form:"page"`
	Size int `json:"size" form:"size"`
	Username string `json:"username" form:"username"`
	Status string `json:"status" form:"status"`
	StartTime string `json:"startTime" form:"startTime"`
	EndTime string `json:"endTime" form:"endTime"`
}

type GetLoginLogByUsernameReq struct {
	Username string `json:"username" form:"username"`
	Page int `json:"page" form:"page"`
	Size int `json:"size" form:"size"`
}

type GetLoginLogByStatusReq struct {
	Status string `json:"status" form:"status"`
	Page int `json:"page" form:"page"`
	Size int `json:"size" form:"size"`
}

type GetLoginLogByTimeRangeReq struct {
	StartTime string `json:"startTime" form:"startTime"`
	EndTime string `json:"endTime" form:"endTime"`
	Page int `json:"page" form:"page"`
	Size int `json:"size" form:"size"`
}

type GetRecentLoginLogReq struct {
	Limit int `json:"limit" form:"limit"`
}

type DeleteLoginLogReq struct {
	Ids string `json:"ids" form:"ids"` // 登录日志ID数组
}
