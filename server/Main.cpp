#include <Api/ApiRouter.hpp>
#include <Application/RecipeApplicationService.hpp>
#include <Application/AdminIngredientApplicationService.hpp>
#include <Application/RuleBasedRecommendationProvider.hpp>
#include <Application/StorageExecutor.hpp>
#include <Application/CookingSessionApplicationService.hpp>
#include <Application/FeedbackApplicationService.hpp>
#include <Application/MealPlanApplicationService.hpp>
#include <Infrastructure/SqliteCookingSessionRepository.hpp>
#include <Infrastructure/SqliteFeedbackRepository.hpp>
#include <Infrastructure/SqliteMealPlanRepository.hpp>
#include <Infrastructure/SqliteAdminIngredientRepository.hpp>
#include <Infrastructure/MigrationRunner.hpp>
#include <Infrastructure/SeedData.hpp>
#include <Infrastructure/SqliteAuthRepository.hpp>
#include <Infrastructure/SqliteAuthService.hpp>
#include <Infrastructure/SqliteAdminRecipeRepository.hpp>
#include <Infrastructure/SqliteIngredientRepository.hpp>
#include <Infrastructure/SqliteRecipeRepository.hpp>
#include <Infrastructure/SqliteDatabase.hpp>
#include <Transport/HttpServer.hpp>

#include <boost/asio/signal_set.hpp>
#include <boost/json/parse.hpp>

#include <algorithm>
#include <cstdint>
#include <filesystem>
#include <fstream>
#include <iostream>
#include <iterator>
#include <memory>
#include <string>
#include <string_view>
#include <thread>
#include <utility>
#include <vector>

