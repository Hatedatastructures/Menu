package menu

import "time"

// ==================== 食材 ====================

type MenuIngredient struct {
	Id              string `json:"id" gorm:"column:Id;type:text;primaryKey"`
	Name            string `json:"name" gorm:"column:Name;type:text;not null;uniqueIndex"`
	Category        string `json:"category" gorm:"column:Category;type:text;not null"`
	DefaultUnit     string `json:"defaultUnit" gorm:"column:DefaultUnit;type:text;not null"`
	IsPantryStaple  bool   `json:"isPantryStaple" gorm:"column:IsPantryStaple;type integer;not null;default:0"`
	SubstituteGroup string `json:"substituteGroup" gorm:"column:SubstituteGroup;type:text;not null;default:''"`
	StoreSkuMapping string `json:"storeSkuMapping" gorm:"column:StoreSkuMapping;type:text"`
	Aliases         []MenuIngredientAlias `json:"aliases" gorm:"foreignKey:IngredientId;references:Id"`
}

func (MenuIngredient) TableName() string { return "Ingredients" }

type MenuIngredientAlias struct {
	Id           int    `json:"-" gorm:"column:Id;type:integer;primaryKey;autoIncrement"`
	IngredientId string `json:"-" gorm:"column:IngredientId;type:text;not null;index"`
	Alias        string `json:"alias" gorm:"column:Alias;type:text;not null;uniqueIndex"`
}

func (MenuIngredientAlias) TableName() string { return "IngredientAliases" }

// ==================== 菜谱 ====================

type MenuRecipe struct {
	Id            string           `json:"id" gorm:"column:Id;type:text;primaryKey"`
	Slug          string           `json:"slug" gorm:"column:Slug;type:text;not null;uniqueIndex"`
	Name          string           `json:"name" gorm:"column:Name;type:text;not null"`
	Cuisine       string           `json:"cuisine" gorm:"column:Cuisine;type:text;not null"`
	Description   string           `json:"description" gorm:"column:Description;type:text;not null"`
	PrepMinutes   int              `json:"prepMinutes" gorm:"column:PrepMinutes;type:integer;not null;default:0"`
	CookMinutes   int              `json:"cookMinutes" gorm:"column:CookMinutes;type:integer;not null;default:0"`
	Servings      int              `json:"servings" gorm:"column:Servings;type:integer;not null;default:1"`
	Difficulty    int              `json:"difficulty" gorm:"column:Difficulty;type:integer;not null;default:1"`
	ImagePath     string           `json:"imagePath" gorm:"column:ImagePath;type:text;not null;default:''"`
	Status        string           `json:"status" gorm:"column:Status;type:text;not null;default:'draft'"`
	AllergensJson string           `json:"-" gorm:"column:AllergensJson;type:text;not null;default:'[]'"`
	CookwareJson  string           `json:"-" gorm:"column:CookwareJson;type:text;not null;default:'[]'"`
	Ingredients   []MenuRecipeIngredient `json:"-" gorm:"foreignKey:RecipeId;references:Id"`
	Steps         []MenuRecipeStep       `json:"-" gorm:"foreignKey:RecipeId;references:Id"`
	CreatedAt     time.Time        `json:"createdAt" gorm:"column:CreatedAt;autoCreateTime"`
	UpdatedAt     time.Time        `json:"updatedAt" gorm:"column:UpdatedAt;autoUpdateTime"`
}

func (MenuRecipe) TableName() string { return "Recipes" }

type MenuRecipeIngredient struct {
	RecipeId      string  `json:"recipeId" gorm:"column:RecipeId;type:text;not null;primaryKey"`
	IngredientId  string  `json:"ingredientId" gorm:"column:IngredientId;type:text;not null;primaryKey"`
	Quantity      float64 `json:"quantity" gorm:"column:Quantity;type:real;not null"`
	Unit          string  `json:"unit" gorm:"column:Unit;type:text;not null"`
	ServingFactor float64 `json:"servingFactor" gorm:"column:ServingFactor;type:real;not null;default:1.0"`
	Preparation   string  `json:"preparation" gorm:"column:Preparation;type:text;not null;default:''"`
	Required      bool    `json:"required" gorm:"column:Required;type:integer;not null;default:1"`
	// 关联食材信息(从 Ingredients 表 JOIN)
	IngredientName    string `json:"ingredientName" gorm:"-"`
	Category          string `json:"category" gorm:"-"`
	DefaultUnit       string `json:"defaultUnit" gorm:"-"`
	IsPantryStaple    bool   `json:"isPantryStaple" gorm:"-"`
}

func (MenuRecipeIngredient) TableName() string { return "RecipeIngredients" }

type MenuRecipeStep struct {
	RecipeId         string `json:"recipeId" gorm:"column:RecipeId;type:text;not null;primaryKey"`
	StepOrder        int    `json:"stepOrder" gorm:"column:StepOrder;type:integer;not null;primaryKey"`
	Title            string `json:"title" gorm:"column:Title;type:text;not null"`
	Instruction      string `json:"instruction" gorm:"column:Instruction;type:text;not null"`
	DurationSeconds  int    `json:"durationSeconds" gorm:"column:DurationSeconds;type:integer;not null;default:0"`
	HasTimer         bool   `json:"hasTimer" gorm:"column:HasTimer;type:integer;not null;default:0"`
}

