#include <gtest/gtest.h>

#include <Application/AuthService.hpp>
#include <Composition/ServerApplication.hpp>
#include <Runtime/ServerConfiguration.hpp>

#include <boost/asio/connect.hpp>
#include <boost/asio/io_context.hpp>
#include <boost/asio/ip/tcp.hpp>
#include <boost/beast/core/flat_buffer.hpp>
#include <boost/beast/http.hpp>
#include <boost/json/parse.hpp>

#include <memory>
#include <chrono>
#include <filesystem>
#include <ranges>
#include <string>
#include <string_view>
#include <thread>
#include <utility>

namespace {

class ApiTestFixture : public ::testing::Test {
protected:
    void SetUp() override {
        const auto Suffix = std::to_string(
            std::chrono::steady_clock::now().time_since_epoch().count());
        Root = std::filesystem::current_path() / ("MenuApiTest-" + Suffix);
        Menu::Runtime::ServerConfiguration Configuration;
        Configuration.Http.Address = "127.0.0.1";
        Configuration.Http.Port = 0;
        Configuration.CorsOrigins = {"https://admin.menu.local"};
        Configuration.DatabasePath = Root / "Data" / "menu.db";
        Configuration.MediaDirectory = Root / "Media";

        auto ApplicationResult = Menu::Server::ServerApplication::Create(
            std::move(Configuration));
        ASSERT_TRUE(ApplicationResult.HasValue())
            << ApplicationResult.ErrorValue().MessageValue();
        Application = std::move(ApplicationResult).Value();
        ASSERT_TRUE(Application->Start().HasValue());
        IoThread = std::thread([this]() { Application->IoContext().run(); });
    }

    void TearDown() override {
        if (Application) {
            Application->Stop();
            Application->IoContext().stop();
        }
        if (IoThread.joinable()) {
            IoThread.join();
        }
        Application.reset();
        std::error_code Error;
        std::filesystem::remove_all(Root, Error);
    }

    [[nodiscard]] boost::beast::http::response<boost::beast::http::string_body> Request(
        boost::beast::http::verb Method,
        std::string_view Target,
        std::string Body = {},
        std::string Origin = {},
        std::string Authorization = {},
        std::string RequestedMethod = {},
        std::string RequestedHeaders = {}) const {
        boost::asio::io_context ClientContext;
        boost::asio::ip::tcp::socket Socket(ClientContext);
        Socket.connect({boost::asio::ip::make_address("127.0.0.1"), Application->LocalPort()});
        boost::beast::http::request<boost::beast::http::string_body> RequestValue(
            Method, Target, 11);
        RequestValue.set(boost::beast::http::field::host, "localhost");
        if (!Origin.empty()) {
            RequestValue.set(boost::beast::http::field::origin, Origin);
        }
        if (!Authorization.empty()) {
            RequestValue.set(boost::beast::http::field::authorization, Authorization);
        }
        if (!RequestedMethod.empty()) {
            RequestValue.set("Access-Control-Request-Method", RequestedMethod);
        }
        if (!RequestedHeaders.empty()) {
            RequestValue.set("Access-Control-Request-Headers", RequestedHeaders);
        }
        if (!Body.empty()) {
            RequestValue.set(boost::beast::http::field::content_type, "application/json");
            RequestValue.body() = std::move(Body);
            RequestValue.prepare_payload();
        }
        boost::beast::http::write(Socket, RequestValue);
        boost::beast::flat_buffer Buffer;
        boost::beast::http::response<boost::beast::http::string_body> Response;
        boost::beast::http::read(Socket, Buffer, Response);
        Socket.shutdown(boost::asio::ip::tcp::socket::shutdown_both);
        return Response;
    }

    [[nodiscard]] static boost::json::value ParseJson(const std::string& Body) {
        boost::system::error_code Error;
        const auto Value = boost::json::parse(Body, Error);
        EXPECT_FALSE(Error);
        return Value;
    }

