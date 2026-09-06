package menu

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	biz_err "shack/internal/error"
	"shack/internal/global"
	menuModel "shack/internal/model/menu"
	req "shack/internal/model/menu/request"
	res "shack/internal/model/menu/response"
)

type AdminService struct{}

// GetAdminStats 获取管理后台统计
func (s *AdminService) GetAdminStats(ctx *gin.Context) (rs res.GetAdminStatsRes, err error) {
	var published, draft, archived, ingredientCount int64

	global.GVA_DB.Model(&menuModel.MenuRecipe{}).Where("Status = ?", "published").Count(&published)
	global.GVA_DB.Model(&menuModel.MenuRecipe{}).Where("Status = ?", "draft").Count(&draft)
	global.GVA_DB.Model(&menuModel.MenuRecipe{}).Where("Status = ?", "archived").Count(&archived)
	global.GVA_DB.Model(&menuModel.MenuIngredient{}).Count(&ingredientCount)

	return res.GetAdminStatsRes{
		PublishedRecipes: int(published),
		DraftRecipes:     int(draft),
		ArchivedRecipes:  int(archived),
		IngredientCount:  int(ingredientCount),
	}, nil
}

// ListAllRecipes 获取全部菜谱(含所有状态)
func (s *AdminService) ListAllRecipes(ctx *gin.Context) error {
	var recipes []menuModel.MenuRecipe
	if err := global.GVA_DB.Order("CreatedAt DESC").Find(&recipes).Error; err != nil {
		return biz_err.New(biz_err.DB_ERROR)
	}

	recipeSvc := &RecipeService{}
	recipeSvc.fillRecipeDetails(recipes)
	list := recipeSvc.toRecipeResList(recipes)
	ctx.Set(string(recipeListKey), list)
	return nil
}

// CreateRecipe 创建菜谱
func (s *AdminService) CreateRecipe(ctx *gin.Context, r req.CreateRecipeReq) (rs res.CreateRecipeRes, err error) {
	recipeId := uuid.New().String()

	allergensJson, _ := json.Marshal(r.Allergens)
	cookwareJson, _ := json.Marshal(r.Cookware)

	recipe := menuModel.MenuRecipe{
		Id:            recipeId,
		Slug:          r.Slug,
		Name:          r.Name,
		Cuisine:       r.Cuisine,
		Description:   r.Description,
		PrepMinutes:   r.PrepMinutes,
		CookMinutes:   r.CookMinutes,
		Servings:      r.Servings,
		Difficulty:    r.Difficulty,
		ImagePath:     r.ImagePath,
		Status:        r.Status,
		AllergensJson: string(allergensJson),
		CookwareJson:  string(cookwareJson),
	}

	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&recipe).Error; err != nil {
			return err
		}

		// 创建食材关联
		for _, ing := range r.Ingredients {
			ri := menuModel.MenuRecipeIngredient{
				RecipeId:      recipeId,
				IngredientId:  ing.IngredientId,
				Quantity:      ing.Quantity,
				Unit:          ing.Unit,
				ServingFactor: ing.ServingFactor,
				Preparation:   ing.Preparation,
				Required:      ing.Required,
			}
			if err := tx.Create(&ri).Error; err != nil {
				return err
			}
		}

		// 创建步骤
		for _, step := range r.Steps {
			s := menuModel.MenuRecipeStep{
				RecipeId:        recipeId,
				StepOrder:       step.StepOrder,
				Title:           step.Title,
				Instruction:     step.Instruction,
				DurationSeconds: step.DurationSeconds,
				HasTimer:        step.HasTimer,
			}
			if err := tx.Create(&s).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return rs, biz_err.New(biz_err.RECIPE_CREATE_FAILED)
	}

	return s.toCreateRecipeRes(recipe, r)
}

