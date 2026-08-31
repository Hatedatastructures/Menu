#pragma once

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
        std::unique_ptr<RecommendationProvider> ProviderValue);

    RecipeApplicationService(const RecipeApplicationService&) = delete;
    RecipeApplicationService& operator=(const RecipeApplicationService&) = delete;
    RecipeApplicationService(RecipeApplicationService&&) noexcept = default;
    RecipeApplicationService& operator=(RecipeApplicationService&&) noexcept = default;

    Foundation::Result<std::vector<Domain::Recipe>> ListPublishedRecipes();

    Foundation::Result<std::optional<Domain::Recipe>> FindPublishedRecipe(
        std::string_view Id);

    Foundation::Result<std::vector<Domain::Recommendation>> RecommendTonight(
        const Domain::RecommendationRequest& Request);

private:
    std::unique_ptr<RecipeRepository> Repository;
    std::unique_ptr<RecommendationProvider> Provider;
};

}  // namespace Menu::Application