    std::thread IoThread;
    std::filesystem::path Root;
    std::unique_ptr<Menu::Server::ServerApplication> Application;
};

}  // namespace

TEST_F(ApiTestFixture, ReturnsHealthAndReadinessFromRealListener) {
    const auto Health = Request(boost::beast::http::verb::get, "/healthz");
    const auto Ready = Request(boost::beast::http::verb::get, "/readyz");

    EXPECT_EQ(static_cast<int>(Health.result()), 200);
    EXPECT_EQ(ParseJson(Health.body()).as_object().at("status").as_string(), "ok");
    EXPECT_EQ(static_cast<int>(Ready.result()), 200);
    EXPECT_EQ(ParseJson(Ready.body()).as_object().at("status").as_string(), "ready");
}

TEST_F(ApiTestFixture, ListsRecipesWithStructuredIngredientsAndSteps) {
    const auto Response = Request(
        boost::beast::http::verb::get, "/api/v1/recipes?cuisine=中餐&limit=3");

    ASSERT_EQ(static_cast<int>(Response.result()), 200);
    const auto Recipes = ParseJson(Response.body()).as_array();
    ASSERT_EQ(Recipes.size(), 3U);
    EXPECT_EQ(Recipes[0].as_object().at("status").as_string(), "published");
    EXPECT_TRUE(Recipes[0].as_object().at("ingredients").is_array());
    EXPECT_TRUE(Recipes[0].as_object().at("steps").is_array());
    const auto Ingredient = Recipes[0].as_object().at("ingredients").as_array()[0].as_object();
    EXPECT_TRUE(Ingredient.contains("ingredientName"));
    EXPECT_TRUE(Ingredient.contains("category"));
    EXPECT_TRUE(Ingredient.contains("defaultUnit"));
    EXPECT_TRUE(Ingredient.contains("isPantryStaple"));
}

TEST_F(ApiTestFixture, ListsIngredientsWithDisplayFields) {
    const auto Response = Request(
        boost::beast::http::verb::get, "/api/v1/ingredients");

    ASSERT_EQ(static_cast<int>(Response.result()), 200);
    const auto Ingredients = ParseJson(Response.body()).as_array();
    ASSERT_FALSE(Ingredients.empty());
    const auto Ingredient = Ingredients[0].as_object();
    EXPECT_TRUE(Ingredient.at("id").is_string());
    EXPECT_FALSE(Ingredient.at("name").as_string().empty());
    EXPECT_FALSE(Ingredient.at("category").as_string().empty());
    EXPECT_FALSE(Ingredient.at("defaultUnit").as_string().empty());
    EXPECT_TRUE(Ingredient.at("isPantryStaple").is_bool());
    EXPECT_TRUE(Ingredient.contains("aliases"));
}

TEST_F(ApiTestFixture, ReturnsRecipeDetailsAndNotFoundError) {
    const auto Found = Request(
        boost::beast::http::verb::get, "/api/v1/recipes/recipe.tomato_egg");
    const auto Missing = Request(
        boost::beast::http::verb::get, "/api/v1/recipes/recipe.does-not-exist");

    ASSERT_EQ(static_cast<int>(Found.result()), 200);
    EXPECT_EQ(ParseJson(Found.body()).as_object().at("id").as_string(), "recipe.tomato_egg");
    ASSERT_EQ(static_cast<int>(Missing.result()), 404);
    EXPECT_EQ(ParseJson(Missing.body()).as_object().at("error").as_object().at("code").as_string(),
              "recipe_not_found");
}

