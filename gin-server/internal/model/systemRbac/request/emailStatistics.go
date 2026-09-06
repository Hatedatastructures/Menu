package request

type GetEmailStatisticsReq struct {
	StartTime string `json:"startTime" form:"startTime"`
	EndTime string `json:"endTime" form:"endTime"`
}
