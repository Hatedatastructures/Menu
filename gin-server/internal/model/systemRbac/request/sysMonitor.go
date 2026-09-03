package request

type GetMonitorTimelineReq struct {
	Type string `json:"type" form:"type"`
	Points int `json:"points" form:"points"`
}
