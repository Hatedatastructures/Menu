//router 解析
package systemRbac

import (
	"github.com/gin-gonic/gin"
	"shack/internal/model/common"
	"shack/internal/utils"
)

type EmailStatisticsRouter struct{}

func (s *EmailStatisticsRouter) InitEmailStatisticsRouter(Router *gin.RouterGroup, apiSet map[string]struct{}) (R gin.IRoutes) {
	emailStatisticsRouter := Router.Group("")
	utils.RegisterApi(emailStatisticsRouter, apiSet, "",
		utils.NewRegisterApiParam(common.HttpGet, "/email/statistics", "获取邮件统计-后台使用", emailStatisticsApi.GetEmailStatisticsHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/email/queue/status", "获取队列状态-后台使用", emailStatisticsApi.GetEmailQueueStatusHandler),
	)
	return emailStatisticsRouter
}
