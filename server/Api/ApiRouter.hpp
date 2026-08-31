#pragma once

#include <Application/AdminRecipeApplicationService.hpp>
#include <Application/AuthService.hpp>
#include <Application/RecipeApplicationService.hpp>
#include <Application/StorageExecutor.hpp>
#include <Transport/Core/HttpHandler.hpp>

#include <atomic>
#include <string>
#include <vector>

namespace Menu::Api {

class ApiRouter final {
public:
    ApiRouter(
        Application::RecipeApplicationService& ServiceValue,
        Application::StorageExecutor& StorageValue,
        Application::AuthService* AuthenticationValue = nullptr,
        Application::AdminRecipeApplicationService* AdminServiceValue = nullptr,
        std::vector<std::string> CorsOriginsValue = {});

    void Handle(
        Transport::HttpRequest Request,
        Transport::HttpResponseCallback Complete);

    void SetReady(bool ReadyValue) noexcept;

private:
    Transport::HttpResponse Route(
        const Transport::HttpRequest& Request,
        const std::string& RequestId);

    Application::RecipeApplicationService& Service;
    Application::StorageExecutor& Storage;
    Application::AuthService* Authentication = nullptr;
    Application::AdminRecipeApplicationService* AdminService = nullptr;
    std::vector<std::string> CorsOrigins;
    std::atomic<bool> Ready = false;
};

}  // namespace Menu::Api
