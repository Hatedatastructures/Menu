const BaseUrl = process.env.MENU_API_BASE_URL ?? "http://127.0.0.1:8080";
const AdminEmail = process.env.MENU_ADMIN_EMAIL ?? "codex-admin-20260901@example.com";
const AdminPassword = process.env.MENU_ADMIN_PASSWORD ?? "Menu-Cook-Local-2026!";

function Assert(Condition, Message) {
    if (!Condition) {
        throw new Error(Message);
    }
}

async function Request(Method, Path, Body = undefined, Token = undefined, ExtraHeaders = {}) {
    const Headers = {Accept: "application/json", ...ExtraHeaders};
    if (Body !== undefined) {
        Headers["Content-Type"] = "application/json";
    }
    if (Token !== undefined) {
        Headers.Authorization = `Bearer ${Token}`;
    }
    const Response = await fetch(`${BaseUrl}${Path}`, {
        method: Method,
        headers: Headers,
        body: Body === undefined ? undefined : JSON.stringify(Body),
    });
    const Text = await Response.text();
    let Json = undefined;
    if (Text.length > 0) {
        try {
            Json = JSON.parse(Text);
        } catch {
            Json = undefined;
        }
    }
    return {Status: Response.status, Headers: Response.headers, Text, Json};
}

function Status(Response, Expected, Name) {
    Assert(Response.Status === Expected, `${Name}: expected ${Expected}, got ${Response.Status}: ${Response.Text}`);
}

function Header(Response, Name) {
    return Response.Headers.get(Name) ?? "";
}

const Result = {};
const Health = await Request("GET", "/healthz");
Status(Health, 200, "healthz");
Assert(Health.Json?.status === "ok", "healthz status is not ok");
Result.Health = true;

const Ready = await Request("GET", "/readyz");
Status(Ready, 200, "readyz");
Assert(Ready.Json?.status === "ready", "readyz status is not ready");
Result.Ready = true;

const RecipesResponse = await Request("GET", "/api/v1/recipes?limit=3");
Status(RecipesResponse, 200, "recipes");
Assert(Array.isArray(RecipesResponse.Json) && RecipesResponse.Json.length > 0, "recipes is empty");
const Recipe = RecipesResponse.Json[0];
Assert(Array.isArray(Recipe.ingredients) && Recipe.ingredients.length > 0, "recipe ingredients are missing");
const RecipeIngredient = Recipe.ingredients[0];
Assert(typeof RecipeIngredient.ingredientName === "string" && RecipeIngredient.ingredientName.length > 0,
    "recipe ingredientName is missing");
Assert(typeof RecipeIngredient.category === "string" && typeof RecipeIngredient.defaultUnit === "string",
    "recipe ingredient display fields are missing");
Result.Recipes = RecipesResponse.Json.length;

const Detail = await Request("GET", `/api/v1/recipes/${encodeURIComponent(Recipe.id)}`);
Status(Detail, 200, "recipe detail");
Assert(Detail.Json?.id === Recipe.id && Array.isArray(Detail.Json.steps), "recipe detail is incomplete");
Result.RecipeDetail = true;

const IngredientsResponse = await Request("GET", "/api/v1/ingredients");
Status(IngredientsResponse, 200, "ingredients");
Assert(Array.isArray(IngredientsResponse.Json) && IngredientsResponse.Json.length > 0, "ingredients is empty");
const Ingredient = IngredientsResponse.Json[0];
Assert(typeof Ingredient.name === "string" && typeof Ingredient.category === "string" &&
    typeof Ingredient.defaultUnit === "string" && typeof Ingredient.isPantryStaple === "boolean",
    "ingredient display fields are missing");
Result.Ingredients = IngredientsResponse.Json.length;

const Recommendation = await Request("POST", "/api/v1/recommendations/tonight", {
    servings: 2,
    availableMinutes: 35,
    cuisines: ["中餐", "西餐", "日系"],
    pantryIngredientIds: ["ingredient.rice"],
    cookware: ["炒锅", "平底锅"],
});
Status(Recommendation, 200, "recommendations");
Assert(Array.isArray(Recommendation.Json) && Recommendation.Json.length > 0 && Recommendation.Json.length <= 3,
    "recommendations count is invalid");
Result.Recommendations = Recommendation.Json.length;

let AdminLogin = await Request("POST", "/api/v1/auth/login", {
    email: AdminEmail,
    password: AdminPassword,
});
if (AdminLogin.Status === 401) {
    AdminLogin = await Request("POST", "/api/v1/auth/register", {
        email: AdminEmail,
        password: AdminPassword,
        displayName: "本地管理验收",
    });
    Status(AdminLogin, 201, "admin register");
} else {
    Status(AdminLogin, 200, "admin login");
}
const AdminAccessToken = AdminLogin.Json?.accessToken;
const AdminRefreshToken = AdminLogin.Json?.refreshToken;
Assert(typeof AdminAccessToken === "string" && typeof AdminRefreshToken === "string", "admin tokens are missing");
Result.AdminLogin = true;

const Refresh = await Request("POST", "/api/v1/auth/refresh", {refreshToken: AdminRefreshToken});
Status(Refresh, 200, "refresh");
Assert(Refresh.Json?.refreshToken && Refresh.Json.refreshToken !== AdminRefreshToken, "refresh token did not rotate");
const ReusedRefresh = await Request("POST", "/api/v1/auth/refresh", {refreshToken: AdminRefreshToken});
Status(ReusedRefresh, 401, "reused refresh token");
Result.RefreshRotation = true;

