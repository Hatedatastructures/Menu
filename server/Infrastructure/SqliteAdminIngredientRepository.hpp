#pragma once

#include <Application/AdminIngredientRepository.hpp>
#include <Infrastructure/SqliteIngredientRepository.hpp>

namespace Menu::Infrastructure {

class SqliteAdminIngredientRepository final
    : public Application::AdminIngredientRepository {
public:
    explicit SqliteAdminIngredientRepository(SqliteDatabase& DatabaseValue)
        : Store(DatabaseValue) {}

    Foundation::Result<std::vector<Domain::Ingredient>> ListAll() override;
    Foundation::Result<Domain::Ingredient> Create(
        const Domain::Ingredient& IngredientValue) override;
    Foundation::Result<Domain::Ingredient> Update(
        std::string_view Id,
        const Domain::Ingredient& IngredientValue) override;
    Foundation::Result<void> Delete(std::string_view Id) override;

private:
    SqliteIngredientRepository Store;
};

}  // namespace Menu::Infrastructure
