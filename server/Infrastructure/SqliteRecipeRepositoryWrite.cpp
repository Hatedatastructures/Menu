#include "Infrastructure/SqliteRecipeRepository.hpp"

#include "Infrastructure/SqliteRecipeRepositorySupport.hpp"
#include "Infrastructure/SqliteWorkflowSupport.hpp"

#include <string_view>
#include <utility>

namespace Menu::Infrastructure {
namespace {

using SqliteRecipeRepositorySupport::DeleteRecipeChildren;
using SqliteRecipeRepositorySupport::InsertIngredientRows;
using SqliteRecipeRepositorySupport::InsertStepRows;
using SqliteRecipeRepositorySupport::SaveRecipeBase;
using SqliteWorkflowSupport::BindText;
using SqliteWorkflowSupport::Prepare;
using SqliteWorkflowSupport::StatementGuard;
using SqliteWorkflowSupport::StorageError;

}  // namespace

Foundation::Result<Domain::Recipe> SqliteRecipeRepository::Create(
    const Domain::Recipe& RecipeValue) {
    const auto BeginResult = Database.BeginTransaction();
    if (!BeginResult.HasValue()) {
        return Foundation::Result<Domain::Recipe>::FromError(BeginResult.ErrorValue());
    }
    const auto BaseResult = SaveRecipeBase(Database, RecipeValue, false);
    if (!BaseResult.HasValue()) {
        Database.Rollback();
        if (sqlite3_errcode(Database.NativeHandle()) == SQLITE_CONSTRAINT) {
            return Foundation::Result<Domain::Recipe>::FromError(
                Foundation::Error(Foundation::ErrorCode::Conflict, "菜谱 ID 或 slug 已存在"));
        }
        return Foundation::Result<Domain::Recipe>::FromError(BaseResult.ErrorValue());
    }
    if (auto Result = InsertIngredientRows(Database, RecipeValue); !Result.HasValue()) {
        Database.Rollback();
        return Foundation::Result<Domain::Recipe>::FromError(Result.ErrorValue());
    }
    if (auto Result = InsertStepRows(Database, RecipeValue); !Result.HasValue()) {
        Database.Rollback();
        return Foundation::Result<Domain::Recipe>::FromError(Result.ErrorValue());
    }
    const auto CommitResult = Database.Commit();
    if (!CommitResult.HasValue()) {
        Database.Rollback();
        return Foundation::Result<Domain::Recipe>::FromError(CommitResult.ErrorValue());
    }
    return RecipeValue;
}

Foundation::Result<Domain::Recipe> SqliteRecipeRepository::Update(
    std::string_view Id,
    const Domain::Recipe& RecipeValue) {
    if (RecipeValue.Id != Id) {
        return Foundation::Result<Domain::Recipe>::FromError(
            Foundation::Error(Foundation::ErrorCode::InvalidArgument, "菜谱 ID 不一致"));
    }
    const auto BeginResult = Database.BeginTransaction();
    if (!BeginResult.HasValue()) {
        return Foundation::Result<Domain::Recipe>::FromError(BeginResult.ErrorValue());
    }
    const auto BaseResult = SaveRecipeBase(Database, RecipeValue, true);
    if (!BaseResult.HasValue()) {
        Database.Rollback();
        return Foundation::Result<Domain::Recipe>::FromError(BaseResult.ErrorValue());
    }
    if (auto Result = DeleteRecipeChildren(Database, Id); !Result.HasValue()) {
        Database.Rollback();
        return Foundation::Result<Domain::Recipe>::FromError(Result.ErrorValue());
    }
    if (auto Result = InsertIngredientRows(Database, RecipeValue); !Result.HasValue()) {
        Database.Rollback();
        return Foundation::Result<Domain::Recipe>::FromError(Result.ErrorValue());
    }
    if (auto Result = InsertStepRows(Database, RecipeValue); !Result.HasValue()) {
        Database.Rollback();
        return Foundation::Result<Domain::Recipe>::FromError(Result.ErrorValue());
    }
    const auto CommitResult = Database.Commit();
    if (!CommitResult.HasValue()) {
        Database.Rollback();
        return Foundation::Result<Domain::Recipe>::FromError(CommitResult.ErrorValue());
    }
    return RecipeValue;
}

Foundation::Result<void> SqliteRecipeRepository::Delete(std::string_view Id) {
    const auto BeginResult = Database.BeginTransaction();
    if (!BeginResult.HasValue()) {
        return BeginResult;
    }
    auto StatementResult = Prepare(Database, "DELETE FROM Recipes WHERE Id = ?;");
    if (!StatementResult.HasValue()) {
        Database.Rollback();
        return Foundation::Result<void>::FromError(StatementResult.ErrorValue());
    }
    StatementGuard Statement = std::move(StatementResult).Value();
    if (auto Result = BindText(Statement.Get(), 1, Id); !Result.HasValue()) {
        Database.Rollback();
        return Result;
    }
    if (sqlite3_step(Statement.Get()) != SQLITE_DONE) {
        Database.Rollback();
        return Foundation::Result<void>::FromError(StorageError(Database.NativeHandle()));
    }
    if (sqlite3_changes(Database.NativeHandle()) == 0) {
        Database.Rollback();
        return Foundation::Result<void>::FromError(
            Foundation::Error(Foundation::ErrorCode::NotFound, "菜谱不存在"));
    }
    const auto CommitResult = Database.Commit();
    if (!CommitResult.HasValue()) {
        Database.Rollback();
        return CommitResult;
    }
    return Foundation::Result<void>();
}

}  // namespace Menu::Infrastructure
