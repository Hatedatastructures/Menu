//router 解析
package systemRbac

import (
	"github.com/gin-gonic/gin"
	"shack/internal/middleware"
	"shack/internal/model/common"
	"shack/internal/utils"
)

type OprationRouter struct{}

func (s *OprationRouter) InitOprationRouter(Router *gin.RouterGroup, apiSet map[string]struct{}) (R gin.IRoutes) {
	oprationRouter := Router.Group("").Use(middleware.OperationRecord())
	utils.RegisterApi(oprationRouter, apiSet, "",
		utils.NewRegisterApiParam(common.HttpDelete, "/operation/clear", "清空所有操作记录-后台使用", oprationApi.ClearOperationRecordHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/operation/statistics", "获取操作统计信息-后台使用", oprationApi.GetOperationStatisticsHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/operation/user/:userId", "根据用户ID获取操作记录-后台使用", oprationApi.GetOperationRecordsByUserIdHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/operation/range", "按时间范围查询操作记录-后台使用", oprationApi.GetOperationRecordsByTimeRangeHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/operation/recent", "获取最近N条记录-后台使用", oprationApi.GetRecentOperationRecordsHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/operation/errors", "获取异常请求记录-后台使用", oprationApi.GetErrorRecordsHandler),
		utils.NewRegisterApiParam(common.HttpPost, "/operation/expired/delete", "批量删除过期记录-后台使用", oprationApi.DeleteExpiredRecordsHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/operation/method/:method", "按请求方法查询记录-后台使用", oprationApi.GetOperationRecordsByMethodHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/operation/path", "按请求路径查询记录-后台使用", oprationApi.GetOperationRecordsByPathHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/operation/today/count", "统计今日操作数-后台使用", oprationApi.CountTodayOperationsHandler),
	)
	return oprationRouter
}
