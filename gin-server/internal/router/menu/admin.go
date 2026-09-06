//router 解析
package menu

import (
	"github.com/gin-gonic/gin"
	"shack/internal/model/common"
	"shack/internal/utils"
)

type AdminRouter struct{}

func (s *AdminRouter) InitAdminRouter(Router *gin.RouterGroup, apiSet map[string]struct{}) (R gin.IRoutes) {
	adminRouter := Router.Group("")
	utils.RegisterApi(adminRouter, apiSet, "",
		utils.NewRegisterApiParam(common.HttpGet, "/api/v1/admin/stats", "获取管理后台统计", adminApi.GetAdminStatsHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/api/v1/admin/recipes", "获取全部菜谱(含所有状态)", adminApi.ListAllRecipesHandler),
		utils.NewRegisterApiParam(common.HttpPost, "/api/v1/admin/recipes", "创建菜谱", adminApi.CreateRecipeHandler),
		utils.NewRegisterApiParam(common.HttpPut, "/api/v1/admin/recipes/:id", "更新菜谱", adminApi.UpdateRecipeHandler),
		utils.NewRegisterApiParam(common.HttpDelete, "/api/v1/admin/recipes/:id", "删除菜谱", adminApi.DeleteRecipeHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/api/v1/admin/ingredients", "获取全部食材", adminApi.ListAllIngredientsHandler),
		utils.NewRegisterApiParam(common.HttpPost, "/api/v1/admin/ingredients", "创建食材", adminApi.CreateIngredientHandler),
		utils.NewRegisterApiParam(common.HttpPut, "/api/v1/admin/ingredients/:id", "更新食材", adminApi.UpdateIngredientHandler),
		utils.NewRegisterApiParam(common.HttpDelete, "/api/v1/admin/ingredients/:id", "删除食材", adminApi.DeleteIngredientHandler),
	)
	return adminRouter
}
