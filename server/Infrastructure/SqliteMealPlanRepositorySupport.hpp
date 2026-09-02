#pragma once

#include <Foundation/Result.hpp>
#include <Infrastructure/SqliteDatabase.hpp>

#include <string>
#include <string_view>

namespace Menu::Infrastructure::SqliteMealPlanRepositorySupport {

Foundation::Result<std::string> FindRecipeName(
    SqliteDatabase& Database,
    std::string_view RecipeId);

Foundation::Result<void> StepDone(
    SqliteDatabase& Database,
    sqlite3_stmt* Statement);

}  // namespace Menu::Infrastructure::SqliteMealPlanRepositorySupport
