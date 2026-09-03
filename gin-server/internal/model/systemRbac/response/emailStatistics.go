package response

type GetEmailStatisticsRes struct {
	StartTime string `json:"startTime"`
	EndTime string `json:"endTime"`
	TotalSent int `json:"totalSent"` // 总发送量
	TotalSuccess int `json:"totalSuccess"` // 总成功量
	TotalFailed int `json:"totalFailed"` // 总失败量
	SuccessRate string `json:"successRate"` // 成功率
	ByType map[string]int `json:"byType"` // 按类型统计
	ByDate []GetEmailStatisticsResBydate `json:"byDate"` // []
}

type GetEmailStatisticsResBydate struct {
	Date string `json:"date"` // 日期
	Count int `json:"count"` // 数量
	Success int `json:"success"` // 成功数
	Failed int `json:"failed"` // 失败数
}

type GetEmailQueueStatusRes struct {
	Pending int `json:"pending"` // 待处理数量
	Processing int `json:"processing"` // 处理中数量
	Workers int `json:"workers"` // Worker数量
	AvgProcessTime string `json:"avgProcessTime"` // 平均处理时间
}
