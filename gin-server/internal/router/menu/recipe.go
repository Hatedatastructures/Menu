//router 解析
package menu

import (
	"github.com/gin-gonic/gin"
	"shack/internal/model/common"
	"shack/internal/utils"
)

type RecipeRouter struct{}

func (s *RecipeRouter) InitRecipeRouter(Router *gin.RouterGroup, apiSet map[string]struct{}) (R gin.IRoutes) {
	recipeRouter := Router.Group("")
	utils.RegisterApi(recipeRouter, apiSet, "",
		utils.NewRegisterApiParam(common.HttpGet, "/api/v1/recipes", "菜谱公开接口模块 对应 C++ server RecipeRoutes (recipes/ingredients/recommendations) 获取已发布菜谱列表(支持筛选)", recipeApi.ListPublishedRecipesHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/api/v1/ingredients", "获取全部食材列表", recipeApi.ListIngredientsHandler),
		utils.NewRegisterApiParam(common.HttpGet, "/api/v1/recipes/:id", "根据ID获取单个已发布菜谱", recipeApi.GetPublishedRecipeHandler),
		utils.NewRegisterApiParam(common.HttpPost, "/api/v1/recommendations/tonight", "今晚推荐菜谱", recipeApi.RecommendTonightHandler),
	)
	return recipeRouter
}
