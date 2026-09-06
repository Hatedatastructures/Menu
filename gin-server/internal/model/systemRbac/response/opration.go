package response

type GetOperationStatisticsRes struct {
	Total int64 `json:"total"` // 总数
	TodayCount int64 `json:"todayCount"` // 今日数量
	ErrorCount int64 `json:"errorCount"` // 异常数量
	AvgLatency float64 `json:"avgLatency"` // 平均延迟(ms)
	MethodStats []GetOperationStatisticsResMethodstat `json:"methodStats"` // []
	PathStats []GetOperationStatisticsResPathstat `json:"pathStats"` // []
}

type GetOperationStatisticsResMethodstat struct {
	Method string `json:"method"` // 请求方法
	Count int64 `json:"count"` // 数量
}

type GetOperationStatisticsResPathstat struct {
	Path string `json:"path"` // 请求路径
	Count int64 `json:"count"` // 数量
}

type GetOperationRecordsByUserIdRes struct {
	Page int `json:"page"`
	Size int `json:"size"`
	List []GetOperationRecordsByUserIdResList `json:"list"` // []
	Total int64 `json:"total"` // 总数
}

type GetOperationRecordsByUserIdResList struct {
	Id int64 `json:"id"` // ID
	UserId int `json:"userId"` // 用户ID
	Ip string `json:"ip"` // 请求IP
	Method string `json:"method"` // 请求方法
	Path string `json:"path"` // 请求路径
	Status int `json:"status"` // 状态
	Latency int64 `json:"latency"` // 延迟(ns)
	ErrorMessage string `json:"errorMessage"` // 错误信息
	CreatedAt string `json:"createdAt"` // 创建时间
}

type GetOperationRecordsByTimeRangeRes struct {
	StartTime string `json:"startTime"`
	EndTime string `json:"endTime"`
	Page int `json:"page"`
	Size int `json:"size"`
	List []GetOperationRecordsByTimeRangeResList `json:"list"` // []
	Total int64 `json:"total"` // 总数
}

type GetOperationRecordsByTimeRangeResList struct {
	Id int64 `json:"id"` // ID
	UserId int `json:"userId"` // 用户ID
	Ip string `json:"ip"` // 请求IP
	Method string `json:"method"` // 请求方法
	Path string `json:"path"` // 请求路径
	Status int `json:"status"` // 状态
	Latency int64 `json:"latency"` // 延迟(ns)
	ErrorMessage string `json:"errorMessage"` // 错误信息
	CreatedAt string `json:"createdAt"` // 创建时间
}

type GetRecentOperationRecordsRes struct {
	Limit int `json:"limit"`
	List []GetRecentOperationRecordsResList `json:"list"` // []
}

type GetRecentOperationRecordsResList struct {
	Id int64 `json:"id"` // ID
	UserId int `json:"userId"` // 用户ID
	Ip string `json:"ip"` // 请求IP
	Method string `json:"method"` // 请求方法
	Path string `json:"path"` // 请求路径
	Status int `json:"status"` // 状态
	Latency int64 `json:"latency"` // 延迟(ns)
	ErrorMessage string `json:"errorMessage"` // 错误信息
	CreatedAt string `json:"createdAt"` // 创建时间
}

type GetErrorRecordsRes struct {
	Page int `json:"page"`
	Size int `json:"size"`
	List []GetErrorRecordsResList `json:"list"` // []
	Total int64 `json:"total"` // 总数
}

type GetErrorRecordsResList struct {
	Id int64 `json:"id"` // ID
	UserId int `json:"userId"` // 用户ID
	Ip string `json:"ip"` // 请求IP
	Method string `json:"method"` // 请求方法
	Path string `json:"path"` // 请求路径
	Status int `json:"status"` // 状态
	ErrorMessage string `json:"errorMessage"` // 错误信息
	CreatedAt string `json:"createdAt"` // 创建时间
}

type GetOperationRecordsByMethodRes struct {
	Page int `json:"page"`
	Size int `json:"size"`
	List []GetOperationRecordsByMethodResList `json:"list"` // []
	Total int64 `json:"total"` // 总数
}

type GetOperationRecordsByMethodResList struct {
	Id int64 `json:"id"` // ID
	UserId int `json:"userId"` // 用户ID
	Ip string `json:"ip"` // 请求IP
	Method string `json:"method"` // 请求方法
	Path string `json:"path"` // 请求路径
	Status int `json:"status"` // 状态
	Latency int64 `json:"latency"` // 延迟(ns)
	ErrorMessage string `json:"errorMessage"` // 错误信息
	CreatedAt string `json:"createdAt"` // 创建时间
}

type GetOperationRecordsByPathRes struct {
	Path string `json:"path"`
	Page int `json:"page"`
	Size int `json:"size"`
	List []GetOperationRecordsByPathResList `json:"list"` // []
	Total int64 `json:"total"` // 总数
}

type GetOperationRecordsByPathResList struct {
	Id int64 `json:"id"` // ID
	UserId int `json:"userId"` // 用户ID
	Ip string `json:"ip"` // 请求IP
	Method string `json:"method"` // 请求方法
	Path string `json:"path"` // 请求路径
	Status int `json:"status"` // 状态
	Latency int64 `json:"latency"` // 延迟(ns)
	ErrorMessage string `json:"errorMessage"` // 错误信息
	CreatedAt string `json:"createdAt"` // 创建时间
}

type CountTodayOperationsRes struct {
	Count int64 `json:"count"` // 今日操作总数
	SuccessCount int64 `json:"successCount"` // 成功数
	ErrorCount int64 `json:"errorCount"` // 异常数
}
