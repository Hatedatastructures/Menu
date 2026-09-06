//router 解析
package systemRbac

import (
	"github.com/gin-gonic/gin"
	"shack/internal/model/common"
	"shack/internal/utils"
)

type ApiRouter struct{}

func (s *ApiRouter) InitApiRouter(Router *gin.RouterGroup, apiSet map[string]struct{}) (R gin.IRoutes) {
	apiRouter := Router.Group("")
	utils.RegisterApi(apiRouter, apiSet, "",
		utils.NewRegisterApiParam(common.HttpPost, "/api/apis", "创建API-后台使用", apiApi.CreateApiHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/api/apis", "API列表-后台使用", apiApi.GetApiListHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/api/apis/:id", "API详情-后台使用", apiApi.GetApiDetailHandler),
		utils.NewRegisterApiParam(common.HttpPut, "/api/apis/:id", "更新API-后台使用", apiApi.UpdateApiHandler),
		utils.NewRegisterApiParam(common.HttpDelete, "/api/apis/:id", "删除API-后台使用", apiApi.DeleteApiHandler),
	)
	return apiRouter
}
