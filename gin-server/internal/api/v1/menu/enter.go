package menu

import "shack/internal/service"

type ApiGroup struct {
	AdminApi
	AuthApi
	HealthApi
	RecipeApi
	WorkflowApi
}

var (
	adminService = service.ServiceGroupApp.MenuServiceGroup.AdminService
	authService = service.ServiceGroupApp.MenuServiceGroup.AuthService
	healthService = service.ServiceGroupApp.MenuServiceGroup.HealthService
	recipeService = service.ServiceGroupApp.MenuServiceGroup.RecipeService
	workflowService = service.ServiceGroupApp.MenuServiceGroup.WorkflowService
)