TEST_F(ApiTestFixture, RejectsInvalidQueryAndRecommendationBody) {
    const auto Query = Request(
        boost::beast::http::verb::get, "/api/v1/recipes?limit=101");
    const auto Body = Request(
        boost::beast::http::verb::post,
        "/api/v1/recommendations/tonight",
        "{\"servings\":0,\"availableMinutes\":30}");

    ASSERT_EQ(static_cast<int>(Query.result()), 400);
    EXPECT_EQ(ParseJson(Query.body()).as_object().at("error").as_object().at("code").as_string(),
              "invalid_query");
    ASSERT_EQ(static_cast<int>(Body.result()), 400);
    EXPECT_EQ(ParseJson(Body.body()).as_object().at("error").as_object().at("code").as_string(),
              "invalid_body");
}

TEST_F(ApiTestFixture, RejectsRecommendationBodyWithMissingRequiredArrays) {
    const auto Response = Request(
        boost::beast::http::verb::post,
        "/api/v1/recommendations/tonight",
        "{\"servings\":2,\"availableMinutes\":30}");

    ASSERT_EQ(static_cast<int>(Response.result()), 400) << Response.body();
    EXPECT_EQ(ParseJson(Response.body()).as_object().at("error").as_object().at("code").as_string(),
              "invalid_body");
}

TEST_F(ApiTestFixture, ReturnsAtMostThreeTonightRecommendations) {
    const auto Response = Request(
        boost::beast::http::verb::post,
        "/api/v1/recommendations/tonight",
        "{\"servings\":2,\"availableMinutes\":35,\"cuisines\":[\"中餐\"],"
        "\"allergies\":[],\"pantryIngredientIds\":[\"ingredient.rice\"],"
        "\"cookware\":[\"炒锅\"],\"recentRecipeIds\":[]}");

    ASSERT_EQ(static_cast<int>(Response.result()), 200);
    const auto Recommendations = ParseJson(Response.body()).as_array();
    ASSERT_LE(Recommendations.size(), 3U);
    ASSERT_FALSE(Recommendations.empty());
    EXPECT_TRUE(Recommendations[0].as_object().at("recipe").is_object());
}

TEST_F(ApiTestFixture, RejectsOriginOutsideConfiguredCorsAllowlist) {
    const auto Response = Request(
        boost::beast::http::verb::get,
        "/healthz",
        {},
        "https://untrusted.example");

    ASSERT_EQ(static_cast<int>(Response.result()), 403);
    EXPECT_EQ(ParseJson(Response.body()).as_object().at("error").as_object().at("code").as_string(),
              "cors_denied");
}

TEST_F(ApiTestFixture, AllowsAdminPreflightMethodsAndAuthorizationHeader) {
    const auto Response = Request(
        boost::beast::http::verb::options,
        "/api/v1/admin/recipes",
        {},
        "https://admin.menu.local",
        {},
        "PATCH",
        "Authorization, Content-Type");

    ASSERT_EQ(static_cast<int>(Response.result()), 204);
    EXPECT_EQ(Response["Access-Control-Allow-Origin"], "https://admin.menu.local");
    EXPECT_NE(Response["Vary"].find("Origin"), std::string::npos);
    EXPECT_NE(Response["Access-Control-Allow-Methods"].find("PATCH"), std::string::npos);
    EXPECT_NE(Response["Access-Control-Allow-Methods"].find("DELETE"), std::string::npos);
    EXPECT_NE(Response["Access-Control-Allow-Headers"].find("Authorization"), std::string::npos);
}

