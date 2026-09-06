//router 解析
package systemRbac

import (
	"github.com/gin-gonic/gin"
	"shack/internal/model/common"
	"shack/internal/utils"
)

type ChatRouter struct{}

func (s *ChatRouter) InitChatRouter(Router *gin.RouterGroup, apiSet map[string]struct{}) (R gin.IRoutes) {
	chatRouter := Router.Group("")
	utils.RegisterApi(chatRouter, apiSet, "",
		utils.NewRegisterApiParam(common.HttpPost, "/sys/chat/send", "发送私聊消息-前台使用", chatApi.SendChatMessageHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/sys/chat/history/:targetId", "获取聊天记录-前台使用", chatApi.GetChatHistoryHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/sys/chat/contacts", "获取最近联系人-前台使用", chatApi.GetRecentContactsHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/sys/chat/users", "获取用户列表(用于选择聊天对象)-前台使用", chatApi.GetChatUsersHandler),
		utils.NewRegisterApiParam(common.HttpPost, "/sys/chat/read/:senderId", "标记消息为已读-前台使用", chatApi.MarkChatAsReadHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/sys/chat/unread-count", "获取未读消息数量-前台使用", chatApi.GetChatUnreadCountHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/sys/chat/online/:userId", "检查用户是否在线-前台使用", chatApi.CheckUserOnlineHandler),
		utils.NewRegisterApiParam(common.HttpDelete, "/sys/chat/clear/:targetId", "清空聊天记录-前台使用", chatApi.ClearChatHistoryHandler),
		utils.NewRegisterApiParam(common.HttpPost, "/sys/chat/block/:targetId", "拉黑用户-前台使用", chatApi.BlockUserHandler),
		utils.NewRegisterApiParam(common.HttpDelete, "/sys/chat/block/:targetId", "取消拉黑-前台使用", chatApi.UnblockUserHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/sys/chat/blacklist", "获取黑名单列表-前台使用", chatApi.GetBlacklistHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/sys/chat/blocked/:targetId", "检查是否拉黑-前台使用", chatApi.CheckIsBlockedHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/sys/chat/stats", "获取消息统计-前台使用", chatApi.GetMessageStatsHandler),
	)
	return chatRouter
}
