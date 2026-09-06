package request

type GetNoticePageReq struct {
	Page int `json:"page" form:"page"`
	PageSize int `json:"pageSize" form:"pageSize"`
	Title string `json:"title" form:"title"`
	NoticeType int `json:"noticeType" form:"noticeType"`
	Status int `json:"status" form:"status"`
}

type GetMyNoticesReq struct {
	Page int `json:"page" form:"page"`
	PageSize int `json:"pageSize" form:"pageSize"`
	IsRead int `json:"isRead" form:"isRead"`
}

type GetNoticeDetailReq struct {
	Id int `json:"id" form:"id"`
}

type CreateNoticeReq struct {
	Title string `json:"title" form:"title"` // 通知标题
	Content string `json:"content" form:"content"` // 通知内容
	NoticeType int `json:"noticeType" form:"noticeType"` // 通知类型(1通知2公告)
	Channels string `json:"channels" form:"channels"` // 推送渠道(JSON数组)
	TargetType int `json:"targetType" form:"targetType"` // 推送对象类型(1指定用户2按部门3全部)
	TargetIds string `json:"targetIds" form:"targetIds"` // 推送对象ID(JSON数组)
}

type UpdateNoticeReq struct {
	Id int64 `json:"id" form:"id"` // 通知ID
	Title string `json:"title" form:"title"` // 通知标题
	Content string `json:"content" form:"content"` // 通知内容
	NoticeType int `json:"noticeType" form:"noticeType"` // 通知类型
	Channels string `json:"channels" form:"channels"` // 推送渠道
	TargetType int `json:"targetType" form:"targetType"` // 推送对象类型
	TargetIds string `json:"targetIds" form:"targetIds"` // 推送对象ID
	Status int `json:"status" form:"status"` // 状态
}

type DeleteNoticeReq struct {
	Id int `json:"id" form:"id"`
}

type PublishNoticeReq struct {
	Id int `json:"id" form:"id"`
}

type MarkNoticeAsReadReq struct {
	Id int `json:"id" form:"id"`
}

type GetNoticeSendLogsReq struct {
	Id int `json:"id" form:"id"`
}

type RetryNoticeSendReq struct {
	Id int `json:"id" form:"id"`
	Channel string `json:"channel" form:"channel"`
}
