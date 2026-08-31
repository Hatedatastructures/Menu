#pragma once

#include <Foundation/Result.hpp>
#include <Domain/Ingredient.hpp>
#include <Domain/Recipe.hpp>

#include <string>
#include <string_view>
#include <vector>

namespace Menu::Domain {

struct IngredientAvailability {
    std::size_t AvailableCount = 0;
    std::size_t MissingCount = 0;
    std::vector<std::string> AvailableIngredientIds;
    std::vector<std::string> MissingIngredientIds;
};

class RecipeRules final {
public:
    static Foundation::Result<std::vector<RecipeIngredient>> ScaleRecipeIngredients(
        const Recipe& RecipeValue,
        int TargetServings);

    static Foundation::Result<void> ValidateRecipe(const Recipe& RecipeValue);

    static Foundation::Result<std::string> NormalizeIngredientAlias(
        std::string_view Alias,
        const std::vector<Ingredient>& Ingredients);

    static IngredientAvailability BuildIngredientAvailability(
        const Recipe& RecipeValue,
        const std::vector<std::string>& PantryIngredientIds);
};

}  // namespace Menu::Domain
