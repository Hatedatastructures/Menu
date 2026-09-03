//router 解析
package gen

import (
	"github.com/gin-gonic/gin"
	"shack/internal/model/common"
	"shack/internal/utils"
)

type ApikeyRouter struct{}

func (s *ApikeyRouter) InitApikeyRouter(Router *gin.RouterGroup, apiSet map[string]struct{}) (R gin.IRoutes) {
	apikeyRouter := Router.Group("")
	utils.RegisterApi(apikeyRouter, apiSet, "",
		utils.NewRegisterApiParam(common.HttpGet, "/gen/apikey", "获取我的API Key配置", apikeyApi.GetMyApiKeyHandler),
		utils.NewRegisterApiParam(common.HttpPost, "/gen/apikey", "保存/更新API Key", apikeyApi.SaveApiKeyHandler),
		utils.NewRegisterApiParam(common.HttpDelete, "/gen/apikey", "删除API Key", apikeyApi.DeleteApiKeyHandler),
		utils.NewRegisterApiParam(common.HttpPost, "/gen/apikey/test", "测试API Key是否可用", apikeyApi.TestApiKeyHandler),
	)
	return apikeyRouter
}
