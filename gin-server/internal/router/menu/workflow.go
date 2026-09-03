//router 解析
package menu

import (
	"github.com/gin-gonic/gin"
	"shack/internal/model/common"
	"shack/internal/utils"
)

type WorkflowRouter struct{}

func (s *WorkflowRouter) InitWorkflowRouter(Router *gin.RouterGroup, apiSet map[string]struct{}) (R gin.IRoutes) {
	workflowRouter := Router.Group("")
	utils.RegisterApi(workflowRouter, apiSet, "",
		utils.NewRegisterApiParam(common.HttpGet, "/api/v1/plans", "工作流模块(需登录) 对应 C++ server WorkflowRoutes (plans/cooking-sessions/feedback) 获取本周餐食计划", workflowApi.ListPlansHandler),
		utils.NewRegisterApiParam(common.HttpPost, "/api/v1/plans", "保存餐食计划", workflowApi.SavePlanHandler),
		utils.NewRegisterApiParam(common.HttpPost, "/api/v1/cooking-sessions", "创建做饭会话", workflowApi.CreateCookingSessionHandler),
		utils.NewRegisterApiParam(common.HttpPut, "/api/v1/cooking-sessions/:id", "更新做饭会话进度", workflowApi.UpdateCookingSessionHandler),
		utils.NewRegisterApiParam(common.HttpPost, "/api/v1/feedback", "提交做饭反馈", workflowApi.SubmitFeedbackHandler),
	)
	return workflowRouter
}
