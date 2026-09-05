package menu

import (
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	biz_err "shack/internal/error"
	"shack/internal/global"
	menuModel "shack/internal/model/menu"
	req "shack/internal/model/menu/request"
	res "shack/internal/model/menu/response"
)

type WorkflowService struct{}

type wfContextKey string

const planListKey wfContextKey = "planList"

// ListPlans 获取本周餐食计划
func (s *WorkflowService) ListPlans(ctx *gin.Context, r req.ListPlansReq) error {
	userId := ctx.GetString("userId")
	if userId == "" {
		return biz_err.New(biz_err.UNAUTHORIZED)
	}

	db := global.GVA_DB.Model(&menuModel.MenuMealPlan{}).Where("UserId = ?", userId)

	if r.From != "" {
		db = db.Where("PlanDate >= ?", r.From)
	}
	if r.To != "" {
		db = db.Where("PlanDate <= ?", r.To)
	}

	var plans []menuModel.MenuMealPlan
	if err := db.Order("PlanDate ASC").Find(&plans).Error; err != nil {
		return biz_err.New(biz_err.DB_ERROR)
	}

	// 填充每个计划的 items
	for i := range plans {
		s.fillPlanItems(&plans[i])
	}

	list := s.toPlanResList(plans)
	ctx.Set(string(planListKey), list)
	return nil
}

// SavePlan 保存餐食计划
func (s *WorkflowService) SavePlan(ctx *gin.Context, r req.SavePlanReq) (rs res.SavePlanRes, err error) {
	userId := ctx.GetString("userId")
	if userId == "" {
		return rs, biz_err.New(biz_err.UNAUTHORIZED)
	}

	var plan menuModel.MenuMealPlan

	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		// 查找已有计划
		if err := tx.Where("UserId = ? AND PlanDate = ?", userId, r.PlanDate).First(&plan).Error; err != nil {
			// 不存在则创建
			plan = menuModel.MenuMealPlan{
				Id:       uuid.New().String(),
				UserId:   userId,
				PlanDate: r.PlanDate,
			}
			if err := tx.Create(&plan).Error; err != nil {
				return err
			}
		}

		// 删除旧的计划项
		tx.Where("MealPlanId = ?", plan.Id).Delete(&menuModel.MenuMealPlanItem{})

		// 创建新的计划项
		for _, item := range r.Items {
			mpi := menuModel.MenuMealPlanItem{
				Id:         uuid.New().String(),
				MealPlanId: plan.Id,
				RecipeId:   item.RecipeId,
				Servings:   item.Servings,
				SortOrder:  item.SortOrder,
			}
			if err := tx.Create(&mpi).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return rs, biz_err.New(biz_err.PLAN_SAVE_FAILED)
	}

	// 重新查询完整数据
	s.fillPlanItems(&plan)
	return s.toPlanRes(plan), nil
}

// CreateCookingSession 创建做饭会话
func (s *WorkflowService) CreateCookingSession(ctx *gin.Context, r req.CreateCookingSessionReq) (rs res.CreateCookingSessionRes, err error) {
	userId := ctx.GetString("userId")
	if userId == "" {
		return rs, biz_err.New(biz_err.UNAUTHORIZED)
	}

	now := time.Now().Format(time.RFC3339)
	session := menuModel.MenuCookingSession{
		Id:               uuid.New().String(),
		UserId:           userId,
		RecipeId:         r.RecipeId,
		CurrentStepOrder: 1,
		State:            "active",
		StartedAt:        now,
		UpdatedAt:        now,
	}

	if err := global.GVA_DB.Create(&session).Error; err != nil {
		return rs, biz_err.New(biz_err.SESSION_CREATE_FAILED)
	}

	return res.CreateCookingSessionRes{
		Id:               session.Id,
		RecipeId:         session.RecipeId,
		CurrentStepOrder: session.CurrentStepOrder,
		State:            session.State,
		StartedAt:        session.StartedAt,
		UpdatedAt:        session.UpdatedAt,
	}, nil
}

