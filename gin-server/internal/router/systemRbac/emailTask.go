//router 解析
package systemRbac

import (
	"github.com/gin-gonic/gin"
	"shack/internal/model/common"
	"shack/internal/utils"
)

type EmailTaskRouter struct{}

func (s *EmailTaskRouter) InitEmailTaskRouter(Router *gin.RouterGroup, apiSet map[string]struct{}) (R gin.IRoutes) {
	emailTaskRouter := Router.Group("")
	utils.RegisterApi(emailTaskRouter, apiSet, "",
		utils.NewRegisterApiParam(common.HttpPost, "/email/task", "创建邮件任务-后台使用", emailTaskApi.CreateEmailTaskHandler),
		utils.NewRegisterApiParam(common.HttpPut, "/email/task/:id", "更新邮件任务-后台使用", emailTaskApi.UpdateEmailTaskHandler),
		utils.NewRegisterApiParam(common.HttpDelete, "/email/task/:id", "删除邮件任务-后台使用", emailTaskApi.DeleteEmailTaskHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/email/task/list", "获取邮件任务列表-后台使用", emailTaskApi.GetEmailTaskListHandler),
		utils.NewRegisterApiParam(common.HttpPost, "/email/task/start/:id", "启动邮件任务-后台使用", emailTaskApi.StartEmailTaskHandler),
		utils.NewRegisterApiParam(common.HttpPost, "/email/task/stop/:id", "停止邮件任务-后台使用", emailTaskApi.StopEmailTaskHandler),
		utils.NewRegisterApiParam(common.HttpPost, "/email/task/run/:id", "手动执行邮件任务-后台使用", emailTaskApi.RunEmailTaskHandler),
	)
	return emailTaskRouter
}
