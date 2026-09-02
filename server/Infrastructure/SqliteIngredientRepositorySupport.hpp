#pragma once

#include <Domain/Ingredient.hpp>
#include <Foundation/Result.hpp>
#include <Infrastructure/SqliteDatabase.hpp>

#include <string>
#include <string_view>
#include <vector>

namespace Menu::Infrastructure::SqliteIngredientRepositorySupport {

Foundation::Result<std::vector<std::string>> ReadAliases(
    SqliteDatabase& Database,
    std::string_view IngredientId);

Foundation::Result<Domain::Ingredient> ReadIngredient(
    SqliteDatabase& Database,
    std::string_view IngredientId);

Foundation::Result<void> InsertAliases(
    SqliteDatabase& Database,
    const Domain::Ingredient& IngredientValue);

}  // namespace Menu::Infrastructure::SqliteIngredientRepositorySupport
