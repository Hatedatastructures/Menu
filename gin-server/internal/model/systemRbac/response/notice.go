package response


import "time"

type GetNoticePageRes struct {

	Title string `json:"title"`
	NoticeType int `json:"noticeType"`
	Status int `json:"status"`
	List []GetNoticePageResList `json:"list"` // []
	Total int64 `json:"total"` // 总数
	Page int `json:"page"` // 当前页
	PageSize int `json:"pageSize"` // 每页数量
}

type GetNoticePageResList struct {
	Id int64 `json:"id"` // 通知ID
	Title string `json:"title"` // 通知标题
	Content string `json:"content"` // 通知内容
	NoticeType int `json:"noticeType"` // 通知类型
	Channels string `json:"channels"` // 推送渠道
	TargetType int `json:"targetType"` // 推送对象类型
	TargetIds string `json:"targetIds"` // 推送对象ID
	Status int `json:"status"` // 状态
	CreateBy int64 `json:"createBy"` // 创建者ID
	CreateName string `json:"createName"` // 创建者名称
	CreateTime time.Time `json:"createTime"` // 创建时间
}

type GetMyNoticesRes struct {

	IsRead int `json:"isRead"`
	List []GetMyNoticesResList `json:"list"` // []
	Total int64 `json:"total"` // 总数
	Page int `json:"page"` // 当前页
	PageSize int `json:"pageSize"` // 每页数量
}

type GetMyNoticesResList struct {
	Id int64 `json:"id"` // 通知ID
	Title string `json:"title"` // 通知标题
	Content string `json:"content"` // 通知内容
	NoticeType int `json:"noticeType"` // 通知类型
	Status int `json:"status"` // 状态
	CreateTime time.Time `json:"createTime"` // 创建时间
}

type GetNoticeDetailRes struct {
	Id int64 `json:"id"` // 通知ID
	Title string `json:"title"` // 通知标题
	Content string `json:"content"` // 通知内容
	NoticeType int `json:"noticeType"` // 通知类型
	Channels string `json:"channels"` // 推送渠道
	TargetType int `json:"targetType"` // 推送对象类型
	TargetIds string `json:"targetIds"` // 推送对象ID
	Status int `json:"status"` // 状态
	CreateBy int64 `json:"createBy"` // 创建者ID
	CreateName string `json:"createName"` // 创建者名称
	CreateTime time.Time `json:"createTime"` // 创建时间
}

type GetNoticeUnreadCountRes struct {
	Count int `json:"count"` // 未读数量
}

type GetNoticeChannelsRes struct {
	List []GetNoticeChannelsResList `json:"list"` // []
}

type GetNoticeChannelsResList struct {
	Channel string `json:"channel"` // 渠道标识
	Name string `json:"name"` // 渠道名称
	Enabled bool `json:"enabled"` // 是否启用
}

type GetNoticeSendLogsRes struct {
	List []GetNoticeSendLogsResList `json:"list"` // []
}

type GetNoticeSendLogsResList struct {
	Id int64 `json:"id"` // 记录ID
	NoticeId int64 `json:"noticeId"` // 通知ID
	Channel string `json:"channel"` // 推送渠道
	TargetId int64 `json:"targetId"` // 推送目标ID
	Status int `json:"status"` // 状态(0失败1成功)
	ErrorMsg string `json:"errorMsg"` // 错误信息
	SendTime time.Time `json:"sendTime"` // 发送时间
}
