package menu_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"shack/internal/api/v1/menu"
	"shack/internal/global"
	menuModel "shack/internal/model/menu"

	"github.com/glebarez/sqlite"
)

// ==================== 测试基础设施 ====================

var testDB *gorm.DB
var testRouter *gin.Engine

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	// 初始化测试数据库(内存模式)
	var err error
	testDB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("failed to open test db: %v", err))
	}

	// AutoMigrate 所有 menu 模型
	if err := testDB.AutoMigrate(
		&menuModel.MenuUser{},
		&menuModel.MenuRefreshToken{},
		&menuModel.MenuIngredient{},
		&menuModel.MenuIngredientAlias{},
		&menuModel.MenuRecipe{},
		&menuModel.MenuRecipeIngredient{},
		&menuModel.MenuRecipeStep{},
		&menuModel.MenuMealPlan{},
		&menuModel.MenuMealPlanItem{},
		&menuModel.MenuCookingSession{},
		&menuModel.MenuFeedback{},
	); err != nil {
		panic(fmt.Sprintf("failed to migrate: %v", err))
	}

	// 设置全局 DB
	global.GVA_DB = testDB

	// 初始化测试路由
	testRouter = setupTestRouter()

	m.Run()
}

func setupTestRouter() *gin.Engine {
	r := gin.New()

	// 注册所有 menu 路由(绕过 utils.RegisterApi 的 apiSet 检查)
	apiGroup := r.Group("")

	// Health (公开)
	apiGroup.GET("/healthz", (&menu.HealthApi{}).HealthzHandler)
	apiGroup.GET("/readyz", (&menu.HealthApi{}).ReadyzHandler)

	// Auth (公开)
	apiGroup.POST("/api/v1/auth/register", (&menu.AuthApi{}).RegisterHandler)
	apiGroup.POST("/api/v1/auth/login", (&menu.AuthApi{}).LoginHandler)
	apiGroup.POST("/api/v1/auth/refresh", (&menu.AuthApi{}).RefreshTokenHandler)

	// Recipe (公开)
	apiGroup.GET("/api/v1/recipes", (&menu.RecipeApi{}).ListPublishedRecipesHandler)
	apiGroup.GET("/api/v1/ingredients", (&menu.RecipeApi{}).ListIngredientsHandler)
	apiGroup.GET("/api/v1/recipes/:id", (&menu.RecipeApi{}).GetPublishedRecipeHandler)
	apiGroup.POST("/api/v1/recommendations/tonight", (&menu.RecipeApi{}).RecommendTonightHandler)

	// Admin (需认证)
	apiGroup.GET("/api/v1/admin/stats", (&menu.AdminApi{}).GetAdminStatsHandler)
	apiGroup.GET("/api/v1/admin/recipes", (&menu.AdminApi{}).ListAllRecipesHandler)
	apiGroup.POST("/api/v1/admin/recipes", (&menu.AdminApi{}).CreateRecipeHandler)
	apiGroup.PUT("/api/v1/admin/recipes/:id", (&menu.AdminApi{}).UpdateRecipeHandler)
	apiGroup.DELETE("/api/v1/admin/recipes/:id", (&menu.AdminApi{}).DeleteRecipeHandler)
	apiGroup.GET("/api/v1/admin/ingredients", (&menu.AdminApi{}).ListAllIngredientsHandler)
	apiGroup.POST("/api/v1/admin/ingredients", (&menu.AdminApi{}).CreateIngredientHandler)
	apiGroup.PUT("/api/v1/admin/ingredients/:id", (&menu.AdminApi{}).UpdateIngredientHandler)
	apiGroup.DELETE("/api/v1/admin/ingredients/:id", (&menu.AdminApi{}).DeleteIngredientHandler)

	// Workflow (需认证,注入测试 userId)
	authHandler := func(h func(*gin.Context)) gin.HandlerFunc {
		return func(c *gin.Context) {
			c.Set("userId", "test-user-001")
			h(c)
		}
	}
	apiGroup.GET("/api/v1/plans", authHandler((&menu.WorkflowApi{}).ListPlansHandler))
	apiGroup.POST("/api/v1/plans", authHandler((&menu.WorkflowApi{}).SavePlanHandler))
	apiGroup.POST("/api/v1/cooking-sessions", authHandler((&menu.WorkflowApi{}).CreateCookingSessionHandler))
	apiGroup.PUT("/api/v1/cooking-sessions/:id", authHandler((&menu.WorkflowApi{}).UpdateCookingSessionHandler))
	apiGroup.POST("/api/v1/feedback", authHandler((&menu.WorkflowApi{}).SubmitFeedbackHandler))

	return r
}

// ==================== 测试辅助函数 ====================

type apiResponse struct {
	Code      int             `json:"code"`
	Msg       string          `json:"msg"`
	Data      json.RawMessage `json:"data"`
	RequestId string          `json:"requestId"`
	TimeStamp int64           `json:"timeStamp"`
}

func doRequest(method, path string, body interface{}) (*httptest.ResponseRecorder, apiResponse) {
	var reqBody *bytes.Buffer
	if body != nil {
		jsonBytes, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(jsonBytes)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}

	req, _ := http.NewRequest(method, path, reqBody)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	var resp apiResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	return w, resp
}

func doRequestWithQuery(method, path string) (*httptest.ResponseRecorder, apiResponse) {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	var resp apiResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	return w, resp
}

// ==================== Health 模块测试 ====================

