//router 解析
package gen

import (
	"github.com/gin-gonic/gin"
	"shack/internal/model/common"
	"shack/internal/utils"
)

type HistoryRouter struct{}

func (s *HistoryRouter) InitHistoryRouter(Router *gin.RouterGroup, apiSet map[string]struct{}) (R gin.IRoutes) {
	historyRouter := Router.Group("")
	utils.RegisterApi(historyRouter, apiSet, "",
		utils.NewRegisterApiParam(common.HttpGet, "/gen/history", "获取生成历史列表", historyApi.GetHistoryListHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/gen/history/:id", "获取历史详情(含答案解析)", historyApi.GetHistoryDetailHandler),
		utils.NewRegisterApiParam(common.HttpDelete, "/gen/history/:id", "删除历史记录", historyApi.DeleteHistoryHandler),
	)
	return historyRouter
}
