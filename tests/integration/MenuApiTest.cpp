#include <gtest/gtest.h>

#include <Api/ApiRouter.hpp>
#include <Application/RecipeApplicationService.hpp>
#include <Application/RuleBasedRecommendationProvider.hpp>
#include <Application/StorageExecutor.hpp>
#include <Infrastructure/MigrationRunner.hpp>
#include <Infrastructure/SeedData.hpp>
#include <Infrastructure/SqliteIngredientRepository.hpp>
#include <Infrastructure/SqliteRecipeRepository.hpp>
#include <TestFixtures/DatabaseFixtures.hpp>
#include <Transport/HttpServer.hpp>

#include <boost/asio/connect.hpp>
#include <boost/asio/io_context.hpp>
#include <boost/asio/ip/tcp.hpp>
#include <boost/beast/core/flat_buffer.hpp>
#include <boost/beast/http.hpp>
#include <boost/json/parse.hpp>

#include <memory>
#include <string>
#include <string_view>
#include <thread>
#include <utility>

namespace {

class ApiTestFixture : public ::testing::Test {
protected:
    void SetUp() override {
        Database = std::make_unique<Menu::Tests::DatabaseFixtures::TemporaryDatabase>(
            Menu::Tests::DatabaseFixtures::OpenTemporary());
        ASSERT_TRUE(Menu::Infrastructure::MigrationRunner::Apply(*Database).HasValue());
        ASSERT_TRUE(Menu::Infrastructure::SeedData::InsertIfEmpty(*Database).HasValue());

        Repository = std::make_unique<Menu::Infrastructure::SqliteRecipeRepository>(*Database);
        IngredientStore = std::make_unique<Menu::Infrastructure::SqliteIngredientRepository>(*Database);
        Service = std::make_unique<Menu::Application::RecipeApplicationService>(
            std::move(Repository),
            std::make_unique<Menu::Application::RuleBasedRecommendationProvider>(),
            std::move(IngredientStore));
        Storage = std::make_unique<Menu::Application::StorageExecutor>(2);
        Router = std::make_unique<Menu::Api::ApiRouter>(*Service, *Storage);
        Router->SetReady(true);

        Menu::Transport::HttpServerOptions Options;
        Options.Port = 0;
        Server = std::make_unique<Menu::Transport::HttpServer>(IoContext, Options);
        const auto StartResult = Server->Start(
            [this](Menu::Transport::HttpRequest Request,
                   Menu::Transport::HttpResponseCallback Complete) {
                Router->Handle(std::move(Request), std::move(Complete));
            });
        ASSERT_TRUE(StartResult.HasValue());
        IoThread = std::thread([this]() {
            IoContext.run();
        });
    }

    void TearDown() override {
        if (Server) {
            Server->Stop();
        }
        if (Storage) {
            Storage->Shutdown();
        }
        IoContext.stop();
        if (IoThread.joinable()) {
            IoThread.join();
        }
        Router.reset();
        Service.reset();
        Database.reset();
    }

    [[nodiscard]] boost::beast::http::response<boost::beast::http::string_body> Request(
        boost::beast::http::verb Method,
        std::string_view Target,
        std::string Body = {}) const {
        boost::asio::io_context ClientContext;
        boost::asio::ip::tcp::socket Socket(ClientContext);
        Socket.connect({boost::asio::ip::make_address("127.0.0.1"), Server->LocalPort()});
        boost::beast::http::request<boost::beast::http::string_body> RequestValue(
            Method, Target, 11);
        RequestValue.set(boost::beast::http::field::host, "localhost");
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

    boost::asio::io_context IoContext;
    std::thread IoThread;
    std::unique_ptr<Menu::Tests::DatabaseFixtures::TemporaryDatabase> Database;
    std::unique_ptr<Menu::Infrastructure::SqliteRecipeRepository> Repository;
    std::unique_ptr<Menu::Infrastructure::SqliteIngredientRepository> IngredientStore;
    std::unique_ptr<Menu::Application::RecipeApplicationService> Service;
    std::unique_ptr<Menu::Application::StorageExecutor> Storage;
    std::unique_ptr<Menu::Api::ApiRouter> Router;
    std::unique_ptr<Menu::Transport::HttpServer> Server;
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

TEST_F(ApiTestFixture, ReturnsAtMostThreeTonightRecommendations) {
    const auto Response = Request(
        boost::beast::http::verb::post,
        "/api/v1/recommendations/tonight",
        "{\"servings\":2,\"availableMinutes\":35,\"cuisines\":[\"中餐\"],"
        "\"pantryIngredientIds\":[\"ingredient.rice\"],\"cookware\":[\"炒锅\"]}");

    ASSERT_EQ(static_cast<int>(Response.result()), 200);
    const auto Recommendations = ParseJson(Response.body()).as_array();
    ASSERT_LE(Recommendations.size(), 3U);
    ASSERT_FALSE(Recommendations.empty());
    EXPECT_TRUE(Recommendations[0].as_object().at("recipe").is_object());
}
