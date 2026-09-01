#include "Application/AdminIngredientApplicationService.hpp"

#include <set>
#include <utility>

namespace Menu::Application {
namespace {

Foundation::Error InvalidIngredient(std::string Message) {
    return Foundation::Error(Foundation::ErrorCode::InvalidArgument, std::move(Message));
}

bool HasControlCharacter(std::string_view Value) {
    for (unsigned char Character : Value) {
        if (Character < 0x20U || Character == 0x7FU) {
            return true;
        }
    }
    return false;
}

Foundation::Result<void> ValidateIngredient(const Domain::Ingredient& Ingredient) {
    if (Ingredient.Id.empty() || Ingredient.Id.size() > 128U ||
        Ingredient.Name.empty() || Ingredient.Name.size() > 160U ||
        Ingredient.Category.empty() || Ingredient.Category.size() > 64U ||
        Ingredient.DefaultUnit.empty() || Ingredient.DefaultUnit.size() > 16U ||
        Ingredient.SubstituteGroup.size() > 64U ||
        Ingredient.StoreSkuMapping.size() > 256U || HasControlCharacter(Ingredient.Id) ||
        HasControlCharacter(Ingredient.Name) || HasControlCharacter(Ingredient.Category)) {
        return Foundation::Result<void>::FromError(
            InvalidIngredient("食材基础字段无效"));
    }
    static const std::set<std::string> Units = {
        "g", "kg", "ml", "l", "个", "份", "勺", "茶匙"};
    if (!Units.contains(Ingredient.DefaultUnit) || Ingredient.Aliases.size() > 16U) {
        return Foundation::Result<void>::FromError(
            InvalidIngredient("食材单位或别名无效"));
    }
    for (const std::string& Alias : Ingredient.Aliases) {
        if (Alias.empty() || Alias.size() > 64U || HasControlCharacter(Alias)) {
            return Foundation::Result<void>::FromError(
                InvalidIngredient("食材别名无效"));
        }
    }
    return Foundation::Result<void>();
}

}  // namespace

AdminIngredientApplicationService::AdminIngredientApplicationService(
    std::unique_ptr<AdminIngredientRepository> RepositoryValue)
    : Repository(std::move(RepositoryValue)) {}

Foundation::Result<std::vector<Domain::Ingredient>>
AdminIngredientApplicationService::ListAll() {
    if (!Repository) {
        return Foundation::Result<std::vector<Domain::Ingredient>>::FromError(
            Foundation::Error(Foundation::ErrorCode::StorageUnavailable, "食材管理仓储不可用"));
    }
    return Repository->ListAll();
}

Foundation::Result<Domain::Ingredient> AdminIngredientApplicationService::Create(
    const Domain::Ingredient& IngredientValue) {
    if (!Repository) {
        return Foundation::Result<Domain::Ingredient>::FromError(
            Foundation::Error(Foundation::ErrorCode::StorageUnavailable, "食材管理仓储不可用"));
    }
    const auto Validation = ValidateIngredient(IngredientValue);
    if (!Validation.HasValue()) {
        return Foundation::Result<Domain::Ingredient>::FromError(Validation.ErrorValue());
    }
    return Repository->Create(IngredientValue);
}

Foundation::Result<Domain::Ingredient> AdminIngredientApplicationService::Update(
    std::string_view Id,
    const Domain::Ingredient& IngredientValue) {
    if (!Repository) {
        return Foundation::Result<Domain::Ingredient>::FromError(
            Foundation::Error(Foundation::ErrorCode::StorageUnavailable, "食材管理仓储不可用"));
    }
    if (Id != IngredientValue.Id) {
        return Foundation::Result<Domain::Ingredient>::FromError(
            InvalidIngredient("食材 ID 不一致"));
    }
    const auto Validation = ValidateIngredient(IngredientValue);
    if (!Validation.HasValue()) {
        return Foundation::Result<Domain::Ingredient>::FromError(Validation.ErrorValue());
    }
    return Repository->Update(Id, IngredientValue);
}

Foundation::Result<void> AdminIngredientApplicationService::Delete(std::string_view Id) {
    if (!Repository) {
        return Foundation::Result<void>::FromError(
            Foundation::Error(Foundation::ErrorCode::StorageUnavailable, "食材管理仓储不可用"));
    }
    return Repository->Delete(Id);
}

}  // namespace Menu::Application
