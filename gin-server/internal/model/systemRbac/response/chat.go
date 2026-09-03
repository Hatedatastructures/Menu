package response


import "time"

type SendChatMessageRes struct {
	Id int64 `json:"id"` // 消息ID
	SenderId int64 `json:"senderId"` // 发送者ID
	SenderName string `json:"senderName"` // 发送者名称
	ReceiverId int64 `json:"receiverId"` // 接收者ID
	Content string `json:"content"` // 消息内容
	MsgType int `json:"msgType"` // 消息类型
	IsRead int `json:"isRead"` // 是否已读
	SendTime time.Time `json:"sendTime"` // 发送时间
}

type GetChatHistoryRes struct {

	List []GetChatHistoryResList `json:"list"` // []
	Total int64 `json:"total"` // 总数
	Page int `json:"page"` // 当前页
	PageSize int `json:"pageSize"` // 每页数量
}

type GetChatHistoryResList struct {
	Id int64 `json:"id"` // 消息ID
	SenderId int64 `json:"senderId"` // 发送者ID
	SenderName string `json:"senderName"` // 发送者名称
	SenderAvatar string `json:"senderAvatar"` // 发送者头像
	ReceiverId int64 `json:"receiverId"` // 接收者ID
	Content string `json:"content"` // 消息内容
	MsgType int `json:"msgType"` // 消息类型
	IsRead int `json:"isRead"` // 是否已读
	SendTime time.Time `json:"sendTime"` // 发送时间
}

type GetRecentContactsRes struct {
	List []GetRecentContactsResList `json:"list"` // []
}

type GetRecentContactsResList struct {
	Id int64 `json:"id"` // 消息ID
	SenderId int64 `json:"senderId"` // 发送者ID
	SenderName string `json:"senderName"` // 发送者名称
	ReceiverId int64 `json:"receiverId"` // 接收者ID
	Content string `json:"content"` // 消息内容
	MsgType int `json:"msgType"` // 消息类型
	IsRead int `json:"isRead"` // 是否已读
	SendTime time.Time `json:"sendTime"` // 发送时间
}

type GetChatUsersRes struct {
	List []GetChatUsersResList `json:"list"` // []
}

type GetChatUsersResList struct {
	Id int64 `json:"id"` // 用户ID
	Username string `json:"username"` // 用户名
	Nickname string `json:"nickname"` // 昵称
	Avatar string `json:"avatar"` // 头像
	LastMessage string `json:"lastMessage"` // 最新消息
	LastMessageTime time.Time `json:"lastMessageTime"` // 最新消息时间
	IsBlocked bool `json:"isBlocked"` // 是否已拉黑
}

type GetChatUnreadCountRes struct {
	Count int `json:"count"` // 未读数量
}

type CheckUserOnlineRes struct {
	Online bool `json:"online"` // 是否在线
}

type GetBlacklistRes struct {
	List []GetBlacklistResList `json:"list"` // []
}

type GetBlacklistResList struct {
	Id int64 `json:"id"` // ID
	UserId int64 `json:"userId"` // 用户ID
	BlockedUserId int64 `json:"blockedUserId"` // 被拉黑用户ID
	BlockedUserName string `json:"blockedUserName"` // 被拉黑用户名
	BlockedUserAvatar string `json:"blockedUserAvatar"` // 被拉黑用户头像
	CreateTime time.Time `json:"createTime"` // 创建时间
}

type CheckIsBlockedRes struct {
	Blocked bool `json:"blocked"` // 是否拉黑
}

type GetMessageStatsRes struct {
	ChatCount int `json:"chatCount"` // 聊天未读数
}
