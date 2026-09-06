//router 解析
package systemRbac

import (
	"github.com/gin-gonic/gin"
	"shack/internal/model/common"
	"shack/internal/utils"
)

type AiConfigRouter struct{}

func (s *AiConfigRouter) InitAiConfigRouter(Router *gin.RouterGroup, apiSet map[string]struct{}) (R gin.IRoutes) {
	aiConfigRouter := Router.Group("")
	utils.RegisterApi(aiConfigRouter, apiSet, "",
		utils.NewRegisterApiParam(common.HttpPost, "/api/ai/config", "创建AI配置-后台使用", aiConfigApi.CreateAiConfigHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/api/ai/config", "获取AI配置列表-后台使用", aiConfigApi.GetAiConfigListHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/api/ai/config/:id", "获取AI配置详情-后台使用", aiConfigApi.GetAiConfigHandler),
		utils.NewRegisterApiParam(common.HttpPut, "/api/ai/config/:id", "更新AI配置-后台使用", aiConfigApi.UpdateAiConfigHandler),
		utils.NewRegisterApiParam(common.HttpDelete, "/api/ai/config/:id", "删除AI配置-后台使用", aiConfigApi.DeleteAiConfigHandler),
		utils.NewRegisterApiParam(common.HttpDelete, "/api/ai/config", "批量删除AI配置-后台使用", aiConfigApi.BatchDeleteAiConfigHandler),
		utils.NewRegisterApiParam(common.HttpPost, "/api/ai/test", "测试AI接口-后台使用", aiConfigApi.TestAiConfigHandler),
		utils.NewRegisterApiParam(common.HttpPost, "/api/ai/chat", "AI对话调用-后台使用", aiConfigApi.AiChatHandler),
		utils.NewRegisterApiParam(common.HttpPut, "/api/ai/config/:id/status", "切换AI配置启用状态-后台使用", aiConfigApi.ToggleAiConfigStatusHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/api/ai/config/enabled", "获取启用的AI配置列表-前台使用", aiConfigApi.GetEnabledAiConfigsHandler),
	)
	return aiConfigRouter
}