TEST_F(ApiTestFixture, RegistersLogsInAndRotatesRefreshToken) {
    const auto Registered = Request(
        boost::beast::http::verb::post,
        "/api/v1/auth/register",
        "{\"email\":\"cook@example.com\",\"password\":\"Menu-Cook-Password-2026\","
        "\"displayName\":\"厨房\"}");
    ASSERT_EQ(static_cast<int>(Registered.result()), 201);
    const auto RegisteredBody = ParseJson(Registered.body()).as_object();
    const std::string RegisteredRefresh =
        std::string(RegisteredBody.at("refreshToken").as_string().c_str());
    EXPECT_FALSE(RegisteredBody.at("accessToken").as_string().empty());
    EXPECT_TRUE(RegisteredBody.at("user").as_object().at("isAdmin").as_bool());

    const auto LoggedIn = Request(
        boost::beast::http::verb::post,
        "/api/v1/auth/login",
        "{\"email\":\"cook@example.com\",\"password\":\"Menu-Cook-Password-2026\"}");
    ASSERT_EQ(static_cast<int>(LoggedIn.result()), 200);
    const auto LoggedInBody = ParseJson(LoggedIn.body()).as_object();
    EXPECT_FALSE(LoggedInBody.at("accessToken").as_string().empty());

    const auto Refreshed = Request(
        boost::beast::http::verb::post,
        "/api/v1/auth/refresh",
        "{\"refreshToken\":\"" + RegisteredRefresh + "\"}");
    ASSERT_EQ(static_cast<int>(Refreshed.result()), 200);
    const auto RefreshedBody = ParseJson(Refreshed.body()).as_object();
    EXPECT_NE(RefreshedBody.at("refreshToken").as_string(), RegisteredRefresh);

    const auto Reused = Request(
        boost::beast::http::verb::post,
        "/api/v1/auth/refresh",
        "{\"refreshToken\":\"" + RegisteredRefresh + "\"}");
    EXPECT_EQ(static_cast<int>(Reused.result()), 401);
}

TEST_F(ApiTestFixture, RejectsInvalidLoginWithoutRevealingAccountState) {
    const auto Response = Request(
        boost::beast::http::verb::post,
        "/api/v1/auth/login",
        "{\"email\":\"missing@example.com\",\"password\":\"Wrong-Password-2026\"}");

    ASSERT_EQ(static_cast<int>(Response.result()), 401) << Response.body();
    EXPECT_EQ(ParseJson(Response.body()).as_object().at("error").as_object().at("code").as_string(),
              "authentication_failed");
}

TEST_F(ApiTestFixture, RejectsOversizedLoginPasswordAtApiBoundary) {
    const auto Registered = Request(
        boost::beast::http::verb::post,
        "/api/v1/auth/register",
        "{\"email\":\"bounded@example.com\",\"password\":\"Menu-Cook-Password-2026\","
        "\"displayName\":\"厨房\"}");
    ASSERT_EQ(static_cast<int>(Registered.result()), 201);

    const std::string LongPassword(129, 'x');
    const auto DirectResult = Application->Authentication().Login(
        "bounded@example.com", LongPassword);
    ASSERT_FALSE(DirectResult.HasValue());
    EXPECT_EQ(DirectResult.ErrorValue().CodeValue(),
              Menu::Foundation::ErrorCode::AuthenticationFailed);
    const auto Response = Request(
        boost::beast::http::verb::post,
        "/api/v1/auth/login",
        "{\"email\":\"bounded@example.com\",\"password\":\"" + LongPassword + "\"}");

    ASSERT_EQ(static_cast<int>(Response.result()), 400) << Response.body();
    EXPECT_EQ(ParseJson(Response.body()).as_object().at("error").as_object().at("code").as_string(),
              "invalid_body");
}

TEST_F(ApiTestFixture, RejectsWhitespaceAndControlCharactersInDisplayName) {
    const auto WhitespaceResponse = Request(
        boost::beast::http::verb::post,
        "/api/v1/auth/register",
        "{\"email\":\"blank@example.com\",\"password\":\"Menu-Cook-Password-2026\","
        "\"displayName\":\"   \"}");
    const auto ControlResponse = Request(
        boost::beast::http::verb::post,
        "/api/v1/auth/register",
        "{\"email\":\"control@example.com\",\"password\":\"Menu-Cook-Password-2026\","
        "\"displayName\":\"\\u0001\"}");

    EXPECT_EQ(static_cast<int>(WhitespaceResponse.result()), 400);
    EXPECT_EQ(static_cast<int>(ControlResponse.result()), 400);
}

