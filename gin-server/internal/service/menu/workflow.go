package menu

import (
	"github.com/gin-gonic/gin"
	req "shack/internal/model/menu/request"
	res "shack/internal/model/menu/response"
)

type WorkflowService struct{}

// 工作流模块(需登录) 对应 C++ server WorkflowRoutes (plans/cooking-sessions/feedback) 获取本周餐食计划
// Auth: shack 
// Github: https://github.com/hishack
// Time : 2026年09月02日 23:41:50
func (s *WorkflowService) ListPlans(
	ctx *gin.Context,
	r req.ListPlansReq,
) (err error) {
	return nil
}

// 保存餐食计划
// Auth: shack 
// Github: https://github.com/hishack
// Time : 2026年09月02日 23:41:50
func (s *WorkflowService) SavePlan(
	ctx *gin.Context,
	r req.SavePlanReq,
) (rs res.SavePlanRes, err error) {
	return rs, nil
}

// 创建做饭会话
// Auth: shack 
// Github: https://github.com/hishack
// Time : 2026年09月02日 23:41:50
func (s *WorkflowService) CreateCookingSession(
	ctx *gin.Context,
	r req.CreateCookingSessionReq,
) (rs res.CreateCookingSessionRes, err error) {
	return rs, nil
}

// 更新做饭会话进度
// Auth: shack 
// Github: https://github.com/hishack
// Time : 2026年09月02日 23:41:50
func (s *WorkflowService) UpdateCookingSession(
	ctx *gin.Context,
	r req.UpdateCookingSessionReq,
) (rs res.UpdateCookingSessionRes, err error) {
	return rs, nil
}

// 提交做饭反馈
// Auth: shack 
// Github: https://github.com/hishack
// Time : 2026年09月02日 23:41:50
func (s *WorkflowService) SubmitFeedback(
	ctx *gin.Context,
	r req.SubmitFeedbackReq,
) (rs res.SubmitFeedbackRes, err error) {
	return rs, nil
}