func TestHealthz(t *testing.T) {
	w, resp := doRequestWithQuery(http.MethodGet, "/healthz")
	if w.Code != 200 {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
	if resp.Code != 0 {
		t.Fatalf("expected business code 0, got %d", resp.Code)
	}

	var data struct {
		Status string `json:"status"`
	}
	json.Unmarshal(resp.Data, &data)
	if data.Status != "ok" {
		t.Fatalf("expected status 'ok', got '%s'", data.Status)
	}
	t.Logf("✓ Healthz: status=%s", data.Status)
}

func TestReadyz(t *testing.T) {
	w, resp := doRequestWithQuery(http.MethodGet, "/readyz")
	if w.Code != 200 {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var data struct {
		Status string `json:"status"`
	}
	json.Unmarshal(resp.Data, &data)
	if data.Status != "ready" {
		t.Fatalf("expected status 'ready', got '%s'", data.Status)
	}
	t.Logf("✓ Readyz: status=%s", data.Status)
}

// ==================== Auth 模块测试 ====================

func TestRegister(t *testing.T) {
	w, resp := doRequest(http.MethodPost, "/api/v1/auth/register", map[string]interface{}{
		"email":       "test@example.com",
		"password":    "password123",
		"displayName": "测试用户",
	})
	if w.Code != 200 {
		t.Fatalf("expected status 200, got %d, body: %s", w.Code, w.Body.String())
	}
	if resp.Code != 0 {
		t.Fatalf("expected business code 0, got %d: %s", resp.Code, resp.Msg)
	}

	var data struct {
		User struct {
			Id          string `json:"id"`
			Email       string `json:"email"`
			DisplayName string `json:"displayName"`
			IsAdmin     bool   `json:"isAdmin"`
		} `json:"user"`
		AccessToken               string `json:"accessToken"`
		RefreshToken              string `json:"refreshToken"`
		AccessTokenExpiresInSeconds int    `json:"accessTokenExpiresInSeconds"`
	}
	json.Unmarshal(resp.Data, &data)

	if data.User.Email != "test@example.com" {
		t.Fatalf("expected email 'test@example.com', got '%s'", data.User.Email)
	}
	if data.User.DisplayName != "测试用户" {
		t.Fatalf("expected displayName '测试用户', got '%s'", data.User.DisplayName)
	}
	if data.User.Id == "" {
		t.Fatal("expected non-empty user id")
	}
	if data.AccessToken == "" {
		t.Fatal("expected non-empty accessToken")
	}
	if data.RefreshToken == "" {
		t.Fatal("expected non-empty refreshToken")
	}
	if data.AccessTokenExpiresInSeconds != 3600 {
		t.Fatalf("expected accessTokenExpiresInSeconds=3600, got %d", data.AccessTokenExpiresInSeconds)
	}
	t.Logf("✓ Register: userId=%s, email=%s", data.User.Id, data.User.Email)
}

func TestRegisterDuplicate(t *testing.T) {
	_, resp := doRequest(http.MethodPost, "/api/v1/auth/register", map[string]interface{}{
		"email":       "test@example.com",
		"password":    "password123",
		"displayName": "重复用户",
	})
	// 业务失败时 handler 返回 HTTP 201
	if resp.Code == 0 {
		t.Fatal("expected error for duplicate email, got success")
	}
	t.Logf("✓ RegisterDuplicate: code=%d, msg=%s", resp.Code, resp.Msg)
}

func TestLogin(t *testing.T) {
	w, resp := doRequest(http.MethodPost, "/api/v1/auth/login", map[string]interface{}{
		"email":    "test@example.com",
		"password": "password123",
	})
	if w.Code != 200 {
		t.Fatalf("expected status 200, got %d, body: %s", w.Code, w.Body.String())
	}
	if resp.Code != 0 {
		t.Fatalf("expected business code 0, got %d: %s", resp.Code, resp.Msg)
	}

	var data struct {
		User struct {
			Email string `json:"email"`
		} `json:"user"`
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
	}
	json.Unmarshal(resp.Data, &data)

	if data.User.Email != "test@example.com" {
		t.Fatalf("expected email 'test@example.com', got '%s'", data.User.Email)
	}
	if data.AccessToken == "" {
		t.Fatal("expected non-empty accessToken")
	}
	t.Logf("✓ Login: email=%s, tokenLen=%d", data.User.Email, len(data.AccessToken))
}

func TestLoginWrongPassword(t *testing.T) {
	_, resp := doRequest(http.MethodPost, "/api/v1/auth/login", map[string]interface{}{
		"email":    "test@example.com",
		"password": "wrongpassword",
	})
	if resp.Code == 0 {
		t.Fatal("expected error for wrong password, got success")
	}
	t.Logf("✓ LoginWrongPassword: code=%d, msg=%s", resp.Code, resp.Msg)
}

// 获取刷新令牌(复用登录结果)
func getRefreshToken(t *testing.T) string {
	t.Helper()
	_, resp := doRequest(http.MethodPost, "/api/v1/auth/login", map[string]interface{}{
		"email":    "test@example.com",
		"password": "password123",
	})
	if resp.Code != 0 {
		t.Fatalf("login failed: %s", resp.Msg)
	}
	var data struct {
		RefreshToken string `json:"refreshToken"`
	}
	json.Unmarshal(resp.Data, &data)
	return data.RefreshToken
}

func TestRefreshToken(t *testing.T) {
	refreshToken := getRefreshToken(t)

	w, resp := doRequest(http.MethodPost, "/api/v1/auth/refresh", map[string]interface{}{
		"refreshToken": refreshToken,
	})
	if w.Code != 200 {
		t.Fatalf("expected status 200, got %d, body: %s", w.Code, w.Body.String())
	}
	if resp.Code != 0 {
		t.Fatalf("expected business code 0, got %d: %s", resp.Code, resp.Msg)
	}

	var data struct {
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
	}
	json.Unmarshal(resp.Data, &data)

	if data.AccessToken == "" {
		t.Fatal("expected new accessToken")
	}
	if data.RefreshToken == "" {
		t.Fatal("expected new refreshToken")
	}
	t.Logf("✓ RefreshToken: new token obtained")
}

// ==================== Recipe 模块测试 ====================

func TestListPublishedRecipesEmpty(t *testing.T) {
	w, resp := doRequestWithQuery(http.MethodGet, "/api/v1/recipes")
	if w.Code != 200 {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
	if resp.Code != 0 {
		t.Fatalf("expected business code 0, got %d", resp.Code)
	}
	t.Log("✓ ListPublishedRecipes (empty): returned successfully")
}

func TestListIngredientsEmpty(t *testing.T) {
	w, resp := doRequestWithQuery(http.MethodGet, "/api/v1/ingredients")
	if w.Code != 200 {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
	if resp.Code != 0 {
		t.Fatalf("expected business code 0, got %d", resp.Code)
	}
	t.Log("✓ ListIngredients (empty): returned successfully")
}

func TestRecommendTonight(t *testing.T) {
	w, resp := doRequest(http.MethodPost, "/api/v1/recommendations/tonight", map[string]interface{}{
		"servings":         2,
		"availableMinutes": 60,
		"cuisines":         []string{"中餐"},
		"allergies":        []string{},
		"pantryIngredientIds": []string{},
		"cookware":         []string{},
		"recentRecipeIds":  []string{},
	})
	if w.Code != 200 {
		t.Fatalf("expected status 200, got %d, body: %s", w.Code, w.Body.String())
	}
	if resp.Code != 0 {
		t.Fatalf("expected business code 0, got %d: %s", resp.Code, resp.Msg)
	}
	t.Log("✓ RecommendTonight: returned successfully")
}

// ==================== Admin 模块测试 ====================

func TestGetAdminStats(t *testing.T) {
	w, resp := doRequestWithQuery(http.MethodGet, "/api/v1/admin/stats")
	if w.Code != 200 {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
	if resp.Code != 0 {
		t.Fatalf("expected business code 0, got %d: %s", resp.Code, resp.Msg)
	}

	var data struct {
		PublishedRecipes int `json:"publishedRecipes"`
		DraftRecipes     int `json:"draftRecipes"`
		ArchivedRecipes  int `json:"archivedRecipes"`
		IngredientCount  int `json:"ingredientCount"`
	}
	json.Unmarshal(resp.Data, &data)

	if data.PublishedRecipes != 0 || data.IngredientCount != 0 {
		t.Fatalf("expected all zeros, got %+v", data)
	}
	t.Logf("✓ GetAdminStats: %+v", data)
}

func TestAdminRecipeCRUD(t *testing.T) {
	// 0. 先创建测试食材(菜谱引用的食材必须存在)
	_, tomatoResp := doRequest(http.MethodPost, "/api/v1/admin/ingredients", map[string]interface{}{
		"name": "番茄-CRUD测试", "category": "蔬菜", "defaultUnit": "g",
	})
	var tomatoCreated struct{ Id string }
	json.Unmarshal(tomatoResp.Data, &tomatoCreated)

	_, eggResp := doRequest(http.MethodPost, "/api/v1/admin/ingredients", map[string]interface{}{
		"name": "鸡蛋-CRUD测试", "category": "蛋类", "defaultUnit": "个",
	})
	var eggCreated struct{ Id string }
	json.Unmarshal(eggResp.Data, &eggCreated)

	// 1. 创建菜谱
	w, resp := doRequest(http.MethodPost, "/api/v1/admin/recipes", map[string]interface{}{
		"slug":        "test-stir-fry",
		"name":        "番茄炒蛋",
		"cuisine":     "中餐",
		"description": "经典的家常菜",
		"prepMinutes": 5,
		"cookMinutes": 10,
		"servings":    2,
		"difficulty":  1,
		"imagePath":   "/images/tomato-egg.jpg",
		"status":      "published",
		"allergens":   []string{"鸡蛋"},
		"cookware":    []string{"炒锅"},
		"ingredients": []map[string]interface{}{
			{"ingredientId": tomatoCreated.Id, "quantity": 200, "unit": "g", "required": true, "servingFactor": 1.0, "preparation": "切块"},
			{"ingredientId": eggCreated.Id, "quantity": 3, "unit": "个", "required": true, "servingFactor": 1.0, "preparation": "打散"},
		},
		"steps": []map[string]interface{}{
			{"stepOrder": 1, "title": "准备食材", "instruction": "番茄切块,鸡蛋打散", "durationSeconds": 300, "hasTimer": false},
			{"stepOrder": 2, "title": "炒蛋", "instruction": "热油炒蛋至凝固盛出", "durationSeconds": 120, "hasTimer": true},
			{"stepOrder": 3, "title": "炒番茄", "instruction": "炒番茄出汁后加入鸡蛋", "durationSeconds": 180, "hasTimer": true},
		},
	})
	if w.Code != 200 || resp.Code != 0 {
		t.Fatalf("create recipe failed: status=%d, code=%d, msg=%s, body=%s", w.Code, resp.Code, resp.Msg, w.Body.String())
	}

	var created struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	}
	json.Unmarshal(resp.Data, &created)
	if created.Name != "番茄炒蛋" {
		t.Fatalf("expected name '番茄炒蛋', got '%s'", created.Name)
	}
	t.Logf("✓ CreateRecipe: id=%s, name=%s", created.Id, created.Name)

	recipeId := created.Id

	// 2. 获取菜谱列表
	w, resp = doRequestWithQuery(http.MethodGet, "/api/v1/admin/recipes")
	if w.Code != 200 || resp.Code != 0 {
		t.Fatalf("list recipes failed: %s", resp.Msg)
	}
	var recipeList []struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	}
	json.Unmarshal(resp.Data, &recipeList)
	if len(recipeList) != 1 {
		t.Fatalf("expected 1 recipe, got %d", len(recipeList))
	}
	t.Logf("✓ ListAllRecipes: count=%d", len(recipeList))

	// 3. 更新菜谱
	w, resp = doRequest(http.MethodPut, "/api/v1/admin/recipes/"+recipeId, map[string]interface{}{
		"slug":        "test-stir-fry-v2",
		"name":        "番茄炒蛋(改良版)",
		"cuisine":     "中餐",
		"description": "改良版家常菜",
		"prepMinutes": 5,
		"cookMinutes": 12,
		"servings":    3,
		"difficulty":  2,
		"imagePath":   "/images/tomato-egg-v2.jpg",
		"status":      "published",
		"allergens":   []string{"鸡蛋"},
		"cookware":    []string{"炒锅"},
		"ingredients": []map[string]interface{}{
			{"ingredientId": tomatoCreated.Id, "quantity": 300, "unit": "g", "required": true, "servingFactor": 1.0, "preparation": "切块"},
		},
		"steps": []map[string]interface{}{
			{"stepOrder": 1, "title": "准备", "instruction": "切食材", "durationSeconds": 300, "hasTimer": false},
		},
	})
	if w.Code != 200 || resp.Code != 0 {
		t.Fatalf("update recipe failed: status=%d, code=%d, msg=%s", w.Code, resp.Code, resp.Msg)
	}
	t.Logf("✓ UpdateRecipe: id=%s", recipeId)

	// 4. 获取单个菜谱(公开接口)
	w, resp = doRequestWithQuery(http.MethodGet, "/api/v1/recipes/"+recipeId)
	if w.Code != 200 || resp.Code != 0 {
		t.Fatalf("get recipe failed: %s", resp.Msg)
	}
	var detail struct {
		Id         string `json:"id"`
		Name       string `json:"name"`
		TotalMinutes int  `json:"totalMinutes"`
		Ingredients []struct {
			IngredientId   string  `json:"ingredientId"`
			Quantity       float64 `json:"quantity"`
			IngredientName string  `json:"ingredientName"`
		} `json:"ingredients"`
		Steps []struct {
			StepOrder int    `json:"stepOrder"`
			Title     string `json:"title"`
		} `json:"steps"`
	}
	json.Unmarshal(resp.Data, &detail)
	if detail.Name != "番茄炒蛋(改良版)" {
		t.Fatalf("expected updated name, got '%s'", detail.Name)
	}
	if detail.TotalMinutes != 17 {
		t.Fatalf("expected totalMinutes=17, got %d", detail.TotalMinutes)
	}
	if len(detail.Ingredients) != 1 {
		t.Fatalf("expected 1 ingredient, got %d", len(detail.Ingredients))
	}
	if len(detail.Steps) != 1 {
		t.Fatalf("expected 1 step, got %d", len(detail.Steps))
	}
	t.Logf("✓ GetPublishedRecipe: name=%s, totalMin=%d, ings=%d, steps=%d",
		detail.Name, detail.TotalMinutes, len(detail.Ingredients), len(detail.Steps))

	// 5. 删除菜谱
	w, resp = doRequestWithQuery(http.MethodDelete, "/api/v1/admin/recipes/"+recipeId)
	if w.Code != 200 || resp.Code != 0 {
		t.Fatalf("delete recipe failed: %s", resp.Msg)
	}
	t.Logf("✓ DeleteRecipe: id=%s", recipeId)

	// 6. 确认删除后查不到
	w, resp = doRequestWithQuery(http.MethodGet, "/api/v1/recipes/"+recipeId)
	if resp.Code == 0 {
		t.Fatal("expected 404 after delete, got success")
	}
	t.Log("✓ VerifyDelete: recipe not found after deletion")
}

func TestAdminIngredientCRUD(t *testing.T) {
	// 1. 创建食材
	w, resp := doRequest(http.MethodPost, "/api/v1/admin/ingredients", map[string]interface{}{
		"name":            "测试食材-番茄",
		"aliases":         []string{"西红柿", "Tomato"},
		"category":        "蔬菜",
		"defaultUnit":     "g",
		"isPantryStaple":  false,
		"substituteGroup": "",
		"storeSkuMapping": "",
	})
	if w.Code != 200 || resp.Code != 0 {
		t.Fatalf("create ingredient failed: %s, body: %s", resp.Msg, w.Body.String())
	}

	var created struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	}
	json.Unmarshal(resp.Data, &created)
	t.Logf("✓ CreateIngredient: id=%s, name=%s", created.Id, created.Name)

	ingredientId := created.Id

	// 2. 获取食材列表
	w, resp = doRequestWithQuery(http.MethodGet, "/api/v1/admin/ingredients")
	if w.Code != 200 || resp.Code != 0 {
		t.Fatalf("list ingredients failed: %s", resp.Msg)
	}
	var list []struct {
		Id   string   `json:"id"`
		Name string   `json:"name"`
	}
	json.Unmarshal(resp.Data, &list)
	if len(list) < 1 {
		t.Fatalf("expected at least 1 ingredient, got %d", len(list))
	}
	t.Logf("✓ ListAllIngredients: count=%d", len(list))

	// 3. 更新食材
	w, resp = doRequest(http.MethodPut, "/api/v1/admin/ingredients/"+ingredientId, map[string]interface{}{
		"name":            "番茄(有机)",
		"aliases":         []string{"有机番茄"},
		"category":        "蔬菜",
		"defaultUnit":     "g",
		"isPantryStaple":  false,
		"substituteGroup": "",
	})
	if w.Code != 200 || resp.Code != 0 {
		t.Fatalf("update ingredient failed: %s", resp.Msg)
	}
	t.Logf("✓ UpdateIngredient: id=%s", ingredientId)

	// 4. 删除食材
	w, resp = doRequestWithQuery(http.MethodDelete, "/api/v1/admin/ingredients/"+ingredientId)
	if w.Code != 200 || resp.Code != 0 {
		t.Fatalf("delete ingredient failed: %s", resp.Msg)
	}
	t.Logf("✓ DeleteIngredient: id=%s", ingredientId)
}

// ==================== Recipe 公开接口 + 筛选测试 ====================

func TestRecipeWithIngredients(t *testing.T) {
	// 先创建食材
	w, resp := doRequest(http.MethodPost, "/api/v1/admin/ingredients", map[string]interface{}{
		"name":           "鸡蛋-食谱测试",
		"category":       "蛋类",
		"defaultUnit":    "个",
		"isPantryStaple": true,
	})
	if resp.Code != 0 {
		t.Fatalf("create ingredient failed: %s", resp.Msg)
	}
	var ingCreated struct {
		Id string `json:"id"`
	}
	json.Unmarshal(resp.Data, &ingCreated)

	// 创建带食材的菜谱
	w, resp = doRequest(http.MethodPost, "/api/v1/admin/recipes", map[string]interface{}{
		"slug":        "scrambled-eggs",
		"name":        "炒鸡蛋",
		"cuisine":     "中餐",
		"description": "简单炒蛋",
		"prepMinutes": 2,
		"cookMinutes": 5,
		"servings":    1,
		"difficulty":  1,
		"status":      "published",
		"ingredients": []map[string]interface{}{
			{"ingredientId": ingCreated.Id, "quantity": 3, "unit": "个", "required": true, "servingFactor": 1.0},
		},
		"steps": []map[string]interface{}{
			{"stepOrder": 1, "title": "打蛋", "instruction": "搅拌均匀", "durationSeconds": 60, "hasTimer": false},
			{"stepOrder": 2, "title": "炒", "instruction": "热油快炒", "durationSeconds": 180, "hasTimer": true},
		},
	})
	if resp.Code != 0 {
		t.Fatalf("create recipe failed: %s", resp.Msg)
	}
	var recipeCreated struct {
		Id string `json:"id"`
	}
	json.Unmarshal(resp.Data, &recipeCreated)

	// 通过公开接口获取
	w, resp = doRequestWithQuery(http.MethodGet, "/api/v1/recipes/"+recipeCreated.Id)
	if w.Code != 200 || resp.Code != 0 {
		t.Fatalf("get recipe failed: %s", resp.Msg)
	}

	var recipe struct {
		Id           string `json:"id"`
		Name         string `json:"name"`
		TotalMinutes int    `json:"totalMinutes"`
		Ingredients  []struct {
			IngredientId   string  `json:"ingredientId"`
			IngredientName string  `json:"ingredientName"`
			Quantity       float64 `json:"quantity"`
			Unit           string  `json:"unit"`
			IsPantryStaple bool    `json:"isPantryStaple"`
		} `json:"ingredients"`
		Steps []struct {
			StepOrder       int    `json:"stepOrder"`
			Title           string `json:"title"`
			DurationSeconds int    `json:"durationSeconds"`
			HasTimer        bool   `json:"hasTimer"`
		} `json:"steps"`
	}
	json.Unmarshal(resp.Data, &recipe)

	if recipe.Name != "炒鸡蛋" {
		t.Fatalf("expected name '炒鸡蛋', got '%s'", recipe.Name)
	}
	if recipe.TotalMinutes != 7 {
		t.Fatalf("expected totalMinutes=7, got %d", recipe.TotalMinutes)
	}
	if len(recipe.Ingredients) != 1 {
		t.Fatalf("expected 1 ingredient, got %d", len(recipe.Ingredients))
	}
	if recipe.Ingredients[0].IngredientName != "鸡蛋-食谱测试" {
		t.Fatalf("expected ingredientName '鸡蛋-食谱测试', got '%s'", recipe.Ingredients[0].IngredientName)
	}
	if recipe.Ingredients[0].Quantity != 3 {
		t.Fatalf("expected quantity=3, got %f", recipe.Ingredients[0].Quantity)
	}
	if len(recipe.Steps) != 2 {
		t.Fatalf("expected 2 steps, got %d", len(recipe.Steps))
	}
	if recipe.Steps[1].HasTimer != true {
		t.Fatal("expected step 2 hasTimer=true")
	}
	t.Logf("✓ RecipeDetail with ingredients & steps: name=%s, ings=%d, steps=%d",
		recipe.Name, len(recipe.Ingredients), len(recipe.Steps))
}

func TestRecipeFilterByCuisine(t *testing.T) {
	// 创建菜谱(已在 TestRecipeWithIngredients 中创建了"中餐"菜谱)
	// 创建另一个菜谱
	w, resp := doRequest(http.MethodPost, "/api/v1/admin/recipes", map[string]interface{}{
		"slug":        "pasta",
		"name":        "意大利面",
		"cuisine":     "意大利",
		"description": "pasta",
		"prepMinutes": 10,
		"cookMinutes": 20,
		"servings":    2,
		"difficulty":  2,
		"status":      "published",
	})
	if resp.Code != 0 {
		t.Fatalf("create recipe failed: %s", resp.Msg)
	}

	// 按菜系筛选
	w, resp = doRequestWithQuery(http.MethodGet, "/api/v1/recipes?cuisine=中餐")
	if w.Code != 200 || resp.Code != 0 {
		t.Fatalf("filter by cuisine failed: %s", resp.Msg)
	}
	var list []struct {
		Name    string `json:"name"`
		Cuisine string `json:"cuisine"`
	}
	json.Unmarshal(resp.Data, &list)
	for _, r := range list {
		if r.Cuisine != "中餐" {
			t.Fatalf("expected cuisine '中餐', got '%s' for recipe '%s'", r.Cuisine, r.Name)
		}
	}
	t.Logf("✓ FilterByCuisine: found %d recipes with cuisine=中餐", len(list))

	// 按时间筛选
	w, resp = doRequestWithQuery(http.MethodGet, "/api/v1/recipes?maxMinutes=10")
	if w.Code != 200 || resp.Code != 0 {
		t.Fatalf("filter by maxMinutes failed: %s", resp.Msg)
	}
	var filtered []struct {
		Name        string `json:"name"`
		TotalMinutes int   `json:"totalMinutes"`
	}
	json.Unmarshal(resp.Data, &filtered)
	for _, r := range filtered {
		if r.TotalMinutes > 10 {
			t.Fatalf("expected totalMinutes<=10, got %d for '%s'", r.TotalMinutes, r.Name)
		}
	}
	t.Logf("✓ FilterByMaxMinutes: found %d recipes with totalMinutes<=10", len(filtered))
}

// ==================== Workflow 模块测试 ====================

func TestCookingSessionFlow(t *testing.T) {
	// 先创建一个菜谱
	w, resp := doRequest(http.MethodPost, "/api/v1/admin/recipes", map[string]interface{}{
		"slug":        "test-session-recipe",
		"name":        "测试菜谱",
		"cuisine":     "中餐",
		"description": "用于会话测试",
		"prepMinutes": 5,
		"cookMinutes": 10,
		"servings":    2,
		"difficulty":  1,
		"status":      "published",
		"steps": []map[string]interface{}{
			{"stepOrder": 1, "title": "步骤1", "instruction": "操作1", "durationSeconds": 60, "hasTimer": false},
			{"stepOrder": 2, "title": "步骤2", "instruction": "操作2", "durationSeconds": 120, "hasTimer": true},
		},
	})
	if resp.Code != 0 {
		t.Fatalf("create recipe failed: %s", resp.Msg)
	}
	var recipeCreated struct {
		Id string `json:"id"`
	}
	json.Unmarshal(resp.Data, &recipeCreated)

	// 创建做饭会话
	w, resp = doRequest(http.MethodPost, "/api/v1/cooking-sessions", map[string]interface{}{
		"recipeId": recipeCreated.Id,
	})
	if w.Code != 200 || resp.Code != 0 {
		t.Fatalf("create session failed: status=%d, code=%d, msg=%s, body=%s", w.Code, resp.Code, resp.Msg, w.Body.String())
	}

	var session struct {
		Id               string `json:"id"`
		RecipeId         string `json:"recipeId"`
		CurrentStepOrder int    `json:"currentStepOrder"`
		State            string `json:"state"`
		StartedAt        string `json:"startedAt"`
	}
	json.Unmarshal(resp.Data, &session)

	if session.RecipeId != recipeCreated.Id {
		t.Fatalf("expected recipeId=%s, got %s", recipeCreated.Id, session.RecipeId)
	}
	if session.CurrentStepOrder != 1 {
		t.Fatalf("expected currentStepOrder=1, got %d", session.CurrentStepOrder)
	}
	if session.State != "active" {
		t.Fatalf("expected state='active', got '%s'", session.State)
	}
	t.Logf("✓ CreateCookingSession: id=%s, recipeId=%s", session.Id, session.RecipeId)

	// 更新会话进度
	w, resp = doRequest(http.MethodPut, "/api/v1/cooking-sessions/"+session.Id, map[string]interface{}{
		"currentStepOrder": 2,
		"state":            "active",
	})
	if w.Code != 200 || resp.Code != 0 {
		t.Fatalf("update session failed: %s", resp.Msg)
	}

	var updated struct {
		CurrentStepOrder int    `json:"currentStepOrder"`
		State            string `json:"state"`
	}
	json.Unmarshal(resp.Data, &updated)
	if updated.CurrentStepOrder != 2 {
		t.Fatalf("expected currentStepOrder=2, got %d", updated.CurrentStepOrder)
	}
	t.Logf("✓ UpdateCookingSession: step=%d, state=%s", updated.CurrentStepOrder, updated.State)
}

func TestSubmitFeedback(t *testing.T) {
	// 先创建菜谱
	w, resp := doRequest(http.MethodPost, "/api/v1/admin/recipes", map[string]interface{}{
		"slug":        "feedback-recipe",
		"name":        "反馈测试菜谱",
		"cuisine":     "中餐",
		"description": "用于反馈测试",
		"prepMinutes": 5,
		"cookMinutes": 10,
		"servings":    2,
		"difficulty":  1,
		"status":      "published",
	})
	if resp.Code != 0 {
		t.Fatalf("create recipe failed: %s", resp.Msg)
	}
	var recipeCreated struct {
		Id string `json:"id"`
	}
	json.Unmarshal(resp.Data, &recipeCreated)

	// 提交反馈
	w, resp = doRequest(http.MethodPost, "/api/v1/feedback", map[string]interface{}{
		"recipeId": recipeCreated.Id,
		"outcome":  "success",
		"tags":     []string{"好吃", "简单"},
		"comment":  "很好吃的家常菜",
	})
	if w.Code != 200 || resp.Code != 0 {
		t.Fatalf("submit feedback failed: status=%d, code=%d, msg=%s, body=%s", w.Code, resp.Code, resp.Msg, w.Body.String())
	}

	var feedback struct {
		Id        string   `json:"id"`
		RecipeId  string   `json:"recipeId"`
		Outcome   string   `json:"outcome"`
		Tags      []string `json:"tags"`
		Comment   string   `json:"comment"`
		CreatedAt string   `json:"createdAt"`
	}
	json.Unmarshal(resp.Data, &feedback)

	if feedback.RecipeId != recipeCreated.Id {
		t.Fatalf("expected recipeId=%s, got %s", recipeCreated.Id, feedback.RecipeId)
	}
	if feedback.Outcome != "success" {
		t.Fatalf("expected outcome='success', got '%s'", feedback.Outcome)
	}
	if len(feedback.Tags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(feedback.Tags))
	}
	if feedback.Comment != "很好吃的家常菜" {
		t.Fatalf("expected comment '很好吃的家常菜', got '%s'", feedback.Comment)
	}
	t.Logf("✓ SubmitFeedback: id=%s, outcome=%s, tags=%v", feedback.Id, feedback.Outcome, feedback.Tags)
}

func TestPlanFlow(t *testing.T) {
	// 先创建菜谱
	w, resp := doRequest(http.MethodPost, "/api/v1/admin/recipes", map[string]interface{}{
		"slug":        "plan-recipe-1",
		"name":        "计划菜谱1",
		"cuisine":     "中餐",
		"description": "用于计划测试",
		"prepMinutes": 5,
		"cookMinutes": 10,
		"servings":    2,
		"difficulty":  1,
		"status":      "published",
	})
	if resp.Code != 0 {
		t.Fatalf("create recipe failed: %s", resp.Msg)
	}
	var recipe1 struct {
		Id string `json:"id"`
	}
	json.Unmarshal(resp.Data, &recipe1)

	// 保存计划
	w, resp = doRequest(http.MethodPost, "/api/v1/plans", map[string]interface{}{
		"planDate": "2026-09-05",
		"items": []map[string]interface{}{
			{"recipeId": recipe1.Id, "servings": 2, "sortOrder": 0},
		},
	})
	if w.Code != 200 || resp.Code != 0 {
		t.Fatalf("save plan failed: status=%d, code=%d, msg=%s, body=%s", w.Code, resp.Code, resp.Msg, w.Body.String())
	}

	var plan struct {
		Id       string `json:"id"`
		PlanDate string `json:"planDate"`
		Items    []struct {
			Id         string `json:"id"`
			RecipeId   string `json:"recipeId"`
			RecipeName string `json:"recipeName"`
			Servings   int    `json:"servings"`
		} `json:"items"`
		CombinedIngredients []struct {
			IngredientId   string  `json:"ingredientId"`
			IngredientName string  `json:"ingredientName"`
			Quantity       float64 `json:"quantity"`
		} `json:"combinedIngredients"`
	}
	json.Unmarshal(resp.Data, &plan)

	if plan.PlanDate != "2026-09-05" {
		t.Fatalf("expected planDate='2026-09-05', got '%s'", plan.PlanDate)
	}
	if len(plan.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(plan.Items))
	}
	t.Logf("✓ SavePlan: id=%s, date=%s, items=%d", plan.Id, plan.PlanDate, len(plan.Items))

	// 查询计划
	w, resp = doRequestWithQuery(http.MethodGet, "/api/v1/plans?from=2026-09-01&to=2026-09-30")
	if w.Code != 200 || resp.Code != 0 {
		t.Fatalf("list plans failed: %s", resp.Msg)
	}
	var planList []struct {
		Id       string `json:"id"`
		PlanDate string `json:"planDate"`
	}
	json.Unmarshal(resp.Data, &planList)
	if len(planList) != 1 {
		t.Fatalf("expected 1 plan, got %d", len(planList))
	}
	t.Logf("✓ ListPlans: count=%d", len(planList))
}

