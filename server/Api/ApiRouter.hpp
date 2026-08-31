#pragma once

#include <Application/RecipeApplicationService.hpp>
#include <Application/StorageExecutor.hpp>
#include <Transport/Core/HttpHandler.hpp>

#include <atomic>
#include <string>

namespace Menu::Api {

class ApiRouter final {
public:
    ApiRouter(
        Application::RecipeApplicationService& ServiceValue,
        Application::StorageExecutor& StorageValue);

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
    std::atomic<bool> Ready = false;
};

}  // namespace Menu::Api
