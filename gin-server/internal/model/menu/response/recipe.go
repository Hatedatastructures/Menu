package response

type GetPublishedRecipeRes struct {
	Id string `json:"id"` // 菜谱ID
	Slug string `json:"slug"` // 菜谱别名
	Name string `json:"name"` // 菜谱名称
	Cuisine string `json:"cuisine"` // 菜系
	Description string `json:"description"` // 描述
	PrepMinutes int `json:"prepMinutes"` // 准备时间(分钟)
	CookMinutes int `json:"cookMinutes"` // 烹饪时间(分钟)
	TotalMinutes int `json:"totalMinutes"` // 总时间(分钟)
	Servings int `json:"servings"` // 份数
	Difficulty int `json:"difficulty"` // 难度(1-5)
	ImagePath string `json:"imagePath"` // 图片路径
	Status string `json:"status"` // 状态
	Allergens []string `json:"allergens"` // 过敏原列表
	Cookware []string `json:"cookware"` // 所需厨具
	Ingredients []GetPublishedRecipeResIngredient `json:"ingredients"`
	Steps []GetPublishedRecipeResStep `json:"steps"`
}

type GetPublishedRecipeResIngredient struct {
	IngredientId string `json:"ingredientId"` // 食材ID
	Quantity float64 `json:"quantity"` // 用量
	Unit string `json:"unit"` // 单位
	Required bool `json:"required"` // 是否必需
	ServingFactor float64 `json:"servingFactor"` // 份数系数
	Preparation string `json:"preparation"` // 处理方式
	IngredientName string `json:"ingredientName"` // 食材名称
	Category string `json:"category"` // 食材分类
	DefaultUnit string `json:"defaultUnit"` // 默认单位
	IsPantryStaple bool `json:"isPantryStaple"` // 是否常备食材
}

type GetPublishedRecipeResStep struct {
	StepOrder int `json:"stepOrder"` // 步骤序号
	Title string `json:"title"` // 步骤标题
	Instruction string `json:"instruction"` // 操作说明
	DurationSeconds int `json:"durationSeconds"` // 持续时间(秒)
	HasTimer bool `json:"hasTimer"` // 是否需要计时
}
