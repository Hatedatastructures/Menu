#pragma once

#include <Domain/Recipe.hpp>
#include <Foundation/Result.hpp>
#include <Infrastructure/SqliteDatabase.hpp>

#include <string>
#include <string_view>
#include <cstdint>
#include <vector>

namespace Menu::Infrastructure::SqliteWorkflowSupport {

class StatementGuard final {
public:
    explicit StatementGuard(sqlite3_stmt* StatementValue);
    StatementGuard(const StatementGuard&) = delete;
    StatementGuard& operator=(const StatementGuard&) = delete;
    StatementGuard(StatementGuard&& Other) noexcept;
    StatementGuard& operator=(StatementGuard&& Other) noexcept;
    ~StatementGuard();

    [[nodiscard]] sqlite3_stmt* Get() const noexcept;

private:
    void Reset() noexcept;

    sqlite3_stmt* Statement = nullptr;
};

Foundation::Result<StatementGuard> Prepare(
    SqliteDatabase& Database,
    std::string_view Sql);

Foundation::Result<void> BindText(
    sqlite3_stmt* Statement,
    int Index,
    std::string_view Value);

Foundation::Result<void> BindInteger(
    sqlite3_stmt* Statement,
    int Index,
    std::int64_t Value);

std::string ColumnText(sqlite3_stmt* Statement, int Column);

Foundation::Error StorageError(sqlite3* Handle);

Foundation::Result<std::vector<Domain::RecipeIngredient>> ReadRecipeIngredients(
    SqliteDatabase& Database,
    std::string_view RecipeId,
    int RecipeServings,
    int RequestedServings);

std::string CurrentTimestamp();

}  // namespace Menu::Infrastructure::SqliteWorkflowSupport
