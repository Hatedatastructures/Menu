package menu

import api "shack/internal/api/v1"

type RouterGroup struct {
	AdminRouter
	AuthRouter
	HealthRouter
	RecipeRouter
	WorkflowRouter
}

var (
	adminApi = api.ApiGroupApp.MenuApiGroup.AdminApi
	authApi = api.ApiGroupApp.MenuApiGroup.AuthApi
	healthApi = api.ApiGroupApp.MenuApiGroup.HealthApi
	recipeApi = api.ApiGroupApp.MenuApiGroup.RecipeApi
	workflowApi = api.ApiGroupApp.MenuApiGroup.WorkflowApi
)
