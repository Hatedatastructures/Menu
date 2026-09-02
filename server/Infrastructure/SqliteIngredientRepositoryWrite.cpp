#include "Infrastructure/SqliteIngredientRepository.hpp"

#include "Infrastructure/SqliteIngredientRepositorySupport.hpp"
#include "Infrastructure/SqliteWorkflowSupport.hpp"

#include <string_view>
#include <utility>

namespace Menu::Infrastructure {
namespace {

using SqliteIngredientRepositorySupport::InsertAliases;
using SqliteIngredientRepositorySupport::ReadIngredient;
using SqliteWorkflowSupport::BindInteger;
using SqliteWorkflowSupport::BindText;
using SqliteWorkflowSupport::Prepare;
using SqliteWorkflowSupport::StatementGuard;
using SqliteWorkflowSupport::StorageError;

}  // namespace

Foundation::Result<Domain::Ingredient> SqliteIngredientRepository::Create(
    const Domain::Ingredient& IngredientValue) {
    const auto BeginResult = Database.BeginTransaction();
    if (!BeginResult.HasValue()) {
        return Foundation::Result<Domain::Ingredient>::FromError(BeginResult.ErrorValue());
    }
    auto StatementResult = Prepare(
        Database,
        "INSERT INTO Ingredients "
        "(Id, Name, Category, DefaultUnit, IsPantryStaple, SubstituteGroup, StoreSkuMapping) "
        "VALUES (?, ?, ?, ?, ?, ?, ?);");
    if (!StatementResult.HasValue()) {
        Database.Rollback();
        return Foundation::Result<Domain::Ingredient>::FromError(StatementResult.ErrorValue());
    }
    StatementGuard Statement = std::move(StatementResult).Value();
    if (auto Result = BindText(Statement.Get(), 1, IngredientValue.Id); !Result.HasValue()) {
        Database.Rollback();
        return Foundation::Result<Domain::Ingredient>::FromError(Result.ErrorValue());
    }
    if (auto Result = BindText(Statement.Get(), 2, IngredientValue.Name); !Result.HasValue()) {
        Database.Rollback();
        return Foundation::Result<Domain::Ingredient>::FromError(Result.ErrorValue());
    }
    if (auto Result = BindText(Statement.Get(), 3, IngredientValue.Category); !Result.HasValue()) {
        Database.Rollback();
        return Foundation::Result<Domain::Ingredient>::FromError(Result.ErrorValue());
    }
    if (auto Result = BindText(Statement.Get(), 4, IngredientValue.DefaultUnit);
        !Result.HasValue()) {
        Database.Rollback();
        return Foundation::Result<Domain::Ingredient>::FromError(Result.ErrorValue());
    }
    if (auto Result = BindInteger(
            Statement.Get(), 5, IngredientValue.IsPantryStaple ? 1 : 0);
        !Result.HasValue()) {
        Database.Rollback();
        return Foundation::Result<Domain::Ingredient>::FromError(Result.ErrorValue());
    }
    if (auto Result = BindText(Statement.Get(), 6, IngredientValue.SubstituteGroup);
        !Result.HasValue()) {
        Database.Rollback();
        return Foundation::Result<Domain::Ingredient>::FromError(Result.ErrorValue());
    }
    if (auto Result = BindText(Statement.Get(), 7, IngredientValue.StoreSkuMapping);
        !Result.HasValue()) {
        Database.Rollback();
        return Foundation::Result<Domain::Ingredient>::FromError(Result.ErrorValue());
    }
    const int StepResult = sqlite3_step(Statement.Get());
    if (StepResult != SQLITE_DONE) {
        Database.Rollback();
        if (StepResult == SQLITE_CONSTRAINT) {
            return Foundation::Result<Domain::Ingredient>::FromError(
                Foundation::Error(Foundation::ErrorCode::Conflict, "食材已存在"));
        }
        return Foundation::Result<Domain::Ingredient>::FromError(
            StorageError(Database.NativeHandle()));
    }
    const auto AliasResult = InsertAliases(Database, IngredientValue);
    if (!AliasResult.HasValue()) {
        Database.Rollback();
        if (sqlite3_errcode(Database.NativeHandle()) == SQLITE_CONSTRAINT) {
            return Foundation::Result<Domain::Ingredient>::FromError(
                Foundation::Error(Foundation::ErrorCode::Conflict, "食材别名已存在"));
        }
        return Foundation::Result<Domain::Ingredient>::FromError(AliasResult.ErrorValue());
    }
    const auto CommitResult = Database.Commit();
    if (!CommitResult.HasValue()) {
        Database.Rollback();
        return Foundation::Result<Domain::Ingredient>::FromError(CommitResult.ErrorValue());
    }
    return ReadIngredient(Database, IngredientValue.Id);
}

Foundation::Result<Domain::Ingredient> SqliteIngredientRepository::Update(
    std::string_view Id,
    const Domain::Ingredient& IngredientValue) {
    const auto BeginResult = Database.BeginTransaction();
    if (!BeginResult.HasValue()) {
        return Foundation::Result<Domain::Ingredient>::FromError(BeginResult.ErrorValue());
    }
    auto StatementResult = Prepare(
        Database,
        "UPDATE Ingredients SET Name = ?, Category = ?, DefaultUnit = ?, "
        "IsPantryStaple = ?, SubstituteGroup = ?, StoreSkuMapping = ? WHERE Id = ?;");
    if (!StatementResult.HasValue()) {
        Database.Rollback();
        return Foundation::Result<Domain::Ingredient>::FromError(StatementResult.ErrorValue());
    }
    StatementGuard Statement = std::move(StatementResult).Value();
    if (auto Result = BindText(Statement.Get(), 1, IngredientValue.Name); !Result.HasValue()) {
        Database.Rollback();
        return Foundation::Result<Domain::Ingredient>::FromError(Result.ErrorValue());
    }
    if (auto Result = BindText(Statement.Get(), 2, IngredientValue.Category); !Result.HasValue()) {
        Database.Rollback();
        return Foundation::Result<Domain::Ingredient>::FromError(Result.ErrorValue());
    }
    if (auto Result = BindText(Statement.Get(), 3, IngredientValue.DefaultUnit);
        !Result.HasValue()) {
        Database.Rollback();
        return Foundation::Result<Domain::Ingredient>::FromError(Result.ErrorValue());
    }
    if (auto Result = BindInteger(
            Statement.Get(), 4, IngredientValue.IsPantryStaple ? 1 : 0);
        !Result.HasValue()) {
        Database.Rollback();
        return Foundation::Result<Domain::Ingredient>::FromError(Result.ErrorValue());
    }
    if (auto Result = BindText(Statement.Get(), 5, IngredientValue.SubstituteGroup);
        !Result.HasValue()) {
        Database.Rollback();
        return Foundation::Result<Domain::Ingredient>::FromError(Result.ErrorValue());
    }
    if (auto Result = BindText(Statement.Get(), 6, IngredientValue.StoreSkuMapping);
        !Result.HasValue()) {
        Database.Rollback();
        return Foundation::Result<Domain::Ingredient>::FromError(Result.ErrorValue());
    }
    if (auto Result = BindText(Statement.Get(), 7, Id); !Result.HasValue()) {
        Database.Rollback();
        return Foundation::Result<Domain::Ingredient>::FromError(Result.ErrorValue());
    }
    const int StepResult = sqlite3_step(Statement.Get());
    if (StepResult != SQLITE_DONE) {
        Database.Rollback();
        if (StepResult == SQLITE_CONSTRAINT) {
            return Foundation::Result<Domain::Ingredient>::FromError(
                Foundation::Error(Foundation::ErrorCode::Conflict, "食材名称已存在"));
        }
        return Foundation::Result<Domain::Ingredient>::FromError(
            StorageError(Database.NativeHandle()));
    }
    if (sqlite3_changes(Database.NativeHandle()) == 0) {
        Database.Rollback();
        return Foundation::Result<Domain::Ingredient>::FromError(
            Foundation::Error(Foundation::ErrorCode::NotFound, "食材不存在"));
    }
    auto DeleteResult = Prepare(
        Database, "DELETE FROM IngredientAliases WHERE IngredientId = ?;");
    if (!DeleteResult.HasValue()) {
        Database.Rollback();
        return Foundation::Result<Domain::Ingredient>::FromError(DeleteResult.ErrorValue());
    }
    StatementGuard DeleteStatement = std::move(DeleteResult).Value();
    if (auto Result = BindText(DeleteStatement.Get(), 1, Id); !Result.HasValue()) {
        Database.Rollback();
        return Foundation::Result<Domain::Ingredient>::FromError(Result.ErrorValue());
    }
    if (sqlite3_step(DeleteStatement.Get()) != SQLITE_DONE) {
        Database.Rollback();
        return Foundation::Result<Domain::Ingredient>::FromError(
            StorageError(Database.NativeHandle()));
    }
    const auto AliasResult = InsertAliases(Database, IngredientValue);
    if (!AliasResult.HasValue()) {
        Database.Rollback();
        return Foundation::Result<Domain::Ingredient>::FromError(AliasResult.ErrorValue());
    }
    const auto CommitResult = Database.Commit();
    if (!CommitResult.HasValue()) {
        Database.Rollback();
        return Foundation::Result<Domain::Ingredient>::FromError(CommitResult.ErrorValue());
    }
    return ReadIngredient(Database, Id);
}

Foundation::Result<void> SqliteIngredientRepository::Delete(std::string_view Id) {
    const auto BeginResult = Database.BeginTransaction();
    if (!BeginResult.HasValue()) {
        return BeginResult;
    }
    auto StatementResult = Prepare(Database, "DELETE FROM Ingredients WHERE Id = ?;");
    if (!StatementResult.HasValue()) {
        Database.Rollback();
        return Foundation::Result<void>::FromError(StatementResult.ErrorValue());
    }
    StatementGuard Statement = std::move(StatementResult).Value();
    if (auto Result = BindText(Statement.Get(), 1, Id); !Result.HasValue()) {
        Database.Rollback();
        return Result;
    }
    const int StepResult = sqlite3_step(Statement.Get());
    if (StepResult != SQLITE_DONE) {
        Database.Rollback();
        if (StepResult == SQLITE_CONSTRAINT) {
            return Foundation::Result<void>::FromError(
                Foundation::Error(Foundation::ErrorCode::Conflict, "食材已被菜谱使用"));
        }
        return Foundation::Result<void>::FromError(StorageError(Database.NativeHandle()));
    }
    if (sqlite3_changes(Database.NativeHandle()) == 0) {
        Database.Rollback();
        return Foundation::Result<void>::FromError(
            Foundation::Error(Foundation::ErrorCode::NotFound, "食材不存在"));
    }
    const auto CommitResult = Database.Commit();
    if (!CommitResult.HasValue()) {
        Database.Rollback();
        return CommitResult;
    }
    return Foundation::Result<void>();
}

}  // namespace Menu::Infrastructure
