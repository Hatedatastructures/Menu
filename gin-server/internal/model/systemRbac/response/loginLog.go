package response

type GetLoginLogListRes struct {
	Page int `json:"page"`
	Size int `json:"size"`
	Username string `json:"username"`
	Status string `json:"status"`
	StartTime string `json:"startTime"`
	EndTime string `json:"endTime"`
	List []GetLoginLogListResList `json:"list"` // []
	Total int64 `json:"total"` // 总数
}

type GetLoginLogListResList struct {
	Id int64 `json:"id"` // ID
	Username string `json:"username"` // 用户名
	Ip string `json:"ip"` // IP地址
	Location string `json:"location"` // 登录地点
	Browser string `json:"browser"` // 浏览器
	Os string `json:"os"` // 操作系统
	Status string `json:"status"` // 状态
	LoginTime string `json:"loginTime"` // 登录时间
}

type GetLoginLogByUsernameRes struct {
	Page int `json:"page"`
	Size int `json:"size"`
	List []GetLoginLogByUsernameResList `json:"list"` // []
	Total int64 `json:"total"` // 总数
}

type GetLoginLogByUsernameResList struct {
	Id int64 `json:"id"` // ID
	Username string `json:"username"` // 用户名
	Ip string `json:"ip"` // IP地址
	Location string `json:"location"` // 登录地点
	Browser string `json:"browser"` // 浏览器
	Os string `json:"os"` // 操作系统
	Status string `json:"status"` // 状态
	LoginTime string `json:"loginTime"` // 登录时间
}

type GetLoginLogByStatusRes struct {
	Page int `json:"page"`
	Size int `json:"size"`
	List []GetLoginLogByStatusResList `json:"list"` // []
	Total int64 `json:"total"` // 总数
}

type GetLoginLogByStatusResList struct {
	Id int64 `json:"id"` // ID
	Username string `json:"username"` // 用户名
	Ip string `json:"ip"` // IP地址
	Location string `json:"location"` // 登录地点
	Browser string `json:"browser"` // 浏览器
	Os string `json:"os"` // 操作系统
	Status string `json:"status"` // 状态
	LoginTime string `json:"loginTime"` // 登录时间
}

type GetLoginLogByTimeRangeRes struct {
	StartTime string `json:"startTime"`
	EndTime string `json:"endTime"`
	Page int `json:"page"`
	Size int `json:"size"`
	List []GetLoginLogByTimeRangeResList `json:"list"` // []
	Total int64 `json:"total"` // 总数
}

type GetLoginLogByTimeRangeResList struct {
	Id int64 `json:"id"` // ID
	Username string `json:"username"` // 用户名
	Ip string `json:"ip"` // IP地址
	Location string `json:"location"` // 登录地点
	Browser string `json:"browser"` // 浏览器
	Os string `json:"os"` // 操作系统
	Status string `json:"status"` // 状态
	LoginTime string `json:"loginTime"` // 登录时间
}

type GetRecentLoginLogRes struct {
	Limit int `json:"limit"`
	List []GetRecentLoginLogResList `json:"list"` // []
}

type GetRecentLoginLogResList struct {
	Id int64 `json:"id"` // ID
	Username string `json:"username"` // 用户名
	Ip string `json:"ip"` // IP地址
	Location string `json:"location"` // 登录地点
	Browser string `json:"browser"` // 浏览器
	Os string `json:"os"` // 操作系统
	Status string `json:"status"` // 状态
	LoginTime string `json:"loginTime"` // 登录时间
}

type GetLoginStatisticsRes struct {
	Total int64 `json:"total"` // 总数
	TodayCount int64 `json:"todayCount"` // 今日登录数
	SuccessCount int64 `json:"successCount"` // 成功数
	FailCount int64 `json:"failCount"` // 失败数
}