// UpdateCookingSession 更新做饭会话进度
func (s *WorkflowService) UpdateCookingSession(ctx *gin.Context, r req.UpdateCookingSessionReq) (rs res.UpdateCookingSessionRes, err error) {
	userId := ctx.GetString("userId")
	if userId == "" {
		return rs, biz_err.New(biz_err.UNAUTHORIZED)
	}

	sessionId := ctx.Param("id")
	if sessionId == "" {
		return rs, biz_err.New(biz_err.PARAM_MISSING)
	}

	var session menuModel.MenuCookingSession
	if err := global.GVA_DB.Where("Id = ? AND UserId = ?", sessionId, userId).First(&session).Error; err != nil {
		return rs, biz_err.New(biz_err.SESSION_NOT_FOUND)
	}

	now := time.Now().Format(time.RFC3339)
	updates := map[string]interface{}{
		"CurrentStepOrder": r.CurrentStepOrder,
		"State":            r.State,
		"UpdatedAt":        now,
	}

	if err := global.GVA_DB.Model(&session).Updates(updates).Error; err != nil {
		return rs, biz_err.New(biz_err.SESSION_UPDATE_FAILED)
	}

	return res.UpdateCookingSessionRes{
		Id:               session.Id,
		RecipeId:         session.RecipeId,
		CurrentStepOrder: r.CurrentStepOrder,
		State:            r.State,
		StartedAt:        session.StartedAt,
		UpdatedAt:        now,
	}, nil
}

// SubmitFeedback 提交做饭反馈
func (s *WorkflowService) SubmitFeedback(ctx *gin.Context, r req.SubmitFeedbackReq) (rs res.SubmitFeedbackRes, err error) {
	userId := ctx.GetString("userId")
	if userId == "" {
		return rs, biz_err.New(biz_err.UNAUTHORIZED)
	}

	tagsJson, _ := json.Marshal(r.Tags)

	feedback := menuModel.MenuFeedback{
		Id:       uuid.New().String(),
		UserId:   userId,
		RecipeId: r.RecipeId,
		Outcome:  r.Outcome,
		TagsJson: string(tagsJson),
		Comment:  r.Comment,
	}

	if err := global.GVA_DB.Create(&feedback).Error; err != nil {
		return rs, biz_err.New(biz_err.FEEDBACK_CREATE_FAILED)
	}

	return res.SubmitFeedbackRes{
		Id:        feedback.Id,
		RecipeId:  feedback.RecipeId,
		Outcome:   feedback.Outcome,
		Tags:      r.Tags,
		Comment:   feedback.Comment,
		CreatedAt: feedback.CreatedAt.Format(time.RFC3339),
	}, nil
}

// ==================== 内部辅助 ====================

// fillPlanItems 填充计划项详情
func (s *WorkflowService) fillPlanItems(plan *menuModel.MenuMealPlan) {
	var items []menuModel.MenuMealPlanItem
	global.GVA_DB.Where("MealPlanId = ?", plan.Id).Order("SortOrder ASC").Find(&items)

	for i := range items {
		// 查询菜谱名和图片
		var recipe menuModel.MenuRecipe
		if err := global.GVA_DB.Where("Id = ?", items[i].RecipeId).First(&recipe).Error; err == nil {
			items[i].RecipeName = recipe.Name
			items[i].ImagePath = recipe.ImagePath
		}

		// 查询菜谱食材(分开查询避免 JOIN 问题)
		var recipeIngs []menuModel.MenuRecipeIngredient
		global.GVA_DB.Where("RecipeId = ?", items[i].RecipeId).Find(&recipeIngs)

		// 查询关联食材信息
		ingIds := make([]string, 0, len(recipeIngs))
		for _, ri := range recipeIngs {
			ingIds = append(ingIds, ri.IngredientId)
		}
		ingInfoMap := make(map[string]menuModel.MenuIngredient)
		if len(ingIds) > 0 {
			var ings []menuModel.MenuIngredient
			global.GVA_DB.Where("Id IN ?", ingIds).Find(&ings)
			for _, ing := range ings {
				ingInfoMap[ing.Id] = ing
			}
		}

		ings := make([]menuModel.MenuPlanIngredient, len(recipeIngs))
		for j, ri := range recipeIngs {
			qty := ri.Quantity * float64(items[i].Servings)
			var ingName, category, defaultUnit string
			var isPantryStaple bool
			if info, ok := ingInfoMap[ri.IngredientId]; ok {
				ingName = info.Name
				category = info.Category
				defaultUnit = info.DefaultUnit
				isPantryStaple = info.IsPantryStaple
			}
			ings[j] = menuModel.MenuPlanIngredient{
				IngredientId:   ri.IngredientId,
				IngredientName: ingName,
				Category:       category,
				DefaultUnit:    defaultUnit,
				IsPantryStaple: isPantryStaple,
				Quantity:       qty,
				Unit:           ri.Unit,
				Required:       ri.Required,
				Preparation:    ri.Preparation,
			}
		}
		items[i].Ingredients = ings
	}

	plan.Items = items

	// 合并所有食材
	combined := s.combineIngredients(items)
	plan.CombinedIngredients = combined
}

