package menu

import (
	req "shack/internal/model/menu/request"
	_ "shack/internal/model/menu/response"
	"shack/internal/utils/validator"
	"shack/internal/vo"

	"github.com/gin-gonic/gin"
)

type WorkflowApi struct{}

// ListPlansHandler
// @Tags menuworkflowApi
// @Summary ListPlansHandler 工作流模块(需登录) 对应 C++ server WorkflowRoutes (plans/cooking-sessions/feedback) 获取本周餐食计划
// @Description ListPlansHandler 工作流模块(需登录) 对应 C++ server WorkflowRoutes (plans/cooking-sessions/feedback) 获取本周餐食计划
// @Security ApiKeyAuth
// @Param data body req.ListPlansReq true "请求参数"
// @Success 200 {object} vo.Result{}
// @Router /api/v1/plans [GET]
func (s *WorkflowApi) ListPlansHandler(c *gin.Context) {
	var req req.ListPlansReq
	// query 参数
	{
		val := c.Query("from")
		if val != "" {
			req.From = val

		}
	}
	{
		val := c.Query("to")
		if val != "" {
			req.To = val

		}
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	err := workflowService.ListPlans(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, _ := c.Get("planList")
	c.JSON(200, vo.Success(c, data))
}

// SavePlanHandler
// @Tags menuworkflowApi
// @Summary SavePlanHandler 保存餐食计划
// @Description SavePlanHandler 保存餐食计划
// @Security ApiKeyAuth
// @Param data body req.SavePlanReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.SavePlanRes}
// @Router /api/v1/plans [POST]
func (s *WorkflowApi) SavePlanHandler(c *gin.Context) {
	var req req.SavePlanReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := workflowService.SavePlan(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// CreateCookingSessionHandler
// @Tags menuworkflowApi
// @Summary CreateCookingSessionHandler 创建做饭会话
// @Description CreateCookingSessionHandler 创建做饭会话
// @Security ApiKeyAuth
// @Param data body req.CreateCookingSessionReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.CreateCookingSessionRes}
// @Router /api/v1/cooking-sessions [POST]
func (s *WorkflowApi) CreateCookingSessionHandler(c *gin.Context) {
	var req req.CreateCookingSessionReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := workflowService.CreateCookingSession(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// UpdateCookingSessionHandler
// @Tags menuworkflowApi
// @Summary UpdateCookingSessionHandler 更新做饭会话进度
// @Description UpdateCookingSessionHandler 更新做饭会话进度
// @Security ApiKeyAuth
// @Param data body req.UpdateCookingSessionReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.UpdateCookingSessionRes}
// @Router /api/v1/cooking-sessions/:id [PUT]
func (s *WorkflowApi) UpdateCookingSessionHandler(c *gin.Context) {
	var req req.UpdateCookingSessionReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := workflowService.UpdateCookingSession(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// SubmitFeedbackHandler
// @Tags menuworkflowApi
// @Summary SubmitFeedbackHandler 提交做饭反馈
// @Description SubmitFeedbackHandler 提交做饭反馈
// @Security ApiKeyAuth
// @Param data body req.SubmitFeedbackReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.SubmitFeedbackRes}
// @Router /api/v1/feedback [POST]
func (s *WorkflowApi) SubmitFeedbackHandler(c *gin.Context) {
	var req req.SubmitFeedbackReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := workflowService.SubmitFeedback(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}