TEST_F(ApiTestFixture, AdminCanCreatePublishAndDeleteRecipeWithAccessToken) {
    const auto Registered = Request(
        boost::beast::http::verb::post,
        "/api/v1/auth/register",
        "{\"email\":\"admin-crud@example.com\",\"password\":\"Menu-Cook-Password-2026\","
        "\"displayName\":\"管理员\"}");
    ASSERT_EQ(static_cast<int>(Registered.result()), 201);
    const auto RegisteredBody = ParseJson(Registered.body()).as_object();
    const std::string AccessToken =
        std::string(RegisteredBody.at("accessToken").as_string().c_str());
    const std::string Authorization = "Bearer " + AccessToken;
    const std::string RecipeBody =
        "{\"id\":\"recipe.admin-test\",\"slug\":\"admin-test\","
        "\"name\":\"管理员测试菜\",\"cuisine\":\"中餐\","
        "\"description\":\"用于管理流程测试\",\"prepMinutes\":5,"
        "\"cookMinutes\":10,\"servings\":2,\"difficulty\":1,"
        "\"imagePath\":\"assets/media/admin-test.png\",\"status\":\"draft\","
        "\"allergens\":[],\"cookware\":[\"炒锅\"],\"ingredients\":[{"
        "\"ingredientId\":\"ingredient.rice\",\"quantity\":200,\"unit\":\"g\","
        "\"required\":true,\"servingFactor\":1,\"preparation\":\"\"}],"
        "\"steps\":[{\"stepOrder\":1,\"title\":\"准备\","
        "\"instruction\":\"准备米饭\",\"durationSeconds\":60,\"hasTimer\":false}]}";
    const auto Created = Request(
        boost::beast::http::verb::post,
        "/api/v1/admin/recipes",
        RecipeBody,
        {},
        Authorization);
    ASSERT_EQ(static_cast<int>(Created.result()), 201);

    const auto AdminList = Request(
        boost::beast::http::verb::get,
        "/api/v1/admin/recipes",
        {},
        {},
        Authorization);
    ASSERT_EQ(static_cast<int>(AdminList.result()), 200);
    EXPECT_TRUE(AdminList.body().find("recipe.admin-test") != std::string::npos);

    std::string PublishedBody = RecipeBody;
    const std::string DraftStatus = "\"status\":\"draft\"";
    const auto DraftPosition = PublishedBody.find(DraftStatus);
    ASSERT_NE(DraftPosition, std::string::npos);
    PublishedBody.replace(
        DraftPosition, DraftStatus.size(), "\"status\":\"published\"");
    const auto Updated = Request(
        boost::beast::http::verb::patch,
        "/api/v1/admin/recipes/recipe.admin-test",
        PublishedBody,
        {},
        Authorization);
    ASSERT_EQ(static_cast<int>(Updated.result()), 200);
    EXPECT_EQ(ParseJson(Updated.body()).as_object().at("status").as_string(), "published");

    const auto Deleted = Request(
        boost::beast::http::verb::delete_,
        "/api/v1/admin/recipes/recipe.admin-test",
        {},
        {},
        Authorization);
    EXPECT_EQ(static_cast<int>(Deleted.result()), 204);
}

TEST_F(ApiTestFixture, ProtectsAdminRecipesFromMissingAndNonAdminCredentials) {
    const auto AdminRegistered = Request(
        boost::beast::http::verb::post,
        "/api/v1/auth/register",
        "{\"email\":\"owner@example.com\",\"password\":\"Menu-Cook-Password-2026\","
        "\"displayName\":\"管理员\"}");
    ASSERT_EQ(static_cast<int>(AdminRegistered.result()), 201);
    const auto UserRegistered = Request(
        boost::beast::http::verb::post,
        "/api/v1/auth/register",
        "{\"email\":\"member@example.com\",\"password\":\"Menu-Cook-Password-2026\","
        "\"displayName\":\"成员\"}");
    ASSERT_EQ(static_cast<int>(UserRegistered.result()), 201);
    const auto UserBody = ParseJson(UserRegistered.body()).as_object();
    const std::string UserAuthorization =
        "Bearer " + std::string(UserBody.at("accessToken").as_string().c_str());

    const auto Missing = Request(
        boost::beast::http::verb::get,
        "/api/v1/admin/recipes");
    const auto NonAdmin = Request(
        boost::beast::http::verb::get,
        "/api/v1/admin/recipes",
        {},
        {},
        UserAuthorization);

    EXPECT_EQ(static_cast<int>(Missing.result()), 401);
    EXPECT_EQ(static_cast<int>(NonAdmin.result()), 403);
}

