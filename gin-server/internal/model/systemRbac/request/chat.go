package request

type SendChatMessageReq struct {
	ReceiverId int64 `json:"receiverId" form:"receiverId"` // 接收者ID
	Content string `json:"content" form:"content"` // 消息内容
	MsgType int `json:"msgType" form:"msgType"` // 消息类型(1文本2图片3文件)
}

type GetChatHistoryReq struct {
	TargetId int `json:"targetId" form:"targetId"`
	Page int `json:"page" form:"page"`
	PageSize int `json:"pageSize" form:"pageSize"`
}

type MarkChatAsReadReq struct {
	SenderId int `json:"senderId" form:"senderId"`
}

type CheckUserOnlineReq struct {
	UserId int `json:"userId" form:"userId"`
}

type ClearChatHistoryReq struct {
	TargetId int `json:"targetId" form:"targetId"`
}

type BlockUserReq struct {
	TargetId int `json:"targetId" form:"targetId"`
}

type UnblockUserReq struct {
	TargetId int `json:"targetId" form:"targetId"`
}

type CheckIsBlockedReq struct {
	TargetId int `json:"targetId" form:"targetId"`
}
