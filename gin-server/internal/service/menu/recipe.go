package menu

import (
	"encoding/json"
	"math"

	"github.com/gin-gonic/gin"

	biz_err "shack/internal/error"
	"shack/internal/global"
	menuModel "shack/internal/model/menu"
	req "shack/internal/model/menu/request"
	res "shack/internal/model/menu/response"
)

type RecipeService struct{}

type contextKey string

const recipeListKey contextKey = "recipeList"
const recipeRecommendKey contextKey = "recipeRecommend"

// ListPublishedRecipes 获取已发布菜谱列表(支持筛选)
func (s *RecipeService) ListPublishedRecipes(ctx *gin.Context, r req.ListPublishedRecipesReq) error {
	db := global.GVA_DB.Model(&menuModel.MenuRecipe{}).Where("Status = ?", "published")

	if r.Cuisine != "" {
		db = db.Where("Cuisine = ?", r.Cuisine)
	}
	if r.MaxMinutes > 0 {
		db = db.Where("PrepMinutes + CookMinutes <= ?", r.MaxMinutes)
	}
	if r.Difficulty > 0 {
		db = db.Where("Difficulty = ?", r.Difficulty)
	}

	limit := 20
	if r.Limit > 0 {
		limit = int(math.Min(float64(r.Limit), 100))
	}

	var recipes []menuModel.MenuRecipe
	if err := db.Order("CreatedAt DESC").Limit(limit).Find(&recipes).Error; err != nil {
		return biz_err.New(biz_err.DB_ERROR)
	}

	// 填充关联数据
	s.fillRecipeDetails(recipes)

	// 转换为响应格式
	list := s.toRecipeResList(recipes)
	ctx.Set(string(recipeListKey), list)
	return nil
}

// ListIngredients 获取全部食材列表
func (s *RecipeService) ListIngredients(ctx *gin.Context) error {
	var ingredients []menuModel.MenuIngredient
	if err := global.GVA_DB.Preload("Aliases").Order("Name ASC").Find(&ingredients).Error; err != nil {
		return biz_err.New(biz_err.DB_ERROR)
	}

	type ingredientItem struct {
		Id              string   `json:"id"`
		Name            string   `json:"name"`
		Aliases         []string `json:"aliases"`
		Category        string   `json:"category"`
		DefaultUnit     string   `json:"defaultUnit"`
		IsPantryStaple  bool     `json:"isPantryStaple"`
		SubstituteGroup string   `json:"substituteGroup"`
		StoreSkuMapping string   `json:"storeSkuMapping"`
	}

	result := make([]ingredientItem, 0, len(ingredients))
	for _, ing := range ingredients {
		aliases := make([]string, 0, len(ing.Aliases))
		for _, a := range ing.Aliases {
			aliases = append(aliases, a.Alias)
		}
		result = append(result, ingredientItem{
			Id:              ing.Id,
			Name:            ing.Name,
			Aliases:         aliases,
			Category:        ing.Category,
			DefaultUnit:     ing.DefaultUnit,
			IsPantryStaple:  ing.IsPantryStaple,
			SubstituteGroup: ing.SubstituteGroup,
			StoreSkuMapping: ing.StoreSkuMapping,
		})
	}
	ctx.Set(string(recipeListKey), result)
	return nil
}

// GetPublishedRecipe 根据ID获取单个已发布菜谱
func (s *RecipeService) GetPublishedRecipe(ctx *gin.Context) (rs res.GetPublishedRecipeRes, err error) {
	recipeId := ctx.Param("id")
	if recipeId == "" {
		return rs, biz_err.New(biz_err.PARAM_MISSING)
	}

	var recipe menuModel.MenuRecipe
	if err := global.GVA_DB.Where("Id = ? AND Status = ?", recipeId, "published").First(&recipe).Error; err != nil {
		return rs, biz_err.New(biz_err.RECIPE_NOT_FOUND)
	}

	s.fillRecipeSingle(&recipe)
	return s.toRecipeRes(recipe), nil
}

