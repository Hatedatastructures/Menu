#include "Composition/ServerApplication.hpp"

#include <Api/ApiDependencies.hpp>
#include <Api/ApiRouter.hpp>
#include <Application/AdminIngredientApplicationService.hpp>
#include <Application/AdminRecipeApplicationService.hpp>
#include <Application/CookingSessionApplicationService.hpp>
#include <Application/FeedbackApplicationService.hpp>
#include <Application/MealPlanApplicationService.hpp>
#include <Application/RecipeApplicationService.hpp>
#include <Application/RuleBasedRecommendationProvider.hpp>
#include <Application/StorageExecutor.hpp>
#include <Infrastructure/MigrationRunner.hpp>
#include <Infrastructure/SeedData.hpp>
#include <Infrastructure/SqliteAdminIngredientRepository.hpp>
#include <Infrastructure/SqliteAdminRecipeRepository.hpp>
#include <Infrastructure/SqliteAuthRepository.hpp>
#include <Infrastructure/SqliteAuthService.hpp>
#include <Infrastructure/SqliteCookingSessionRepository.hpp>
#include <Infrastructure/SqliteDatabase.hpp>
#include <Infrastructure/SqliteFeedbackRepository.hpp>
#include <Infrastructure/SqliteIngredientRepository.hpp>
#include <Infrastructure/SqliteMealPlanRepository.hpp>
#include <Infrastructure/SqliteRecipeRepository.hpp>
#include <Transport/HttpServer.hpp>

#include <boost/asio/signal_set.hpp>
#include <boost/asio/post.hpp>

#include <algorithm>
#include <filesystem>
#include <thread>
#include <utility>
#include <vector>

