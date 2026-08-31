#pragma once

#include <Application/IngredientRepository.hpp>
#include <Infrastructure/SqliteDatabase.hpp>

namespace Menu::Infrastructure {

class SqliteIngredientRepository final : public Application::IngredientRepository {
public:
    explicit SqliteIngredientRepository(SqliteDatabase& DatabaseValue)
        : Database(DatabaseValue) {}

    Foundation::Result<std::vector<Domain::Ingredient>> ListAll() override;

private:
    SqliteDatabase& Database;
};

}  // namespace Menu::Infrastructure