// RecommendTonight 今晚推荐菜谱
func (s *RecipeService) RecommendTonight(ctx *gin.Context, r req.RecommendTonightReq) error {
	var recipes []menuModel.MenuRecipe
	if err := global.GVA_DB.Where("Status = ?", "published").Find(&recipes).Error; err != nil {
		return biz_err.New(biz_err.DB_ERROR)
	}

	s.fillRecipeDetails(recipes)

	type recommendItem struct {
		Recipe                res.GetPublishedRecipeRes `json:"recipe"`
		AvailableIngredientCount int                     `json:"availableIngredientCount"`
		MissingIngredientCount   int                     `json:"missingIngredientCount"`
		Score                    int                     `json:"score"`
	}

	// 将 pantryIngredientIds 转为 set
	pantrySet := make(map[string]bool)
	for _, id := range r.PantryIngredientIds {
		pantrySet[id] = true
	}

	results := make([]recommendItem, 0)
	for _, recipe := range recipes {
		// 时间筛选
		totalMinutes := recipe.PrepMinutes + recipe.CookMinutes
		if r.AvailableMinutes > 0 && totalMinutes > r.AvailableMinutes {
			continue
		}
		// 菜系筛选
		if len(r.Cuisines) > 0 {
			matched := false
			for _, c := range r.Cuisines {
				if c == recipe.Cuisine {
					matched = true
					break
				}
			}
			if !matched {
				continue
			}
		}

		available := 0
		missing := 0
		for _, ing := range recipe.Ingredients {
			if pantrySet[ing.IngredientId] {
				available++
			} else {
				missing++
			}
		}
		score := available * 10
		if missing == 0 {
			score += 50
		}

		results = append(results, recommendItem{
			Recipe:                  s.toRecipeRes(recipe),
			AvailableIngredientCount: available,
			MissingIngredientCount:   missing,
			Score:                    score,
		})
	}

	// 按分数排序(简单冒泡)
	for i := 0; i < len(results); i++ {
		for j := i + 1; j < len(results); j++ {
			if results[j].Score > results[i].Score {
				results[i], results[j] = results[j], results[i]
			}
		}
	}

	// 限制返回数量
	limit := 10
	if r.Servings > 0 && r.Servings < limit {
		limit = r.Servings
	}
	if len(results) > limit {
		results = results[:limit]
	}

	ctx.Set(string(recipeRecommendKey), results)
	return nil
}

// ==================== 内部辅助方法 ====================

// fillRecipeDetails 批量填充菜谱的食材和步骤详情
func (s *RecipeService) fillRecipeDetails(recipes []menuModel.MenuRecipe) {
	if len(recipes) == 0 {
		return
	}

	ids := make([]string, len(recipes))
	for i, r := range recipes {
		ids[i] = r.Id
	}

	// 查询所有菜谱食材(含食材名)
	var recipeIngs []struct {
		RecipeId         string
		IngredientId     string
		Quantity         float64
		Unit             string
		ServingFactor    float64
		Preparation      string
		Required         bool
		IngredientName   string
		Category         string
		DefaultUnit      string
		IsPantryStaple   bool
	}
	global.GVA_DB.Table("RecipeIngredients").
		Select("RecipeIngredients.*, Ingredients.Name as ingredient_name, Ingredients.Category as category, Ingredients.DefaultUnit as default_unit, Ingredients.IsPantryStaple as is_pantry_staple").
		Joins("LEFT JOIN Ingredients ON RecipeIngredients.IngredientId = Ingredients.Id").
		Where("RecipeIngredients.RecipeId IN ?", ids).
		Scan(&recipeIngs)

	// 查询所有步骤
	var steps []menuModel.MenuRecipeStep
	global.GVA_DB.Where("RecipeId IN ?", ids).Order("StepOrder ASC").Find(&steps)

	// 查询所有别名
	var allAliases []menuModel.MenuIngredientAlias
	ingIdSet := make(map[string]bool)
	for _, ri := range recipeIngs {
		ingIdSet[ri.IngredientId] = true
	}
	ingIds := make([]string, 0, len(ingIdSet))
	for id := range ingIdSet {
		ingIds = append(ingIds, id)
	}
	if len(ingIds) > 0 {
		global.GVA_DB.Where("IngredientId IN ?", ingIds).Find(&allAliases)
	}
	_ = allAliases // 别名在 recipe list 中不需要

	// 组装
	ingMap := make(map[string][]menuModel.MenuRecipeIngredient)
	for _, ri := range recipeIngs {
		ingMap[ri.RecipeId] = append(ingMap[ri.RecipeId], menuModel.MenuRecipeIngredient{
			RecipeId:      ri.RecipeId,
			IngredientId:  ri.IngredientId,
			Quantity:      ri.Quantity,
			Unit:          ri.Unit,
			ServingFactor: ri.ServingFactor,
			Preparation:   ri.Preparation,
			Required:      ri.Required,
			IngredientName: ri.IngredientName,
			Category:       ri.Category,
			DefaultUnit:    ri.DefaultUnit,
			IsPantryStaple: ri.IsPantryStaple,
		})
	}
	stepMap := make(map[string][]menuModel.MenuRecipeStep)
	for _, step := range steps {
		stepMap[step.RecipeId] = append(stepMap[step.RecipeId], step)
	}

	for i := range recipes {
		recipes[i].Ingredients = ingMap[recipes[i].Id]
		recipes[i].Steps = stepMap[recipes[i].Id]
	}
}

