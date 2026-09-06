//router 解析
package systemRbac

import (
	"github.com/gin-gonic/gin"
	"shack/internal/model/common"
	"shack/internal/utils"
)

type NoticeRouter struct{}

func (s *NoticeRouter) InitNoticeRouter(Router *gin.RouterGroup, apiSet map[string]struct{}) (R gin.IRoutes) {
	noticeRouter := Router.Group("")
	utils.RegisterApi(noticeRouter, apiSet, "",
		utils.NewRegisterApiParam(common.HttpGet, "/sys/notice/page", "分页查询通知列表-后台使用", noticeApi.GetNoticePageHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/sys/notice/my", "获取当前用户通知列表-前台使用", noticeApi.GetMyNoticesHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/sys/notice/:id", "获取通知详情-前台/后台使用", noticeApi.GetNoticeDetailHandler),
		utils.NewRegisterApiParam(common.HttpPost, "/sys/notice", "创建通知-后台使用", noticeApi.CreateNoticeHandler),
		utils.NewRegisterApiParam(common.HttpPut, "/sys/notice", "更新通知-后台使用", noticeApi.UpdateNoticeHandler),
		utils.NewRegisterApiParam(common.HttpDelete, "/sys/notice/:id", "删除通知-后台使用", noticeApi.DeleteNoticeHandler),
		utils.NewRegisterApiParam(common.HttpPost, "/sys/notice/:id/publish", "发布通知-后台使用", noticeApi.PublishNoticeHandler),
		utils.NewRegisterApiParam(common.HttpPost, "/sys/notice/:id/read", "标记通知为已读-前台使用", noticeApi.MarkNoticeAsReadHandler),
		utils.NewRegisterApiParam(common.HttpPost, "/sys/notice/read-all", "标记所有通知为已读-前台使用", noticeApi.MarkAllNoticesAsReadHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/sys/notice/unread-count", "获取未读通知数量-前台使用", noticeApi.GetNoticeUnreadCountHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/sys/notice/channels", "获取推送渠道-后台使用", noticeApi.GetNoticeChannelsHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/sys/notice/:id/send-logs", "获取通知推送记录-后台使用", noticeApi.GetNoticeSendLogsHandler),
		utils.NewRegisterApiParam(common.HttpPost, "/sys/notice/:id/retry", "重试失败推送-后台使用", noticeApi.RetryNoticeSendHandler),
	)
	return noticeRouter
}