TEST_F(ApiTestFixture, AuthenticatedUserCanManagePlansSessionsAndFeedback) {
    const auto Registered = Request(
        boost::beast::http::verb::post,
        "/api/v1/auth/register",
        "{\"email\":\"workflow@example.com\",\"password\":\"Menu-Cook-Password-2026\","
        "\"displayName\":\"计划用户\"}");
    ASSERT_EQ(static_cast<int>(Registered.result()), 201);
    const auto RegisteredBody = ParseJson(Registered.body()).as_object();
    const std::string Authorization =
        "Bearer " + std::string(RegisteredBody.at("accessToken").as_string().c_str());

    const auto Unauthorized = Request(
        boost::beast::http::verb::get, "/api/v1/plans?from=2026-09-01&to=2026-09-07");
    ASSERT_EQ(static_cast<int>(Unauthorized.result()), 401);

    const auto CreatedPlan = Request(
        boost::beast::http::verb::post,
        "/api/v1/plans",
        "{\"planDate\":\"2026-09-01\",\"items\":["
        "{\"recipeId\":\"recipe.tomato_egg\",\"servings\":2,\"sortOrder\":0},"
        "{\"recipeId\":\"recipe.pasta_tomato\",\"servings\":2,\"sortOrder\":1}]}" ,
        {},
        Authorization);
    ASSERT_EQ(static_cast<int>(CreatedPlan.result()), 201);
    const auto PlanBody = ParseJson(CreatedPlan.body()).as_object();
    EXPECT_EQ(PlanBody.at("planDate").as_string(), "2026-09-01");
    ASSERT_TRUE(PlanBody.at("items").is_array());
    ASSERT_EQ(PlanBody.at("items").as_array().size(), 2U);
    const auto CombinedIngredients = PlanBody.at("combinedIngredients").as_array();
    const auto Tomato = std::ranges::find_if(
        CombinedIngredients,
        [](const boost::json::value& Value) {
            return Value.as_object().at("ingredientId").as_string() == "ingredient.tomato";
        });
    ASSERT_NE(Tomato, CombinedIngredients.end());
    EXPECT_EQ(Tomato->as_object().at("ingredientName").as_string(), "番茄");

    const auto Plans = Request(
        boost::beast::http::verb::get,
        "/api/v1/plans?from=2026-09-01&to=2026-09-07",
        {},
        {},
        Authorization);
    ASSERT_EQ(static_cast<int>(Plans.result()), 200);
    ASSERT_EQ(ParseJson(Plans.body()).as_array().size(), 1U);

    const auto Session = Request(
        boost::beast::http::verb::post,
        "/api/v1/cooking-sessions",
        "{\"recipeId\":\"recipe.tomato_egg\"}",
        {},
        Authorization);
    ASSERT_EQ(static_cast<int>(Session.result()), 201);
    const auto SessionBody = ParseJson(Session.body()).as_object();
    const std::string SessionId = std::string(SessionBody.at("id").as_string().c_str());
    EXPECT_EQ(SessionBody.at("state").as_string(), "active");

    const auto UpdatedSession = Request(
        boost::beast::http::verb::patch,
        "/api/v1/cooking-sessions/" + SessionId,
        "{\"currentStepOrder\":2,\"state\":\"paused\"}",
        {},
        Authorization);
    ASSERT_EQ(static_cast<int>(UpdatedSession.result()), 200);
    EXPECT_EQ(ParseJson(UpdatedSession.body()).as_object().at("state").as_string(), "paused");

    const auto Feedback = Request(
        boost::beast::http::verb::post,
        "/api/v1/feedback",
        "{\"recipeId\":\"recipe.tomato_egg\",\"outcome\":\"made\","
        "\"tags\":[\"下次还想做\"],\"comment\":\"咸淡合适\"}",
        {},
        Authorization);
    ASSERT_EQ(static_cast<int>(Feedback.result()), 201);
    EXPECT_EQ(ParseJson(Feedback.body()).as_object().at("outcome").as_string(), "made");
}