// fillRecipeSingle 填充单个菜谱详情
func (s *RecipeService) fillRecipeSingle(recipe *menuModel.MenuRecipe) {
	s.fillRecipeDetails([]menuModel.MenuRecipe{*recipe})
}

// toRecipeRes 将 GORM 模型转换为响应 DTO
func (s *RecipeService) toRecipeRes(recipe menuModel.MenuRecipe) res.GetPublishedRecipeRes {
	ingredients := make([]res.GetPublishedRecipeResIngredient, len(recipe.Ingredients))
	for i, ing := range recipe.Ingredients {
		ingredients[i] = res.GetPublishedRecipeResIngredient{
			IngredientId:   ing.IngredientId,
			Quantity:       ing.Quantity,
			Unit:           ing.Unit,
			Required:       ing.Required,
			ServingFactor:  ing.ServingFactor,
			Preparation:    ing.Preparation,
			IngredientName: ing.IngredientName,
			Category:       ing.Category,
			DefaultUnit:    ing.DefaultUnit,
			IsPantryStaple: ing.IsPantryStaple,
		}
	}

	steps := make([]res.GetPublishedRecipeResStep, len(recipe.Steps))
	for i, step := range recipe.Steps {
		steps[i] = res.GetPublishedRecipeResStep{
			StepOrder:       step.StepOrder,
			Title:           step.Title,
			Instruction:     step.Instruction,
			DurationSeconds: step.DurationSeconds,
			HasTimer:        step.HasTimer,
		}
	}

	var allergens, cookware []string
	_ = json.Unmarshal([]byte(recipe.AllergensJson), &allergens)
	_ = json.Unmarshal([]byte(recipe.CookwareJson), &cookware)

	return res.GetPublishedRecipeRes{
		Id:            recipe.Id,
		Slug:          recipe.Slug,
		Name:          recipe.Name,
		Cuisine:       recipe.Cuisine,
		Description:   recipe.Description,
		PrepMinutes:   recipe.PrepMinutes,
		CookMinutes:   recipe.CookMinutes,
		TotalMinutes:  recipe.PrepMinutes + recipe.CookMinutes,
		Servings:      recipe.Servings,
		Difficulty:    recipe.Difficulty,
		ImagePath:     recipe.ImagePath,
		Status:        recipe.Status,
		Allergens:     allergens,
		Cookware:      cookware,
		Ingredients:   ingredients,
		Steps:         steps,
	}
}

// toRecipeResList 批量转换
func (s *RecipeService) toRecipeResList(recipes []menuModel.MenuRecipe) []res.GetPublishedRecipeRes {
	list := make([]res.GetPublishedRecipeRes, len(recipes))
	for i, r := range recipes {
		list[i] = s.toRecipeRes(r)
	}
	return list
}
