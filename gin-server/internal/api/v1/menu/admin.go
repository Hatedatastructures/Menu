package menu

import (

	"github.com/gin-gonic/gin"
	req "shack/internal/model/menu/request"
	_ "shack/internal/model/menu/response"
	"shack/internal/vo"
	"shack/internal/utils/validator"

)

type AdminApi struct{}

// GetAdminStatsHandler
// @Tags menuadminApi
// @Summary GetAdminStatsHandler 获取管理后台统计
// @Description GetAdminStatsHandler 获取管理后台统计
// @Security ApiKeyAuth
// @Success 200 {object} vo.Result{data=_.GetAdminStatsRes}
// @Router /api/v1/admin/stats [GET]
func (s *AdminApi) GetAdminStatsHandler(c *gin.Context) {
	data, err := adminService.GetAdminStats(c)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// ListAllRecipesHandler
// @Tags menuadminApi
// @Summary ListAllRecipesHandler 获取全部菜谱(含所有状态)
// @Description ListAllRecipesHandler 获取全部菜谱(含所有状态)
// @Security ApiKeyAuth
// @Success 200 {object} vo.Result{}
// @Router /api/v1/admin/recipes [GET]
func (s *AdminApi) ListAllRecipesHandler(c *gin.Context) {
	err := adminService.ListAllRecipes(c)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, _ := c.Get("recipeList")
	c.JSON(200, vo.Success(c, data))
}

// CreateRecipeHandler
// @Tags menuadminApi
// @Summary CreateRecipeHandler 创建菜谱
// @Description CreateRecipeHandler 创建菜谱
// @Security ApiKeyAuth
// @Param data body req.CreateRecipeReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.CreateRecipeRes}
// @Router /api/v1/admin/recipes [POST]
func (s *AdminApi) CreateRecipeHandler(c *gin.Context) {
	var req req.CreateRecipeReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := adminService.CreateRecipe(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// UpdateRecipeHandler
// @Tags menuadminApi
// @Summary UpdateRecipeHandler 更新菜谱
// @Description UpdateRecipeHandler 更新菜谱
// @Security ApiKeyAuth
// @Param data body req.UpdateRecipeReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.UpdateRecipeRes}
// @Router /api/v1/admin/recipes/:id [PUT]
func (s *AdminApi) UpdateRecipeHandler(c *gin.Context) {
	var req req.UpdateRecipeReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := adminService.UpdateRecipe(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// DeleteRecipeHandler
// @Tags menuadminApi
// @Summary DeleteRecipeHandler 删除菜谱
// @Description DeleteRecipeHandler 删除菜谱
// @Security ApiKeyAuth
// @Success 200 {object} vo.Result{}
// @Router /api/v1/admin/recipes/:id [DELETE]
func (s *AdminApi) DeleteRecipeHandler(c *gin.Context) {
	err := adminService.DeleteRecipe(c)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, "操作成功"))
}

// ListAllIngredientsHandler
// @Tags menuadminApi
// @Summary ListAllIngredientsHandler 获取全部食材
// @Description ListAllIngredientsHandler 获取全部食材
// @Security ApiKeyAuth
// @Success 200 {object} vo.Result{}
// @Router /api/v1/admin/ingredients [GET]
func (s *AdminApi) ListAllIngredientsHandler(c *gin.Context) {
	err := adminService.ListAllIngredients(c)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, _ := c.Get("recipeList")
	c.JSON(200, vo.Success(c, data))
}

// CreateIngredientHandler
// @Tags menuadminApi
// @Summary CreateIngredientHandler 创建食材
// @Description CreateIngredientHandler 创建食材
// @Security ApiKeyAuth
// @Param data body req.CreateIngredientReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.CreateIngredientRes}
// @Router /api/v1/admin/ingredients [POST]
func (s *AdminApi) CreateIngredientHandler(c *gin.Context) {
	var req req.CreateIngredientReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := adminService.CreateIngredient(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// UpdateIngredientHandler
// @Tags menuadminApi
// @Summary UpdateIngredientHandler 更新食材
// @Description UpdateIngredientHandler 更新食材
// @Security ApiKeyAuth
// @Param data body req.UpdateIngredientReq true "请求参数"
// @Success 200 {object} vo.Result{data=_.UpdateIngredientRes}
// @Router /api/v1/admin/ingredients/:id [PUT]
func (s *AdminApi) UpdateIngredientHandler(c *gin.Context) {
	var req req.UpdateIngredientReq
	// body 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	data, err := adminService.UpdateIngredient(c, req)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, data))
}

// DeleteIngredientHandler
// @Tags menuadminApi
// @Summary DeleteIngredientHandler 删除食材
// @Description DeleteIngredientHandler 删除食材
// @Security ApiKeyAuth
// @Success 200 {object} vo.Result{}
// @Router /api/v1/admin/ingredients/:id [DELETE]
func (s *AdminApi) DeleteIngredientHandler(c *gin.Context) {
	err := adminService.DeleteIngredient(c)
	if err != nil {
		c.JSON(201, vo.Fail(c, "", err))
		return
	}
	c.JSON(200, vo.Success(c, "操作成功"))
}