// ==================== API 契约格式验证 ====================

func TestResponseFormatContract(t *testing.T) {
	// 验证所有响应都遵循 vo.Result 格式
	endpoints := []struct {
		method string
		path   string
		body   interface{}
	}{
		{http.MethodGet, "/healthz", nil},
		{http.MethodGet, "/readyz", nil},
		{http.MethodPost, "/api/v1/auth/register", map[string]interface{}{
			"email": "contract@example.com", "password": "test12345", "displayName": "契约测试",
		}},
		{http.MethodGet, "/api/v1/recipes", nil},
		{http.MethodGet, "/api/v1/ingredients", nil},
		{http.MethodPost, "/api/v1/recommendations/tonight", map[string]interface{}{
			"servings": 2, "availableMinutes": 60,
		}},
		{http.MethodGet, "/api/v1/admin/stats", nil},
		{http.MethodGet, "/api/v1/admin/recipes", nil},
		{http.MethodGet, "/api/v1/admin/ingredients", nil},
	}

	for _, ep := range endpoints {
		var w *httptest.ResponseRecorder
		var resp apiResponse
		if ep.body != nil {
			w, resp = doRequest(ep.method, ep.path, ep.body)
		} else {
			w, resp = doRequestWithQuery(ep.method, ep.path)
		}

		// 验证 HTTP 状态码
		if w.Code != 200 && w.Code != 201 {
			t.Errorf("[%s %s] expected HTTP 200/201, got %d", ep.method, ep.path, w.Code)
		}

		// 验证响应结构: 必须有 code, msg, timeStamp
		if w.Body.Len() == 0 {
			t.Errorf("[%s %s] empty response body", ep.method, ep.path)
		}
		if resp.TimeStamp == 0 {
			t.Errorf("[%s %s] missing timeStamp", ep.method, ep.path)
		}

		t.Logf("✓ Contract: %s %s → HTTP %d, code=%d", ep.method, ep.path, w.Code, resp.Code)
	}
}

