//router 解析
package systemRbac

import (
	"github.com/gin-gonic/gin"
	"shack/internal/model/common"
	"shack/internal/utils"
)

type ApiTestRouter struct{}

func (s *ApiTestRouter) InitApiTestRouter(Router *gin.RouterGroup, apiSet map[string]struct{}) (R gin.IRoutes) {
	apiTestRouter := Router.Group("")
	utils.RegisterApi(apiTestRouter, apiSet, "",
		utils.NewRegisterApiParam(common.HttpPost, "/api-test/log", "保存API测试请求记录-前台使用", apiTestApi.SaveApiTestLogHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/api-test/log/list", "获取API测试记录列表-前台使用", apiTestApi.GetApiTestLogListHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/api-test/log/:id", "获取API测试记录详情-前台使用", apiTestApi.GetApiTestLogDetailHandler),
		utils.NewRegisterApiParam(common.HttpDelete, "/api-test/log/:id", "删除API测试记录-前台使用", apiTestApi.DeleteApiTestLogHandler),
		utils.NewRegisterApiParam(common.HttpDelete, "/api-test/log/clear", "清空API测试记录-前台使用", apiTestApi.ClearApiTestLogHandler),
	)
	return apiTestRouter
}
