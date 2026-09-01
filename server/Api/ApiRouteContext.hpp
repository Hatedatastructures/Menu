#pragma once

#include <Application/AdminRecipeApplicationService.hpp>
#include <Application/AdminIngredientApplicationService.hpp>
#include <Application/AuthService.hpp>
#include <Application/CookingSessionApplicationService.hpp>
#include <Application/FeedbackApplicationService.hpp>
#include <Application/MealPlanApplicationService.hpp>
#include <Application/RecipeApplicationService.hpp>

namespace Menu::Api {

struct ApiRouteContext {
    Application::RecipeApplicationService& Service;
    Application::AuthService* Authentication = nullptr;
    Application::AdminRecipeApplicationService* AdminService = nullptr;
    Application::AdminIngredientApplicationService* AdminIngredientService = nullptr;
    Application::MealPlanApplicationService* MealPlans = nullptr;
    Application::CookingSessionApplicationService* CookingSessions = nullptr;
    Application::FeedbackApplicationService* Feedbacks = nullptr;
};

}  // namespace Menu::Api