// UpdateRecipe 更新菜谱
func (s *AdminService) UpdateRecipe(ctx *gin.Context, r req.UpdateRecipeReq) (rs res.UpdateRecipeRes, err error) {
	recipeId := ctx.Param("id")
	if recipeId == "" {
		return rs, biz_err.New(biz_err.PARAM_MISSING)
	}

	var existing menuModel.MenuRecipe
	if err := global.GVA_DB.Where("Id = ?", recipeId).First(&existing).Error; err != nil {
		return rs, biz_err.New(biz_err.RECIPE_NOT_FOUND)
	}

	// 更新字段
	allergensJson, _ := json.Marshal(r.Allergens)
	cookwareJson, _ := json.Marshal(r.Cookware)

	updates := map[string]interface{}{
		"Slug":          r.Slug,
		"Name":          r.Name,
		"Cuisine":       r.Cuisine,
		"Description":   r.Description,
		"PrepMinutes":   r.PrepMinutes,
		"CookMinutes":   r.CookMinutes,
		"Servings":      r.Servings,
		"Difficulty":    r.Difficulty,
		"ImagePath":     r.ImagePath,
		"Status":        r.Status,
		"AllergensJson": string(allergensJson),
		"CookwareJson":  string(cookwareJson),
	}

	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&menuModel.MenuRecipe{}).Where("Id = ?", recipeId).Updates(updates).Error; err != nil {
			return err
		}

		// 删除旧食材关联,重新创建
		tx.Where("RecipeId = ?", recipeId).Delete(&menuModel.MenuRecipeIngredient{})
		for _, ing := range r.Ingredients {
			ri := menuModel.MenuRecipeIngredient{
				RecipeId:      recipeId,
				IngredientId:  ing.IngredientId,
				Quantity:      ing.Quantity,
				Unit:          ing.Unit,
				ServingFactor: ing.ServingFactor,
				Preparation:   ing.Preparation,
				Required:      ing.Required,
			}
			if err := tx.Create(&ri).Error; err != nil {
				return err
			}
		}

		// 删除旧步骤,重新创建
		tx.Where("RecipeId = ?", recipeId).Delete(&menuModel.MenuRecipeStep{})
		for _, step := range r.Steps {
			st := menuModel.MenuRecipeStep{
				RecipeId:        recipeId,
				StepOrder:       step.StepOrder,
				Title:           step.Title,
				Instruction:     step.Instruction,
				DurationSeconds: step.DurationSeconds,
				HasTimer:        step.HasTimer,
			}
			if err := tx.Create(&st).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return rs, biz_err.New(biz_err.RECIPE_UPDATE_FAILED)
	}

	// 重新查询完整数据
	recipeSvc := &RecipeService{}
	var updated menuModel.MenuRecipe
	global.GVA_DB.Where("Id = ?", recipeId).First(&updated)
	recipeSvc.fillRecipeSingle(&updated)

	return s.toUpdateRecipeRes(updated), nil
}

// DeleteRecipe 删除菜谱
func (s *AdminService) DeleteRecipe(ctx *gin.Context) error {
	recipeId := ctx.Param("id")
	if recipeId == "" {
		return biz_err.New(biz_err.PARAM_MISSING)
	}

	result := global.GVA_DB.Where("Id = ?", recipeId).Delete(&menuModel.MenuRecipe{})
	if result.RowsAffected == 0 {
		return biz_err.New(biz_err.RECIPE_NOT_FOUND)
	}
	return nil
}

