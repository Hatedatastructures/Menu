#pragma once

#include <Api/ApiDependencies.hpp>
#include <Transport/Core/HttpHandler.hpp>

#include <atomic>
#include <functional>
#include <string>
#include <vector>

namespace Menu::Api {

class ApiRouter final {
public:
    explicit ApiRouter(ApiDependencies DependenciesValue);

    void Handle(
        Transport::HttpRequest Request,
        Transport::HttpResponseCallback Complete);

    void SetReady(bool ReadyValue) noexcept;
    void SetCompletionDispatcher(
        Application::StorageExecutor::CompletionDispatcher DispatcherValue);

private:
    Transport::HttpResponse Route(
        const Transport::HttpRequest& Request,
        const std::string& RequestId);

    Application::RecipeApplicationService& Service;
    Application::StorageExecutor& Storage;
    Application::AuthService& Authentication;
    Application::AdminRecipeApplicationService& AdminService;
    Application::AdminIngredientApplicationService& AdminIngredientService;
    Application::MealPlanApplicationService& MealPlans;
    Application::CookingSessionApplicationService& CookingSessions;
    Application::FeedbackApplicationService& Feedbacks;
    std::vector<std::string> CorsOrigins;
    Application::StorageExecutor::CompletionDispatcher CompletionDispatcher;
    std::atomic<bool> Ready = false;
};

}  // namespace Menu::Api
