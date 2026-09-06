package request

type ListPublishedRecipesReq struct {
	Cuisine string `json:"cuisine" form:"cuisine"`
	MaxMinutes int `json:"maxMinutes" form:"maxMinutes"`
	Difficulty int `json:"difficulty" form:"difficulty"`
	Limit int `json:"limit" form:"limit"`
}

type RecommendTonightReq struct {
	Servings int `json:"servings" form:"servings"` // 用餐人数
	AvailableMinutes int `json:"availableMinutes" form:"availableMinutes"` // 可用时间(分钟)
	Cuisines []string `json:"cuisines" form:"cuisines"` // 偏好菜系
	Allergies []string `json:"allergies" form:"allergies"` // 过敏原
	PantryIngredientIds []string `json:"pantryIngredientIds" form:"pantryIngredientIds"` // 已有食材ID
	Cookware []string `json:"cookware" form:"cookware"` // 已有厨具
	RecentRecipeIds []string `json:"recentRecipeIds" form:"recentRecipeIds"` // 最近做过的菜谱ID
}