namespace {

struct ServerSettings {
    Menu::Transport::HttpServerOptions Http;
    std::filesystem::path DatabasePath = "server/data/menu.db";
    std::filesystem::path MediaDirectory = "assets/media";
    std::vector<std::string> CorsOrigins;
};

Menu::Foundation::Error InvalidConfiguration(std::string Message) {
    return Menu::Foundation::Error(
        Menu::Foundation::ErrorCode::InvalidArgument, std::move(Message));
}

Menu::Foundation::Result<std::string> ReadString(
    const boost::json::object& Object,
    std::string_view Key,
    std::string DefaultValue,
    std::size_t MaximumLength) {
    const boost::json::value* Value = Object.if_contains(Key);
    if (Value == nullptr) {
        return DefaultValue;
    }
    if (!Value->is_string() || Value->as_string().size() > MaximumLength) {
        return Menu::Foundation::Result<std::string>::FromError(
            InvalidConfiguration("配置字符串无效"));
    }
    return std::string(Value->as_string().c_str());
}

Menu::Foundation::Result<std::int64_t> ReadInteger(
    const boost::json::object& Object,
    std::string_view Key,
    std::int64_t DefaultValue,
    std::int64_t Minimum,
    std::int64_t Maximum) {
    const boost::json::value* Value = Object.if_contains(Key);
    if (Value == nullptr) {
        return DefaultValue;
    }
    if (!Value->is_int64()) {
        return Menu::Foundation::Result<std::int64_t>::FromError(
            InvalidConfiguration("配置数字无效"));
    }
    const std::int64_t Parsed = Value->as_int64();
    if (Parsed < Minimum || Parsed > Maximum) {
        return Menu::Foundation::Result<std::int64_t>::FromError(
            InvalidConfiguration("配置数字超出范围"));
    }
    return Parsed;
}

Menu::Foundation::Result<std::vector<std::string>> ReadStringArray(
    const boost::json::object& Object,
    std::string_view Key) {
    const boost::json::value* Value = Object.if_contains(Key);
    if (Value == nullptr) {
        return std::vector<std::string>();
    }
    if (!Value->is_array() || Value->as_array().size() > 32U) {
        return Menu::Foundation::Result<std::vector<std::string>>::FromError(
            InvalidConfiguration("配置数组无效"));
    }
    std::vector<std::string> Values;
    for (const boost::json::value& Item : Value->as_array()) {
        if (!Item.is_string() || Item.as_string().size() > 256U) {
            return Menu::Foundation::Result<std::vector<std::string>>::FromError(
                InvalidConfiguration("配置数组元素无效"));
        }
        Values.emplace_back(Item.as_string().c_str());
    }
    return Values;
}

Menu::Foundation::Result<ServerSettings> LoadSettings(
    const std::filesystem::path& ConfigPath) {
    std::ifstream Input(ConfigPath, std::ios::binary);
    if (!Input) {
        return Menu::Foundation::Result<ServerSettings>::FromError(
            InvalidConfiguration("无法打开配置文件"));
    }
    const std::string Text{
        std::istreambuf_iterator<char>(Input), std::istreambuf_iterator<char>()};
    if (Text.size() > 128U * 1024U) {
        return Menu::Foundation::Result<ServerSettings>::FromError(
            InvalidConfiguration("配置文件过大"));
    }

    boost::system::error_code ParseError;
    const boost::json::value Parsed = boost::json::parse(Text, ParseError);
    if (ParseError || !Parsed.is_object()) {
        return Menu::Foundation::Result<ServerSettings>::FromError(
            InvalidConfiguration("配置 JSON 无效"));
    }
    const boost::json::object& Object = Parsed.as_object();
    ServerSettings Settings;
    const auto Address = ReadString(Object, "address", Settings.Http.Address, 64U);
    const auto Port = ReadInteger(Object, "port", Settings.Http.Port, 1, 65535);
    const auto DatabasePath = ReadString(
        Object, "databasePath", Settings.DatabasePath.string(), 512U);
    const auto MediaDirectory = ReadString(
        Object, "mediaDirectory", Settings.MediaDirectory.string(), 512U);
    const auto MaxBodyBytes = ReadInteger(
        Object, "maxBodyBytes", Settings.Http.MaxBodyBytes, 1024, 16 * 1024 * 1024);
    const auto MaxTargetBytes = ReadInteger(
        Object, "maxTargetBytes", Settings.Http.MaxTargetBytes, 256, 1024 * 1024);
    const auto TimeoutSeconds = ReadInteger(
        Object, "timeoutSeconds", Settings.Http.TimeoutSeconds, 1, 300);
    const auto CorsOrigins = ReadStringArray(Object, "corsOrigins");
    if (!Address.HasValue() || !Port.HasValue() || !DatabasePath.HasValue() ||
        !MediaDirectory.HasValue() || !MaxBodyBytes.HasValue() ||
        !MaxTargetBytes.HasValue() || !TimeoutSeconds.HasValue() ||
        !CorsOrigins.HasValue()) {
        return Menu::Foundation::Result<ServerSettings>::FromError(
            InvalidConfiguration("配置字段无效"));
    }
    Settings.Http.Address = Address.Value();
    Settings.Http.Port = static_cast<std::uint16_t>(Port.Value());
    Settings.DatabasePath = DatabasePath.Value();
    Settings.MediaDirectory = MediaDirectory.Value();
    Settings.Http.MaxBodyBytes = static_cast<std::size_t>(MaxBodyBytes.Value());
    Settings.Http.MaxTargetBytes = static_cast<std::size_t>(MaxTargetBytes.Value());
    Settings.Http.TimeoutSeconds = static_cast<int>(TimeoutSeconds.Value());
    Settings.CorsOrigins = CorsOrigins.Value();
    return Settings;
}

void ReportFailure(const Menu::Foundation::Error& ErrorValue) {
    std::cerr << "MenuServer startup failed: "
              << ErrorValue.MessageValue() << '\n';
}

}  // namespace

