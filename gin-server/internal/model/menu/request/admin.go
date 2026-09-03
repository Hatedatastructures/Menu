package request

type CreateRecipeReq struct {
	Slug string `json:"slug" form:"slug"` // 菜谱别名
	Name string `json:"name" form:"name" validate:"required,max=128,min=1"` // 菜谱名称
	Cuisine string `json:"cuisine" form:"cuisine"` // 菜系
	Description string `json:"description" form:"description"` // 描述
	PrepMinutes int `json:"prepMinutes" form:"prepMinutes"` // 准备时间(分钟)
	CookMinutes int `json:"cookMinutes" form:"cookMinutes"` // 烹饪时间(分钟)
	Servings int `json:"servings" form:"servings" validate:"required,gte=1"` // 份数
	Difficulty int `json:"difficulty" form:"difficulty" validate:"required,gte=1,lte=5"` // 难度(1-5)
	ImagePath string `json:"imagePath" form:"imagePath"` // 图片路径
	Status string `json:"status" form:"status" validate:"required,oneof=published draft archived"` // 状态
	Allergens []string `json:"allergens" form:"allergens"` // 过敏原列表
	Cookware []string `json:"cookware" form:"cookware"` // 所需厨具
	Ingredients []CreateRecipeReqIngredient `json:"ingredients" form:"ingredients"`
	Steps []CreateRecipeReqStep `json:"steps" form:"steps"`
}

type CreateRecipeReqIngredient struct {
	IngredientId string `json:"ingredientId" form:"ingredientId"` // 食材ID
	Quantity float64 `json:"quantity" form:"quantity"` // 用量
	Unit string `json:"unit" form:"unit"` // 单位
	Required bool `json:"required" form:"required"` // 是否必需
	ServingFactor float64 `json:"servingFactor" form:"servingFactor"` // 份数系数
	Preparation string `json:"preparation" form:"preparation"` // 处理方式
}

type CreateRecipeReqStep struct {
	StepOrder int `json:"stepOrder" form:"stepOrder"` // 步骤序号
	Title string `json:"title" form:"title"` // 步骤标题
	Instruction string `json:"instruction" form:"instruction"` // 操作说明
	DurationSeconds int `json:"durationSeconds" form:"durationSeconds"` // 持续时间(秒)
	HasTimer bool `json:"hasTimer" form:"hasTimer"` // 是否需要计时
}

type UpdateRecipeReq struct {
	Slug string `json:"slug" form:"slug"`
	Name string `json:"name" form:"name"`
	Cuisine string `json:"cuisine" form:"cuisine"`
	Description string `json:"description" form:"description"`
	PrepMinutes int `json:"prepMinutes" form:"prepMinutes"`
	CookMinutes int `json:"cookMinutes" form:"cookMinutes"`
	Servings int `json:"servings" form:"servings"`
	Difficulty int `json:"difficulty" form:"difficulty"`
	ImagePath string `json:"imagePath" form:"imagePath"`
	Status string `json:"status" form:"status"`
	Allergens []string `json:"allergens" form:"allergens"`
	Cookware []string `json:"cookware" form:"cookware"`
	Ingredients []UpdateRecipeReqIngredient `json:"ingredients" form:"ingredients"`
	Steps []UpdateRecipeReqStep `json:"steps" form:"steps"`
}

type UpdateRecipeReqIngredient struct {
	IngredientId string `json:"ingredientId" form:"ingredientId"`
	Quantity float64 `json:"quantity" form:"quantity"`
	Unit string `json:"unit" form:"unit"`
	Required bool `json:"required" form:"required"`
	ServingFactor float64 `json:"servingFactor" form:"servingFactor"`
	Preparation string `json:"preparation" form:"preparation"`
}

type UpdateRecipeReqStep struct {
	StepOrder int `json:"stepOrder" form:"stepOrder"`
	Title string `json:"title" form:"title"`
	Instruction string `json:"instruction" form:"instruction"`
	DurationSeconds int `json:"durationSeconds" form:"durationSeconds"`
	HasTimer bool `json:"hasTimer" form:"hasTimer"`
}

type CreateIngredientReq struct {
	Name string `json:"name" form:"name" validate:"required,max=64,min=1"` // 食材名称
	Aliases []string `json:"aliases" form:"aliases"` // 别名列表
	Category string `json:"category" form:"category"` // 分类
	DefaultUnit string `json:"defaultUnit" form:"defaultUnit"` // 默认单位
	IsPantryStaple bool `json:"isPantryStaple" form:"isPantryStaple"` // 是否常备食材
	SubstituteGroup string `json:"substituteGroup" form:"substituteGroup"` // 替代组
	StoreSkuMapping string `json:"storeSkuMapping" form:"storeSkuMapping"` // 商店SKU映射
}

type UpdateIngredientReq struct {
	Name string `json:"name" form:"name"`
	Aliases []string `json:"aliases" form:"aliases"`
	Category string `json:"category" form:"category"`
	DefaultUnit string `json:"defaultUnit" form:"defaultUnit"`
	IsPantryStaple bool `json:"isPantryStaple" form:"isPantryStaple"`
	SubstituteGroup string `json:"substituteGroup" form:"substituteGroup"`
	StoreSkuMapping string `json:"storeSkuMapping" form:"storeSkuMapping"`
}
