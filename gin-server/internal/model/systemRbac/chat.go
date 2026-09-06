package systemRbac

import "shack/internal/global"

type ChatConversation struct {
	global.GVA_MODEL

	User1ID int64 `json:"user1Id" gorm:"column:user1_id;comment:用户1"`
	User2ID int64 `json:"user2Id" gorm:"column:user2_id;comment:用户2"`

	LastMsg     string `json:"lastMsg" gorm:"column:last_msg;comment:最后一条消息"`
	LastMsgTime *int64 `json:"lastMsgTime" gorm:"column:last_msg_time;comment:最后消息时间"`

	// 未读（分别存）
	User1Unread int `json:"user1Unread" gorm:"column:user1_unread;comment:用户1未读数"`
	User2Unread int `json:"user2Unread" gorm:"column:user2_unread;comment:用户2未读数"`
}

func (ChatConversation) TableName() string {
	return "chat_conversations"
}