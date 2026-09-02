#include "Infrastructure/SeedDataRows.hpp"

#include "Infrastructure/SqliteWorkflowSupport.hpp"

#include <string>
#include <string_view>
#include <utility>

namespace Menu::Infrastructure::SeedDataRows {
namespace {

using SqliteWorkflowSupport::BindInt;
using SqliteWorkflowSupport::BindText;
using SqliteWorkflowSupport::Prepare;
using SqliteWorkflowSupport::StatementGuard;
using SqliteWorkflowSupport::StorageError;

Foundation::Result<void> StepDone(
    SqliteDatabase& Database,
    sqlite3_stmt* Statement) {
    if (sqlite3_step(Statement) != SQLITE_DONE) {
        return Foundation::Result<void>::FromError(StorageError(Database.NativeHandle()));
    }
    return Foundation::Result<void>();
}

}  // namespace

Foundation::Result<void> InsertIngredient(
    SqliteDatabase& Database,
    const Domain::Ingredient& IngredientValue) {
    auto StatementResult = Prepare(
        Database,
        "INSERT OR IGNORE INTO Ingredients "
        "(Id, Name, Category, DefaultUnit, IsPantryStaple, SubstituteGroup, StoreSkuMapping) "
        "VALUES (?, ?, ?, ?, ?, ?, ?);");
    if (!StatementResult.HasValue()) {
        return Foundation::Result<void>::FromError(StatementResult.ErrorValue());
    }
    StatementGuard Statement = std::move(StatementResult).Value();
    const std::string_view TextValues[] = {
        IngredientValue.Id, IngredientValue.Name, IngredientValue.Category,
        IngredientValue.DefaultUnit};
    for (std::size_t Index = 0; Index < std::size(TextValues); ++Index) {
        if (auto Result = BindText(
                Statement.Get(), static_cast<int>(Index + 1), TextValues[Index]);
            !Result.HasValue()) {
            return Result;
        }
    }
    if (auto Result = BindInt(Statement.Get(), 5, IngredientValue.IsPantryStaple ? 1 : 0);
        !Result.HasValue()) {
        return Result;
    }
    if (auto Result = BindText(Statement.Get(), 6, IngredientValue.SubstituteGroup);
        !Result.HasValue()) {
        return Result;
    }
    if (auto Result = BindText(Statement.Get(), 7, IngredientValue.StoreSkuMapping);
        !Result.HasValue()) {
        return Result;
    }
    if (auto Result = StepDone(Database, Statement.Get()); !Result.HasValue()) {
        return Result;
    }
    for (const std::string& Alias : IngredientValue.Aliases) {
        auto AliasStatementResult = Prepare(
            Database,
            "INSERT OR IGNORE INTO IngredientAliases (IngredientId, Alias) VALUES (?, ?);");
        if (!AliasStatementResult.HasValue()) {
            return Foundation::Result<void>::FromError(AliasStatementResult.ErrorValue());
        }
        StatementGuard AliasStatement = std::move(AliasStatementResult).Value();
        if (auto Result = BindText(AliasStatement.Get(), 1, IngredientValue.Id);
            !Result.HasValue()) {
            return Result;
        }
        if (auto Result = BindText(AliasStatement.Get(), 2, Alias); !Result.HasValue()) {
            return Result;
        }
        if (auto Result = StepDone(Database, AliasStatement.Get()); !Result.HasValue()) {
            return Result;
        }
    }
    return Foundation::Result<void>();
}

}  // namespace Menu::Infrastructure::SeedDataRows
