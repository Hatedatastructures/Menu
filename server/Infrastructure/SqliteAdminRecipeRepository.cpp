#include "Infrastructure/SqliteAdminRecipeRepository.hpp"

namespace Menu::Infrastructure {

Foundation::Result<std::vector<Domain::Recipe>>
SqliteAdminRecipeRepository::ListAll() {
    return Store.ListAll();
}

Foundation::Result<Domain::Recipe> SqliteAdminRecipeRepository::Create(
    const Domain::Recipe& RecipeValue) {
    return Store.Create(RecipeValue);
}

Foundation::Result<Domain::Recipe> SqliteAdminRecipeRepository::Update(
    std::string_view Id,
    const Domain::Recipe& RecipeValue) {
    return Store.Update(Id, RecipeValue);
}

Foundation::Result<void> SqliteAdminRecipeRepository::Delete(std::string_view Id) {
    return Store.Delete(Id);
}

}  // namespace Menu::Infrastructure