int main(int ArgumentCount, char** Arguments) {
    std::filesystem::path ConfigPath = "config/Menu.example.json";
    for (int Index = 1; Index < ArgumentCount; ++Index) {
        if (std::string_view(Arguments[Index]) == "--config" && Index + 1 < ArgumentCount) {
            ConfigPath = Arguments[++Index];
        } else {
            std::cerr << "Usage: MenuServer [--config path]\n";
            return 2;
        }
    }

    const auto SettingsResult = LoadSettings(ConfigPath);
    if (!SettingsResult.HasValue()) {
        ReportFailure(SettingsResult.ErrorValue());
        return 2;
    }
    const ServerSettings Settings = SettingsResult.Value();
    std::error_code FileError;
    std::filesystem::create_directories(Settings.DatabasePath.parent_path(), FileError);
    if (FileError) {
        std::cerr << "MenuServer startup failed: database directory unavailable\n";
        return 2;
    }
    std::filesystem::create_directories(Settings.MediaDirectory, FileError);
    if (FileError) {
        std::cerr << "MenuServer startup failed: media directory unavailable\n";
        return 2;
    }

    auto DatabaseResult = Menu::Infrastructure::SqliteDatabase::Open(Settings.DatabasePath);
    if (!DatabaseResult.HasValue()) {
        ReportFailure(DatabaseResult.ErrorValue());
        return 2;
    }
    auto Database = std::make_unique<Menu::Infrastructure::SqliteDatabase>(
        std::move(DatabaseResult).Value());
    const auto MigrationResult = Menu::Infrastructure::MigrationRunner::Apply(*Database);
    if (!MigrationResult.HasValue()) {
        ReportFailure(MigrationResult.ErrorValue());
        return 2;
    }
    const auto SeedResult = Menu::Infrastructure::SeedData::InsertIfEmpty(*Database);
    if (!SeedResult.HasValue()) {
        ReportFailure(SeedResult.ErrorValue());
        return 2;
    }

    auto Repository = std::make_unique<Menu::Infrastructure::SqliteRecipeRepository>(*Database);
    auto IngredientStore = std::make_unique<Menu::Infrastructure::SqliteIngredientRepository>(*Database);
    auto Service = std::make_unique<Menu::Application::RecipeApplicationService>(
        std::move(Repository),
        std::make_unique<Menu::Application::RuleBasedRecommendationProvider>(),
        std::move(IngredientStore));
    Menu::Application::StorageExecutor StorageExecutor(2);
    auto Authentication = std::make_unique<Menu::Infrastructure::SqliteAuthService>(
        std::make_unique<Menu::Infrastructure::SqliteAuthRepository>(*Database));
    auto AdminService = std::make_unique<Menu::Application::AdminRecipeApplicationService>(
        std::make_unique<Menu::Infrastructure::SqliteAdminRecipeRepository>(*Database));
    auto AdminIngredientService =
        std::make_unique<Menu::Application::AdminIngredientApplicationService>(
            std::make_unique<Menu::Infrastructure::SqliteAdminIngredientRepository>(*Database));
    auto MealPlanService = std::make_unique<Menu::Application::MealPlanApplicationService>(
        std::make_unique<Menu::Infrastructure::SqliteMealPlanRepository>(*Database));
    auto CookingSessionService =
        std::make_unique<Menu::Application::CookingSessionApplicationService>(
            std::make_unique<Menu::Infrastructure::SqliteCookingSessionRepository>(*Database));
    auto FeedbackService = std::make_unique<Menu::Application::FeedbackApplicationService>(
        std::make_unique<Menu::Infrastructure::SqliteFeedbackRepository>(*Database));
    Menu::Api::ApiRouter Router(
        *Service,
        StorageExecutor,
        Authentication.get(),
        AdminService.get(),
        Settings.CorsOrigins,
        AdminIngredientService.get(),
        MealPlanService.get(),
        CookingSessionService.get(),
        FeedbackService.get());

    boost::asio::io_context IoContext;
    Menu::Transport::HttpServer Server(IoContext, Settings.Http);
    const auto StartResult = Server.Start(
        [&Router](Menu::Transport::HttpRequest Request,
                  Menu::Transport::HttpResponseCallback Complete) {
            Router.Handle(std::move(Request), std::move(Complete));
        });
    if (!StartResult.HasValue()) {
        ReportFailure(StartResult.ErrorValue());
        return 2;
    }
    Router.SetReady(true);

    boost::asio::signal_set Signals(IoContext, SIGINT, SIGTERM);
    Signals.async_wait([&Server, &IoContext](boost::system::error_code Error, int) {
        if (!Error) {
            Server.Stop();
            IoContext.stop();
        }
    });

    const unsigned int HardwareThreads = std::thread::hardware_concurrency();
    const std::size_t IoThreadCount = std::max<std::size_t>(
        2U, HardwareThreads == 0U ? 2U : HardwareThreads / 2U);
    std::vector<std::thread> IoThreads;
    IoThreads.reserve(IoThreadCount > 0U ? IoThreadCount - 1U : 0U);
    for (std::size_t Index = 1; Index < IoThreadCount; ++Index) {
        IoThreads.emplace_back([&IoContext]() {
            IoContext.run();
        });
    }
    IoContext.run();
    Server.Stop();
    IoContext.stop();
    for (std::thread& IoThread : IoThreads) {
        IoThread.join();
    }
    StorageExecutor.Shutdown();
    return 0;
}
