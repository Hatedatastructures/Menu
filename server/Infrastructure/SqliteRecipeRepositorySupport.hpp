#pragma once

#include <Domain/Recipe.hpp>
#include <Foundation/Result.hpp>
#include <Infrastructure/SqliteDatabase.hpp>

#include <optional>
#include <string>
#include <string_view>
#include <vector>

namespace Menu::Infrastructure::SqliteRecipeRepositorySupport {

Foundation::Result<std::vector<std::string>> ParseStringArray(
    sqlite3* Handle,
    std::string_view JsonText);

std::string SerializeStrings(const std::vector<std::string>& Values);

Foundation::Result<void> StepDone(
    SqliteDatabase& Database,
    sqlite3_stmt* Statement);

Foundation::Result<void> InsertIngredientRows(
    SqliteDatabase& Database,
    const Domain::Recipe& RecipeValue);

Foundation::Result<void> InsertStepRows(
    SqliteDatabase& Database,
    const Domain::Recipe& RecipeValue);

Foundation::Result<void> DeleteRecipeChildren(
    SqliteDatabase& Database,
    std::string_view RecipeId);

Foundation::Result<void> SaveRecipeBase(
    SqliteDatabase& Database,
    const Domain::Recipe& RecipeValue,
    bool Update);

Foundation::Result<Domain::Recipe> ReadRecipeRow(
    SqliteDatabase& Database,
    sqlite3_stmt* Statement);

Foundation::Result<Domain::Recipe> ReadPublishedRecipe(
    SqliteDatabase& Database,
    std::string_view Id);

}  // namespace Menu::Infrastructure::SqliteRecipeRepositorySupport
