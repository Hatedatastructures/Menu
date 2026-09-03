//router 解析
package systemRbac

import (
	"github.com/gin-gonic/gin"
	"shack/internal/model/common"
	"shack/internal/utils"
)

type LoginLogRouter struct{}

func (s *LoginLogRouter) InitLoginLogRouter(Router *gin.RouterGroup, apiSet map[string]struct{}) (R gin.IRoutes) {
	loginLogRouter := Router.Group("")
	utils.RegisterApi(loginLogRouter, apiSet, "",
		utils.NewRegisterApiParam(common.HttpGet, "/loginLog", "获取登录日志列表-后台使用", loginLogApi.GetLoginLogListHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/loginLog/user/:username", "根据用户名查询登录日志-后台使用", loginLogApi.GetLoginLogByUsernameHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/loginLog/status/:status", "根据状态查询登录日志-后台使用", loginLogApi.GetLoginLogByStatusHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/loginLog/range", "按时间范围查询登录日志-后台使用", loginLogApi.GetLoginLogByTimeRangeHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/loginLog/recent", "获取最近N条登录记录-后台使用", loginLogApi.GetRecentLoginLogHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/loginLog/statistics", "获取登录统计信息-后台使用", loginLogApi.GetLoginStatisticsHandler),
		utils.NewRegisterApiParam(common.HttpDelete, "/loginLog/clear", "清空登录日志-后台使用", loginLogApi.ClearLoginLogHandler),
		utils.NewRegisterApiParam(common.HttpPost, "/loginLog/delete", "批量删除登录日志-后台使用", loginLogApi.DeleteLoginLogHandler),
	)
	return loginLogRouter
}