TEST_F(ApiTestFixture, AdminCanManageIngredientsAndReadStats) {
    const auto Registered = Request(
        boost::beast::http::verb::post,
        "/api/v1/auth/register",
        "{\"email\":\"ingredient-admin@example.com\",\"password\":\"Menu-Cook-Password-2026\","
        "\"displayName\":\"食材管理员\"}");
    ASSERT_EQ(static_cast<int>(Registered.result()), 201);
    const auto RegisteredBody = ParseJson(Registered.body()).as_object();
    const std::string Authorization =
        "Bearer " + std::string(RegisteredBody.at("accessToken").as_string().c_str());
    const std::string IngredientBody =
        "{\"id\":\"ingredient.admin_test\",\"name\":\"本地香草\","
        "\"aliases\":[\"香草\"],\"category\":\"香料\",\"defaultUnit\":\"g\","
        "\"isPantryStaple\":false,\"substituteGroup\":\"\",\"storeSkuMapping\":\"\"}";

    const auto Created = Request(
        boost::beast::http::verb::post,
        "/api/v1/admin/ingredients",
        IngredientBody,
        {},
        Authorization);
    ASSERT_EQ(static_cast<int>(Created.result()), 201);
    EXPECT_EQ(ParseJson(Created.body()).as_object().at("name").as_string(), "本地香草");

    const auto AdminList = Request(
        boost::beast::http::verb::get,
        "/api/v1/admin/ingredients",
        {},
        {},
        Authorization);
    ASSERT_EQ(static_cast<int>(AdminList.result()), 200);
    EXPECT_TRUE(AdminList.body().find("ingredient.admin_test") != std::string::npos);

    const std::string UpdatedBody =
        "{\"id\":\"ingredient.admin_test\",\"name\":\"本地香草叶\","
        "\"aliases\":[\"香草叶\"],\"category\":\"香料\",\"defaultUnit\":\"g\","
        "\"isPantryStaple\":true,\"substituteGroup\":\"\",\"storeSkuMapping\":\"sku.future\"}";
    const auto Updated = Request(
        boost::beast::http::verb::patch,
        "/api/v1/admin/ingredients/ingredient.admin_test",
        UpdatedBody,
        {},
        Authorization);
    ASSERT_EQ(static_cast<int>(Updated.result()), 200);
    const auto UpdatedValue = ParseJson(Updated.body()).as_object();
    EXPECT_EQ(UpdatedValue.at("name").as_string(), "本地香草叶");
    EXPECT_TRUE(UpdatedValue.at("isPantryStaple").as_bool());

    const auto Stats = Request(
        boost::beast::http::verb::get,
        "/api/v1/admin/stats",
        {},
        {},
        Authorization);
    ASSERT_EQ(static_cast<int>(Stats.result()), 200);
    EXPECT_GE(ParseJson(Stats.body()).as_object().at("ingredientCount").as_int64(), 20);

    const auto Deleted = Request(
        boost::beast::http::verb::delete_,
        "/api/v1/admin/ingredients/ingredient.admin_test",
        {},
        {},
        Authorization);
    EXPECT_EQ(static_cast<int>(Deleted.result()), 204);
}
