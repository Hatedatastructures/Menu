#pragma once

#include <Foundation/Result.hpp>
#include <Domain/Recipe.hpp>
#include <Domain/RecipeRules.hpp>

#include <string>
#include <vector>

namespace Menu::Domain {

struct RecommendationRequest {
    int Servings = 1;
    int AvailableMinutes = 60;
    std::vector<std::string> Cuisines;
    std::vector<std::string> Allergies;
    std::vector<std::string> PantryIngredientIds;
    std::vector<std::string> Cookware;
    std::vector<std::string> RecentRecipeIds;
};

struct Recommendation {
    Recipe RecipeValue;
    IngredientAvailability Availability;
    int Score = 0;
};

class RuleBasedRecommendation final {
public:
    static Foundation::Result<std::vector<Recommendation>> Rank(
        const RecommendationRequest& Request,
        const std::vector<Recipe>& Recipes);
};

}  // namespace Menu::Domain
