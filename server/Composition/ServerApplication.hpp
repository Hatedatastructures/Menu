#pragma once

#include <Foundation/Result.hpp>
#include <Runtime/ServerConfiguration.hpp>

#include <boost/asio/io_context.hpp>

#include <cstdint>
#include <memory>

namespace Menu::Api {
class ApiRouter;
}

namespace Menu::Application {
class AuthService;
class AdminIngredientApplicationService;
class AdminRecipeApplicationService;
class CookingSessionApplicationService;
class FeedbackApplicationService;
class MealPlanApplicationService;
class RecipeApplicationService;
class StorageExecutor;
}

namespace Menu::Infrastructure {
class SqliteAuthService;
class SqliteDatabase;
}

namespace Menu::Transport {
class HttpServer;
}

namespace Menu::Server {

class ServerApplication final {
public:
    [[nodiscard]] static Foundation::Result<std::unique_ptr<ServerApplication>> Create(
        Runtime::ServerConfiguration Configuration);

    ~ServerApplication();

    ServerApplication(const ServerApplication&) = delete;
    ServerApplication& operator=(const ServerApplication&) = delete;

    [[nodiscard]] Foundation::Result<void> Start();
    void Stop() noexcept;
    void Run();

    [[nodiscard]] boost::asio::io_context& IoContext() noexcept;
    [[nodiscard]] std::uint16_t LocalPort() const noexcept;
    [[nodiscard]] Api::ApiRouter& Router() noexcept;
    [[nodiscard]] Application::AuthService& Authentication() noexcept;

private:
    ServerApplication(
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
        std::unique_ptr<Api::ApiRouter> RouterValue);

    Runtime::ServerConfiguration Configuration;
    boost::asio::io_context Io;
    std::unique_ptr<Infrastructure::SqliteDatabase> Database;
    std::unique_ptr<Application::RecipeApplicationService> Service;
    std::unique_ptr<Infrastructure::SqliteAuthService> AuthenticationService;
    std::unique_ptr<Application::AdminRecipeApplicationService> AdminService;
    std::unique_ptr<Application::AdminIngredientApplicationService> AdminIngredientService;
    std::unique_ptr<Application::MealPlanApplicationService> MealPlanService;
    std::unique_ptr<Application::CookingSessionApplicationService> CookingSessionService;
    std::unique_ptr<Application::FeedbackApplicationService> FeedbackService;
    std::unique_ptr<Application::StorageExecutor> Storage;
    std::unique_ptr<Api::ApiRouter> RouterValue;
    std::unique_ptr<Transport::HttpServer> Server;
    bool Started = false;
};

}  // namespace Menu::Server
