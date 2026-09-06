//router 解析
package gen

import (
	"github.com/gin-gonic/gin"
	"shack/internal/model/common"
	"shack/internal/utils"
)

type GenerateRouter struct{}

func (s *GenerateRouter) InitGenerateRouter(Router *gin.RouterGroup, apiSet map[string]struct{}) (R gin.IRoutes) {
	generateRouter := Router.Group("")
	utils.RegisterApi(generateRouter, apiSet, "",
		utils.NewRegisterApiParam(common.HttpPost, "/gen/generate", "上传PDF并生成答案解析(异步)", generateApi.GenerateHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/gen/generate/:id", "获取生成进度/结果", generateApi.GetGenerateResultHandler),
	)
	return generateRouter
}