// ==================== 前端对接兼容性测试 ====================

func TestFrontendAuthFlow(t *testing.T) {
	// 模拟前端完整的认证流程: 注册 → 登录 → 刷新
	// 1. 注册
	_, regResp := doRequest(http.MethodPost, "/api/v1/auth/register", map[string]interface{}{
		"email":       "frontend@test.com",
		"password":    "mypass123",
		"displayName": "前端用户",
	})
	if regResp.Code != 0 {
		t.Fatalf("register failed: %s", regResp.Msg)
	}

	// 2. 登录
	_, loginResp := doRequest(http.MethodPost, "/api/v1/auth/login", map[string]interface{}{
		"email":    "frontend@test.com",
		"password": "mypass123",
	})
	if loginResp.Code != 0 {
		t.Fatalf("login failed: %s", loginResp.Msg)
	}

	var loginData struct {
		User struct {
			Id          string `json:"id"`
			Email       string `json:"email"`
			DisplayName string `json:"displayName"`
			IsAdmin     bool   `json:"isAdmin"`
		} `json:"user"`
		AccessToken               string `json:"accessToken"`
		RefreshToken              string `json:"refreshToken"`
		AccessTokenExpiresInSeconds int    `json:"accessTokenExpiresInSeconds"`
	}
	json.Unmarshal(loginResp.Data, &loginData)

	// 验证前端需要的所有字段都存在
	if loginData.User.Id == "" {
		t.Fatal("missing user.id")
	}
	if loginData.User.Email != "frontend@test.com" {
		t.Fatalf("wrong email: %s", loginData.User.Email)
	}
	if loginData.AccessToken == "" {
		t.Fatal("missing accessToken")
	}
	if loginData.RefreshToken == "" {
		t.Fatal("missing refreshToken")
	}
	if loginData.AccessTokenExpiresInSeconds == 0 {
		t.Fatal("missing accessTokenExpiresInSeconds")
	}

	// 3. 刷新令牌
	_, refreshResp := doRequest(http.MethodPost, "/api/v1/auth/refresh", map[string]interface{}{
		"refreshToken": loginData.RefreshToken,
	})
	if refreshResp.Code != 0 {
		t.Fatalf("refresh failed: %s", refreshResp.Msg)
	}

	var refreshData struct {
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
	}
	json.Unmarshal(refreshResp.Data, &refreshData)
	if refreshData.AccessToken == "" {
		t.Fatal("refresh: missing accessToken")
	}
	if refreshData.RefreshToken == "" {
		t.Fatal("refresh: missing refreshToken")
	}
	// 新令牌应该与旧的不同
	if refreshData.RefreshToken == loginData.RefreshToken {
		t.Fatal("refresh: new refreshToken should differ from old")
	}

	t.Logf("✓ FrontendAuthFlow: register→login→refresh all passed")
}

