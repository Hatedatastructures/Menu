package systemRbac

import "shack/internal/global"

type ChatMessage struct {
	global.GVA_MODEL

	ConversationID int64 `json:"conversationId" gorm:"column:conversation_id;index;comment:会话ID"`

	SenderID   int64  `json:"senderId" gorm:"column:sender_id;comment:发送者"`
	ReceiverID int64  `json:"receiverId" gorm:"column:receiver_id;comment:接收者"`

	Content string `json:"content" gorm:"column:content;type:text;comment:消息内容"`

	MsgType int `json:"msgType" gorm:"column:msg_type;comment:类型(1文本2图片3文件)"`

	Status int `json:"status" gorm:"column:status;comment:状态(0发送中1成功2失败)"`

	IsRead bool `json:"isRead" gorm:"column:is_read;comment:是否已读"`

	FileURL string `json:"fileUrl" gorm:"column:file_url;comment:文件地址"`
}

func (ChatMessage) TableName() string {
	return "chat_messages"
}