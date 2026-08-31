#include "Application/RecipeApplicationService.hpp"

#include <utility>

namespace Menu::Application {

RecipeApplicationService::RecipeApplicationService(
    std::unique_ptr<RecipeRepository> RepositoryValue,
    std::unique_ptr<RecommendationProvider> ProviderValue)
    : Repository(std::move(RepositoryValue)), Provider(std::move(ProviderValue)) {}

Foundation::Result<std::vector<Domain::Recipe>>
RecipeApplicationService::ListPublishedRecipes() {
    if (!Repository) {
        return Foundation::Result<std::vector<Domain::Recipe>>::FromError(
            Foundation::Error(Foundation::ErrorCode::StorageUnavailable, "菜谱仓储不可用"));
    }
    return Repository->ListPublished();
}

Foundation::Result<std::optional<Domain::Recipe>>
RecipeApplicationService::FindPublishedRecipe(std::string_view Id) {
    if (!Repository) {
        return Foundation::Result<std::optional<Domain::Recipe>>::FromError(
            Foundation::Error(Foundation::ErrorCode::StorageUnavailable, "菜谱仓储不可用"));
    }
    return Repository->FindPublishedById(Id);
}

Foundation::Result<std::vector<Domain::Recommendation>>
RecipeApplicationService::RecommendTonight(const Domain::RecommendationRequest& Request) {
    if (!Provider) {
        return Foundation::Result<std::vector<Domain::Recommendation>>::FromError(
            Foundation::Error(Foundation::ErrorCode::StorageUnavailable, "推荐服务不可用"));
    }

    const auto RecipesResult = ListPublishedRecipes();
    if (!RecipesResult.HasValue()) {
        return Foundation::Result<std::vector<Domain::Recommendation>>::FromError(
            RecipesResult.ErrorValue());
    }
    return Provider->Recommend(Request, RecipesResult.Value());
}

}  // namespace Menu::Application
