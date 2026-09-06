package response

type SavePlanRes struct {
	Id string `json:"id"` // 计划ID
	PlanDate string `json:"planDate"` // 计划日期
	Items []SavePlanResItem `json:"items"`
	CombinedIngredients []SavePlanResCombinedingredient `json:"combinedIngredients"`
}

type SavePlanResItem struct {
	Id string `json:"id"`
	RecipeId string `json:"recipeId"`
	RecipeName string `json:"recipeName"`
	ImagePath string `json:"imagePath"`
	Servings int `json:"servings"`
	SortOrder int `json:"sortOrder"`
	Ingredients []SavePlanResItemIngredient `json:"ingredients"`
}

type SavePlanResItemIngredient struct {
	IngredientId string `json:"ingredientId"`
	IngredientName string `json:"ingredientName"`
	Category string `json:"category"`
	DefaultUnit string `json:"defaultUnit"`
	IsPantryStaple bool `json:"isPantryStaple"`
	Quantity float64 `json:"quantity"`
	Unit string `json:"unit"`
	Required bool `json:"required"`
	Preparation string `json:"preparation"`
}

type SavePlanResCombinedingredient struct {
	IngredientId string `json:"ingredientId"`
	IngredientName string `json:"ingredientName"`
	Category string `json:"category"`
	DefaultUnit string `json:"defaultUnit"`
	IsPantryStaple bool `json:"isPantryStaple"`
	Quantity float64 `json:"quantity"`
	Unit string `json:"unit"`
	Required bool `json:"required"`
	Preparation string `json:"preparation"`
}

type CreateCookingSessionRes struct {
	Id string `json:"id"` // 会话ID
	RecipeId string `json:"recipeId"` // 菜谱ID
	CurrentStepOrder int `json:"currentStepOrder"` // 当前步骤序号
	State string `json:"state"` // 状态
	StartedAt string `json:"startedAt"` // 开始时间
	UpdatedAt string `json:"updatedAt"` // 更新时间
}

type UpdateCookingSessionRes struct {
	Id string `json:"id"` // 会话ID
	RecipeId string `json:"recipeId"` // 菜谱ID
	CurrentStepOrder int `json:"currentStepOrder"` // 当前步骤序号
	State string `json:"state"` // 状态
	StartedAt string `json:"startedAt"` // 开始时间
	UpdatedAt string `json:"updatedAt"` // 更新时间
}

type SubmitFeedbackRes struct {
	Id string `json:"id"` // 反馈ID
	RecipeId string `json:"recipeId"` // 菜谱ID
	Outcome string `json:"outcome"` // 结果
	Tags []string `json:"tags"` // 标签列表
	Comment string `json:"comment"` // 评论
	CreatedAt string `json:"createdAt"` // 创建时间
}