const UserEmail = `api-e2e-${Date.now()}@example.com`;
const UserPassword = "Menu-Cook-Api-2026!";
const Registered = await Request("POST", "/api/v1/auth/register", {
    email: UserEmail,
    password: UserPassword,
    displayName: "厨房验收",
});
Status(Registered, 201, "register");
const UserAccessToken = Registered.Json?.accessToken;
Assert(typeof UserAccessToken === "string", "user access token is missing");
Result.Register = true;

const Plan = await Request("POST", "/api/v1/plans", {
    planDate: "2026-09-01",
    items: [
        {recipeId: Recipe.id, servings: 2, sortOrder: 0},
    ],
}, UserAccessToken);
Status(Plan, 201, "create plan");
Assert(Array.isArray(Plan.Json?.combinedIngredients), "plan combined ingredients are missing");
const Plans = await Request("GET", "/api/v1/plans?from=2026-09-01&to=2026-09-07", undefined, UserAccessToken);
Status(Plans, 200, "list plans");
Assert(Array.isArray(Plans.Json) && Plans.Json.length > 0, "plan list is empty");
Result.Plans = true;

const Session = await Request("POST", "/api/v1/cooking-sessions", {recipeId: Recipe.id}, UserAccessToken);
Status(Session, 201, "create cooking session");
const SessionId = Session.Json?.id;
Assert(typeof SessionId === "string", "cooking session id is missing");
const UpdatedSession = await Request("PATCH", `/api/v1/cooking-sessions/${encodeURIComponent(SessionId)}`, {
    currentStepOrder: 1,
    state: "paused",
}, UserAccessToken);
Status(UpdatedSession, 200, "update cooking session");
Assert(UpdatedSession.Json?.state === "paused", "cooking session state did not update");
const Feedback = await Request("POST", "/api/v1/feedback", {
    recipeId: Recipe.id,
    outcome: "made",
    tags: ["下次还想做"],
    comment: "真实 HTTP 验收",
}, UserAccessToken);
Status(Feedback, 201, "feedback");
Result.CookingWorkflow = true;

const AdminRecipes = await Request("GET", "/api/v1/admin/recipes", undefined, AdminAccessToken);
Status(AdminRecipes, 200, "admin recipes");
const Stats = await Request("GET", "/api/v1/admin/stats", undefined, AdminAccessToken);
Status(Stats, 200, "admin stats");
const IngredientId = `ingredient.api_e2e_${Date.now()}`;
const IngredientBody = {
    id: IngredientId,
    name: "接口验收香草",
    aliases: ["验收香草"],
    category: "香料",
    defaultUnit: "g",
    isPantryStaple: false,
    substituteGroup: "",
    storeSkuMapping: "",
};
const CreatedIngredient = await Request("POST", "/api/v1/admin/ingredients", IngredientBody, AdminAccessToken);
Status(CreatedIngredient, 201, "admin ingredient create");
const UpdatedIngredient = await Request("PATCH", `/api/v1/admin/ingredients/${encodeURIComponent(IngredientId)}`, {
    ...IngredientBody,
    name: "接口验收香草叶",
    isPantryStaple: true,
    storeSkuMapping: "sku.future",
}, AdminAccessToken);
Status(UpdatedIngredient, 200, "admin ingredient update");
Assert(UpdatedIngredient.Json?.name === "接口验收香草叶", "admin ingredient update did not persist");
const DeletedIngredient = await Request("DELETE", `/api/v1/admin/ingredients/${encodeURIComponent(IngredientId)}`, undefined, AdminAccessToken);
Status(DeletedIngredient, 204, "admin ingredient delete");
Result.AdminIngredientCrud = true;

const Preflight = await Request("OPTIONS", "/api/v1/admin/recipes", undefined, undefined, {
    Origin: "http://127.0.0.1:5173",
    "Access-Control-Request-Method": "PATCH",
    "Access-Control-Request-Headers": "Authorization, Content-Type",
});
Status(Preflight, 204, "allowed cors preflight");
Assert(Header(Preflight, "access-control-allow-methods").includes("PATCH") &&
    Header(Preflight, "access-control-allow-methods").includes("DELETE"), "cors methods are incomplete");
Assert(Header(Preflight, "access-control-allow-headers").toLowerCase().includes("authorization"),
    "cors authorization header is missing");
Assert(Header(Preflight, "vary").toLowerCase().includes("origin") &&
    Header(Preflight, "vary").toLowerCase().includes("access-control-request-method") &&
    Header(Preflight, "vary").toLowerCase().includes("access-control-request-headers"),
    "cors vary headers are incomplete");
const DeniedPreflight = await Request("OPTIONS", "/api/v1/admin/recipes", undefined, undefined, {
    Origin: "https://untrusted.example",
    "Access-Control-Request-Method": "PATCH",
    "Access-Control-Request-Headers": "Authorization",
});
Status(DeniedPreflight, 403, "denied cors preflight");
Result.Cors = true;

const LongPassword = "x".repeat(129);
const OversizedLogin = await Request("POST", "/api/v1/auth/login", {
    email: UserEmail,
    password: LongPassword,
});
Status(OversizedLogin, 401, "oversized login password");
const InvalidDisplayName = await Request("POST", "/api/v1/auth/register", {
    email: `invalid-name-${Date.now()}@example.com`,
    password: UserPassword,
    displayName: "   ",
});
Status(InvalidDisplayName, 400, "blank display name");
Result.SecurityValidation = true;

console.log(JSON.stringify(Result));
