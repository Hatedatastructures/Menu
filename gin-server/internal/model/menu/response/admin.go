package response

type GetAdminStatsRes struct {
	PublishedRecipes int `json:"publishedRecipes"` // 已发布菜谱数
	DraftRecipes int `json:"draftRecipes"` // 草稿菜谱数
	ArchivedRecipes int `json:"archivedRecipes"` // 已归档菜谱数
	IngredientCount int `json:"ingredientCount"` // 食材总数
}

type CreateRecipeRes struct {
	Id string `json:"id"` // 菜谱ID
	Slug string `json:"slug"`
	Name string `json:"name"`
	Cuisine string `json:"cuisine"`
	Description string `json:"description"`
	PrepMinutes int `json:"prepMinutes"`
	CookMinutes int `json:"cookMinutes"`
	TotalMinutes int `json:"totalMinutes"`
	Servings int `json:"servings"`
	Difficulty int `json:"difficulty"`
	ImagePath string `json:"imagePath"`
	Status string `json:"status"`
	Allergens []string `json:"allergens"`
	Cookware []string `json:"cookware"`
	Ingredients []CreateRecipeResIngredient `json:"ingredients"`
	Steps []CreateRecipeResStep `json:"steps"`
}

type CreateRecipeResIngredient struct {
	IngredientId string `json:"ingredientId"`
	Quantity float64 `json:"quantity"`
	Unit string `json:"unit"`
	Required bool `json:"required"`
	ServingFactor float64 `json:"servingFactor"`
	Preparation string `json:"preparation"`
	IngredientName string `json:"ingredientName"`
	Category string `json:"category"`
	DefaultUnit string `json:"defaultUnit"`
	IsPantryStaple bool `json:"isPantryStaple"`
}

type CreateRecipeResStep struct {
	StepOrder int `json:"stepOrder"`
	Title string `json:"title"`
	Instruction string `json:"instruction"`
	DurationSeconds int `json:"durationSeconds"`
	HasTimer bool `json:"hasTimer"`
}

type UpdateRecipeRes struct {
	Id string `json:"id"`
	Slug string `json:"slug"`
	Name string `json:"name"`
	Cuisine string `json:"cuisine"`
	Description string `json:"description"`
	PrepMinutes int `json:"prepMinutes"`
	CookMinutes int `json:"cookMinutes"`
	TotalMinutes int `json:"totalMinutes"`
	Servings int `json:"servings"`
	Difficulty int `json:"difficulty"`
	ImagePath string `json:"imagePath"`
	Status string `json:"status"`
	Allergens []string `json:"allergens"`
	Cookware []string `json:"cookware"`
	Ingredients []UpdateRecipeResIngredient `json:"ingredients"`
	Steps []UpdateRecipeResStep `json:"steps"`
}

type UpdateRecipeResIngredient struct {
	IngredientId string `json:"ingredientId"`
	Quantity float64 `json:"quantity"`
	Unit string `json:"unit"`
	Required bool `json:"required"`
	ServingFactor float64 `json:"servingFactor"`
	Preparation string `json:"preparation"`
	IngredientName string `json:"ingredientName"`
	Category string `json:"category"`
	DefaultUnit string `json:"defaultUnit"`
	IsPantryStaple bool `json:"isPantryStaple"`
}

type UpdateRecipeResStep struct {
	StepOrder int `json:"stepOrder"`
	Title string `json:"title"`
	Instruction string `json:"instruction"`
	DurationSeconds int `json:"durationSeconds"`
	HasTimer bool `json:"hasTimer"`
}

type CreateIngredientRes struct {
	Id string `json:"id"` // 食材ID
	Name string `json:"name"`
	Aliases []string `json:"aliases"`
	Category string `json:"category"`
	DefaultUnit string `json:"defaultUnit"`
	IsPantryStaple bool `json:"isPantryStaple"`
	SubstituteGroup string `json:"substituteGroup"`
	StoreSkuMapping string `json:"storeSkuMapping"`
}

type UpdateIngredientRes struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Aliases []string `json:"aliases"`
	Category string `json:"category"`
	DefaultUnit string `json:"defaultUnit"`
	IsPantryStaple bool `json:"isPantryStaple"`
	SubstituteGroup string `json:"substituteGroup"`
	StoreSkuMapping string `json:"storeSkuMapping"`
}
