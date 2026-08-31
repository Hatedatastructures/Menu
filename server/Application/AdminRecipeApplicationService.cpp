#include "Application/AdminRecipeApplicationService.hpp"

#include <Domain/RecipeRules.hpp>

#include <utility>

namespace Menu::Application {

AdminRecipeApplicationService::AdminRecipeApplicationService(
    std::unique_ptr<AdminRecipeRepository> RepositoryValue)
    : Repository(std::move(RepositoryValue)) {}

Foundation::Result<std::vector<Domain::Recipe>>
AdminRecipeApplicationService::ListAllRecipes() {
    if (!Repository) {
        return Foundation::Result<std::vector<Domain::Recipe>>::FromError(
            Foundation::Error(Foundation::ErrorCode::StorageUnavailable, "管理仓储不可用"));
    }
    return Repository->ListAll();
}

Foundation::Result<Domain::Recipe>
AdminRecipeApplicationService::CreateRecipe(const Domain::Recipe& RecipeValue) {
    if (!Repository) {
        return Foundation::Result<Domain::Recipe>::FromError(
            Foundation::Error(Foundation::ErrorCode::StorageUnavailable, "管理仓储不可用"));
    }
    const auto Validation = Domain::RecipeRules::ValidateRecipe(RecipeValue);
    if (!Validation.HasValue()) {
        return Foundation::Result<Domain::Recipe>::FromError(Validation.ErrorValue());
    }
    return Repository->Create(RecipeValue);
}

Foundation::Result<Domain::Recipe>
AdminRecipeApplicationService::UpdateRecipe(
    std::string_view Id,
    const Domain::Recipe& RecipeValue) {
    if (!Repository) {
        return Foundation::Result<Domain::Recipe>::FromError(
            Foundation::Error(Foundation::ErrorCode::StorageUnavailable, "管理仓储不可用"));
    }
    if (Id != RecipeValue.Id) {
        return Foundation::Result<Domain::Recipe>::FromError(
            Foundation::Error(Foundation::ErrorCode::InvalidArgument, "菜谱 ID 不一致"));
    }
    const auto Validation = Domain::RecipeRules::ValidateRecipe(RecipeValue);
    if (!Validation.HasValue()) {
        return Foundation::Result<Domain::Recipe>::FromError(Validation.ErrorValue());
    }
    return Repository->Update(Id, RecipeValue);
}

Foundation::Result<void>
AdminRecipeApplicationService::DeleteRecipe(std::string_view Id) {
    if (!Repository) {
        return Foundation::Result<void>::FromError(
            Foundation::Error(Foundation::ErrorCode::StorageUnavailable, "管理仓储不可用"));
    }
    return Repository->Delete(Id);
}

}  // namespace Menu::Application
