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
    Application::AuthService& Authentication;
    Application::AdminRecipeApplicationService& AdminService;
    Application::AdminIngredientApplicationService& AdminIngredientService;
    Application::MealPlanApplicationService& MealPlans;
    Application::CookingSessionApplicationService& CookingSessions;
    Application::FeedbackApplicationService& Feedbacks;
};

}  // namespace Menu::Api
