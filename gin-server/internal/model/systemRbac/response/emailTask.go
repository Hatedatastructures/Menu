package response

type GetEmailTaskListRes struct {
	Page int `json:"page"`
	PageSize int `json:"pageSize"`
	Status string `json:"status"`
	Total int64 `json:"total"` // 总数
	List []GetEmailTaskListResList `json:"list"` // []
}

type GetEmailTaskListResList struct {
	Id int64 `json:"id"` // 任务ID
	Name string `json:"name"` // 任务名称
	CronExpr string `json:"cronExpr"` // Cron表达式
	Status string `json:"status"` // 任务状态
	NextRunTime string `json:"nextRunTime"` // 下次执行时间
	LastRunTime string `json:"lastRunTime"` // 上次执行时间
	TotalRuns int `json:"totalRuns"` // 总执行次数
	SuccessRate string `json:"successRate"` // 成功率
}

type RunEmailTaskRes struct {
	ExecutionId string `json:"executionId"` // 执行ID
}
