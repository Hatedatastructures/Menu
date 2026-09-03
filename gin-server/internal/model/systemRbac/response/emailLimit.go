package response

type GetEmailLimitLogRes struct {
	Page int `json:"page"`
	PageSize int `json:"pageSize"`
	Email string `json:"email"`
	Ip string `json:"ip"`
	StartTime string `json:"startTime"`
	EndTime string `json:"endTime"`
	Total int64 `json:"total"` // 总数
	List []GetEmailLimitLogResList `json:"list"` // []
}

type GetEmailLimitLogResList struct {
	Id int64 `json:"id"` // 日志ID
	Email string `json:"email"` // 触犯限流的邮箱
	Ip string `json:"ip"` // 触犯限流的IP
	LimitType string `json:"limitType"` // 限流类型: email/ip/global
	Action string `json:"action"` // 触犯限流的操作
	Blocked bool `json:"blocked"` // 是否被阻止
	Reason string `json:"reason"` // 限流原因
	CreatedAt string `json:"createdAt"` // 创建时间
}

type GetEmailBlacklistListRes struct {
	Page int `json:"page"`
	PageSize int `json:"pageSize"`
	Total int64 `json:"total"` // 总数
	List []GetEmailBlacklistListResList `json:"list"` // []
}

type GetEmailBlacklistListResList struct {
	Id int64 `json:"id"` // 黑名单ID
	Email string `json:"email"` // 被拉黑的邮箱
	Reason string `json:"reason"` // 拉黑原因
	Type string `json:"type"` // 拉黑类型
	CreatedAt string `json:"createdAt"` // 创建时间
	CreatedBy string `json:"createdBy"` // 创建人
}
