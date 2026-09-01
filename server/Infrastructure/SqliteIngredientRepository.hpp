#pragma once

#include <Application/AdminIngredientRepository.hpp>
#include <Application/IngredientRepository.hpp>
#include <Infrastructure/SqliteDatabase.hpp>

namespace Menu::Infrastructure {

class SqliteIngredientRepository final
    : public Application::IngredientRepository,
      public Application::AdminIngredientRepository {
public:
    explicit SqliteIngredientRepository(SqliteDatabase& DatabaseValue)
        : Database(DatabaseValue) {}

    Foundation::Result<std::vector<Domain::Ingredient>> ListAll() override;

    Foundation::Result<Domain::Ingredient> Create(
        const Domain::Ingredient& IngredientValue) override;

    Foundation::Result<Domain::Ingredient> Update(
        std::string_view Id,
        const Domain::Ingredient& IngredientValue) override;

    Foundation::Result<void> Delete(std::string_view Id) override;

private:
    SqliteDatabase& Database;
};

}  // namespace Menu::Infrastructure
