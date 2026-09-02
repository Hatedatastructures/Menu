#pragma once

#include <Domain/Ingredient.hpp>
#include <Domain/Recipe.hpp>
#include <Foundation/Result.hpp>
#include <Infrastructure/SqliteDatabase.hpp>

namespace Menu::Infrastructure::SeedDataRows {

Foundation::Result<void> InsertIngredient(
    SqliteDatabase& Database,
    const Domain::Ingredient& IngredientValue);

Foundation::Result<void> InsertRecipe(
    SqliteDatabase& Database,
    const Domain::Recipe& RecipeValue);

}  // namespace Menu::Infrastructure::SeedDataRows
