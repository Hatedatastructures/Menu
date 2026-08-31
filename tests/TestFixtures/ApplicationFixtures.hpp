#pragma once

#include <Application/RecipeApplicationService.hpp>
#include <Application/RuleBasedRecommendationProvider.hpp>
#include <TestFixtures/RecipeFixtures.hpp>

#include <memory>
#include <optional>
#include <utility>
#include <vector>

namespace Menu::Tests::ApplicationFixtures {

class InMemoryRecipeRepository final : public Application::RecipeRepository {
public:
    InMemoryRecipeRepository() {
        auto SafeRecipe = RecipeFixtures::SingleRecipeWithIngredients(
            2, {Domain::RecipeIngredient{"ingredient.rice", 200.0, "g", true}});
        SafeRecipe.Id = "recipe.safe";

        auto AllergyRecipe = SafeRecipe;
        AllergyRecipe.Id = "recipe.allergy";
        AllergyRecipe.Allergens = {"花生"};

        auto SlowRecipe = SafeRecipe;
        SlowRecipe.Id = "recipe.slow";
        SlowRecipe.CookMinutes = 60;

        auto AlternateRecipe = SafeRecipe;
        AlternateRecipe.Id = "recipe.alternate";

        Recipes = {std::move(SafeRecipe), std::move(AllergyRecipe),
                   std::move(SlowRecipe), std::move(AlternateRecipe)};
    }

    Foundation::Result<std::vector<Domain::Recipe>> ListPublished() override {
        return Recipes;
    }

    Foundation::Result<std::optional<Domain::Recipe>> FindPublishedById(
        std::string_view Id) override {
        for (const Domain::Recipe& RecipeValue : Recipes) {
            if (RecipeValue.Id == Id) {
                return std::optional<Domain::Recipe>(RecipeValue);
            }
        }
        return std::optional<Domain::Recipe>();
    }

private:
    std::vector<Domain::Recipe> Recipes;
};

inline Application::RecipeApplicationService WithSeedRecipes() {
    return Application::RecipeApplicationService(
        std::make_unique<InMemoryRecipeRepository>(),
        std::make_unique<Application::RuleBasedRecommendationProvider>());
}

}  // namespace Menu::Tests::ApplicationFixtures
