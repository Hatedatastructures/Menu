package response

type SendEmailRes struct {
	TaskId string `json:"taskId"` // 邮件任务ID
}

type SendBatchEmailRes struct {
	BatchId string `json:"batchId"` // 批次ID
	TotalCount int `json:"totalCount"` // 总收件人数
}
