#pragma once

#include <Domain/Recipe.hpp>

#include <string>
#include <vector>

namespace Menu::Domain {

struct MealPlanItem {
    std::string Id;
    std::string RecipeId;
    std::string RecipeName;
    std::string ImagePath;
    int Servings = 1;
    int SortOrder = 0;
    std::vector<RecipeIngredient> Ingredients;
};

struct MealPlan {
    std::string Id;
    std::string UserId;
    std::string PlanDate;
    std::vector<MealPlanItem> Items;
    std::vector<RecipeIngredient> CombinedIngredients;
};

}  // namespace Menu::Domain
