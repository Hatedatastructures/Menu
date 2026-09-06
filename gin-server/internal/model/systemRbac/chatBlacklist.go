package systemRbac

import "shack/internal/global"

// ChatBlacklist 聊天黑名单
type ChatBlacklist struct {
	global.GVA_MODEL

	UserId         int64 `json:"userId" gorm:"column:user_id;index;not null;comment:用户ID"`
	BlockedUserId  int64 `json:"blockedUserId" gorm:"column:blocked_user_id;index;not null;comment:被拉黑用户ID"`
	BlockedReason  string `json:"blockedReason" gorm:"column:blocked_reason;type:varchar(500);comment:拉黑原因"`
}

func (ChatBlacklist) TableName() string {
	return "chat_blacklists"
}
