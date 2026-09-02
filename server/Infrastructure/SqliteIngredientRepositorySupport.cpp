#include "Infrastructure/SqliteIngredientRepositorySupport.hpp"

#include "Infrastructure/SqliteWorkflowSupport.hpp"

#include <string>
#include <string_view>
#include <utility>

namespace Menu::Infrastructure::SqliteIngredientRepositorySupport {
namespace {

using SqliteWorkflowSupport::BindText;
using SqliteWorkflowSupport::ColumnText;
using SqliteWorkflowSupport::Prepare;
using SqliteWorkflowSupport::StatementGuard;
using SqliteWorkflowSupport::StorageError;

}  // namespace

Foundation::Result<std::vector<std::string>> ReadAliases(
    SqliteDatabase& Database,
    std::string_view IngredientId) {
    auto StatementResult = Prepare(
        Database,
        "SELECT Alias FROM IngredientAliases WHERE IngredientId = ? ORDER BY Alias;");
    if (!StatementResult.HasValue()) {
        return Foundation::Result<std::vector<std::string>>::FromError(
            StatementResult.ErrorValue());
    }
    StatementGuard Statement = std::move(StatementResult).Value();
    if (auto Result = BindText(Statement.Get(), 1, IngredientId); !Result.HasValue()) {
        return Foundation::Result<std::vector<std::string>>::FromError(Result.ErrorValue());
    }
    std::vector<std::string> Aliases;
    while (true) {
        const int StepResult = sqlite3_step(Statement.Get());
        if (StepResult == SQLITE_DONE) {
            break;
        }
        if (StepResult != SQLITE_ROW) {
            return Foundation::Result<std::vector<std::string>>::FromError(
                StorageError(Database.NativeHandle()));
        }
        Aliases.push_back(ColumnText(Statement.Get(), 0));
    }
    return Aliases;
}

Foundation::Result<Domain::Ingredient> ReadIngredient(
    SqliteDatabase& Database,
    std::string_view IngredientId) {
    auto StatementResult = Prepare(
        Database,
        "SELECT Id, Name, Category, DefaultUnit, IsPantryStaple, SubstituteGroup, "
        "StoreSkuMapping FROM Ingredients WHERE Id = ? LIMIT 1;");
    if (!StatementResult.HasValue()) {
        return Foundation::Result<Domain::Ingredient>::FromError(
            StatementResult.ErrorValue());
    }
    StatementGuard Statement = std::move(StatementResult).Value();
    if (auto Result = BindText(Statement.Get(), 1, IngredientId); !Result.HasValue()) {
        return Foundation::Result<Domain::Ingredient>::FromError(Result.ErrorValue());
    }
    const int StepResult = sqlite3_step(Statement.Get());
    if (StepResult == SQLITE_DONE) {
        return Foundation::Result<Domain::Ingredient>::FromError(
            Foundation::Error(Foundation::ErrorCode::NotFound, "食材不存在"));
    }
    if (StepResult != SQLITE_ROW) {
        return Foundation::Result<Domain::Ingredient>::FromError(
            StorageError(Database.NativeHandle()));
    }
    Domain::Ingredient IngredientValue;
    IngredientValue.Id = ColumnText(Statement.Get(), 0);
    IngredientValue.Name = ColumnText(Statement.Get(), 1);
    IngredientValue.Category = ColumnText(Statement.Get(), 2);
    IngredientValue.DefaultUnit = ColumnText(Statement.Get(), 3);
    IngredientValue.IsPantryStaple = sqlite3_column_int(Statement.Get(), 4) != 0;
    IngredientValue.SubstituteGroup = ColumnText(Statement.Get(), 5);
    IngredientValue.StoreSkuMapping = ColumnText(Statement.Get(), 6);
    const auto Aliases = ReadAliases(Database, IngredientValue.Id);
    if (!Aliases.HasValue()) {
        return Foundation::Result<Domain::Ingredient>::FromError(Aliases.ErrorValue());
    }
    IngredientValue.Aliases = Aliases.Value();
    return IngredientValue;
}

Foundation::Result<void> InsertAliases(
    SqliteDatabase& Database,
    const Domain::Ingredient& IngredientValue) {
    for (const std::string& Alias : IngredientValue.Aliases) {
        auto StatementResult = Prepare(
            Database,
            "INSERT INTO IngredientAliases (IngredientId, Alias) VALUES (?, ?);");
        if (!StatementResult.HasValue()) {
            return Foundation::Result<void>::FromError(StatementResult.ErrorValue());
        }
        StatementGuard Statement = std::move(StatementResult).Value();
        if (auto Result = BindText(Statement.Get(), 1, IngredientValue.Id);
            !Result.HasValue()) {
            return Result;
        }
        if (auto Result = BindText(Statement.Get(), 2, Alias); !Result.HasValue()) {
            return Result;
        }
        if (sqlite3_step(Statement.Get()) != SQLITE_DONE) {
            return Foundation::Result<void>::FromError(StorageError(Database.NativeHandle()));
        }
    }
    return Foundation::Result<void>();
}

}  // namespace Menu::Infrastructure::SqliteIngredientRepositorySupport