// combineIngredients 合并所有计划项的食材
func (s *WorkflowService) combineIngredients(items []menuModel.MenuMealPlanItem) []menuModel.MenuPlanIngredient {
	type key struct {
		IngredientId string
		Unit         string
	}
	merged := make(map[key]*menuModel.MenuPlanIngredient)

	for _, item := range items {
		for _, ing := range item.Ingredients {
			k := key{IngredientId: ing.IngredientId, Unit: ing.Unit}
			if existing, ok := merged[k]; ok {
				existing.Quantity += ing.Quantity
			} else {
				cp := ing
				merged[k] = &cp
			}
		}
	}

	result := make([]menuModel.MenuPlanIngredient, 0, len(merged))
	for _, v := range merged {
		result = append(result, *v)
	}
	return result
}

func (s *WorkflowService) toPlanRes(plan menuModel.MenuMealPlan) res.SavePlanRes {
	items := make([]res.SavePlanResItem, len(plan.Items))
	for i, item := range plan.Items {
		ings := make([]res.SavePlanResItemIngredient, len(item.Ingredients))
		for j, ing := range item.Ingredients {
			ings[j] = res.SavePlanResItemIngredient{
				IngredientId:   ing.IngredientId,
				IngredientName: ing.IngredientName,
				Category:       ing.Category,
				DefaultUnit:    ing.DefaultUnit,
				IsPantryStaple: ing.IsPantryStaple,
				Quantity:       ing.Quantity,
				Unit:           ing.Unit,
				Required:       ing.Required,
				Preparation:    ing.Preparation,
			}
		}
		items[i] = res.SavePlanResItem{
			Id:          item.Id,
			RecipeId:    item.RecipeId,
			RecipeName:  item.RecipeName,
			ImagePath:   item.ImagePath,
			Servings:    item.Servings,
			SortOrder:   item.SortOrder,
			Ingredients: ings,
		}
	}

	combined := make([]res.SavePlanResCombinedingredient, len(plan.CombinedIngredients))
	for i, ing := range plan.CombinedIngredients {
		combined[i] = res.SavePlanResCombinedingredient{
			IngredientId:   ing.IngredientId,
			IngredientName: ing.IngredientName,
			Category:       ing.Category,
			DefaultUnit:    ing.DefaultUnit,
			IsPantryStaple: ing.IsPantryStaple,
			Quantity:       ing.Quantity,
			Unit:           ing.Unit,
			Required:       ing.Required,
			Preparation:    ing.Preparation,
		}
	}

	return res.SavePlanRes{
		Id:                 plan.Id,
		PlanDate:           plan.PlanDate,
		Items:              items,
		CombinedIngredients: combined,
	}
}

func (s *WorkflowService) toPlanResList(plans []menuModel.MenuMealPlan) []res.SavePlanRes {
	list := make([]res.SavePlanRes, len(plans))
	for i, p := range plans {
		list[i] = s.toPlanRes(p)
	}
	return list
}