func (MenuRecipeStep) TableName() string { return "RecipeSteps" }

// ==================== 用户 ====================

type MenuUser struct {
	Id            string `json:"id" gorm:"column:Id;type:text;primaryKey"`
	Email         string `json:"email" gorm:"column:Email;type:text;not null;uniqueIndex"`
	PasswordHash  string `json:"-" gorm:"column:PasswordHash;type:text;not null"`
	DisplayName   string `json:"displayName" gorm:"column:DisplayName;type:text;not null"`
	IsAdmin       bool   `json:"isAdmin" gorm:"column:IsAdmin;type:integer;not null;default:0"`
}

func (MenuUser) TableName() string { return "Users" }

// ==================== 刷新令牌 ====================

type MenuRefreshToken struct {
	Id        string       `json:"id" gorm:"column:Id;type:text;primaryKey"`
	UserId    string       `json:"userId" gorm:"column:UserId;type:text;not null;index"`
	TokenHash string       `json:"-" gorm:"column:TokenHash;type:text;not null;uniqueIndex"`
	ExpiresAt time.Time    `json:"expiresAt" gorm:"column:ExpiresAt;type:text;not null"`
	RevokedAt *time.Time   `json:"-" gorm:"column:RevokedAt;type:text"`
}

func (MenuRefreshToken) TableName() string { return "RefreshTokens" }

// ==================== 餐食计划 ====================

type MenuMealPlan struct {
	Id        string                 `json:"id" gorm:"column:Id;type:text;primaryKey"`
	UserId    string                 `json:"userId" gorm:"column:UserId;type:text;not null;index"`
	PlanDate  string                 `json:"planDate" gorm:"column:PlanDate;type:text;not null"`
	Items     []MenuMealPlanItem     `json:"-" gorm:"foreignKey:MealPlanId;references:Id"`
	CombinedIngredients []MenuPlanIngredient `json:"-" gorm:"-"`
	CreatedAt time.Time              `json:"createdAt" gorm:"column:CreatedAt;autoCreateTime"`
}

func (MenuMealPlan) TableName() string { return "MealPlans" }

type MenuMealPlanItem struct {
	Id          string `json:"id" gorm:"column:Id;type:text;primaryKey"`
	MealPlanId  string `json:"-" gorm:"column:MealPlanId;type:text;not null;index"`
	RecipeId    string `json:"recipeId" gorm:"column:RecipeId;type:text;not null"`
	Servings    int    `json:"servings" gorm:"column:Servings;type:integer;not null;default:1"`
	SortOrder   int    `json:"sortOrder" gorm:"column:SortOrder;type:integer;not null;default:0"`
	// 关联信息(运行时填充)
	RecipeName  string                `json:"recipeName" gorm:"-"`
	ImagePath   string                `json:"imagePath" gorm:"-"`
	Ingredients []MenuPlanIngredient `json:"ingredients" gorm:"-"`
}

func (MenuMealPlanItem) TableName() string { return "MealPlanItems" }

// ==================== 做饭会话 ====================

type MenuCookingSession struct {
	Id               string `json:"id" gorm:"column:Id;type:text;primaryKey"`
	UserId           string `json:"userId" gorm:"column:UserId;type:text;not null;index"`
	RecipeId         string `json:"recipeId" gorm:"column:RecipeId;type:text;not null"`
	CurrentStepOrder int    `json:"currentStepOrder" gorm:"column:CurrentStepOrder;type:integer;not null;default:1"`
	State            string `json:"state" gorm:"column:State;type:text;not null;default:'active'"`
	StartedAt        string `json:"startedAt" gorm:"column:StartedAt;type:text;not null"`
	UpdatedAt        string `json:"updatedAt" gorm:"column:UpdatedAt;type:text;not null"`
}

func (MenuCookingSession) TableName() string { return "CookingSessions" }

// ==================== 反馈 ====================

type MenuFeedback struct {
	Id        string `json:"id" gorm:"column:Id;type:text;primaryKey"`
	UserId    string `json:"userId" gorm:"column:UserId;type:text;not null;index"`
	RecipeId  string `json:"recipeId" gorm:"column:RecipeId;type:text;not null"`
	Outcome   string `json:"outcome" gorm:"column:Outcome;type:text;not null"`
	TagsJson  string `json:"-" gorm:"column:TagsJson;type:text;not null;default:'[]'"`
	Comment   string `json:"comment" gorm:"column:Comment;type:text;not null;default:''"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:CreatedAt;autoCreateTime"`
}

func (MenuFeedback) TableName() string { return "Feedback" }

// ==================== 辅助类型 ====================

// MenuPlanIngredient 用于计划中合并食材的展示
type MenuPlanIngredient struct {
	IngredientId   string  `json:"ingredientId"`
	IngredientName string  `json:"ingredientName"`
	Category       string  `json:"category"`
	DefaultUnit    string  `json:"defaultUnit"`
	IsPantryStaple bool    `json:"isPantryStaple"`
	Quantity       float64 `json:"quantity"`
	Unit           string  `json:"unit"`
	Required       bool    `json:"required"`
	Preparation    string  `json:"preparation"`
}
