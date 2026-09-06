package response

type HealthzRes struct {
	Status string `json:"status"` // 服务状态:ok
}

type ReadyzRes struct {
	Status string `json:"status"` // 服务状态:ready / starting
}