func TestFrontendRecipeLoadFlow(t *testing.T) {
	// 模拟前端加载菜谱流程: 获取列表 → 获取详情 → 获取推荐
	// 1. 获取菜谱列表
	_, listResp := doRequestWithQuery(http.MethodGet, "/api/v1/recipes?limit=20")
	if listResp.Code != 0 {
		t.Fatalf("list recipes failed: %s", listResp.Msg)
	}
	t.Log("✓ FrontendRecipeLoad: list loaded")

	// 2. 获取食材列表
	_, ingResp := doRequestWithQuery(http.MethodGet, "/api/v1/ingredients")
	if ingResp.Code != 0 {
		t.Fatalf("list ingredients failed: %s", ingResp.Msg)
	}
	t.Log("✓ FrontendRecipeLoad: ingredients loaded")

	// 3. 获取推荐
	_, recResp := doRequest(http.MethodPost, "/api/v1/recommendations/tonight", map[string]interface{}{
		"servings":         2,
		"availableMinutes": 60,
		"cuisines":         []string{},
		"allergies":        []string{},
		"pantryIngredientIds": []string{},
		"cookware":         []string{},
		"recentRecipeIds":  []string{},
	})
	if recResp.Code != 0 {
		t.Fatalf("recommend failed: %s", recResp.Msg)
	}
	t.Log("✓ FrontendRecipeLoad: recommendations loaded")
}
