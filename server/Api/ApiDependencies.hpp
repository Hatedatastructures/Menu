#pragma once

#include <Application/AdminIngredientApplicationService.hpp>
#include <Application/AdminRecipeApplicationService.hpp>
#include <Application/AuthService.hpp>
#include <Application/CookingSessionApplicationService.hpp>
#include <Application/FeedbackApplicationService.hpp>
#include <Application/MealPlanApplicationService.hpp>
#include <Application/RecipeApplicationService.hpp>
#include <Application/StorageExecutor.hpp>

#include <string>
#include <vector>

namespace Menu::Api {

struct ApiDependencies final {
    Application::RecipeApplicationService& Service;
    Application::StorageExecutor& Storage;
    std::vector<std::string> CorsOrigins;
    Application::AuthService& Authentication;
    Application::AdminRecipeApplicationService& AdminService;
    Application::AdminIngredientApplicationService& AdminIngredientService;
    Application::MealPlanApplicationService& MealPlans;
    Application::CookingSessionApplicationService& CookingSessions;
    Application::FeedbackApplicationService& Feedbacks;
};

}  // namespace Menu::Api