// ListAllIngredients 获取全部食材
func (s *AdminService) ListAllIngredients(ctx *gin.Context) error {
	var ingredients []menuModel.MenuIngredient
	if err := global.GVA_DB.Preload("Aliases").Order("Name ASC").Find(&ingredients).Error; err != nil {
		return biz_err.New(biz_err.DB_ERROR)
	}

	type ingredientRes struct {
		Id              string   `json:"id"`
		Name            string   `json:"name"`
		Aliases         []string `json:"aliases"`
		Category        string   `json:"category"`
		DefaultUnit     string   `json:"defaultUnit"`
		IsPantryStaple  bool     `json:"isPantryStaple"`
		SubstituteGroup string   `json:"substituteGroup"`
		StoreSkuMapping string   `json:"storeSkuMapping"`
	}

	list := make([]ingredientRes, 0, len(ingredients))
	for _, ing := range ingredients {
		aliases := make([]string, 0, len(ing.Aliases))
		for _, a := range ing.Aliases {
			aliases = append(aliases, a.Alias)
		}
		list = append(list, ingredientRes{
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
	ctx.Set(string(recipeListKey), list)
	return nil
}

// CreateIngredient 创建食材
func (s *AdminService) CreateIngredient(ctx *gin.Context, r req.CreateIngredientReq) (rs res.CreateIngredientRes, err error) {
	ingId := "ing." + uuid.New().String()[:12]

	ingredient := menuModel.MenuIngredient{
		Id:              ingId,
		Name:            r.Name,
		Category:        r.Category,
		DefaultUnit:     r.DefaultUnit,
		IsPantryStaple:  r.IsPantryStaple,
		SubstituteGroup: r.SubstituteGroup,
		StoreSkuMapping: r.StoreSkuMapping,
	}

	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&ingredient).Error; err != nil {
			return err
		}
		// 创建别名
		for _, alias := range r.Aliases {
			a := menuModel.MenuIngredientAlias{
				IngredientId: ingId,
				Alias:        alias,
			}
			if err := tx.Create(&a).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return rs, biz_err.New(biz_err.INGREDIENT_CREATE_FAILED)
	}

	return res.CreateIngredientRes{
		Id:              ingredient.Id,
		Name:            ingredient.Name,
		Aliases:         r.Aliases,
		Category:        ingredient.Category,
		DefaultUnit:     ingredient.DefaultUnit,
		IsPantryStaple:  ingredient.IsPantryStaple,
		SubstituteGroup: ingredient.SubstituteGroup,
		StoreSkuMapping: ingredient.StoreSkuMapping,
	}, nil
}

// UpdateIngredient 更新食材
func (s *AdminService) UpdateIngredient(ctx *gin.Context, r req.UpdateIngredientReq) (rs res.UpdateIngredientRes, err error) {
	ingId := ctx.Param("id")
	if ingId == "" {
		return rs, biz_err.New(biz_err.PARAM_MISSING)
	}

	var existing menuModel.MenuIngredient
	if err := global.GVA_DB.Where("Id = ?", ingId).First(&existing).Error; err != nil {
		return rs, biz_err.New(biz_err.INGREDIENT_NOT_FOUND)
	}

	updates := map[string]interface{}{
		"Name":            r.Name,
		"Category":        r.Category,
		"DefaultUnit":     r.DefaultUnit,
		"IsPantryStaple":  r.IsPantryStaple,
		"SubstituteGroup": r.SubstituteGroup,
		"StoreSkuMapping": r.StoreSkuMapping,
	}

	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&menuModel.MenuIngredient{}).Where("Id = ?", ingId).Updates(updates).Error; err != nil {
			return err
		}
		// 更新别名:删除旧的,创建新的
		tx.Where("IngredientId = ?", ingId).Delete(&menuModel.MenuIngredientAlias{})
		for _, alias := range r.Aliases {
			a := menuModel.MenuIngredientAlias{
				IngredientId: ingId,
				Alias:        alias,
			}
			if err := tx.Create(&a).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return rs, biz_err.New(biz_err.INGREDIENT_UPDATE_FAILED)
	}

	return res.UpdateIngredientRes{
		Id:              ingId,
		Name:            r.Name,
		Aliases:         r.Aliases,
		Category:        r.Category,
		DefaultUnit:     r.DefaultUnit,
		IsPantryStaple:  r.IsPantryStaple,
		SubstituteGroup: r.SubstituteGroup,
		StoreSkuMapping: r.StoreSkuMapping,
	}, nil
}

// DeleteIngredient 删除食材
func (s *AdminService) DeleteIngredient(ctx *gin.Context) error {
	ingId := ctx.Param("id")
	if ingId == "" {
		return biz_err.New(biz_err.PARAM_MISSING)
	}

	result := global.GVA_DB.Where("Id = ?", ingId).Delete(&menuModel.MenuIngredient{})
	if result.RowsAffected == 0 {
		return biz_err.New(biz_err.INGREDIENT_NOT_FOUND)
	}
	return nil
}

// ==================== 内部辅助 ====================

func (s *AdminService) toCreateRecipeRes(recipe menuModel.MenuRecipe, r req.CreateRecipeReq) (res.CreateRecipeRes, error) {
	ingredients := make([]res.CreateRecipeResIngredient, len(r.Ingredients))
	for i, ing := range r.Ingredients {
		ingredients[i] = res.CreateRecipeResIngredient{
			IngredientId:  ing.IngredientId,
			Quantity:      ing.Quantity,
			Unit:          ing.Unit,
			Required:      ing.Required,
			ServingFactor: ing.ServingFactor,
			Preparation:   ing.Preparation,
		}
	}
	steps := make([]res.CreateRecipeResStep, len(r.Steps))
	for i, step := range r.Steps {
		steps[i] = res.CreateRecipeResStep{
			StepOrder:       step.StepOrder,
			Title:           step.Title,
			Instruction:     step.Instruction,
			DurationSeconds: step.DurationSeconds,
			HasTimer:        step.HasTimer,
		}
	}
	return res.CreateRecipeRes{
		Id:           recipe.Id,
		Slug:         recipe.Slug,
		Name:         recipe.Name,
		Cuisine:      recipe.Cuisine,
		Description:  recipe.Description,
		PrepMinutes:  recipe.PrepMinutes,
		CookMinutes:  recipe.CookMinutes,
		TotalMinutes: recipe.PrepMinutes + recipe.CookMinutes,
		Servings:     recipe.Servings,
		Difficulty:   recipe.Difficulty,
		ImagePath:    recipe.ImagePath,
		Status:       recipe.Status,
		Allergens:    r.Allergens,
		Cookware:     r.Cookware,
		Ingredients:  ingredients,
		Steps:        steps,
	}, nil
}

func (s *AdminService) toUpdateRecipeRes(recipe menuModel.MenuRecipe) res.UpdateRecipeRes {
	var allergens, cookware []string
	_ = json.Unmarshal([]byte(recipe.AllergensJson), &allergens)
	_ = json.Unmarshal([]byte(recipe.CookwareJson), &cookware)

	ingredients := make([]res.UpdateRecipeResIngredient, len(recipe.Ingredients))
	for i, ing := range recipe.Ingredients {
		ingredients[i] = res.UpdateRecipeResIngredient{
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
	steps := make([]res.UpdateRecipeResStep, len(recipe.Steps))
	for i, step := range recipe.Steps {
		steps[i] = res.UpdateRecipeResStep{
			StepOrder:       step.StepOrder,
			Title:           step.Title,
			Instruction:     step.Instruction,
			DurationSeconds: step.DurationSeconds,
			HasTimer:        step.HasTimer,
		}
	}
	return res.UpdateRecipeRes{
		Id:           recipe.Id,
		Slug:         recipe.Slug,
		Name:         recipe.Name,
		Cuisine:      recipe.Cuisine,
		Description:  recipe.Description,
		PrepMinutes:  recipe.PrepMinutes,
		CookMinutes:  recipe.CookMinutes,
		TotalMinutes: recipe.PrepMinutes + recipe.CookMinutes,
		Servings:     recipe.Servings,
		Difficulty:   recipe.Difficulty,
		ImagePath:    recipe.ImagePath,
		Status:       recipe.Status,
		Allergens:    allergens,
		Cookware:     cookware,
		Ingredients:  ingredients,
		Steps:        steps,
	}
}
