#pragma once

#include <Application/IngredientRepository.hpp>
#include <Application/RecipeRepository.hpp>
#include <Application/RecommendationProvider.hpp>

#include <memory>
#include <optional>
#include <string_view>

namespace Menu::Application {

class RecipeApplicationService final {
public:
    RecipeApplicationService(
        std::unique_ptr<RecipeRepository> RepositoryValue,
        std::unique_ptr<RecommendationProvider> ProviderValue,
        std::unique_ptr<IngredientRepository> IngredientRepositoryValue = nullptr);

    RecipeApplicationService(const RecipeApplicationService&) = delete;
    RecipeApplicationService& operator=(const RecipeApplicationService&) = delete;
    RecipeApplicationService(RecipeApplicationService&&) noexcept = default;
    RecipeApplicationService& operator=(RecipeApplicationService&&) noexcept = default;

    Foundation::Result<std::vector<Domain::Recipe>> ListPublishedRecipes();

    Foundation::Result<std::optional<Domain::Recipe>> FindPublishedRecipe(
        std::string_view Id);

    Foundation::Result<std::vector<Domain::Recommendation>> RecommendTonight(
        const Domain::RecommendationRequest& Request);

    Foundation::Result<std::vector<Domain::Ingredient>> ListIngredients();

private:
    std::unique_ptr<RecipeRepository> Repository;
    std::unique_ptr<RecommendationProvider> Provider;
    std::unique_ptr<IngredientRepository> Ingredients;
};

}  // namespace Menu::Application
