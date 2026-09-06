package request

type ListPlansReq struct {
	From string `json:"from" form:"from"`
	To string `json:"to" form:"to"`
}

type SavePlanReq struct {
	PlanDate string `json:"planDate" form:"planDate" validate:"required,max=10,min=10"` // 计划日期(YYYY-MM-DD)
	Items []SavePlanReqItem `json:"items" form:"items"`
}

type SavePlanReqItem struct {
	RecipeId string `json:"recipeId" form:"recipeId"` // 菜谱ID
	Servings int `json:"servings" form:"servings"` // 份数
	SortOrder int `json:"sortOrder" form:"sortOrder"` // 排序(可选)
}

type CreateCookingSessionReq struct {
	RecipeId string `json:"recipeId" form:"recipeId" validate:"required,max=128,min=1"` // 菜谱ID
}

type UpdateCookingSessionReq struct {
	CurrentStepOrder int `json:"currentStepOrder" form:"currentStepOrder" validate:"required,gte=1,lte=128"` // 当前步骤序号
	State string `json:"state" form:"state" validate:"required,max=16,min=1"` // 状态
}

type SubmitFeedbackReq struct {
	RecipeId string `json:"recipeId" form:"recipeId" validate:"required,max=128,min=1"` // 菜谱ID
	Outcome string `json:"outcome" form:"outcome" validate:"required,max=16,min=1"` // 结果:success/failure/partial
	Tags []string `json:"tags" form:"tags"` // 标签列表
	Comment string `json:"comment" form:"comment" validate:"required,max=1000"` // 评论
}
