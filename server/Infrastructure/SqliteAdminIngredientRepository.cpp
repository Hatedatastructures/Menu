#include "Infrastructure/SqliteAdminIngredientRepository.hpp"

namespace Menu::Infrastructure {

Foundation::Result<std::vector<Domain::Ingredient>>
SqliteAdminIngredientRepository::ListAll() {
    return Store.ListAll();
}

Foundation::Result<Domain::Ingredient> SqliteAdminIngredientRepository::Create(
    const Domain::Ingredient& IngredientValue) {
    return Store.Create(IngredientValue);
}

Foundation::Result<Domain::Ingredient> SqliteAdminIngredientRepository::Update(
    std::string_view Id,
    const Domain::Ingredient& IngredientValue) {
    return Store.Update(Id, IngredientValue);
}

Foundation::Result<void> SqliteAdminIngredientRepository::Delete(std::string_view Id) {
    return Store.Delete(Id);
}

}  // namespace Menu::Infrastructure