namespace Menu::Server {
namespace {

Foundation::Error StartupError(const char* Message) {
    return Foundation::Error(Foundation::ErrorCode::StorageUnavailable, Message);
}

}  // namespace

Foundation::Result<std::unique_ptr<ServerApplication>> ServerApplication::Create(
    Runtime::ServerConfiguration Configuration) {
    std::error_code FileError;
    std::filesystem::create_directories(Configuration.DatabasePath.parent_path(), FileError);
    if (FileError) {
        return Foundation::Result<std::unique_ptr<ServerApplication>>::FromError(
            StartupError("数据库目录不可用"));
    }
    std::filesystem::create_directories(Configuration.MediaDirectory, FileError);
    if (FileError) {
        return Foundation::Result<std::unique_ptr<ServerApplication>>::FromError(
            StartupError("媒体目录不可用"));
    }

    auto DatabaseResult = Infrastructure::SqliteDatabase::Open(Configuration.DatabasePath);
    if (!DatabaseResult.HasValue()) {
        return Foundation::Result<std::unique_ptr<ServerApplication>>::FromError(
            std::move(DatabaseResult).ErrorValue());
    }
    auto Database = std::make_unique<Infrastructure::SqliteDatabase>(
        std::move(DatabaseResult).Value());
    const auto MigrationResult = Infrastructure::MigrationRunner::Apply(*Database);
    if (!MigrationResult.HasValue()) {
        return Foundation::Result<std::unique_ptr<ServerApplication>>::FromError(
            std::move(MigrationResult).ErrorValue());
    }
    const auto SeedResult = Infrastructure::SeedData::InsertIfEmpty(*Database);
    if (!SeedResult.HasValue()) {
        return Foundation::Result<std::unique_ptr<ServerApplication>>::FromError(
            std::move(SeedResult).ErrorValue());
    }

    auto RecipeRepository = std::make_unique<Infrastructure::SqliteRecipeRepository>(*Database);
    auto IngredientRepository =
        std::make_unique<Infrastructure::SqliteIngredientRepository>(*Database);
    auto Service = std::make_unique<Application::RecipeApplicationService>(
        std::move(RecipeRepository),
        std::make_unique<Application::RuleBasedRecommendationProvider>(),
        std::move(IngredientRepository));
    auto Authentication = std::make_unique<Infrastructure::SqliteAuthService>(
        std::make_unique<Infrastructure::SqliteAuthRepository>(*Database));
    auto AdminService = std::make_unique<Application::AdminRecipeApplicationService>(
        std::make_unique<Infrastructure::SqliteAdminRecipeRepository>(*Database));
    auto AdminIngredientService =
        std::make_unique<Application::AdminIngredientApplicationService>(
            std::make_unique<Infrastructure::SqliteAdminIngredientRepository>(*Database));
    auto MealPlanService = std::make_unique<Application::MealPlanApplicationService>(
        std::make_unique<Infrastructure::SqliteMealPlanRepository>(*Database));
    auto CookingService = std::make_unique<Application::CookingSessionApplicationService>(
        std::make_unique<Infrastructure::SqliteCookingSessionRepository>(*Database));
    auto FeedbackService = std::make_unique<Application::FeedbackApplicationService>(
        std::make_unique<Infrastructure::SqliteFeedbackRepository>(*Database));
    // Repositories share one SQLite connection, so one worker owns all storage calls.
    auto Storage = std::make_unique<Application::StorageExecutor>(1);
    Api::ApiDependencies Dependencies{
        *Service,
        *Storage,
        Configuration.CorsOrigins,
        *Authentication,
        *AdminService,
        *AdminIngredientService,
        *MealPlanService,
        *CookingService,
        *FeedbackService};
    auto Router = std::make_unique<Api::ApiRouter>(std::move(Dependencies));

    return std::unique_ptr<ServerApplication>(new ServerApplication(
        std::move(Configuration),
        std::move(Database),
        std::move(Service),
        std::move(Authentication),
        std::move(AdminService),
        std::move(AdminIngredientService),
        std::move(MealPlanService),
        std::move(CookingService),
        std::move(FeedbackService),
        std::move(Storage),
        std::move(Router)));
}

ServerApplication::ServerApplication(
    Runtime::ServerConfiguration ConfigurationValue,
    std::unique_ptr<Infrastructure::SqliteDatabase> DatabaseValue,
    std::unique_ptr<Application::RecipeApplicationService> ServiceValue,
    std::unique_ptr<Infrastructure::SqliteAuthService> AuthenticationValue,
    std::unique_ptr<Application::AdminRecipeApplicationService> AdminServiceValue,
    std::unique_ptr<Application::AdminIngredientApplicationService> AdminIngredientServiceValue,
    std::unique_ptr<Application::MealPlanApplicationService> MealPlanServiceValue,
    std::unique_ptr<Application::CookingSessionApplicationService> CookingSessionServiceValue,
    std::unique_ptr<Application::FeedbackApplicationService> FeedbackServiceValue,
    std::unique_ptr<Application::StorageExecutor> StorageValue,
    std::unique_ptr<Api::ApiRouter> RouterValueInput)
    : Configuration(std::move(ConfigurationValue)),
      Database(std::move(DatabaseValue)),
      Service(std::move(ServiceValue)),
      AuthenticationService(std::move(AuthenticationValue)),
      AdminService(std::move(AdminServiceValue)),
      AdminIngredientService(std::move(AdminIngredientServiceValue)),
      MealPlanService(std::move(MealPlanServiceValue)),
      CookingSessionService(std::move(CookingSessionServiceValue)),
      FeedbackService(std::move(FeedbackServiceValue)),
      Storage(std::move(StorageValue)),
      RouterValue(std::move(RouterValueInput)),
      Server(std::make_unique<Transport::HttpServer>(Io, Configuration.Http)) {}

ServerApplication::~ServerApplication() {
    Stop();
    if (Storage) {
        Storage->Shutdown();
    }
}

Foundation::Result<void> ServerApplication::Start() {
    if (Started) {
        return Foundation::Result<void>::FromError(
            Foundation::Error(Foundation::ErrorCode::Conflict, "服务已启动"));
    }
    RouterValue->SetCompletionDispatcher([this](std::function<void()> Handler) {
        boost::asio::post(Io, std::move(Handler));
    });
    const auto Result = Server->Start(
        [this](Transport::HttpRequest Request,
               Transport::HttpResponseCallback Complete) {
            RouterValue->Handle(std::move(Request), std::move(Complete));
        });
    if (!Result.HasValue()) {
        return Result;
    }
    RouterValue->SetReady(true);
    Started = true;
    return {};
}

void ServerApplication::Stop() noexcept {
    if (!Started) {
        return;
    }
    RouterValue->SetReady(false);
    Server->Stop();
    Io.stop();
    Started = false;
}

void ServerApplication::Run() {
    boost::asio::signal_set Signals(Io, SIGINT, SIGTERM);
    Signals.async_wait([this](boost::system::error_code Error, int) {
        if (!Error) {
            Stop();
        }
    });

    const unsigned int HardwareThreads = std::thread::hardware_concurrency();
    const std::size_t IoThreadCount = std::max<std::size_t>(
        2U, HardwareThreads == 0U ? 2U : HardwareThreads / 2U);
    std::vector<std::thread> IoThreads;
    IoThreads.reserve(IoThreadCount > 0U ? IoThreadCount - 1U : 0U);
    for (std::size_t Index = 1; Index < IoThreadCount; ++Index) {
        IoThreads.emplace_back([this]() { Io.run(); });
    }
    Io.run();
    Stop();
    for (std::thread& IoThread : IoThreads) {
        IoThread.join();
    }
}

boost::asio::io_context& ServerApplication::IoContext() noexcept {
    return Io;
}

std::uint16_t ServerApplication::LocalPort() const noexcept {
    return Server == nullptr ? 0 : Server->LocalPort();
}

Api::ApiRouter& ServerApplication::Router() noexcept {
    return *RouterValue;
}

Application::AuthService& ServerApplication::Authentication() noexcept {
    return *AuthenticationService;
}

}  // namespace Menu::Server
