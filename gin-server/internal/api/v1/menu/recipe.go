package menu

import (

	"github.com/gin-gonic/gin"
	"strconv"
	req "shack/internal/model/menu/request"
	_ "shack/internal/model/menu/response"
	"shack/internal/vo"
	"shack/internal/utils/validator"

)

type RecipeApi struct{}

// ListPublishedRecipesHandler
// @Tags menurecipeApi
// @Summary ListPublishedRecipesHandler 菜谱公开接口模块 对应 C++ server RecipeRoutes (recipes/ingredients/recommendations) 获取已发布菜谱列表(支持筛选)
// @Description ListPublishedRecipesHandler 菜谱公开接口模块 对应 C++ server RecipeRoutes (recipes/ingredients/recommendations) 获取已发布菜谱列表(支持筛选)
// @Param data body req.ListPublishedRecipesReq true "请求参数"
// @Success 200 {object} vo.Result{}
// @Router /api/v1/recipes [GET]
func (s *RecipeApi) ListPublishedRecipesHandler(c *gin.Context) {
	var req req.ListPublishedRecipesReq
	// query 参数
	{
		val := c.Query("cuisine")
		if val != "" {
			req.Cuisine = val

		}
	}
	{
		val := c.Query("maxMinutes")
		if val != "" {
			parsed, err := strconv.Atoi(val)
			if err != nil {
				c.JSON(201, vo.Fail(c, "", err))
				return
			}
			req.MaxMinutes = parsed

		}
	}
	{
		val := c.Query("difficulty")
		if val != "" {
			parsed, err := strconv.Atoi(val)
			if err != nil {
				c.JSON(201, vo.Fail(c, "", err))
				return
			}
			req.Difficulty = parsed

		}
	}
	{
		val := c.Query("limit")
		if val != "" {
			parsed, err := strconv.Atoi(val)
			if err != nil {
				c.JSON(201, vo.Fail(c, "", err))
				return
			}
			req.Limit = parsed

		}
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	err := recipeService.ListPublishedRecipes(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, "操作成功"))
}

// ListIngredientsHandler
// @Tags menurecipeApi
// @Summary ListIngredientsHandler 获取全部食材列表
// @Description ListIngredientsHandler 获取全部食材列表
// @Success 200 {object} vo.Result{}
// @Router /api/v1/ingredients [GET]
func (s *RecipeApi) ListIngredientsHandler(c *gin.Context) {
	err := recipeService.ListIngredients(c)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, "操作成功"))
}

// GetPublishedRecipeHandler
// @Tags menurecipeApi
// @Summary GetPublishedRecipeHandler 根据ID获取单个已发布菜谱
// @Description GetPublishedRecipeHandler 根据ID获取单个已发布菜谱
// @Success 200 {object} vo.Result{data=_.GetPublishedRecipeRes}
// @Router /api/v1/recipes/:id [GET]
func (s *RecipeApi) GetPublishedRecipeHandler(c *gin.Context) {
	data, err := recipeService.GetPublishedRecipe(c)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// RecommendTonightHandler
// @Tags menurecipeApi
// @Summary RecommendTonightHandler 今晚推荐菜谱
// @Description RecommendTonightHandler 今晚推荐菜谱
// @Param data body req.RecommendTonightReq true "请求参数"
// @Success 200 {object} vo.Result{}
// @Router /api/v1/recommendations/tonight [POST]
func (s *RecipeApi) RecommendTonightHandler(c *gin.Context) {
	var req req.RecommendTonightReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	err := recipeService.RecommendTonight(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, "操作成功"))
}