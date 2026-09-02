#pragma once

#include <Application/AdminRecipeRepository.hpp>
#include <Infrastructure/SqliteRecipeRepository.hpp>

namespace Menu::Infrastructure {

class SqliteAdminRecipeRepository final : public Application::AdminRecipeRepository {
public:
    explicit SqliteAdminRecipeRepository(SqliteDatabase& DatabaseValue)
        : Store(DatabaseValue) {}

    Foundation::Result<std::vector<Domain::Recipe>> ListAll() override;
    Foundation::Result<Domain::Recipe> Create(
        const Domain::Recipe& RecipeValue) override;
    Foundation::Result<Domain::Recipe> Update(
        std::string_view Id,
        const Domain::Recipe& RecipeValue) override;
    Foundation::Result<void> Delete(std::string_view Id) override;

private:
    SqliteRecipeRepository Store;
};

}  // namespace Menu::Infrastructure
