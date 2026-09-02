#include "Infrastructure/SqliteMealPlanRepositorySupport.hpp"

#include "Infrastructure/SqliteWorkflowSupport.hpp"

#include <string>
#include <string_view>
#include <utility>

namespace Menu::Infrastructure::SqliteMealPlanRepositorySupport {
namespace {

using SqliteWorkflowSupport::BindText;
using SqliteWorkflowSupport::ColumnText;
using SqliteWorkflowSupport::Prepare;
using SqliteWorkflowSupport::StatementGuard;
using SqliteWorkflowSupport::StorageError;

}  // namespace

Foundation::Result<std::string> FindRecipeName(
    SqliteDatabase& Database,
    std::string_view RecipeId) {
    auto StatementResult = Prepare(
        Database, "SELECT Name FROM Recipes WHERE Id = ? LIMIT 1;");
    if (!StatementResult.HasValue()) {
        return Foundation::Result<std::string>::FromError(StatementResult.ErrorValue());
    }
    StatementGuard Statement = std::move(StatementResult).Value();
    if (auto Result = BindText(Statement.Get(), 1, RecipeId); !Result.HasValue()) {
        return Foundation::Result<std::string>::FromError(Result.ErrorValue());
    }
    const int StepResult = sqlite3_step(Statement.Get());
    if (StepResult == SQLITE_DONE) {
        return Foundation::Result<std::string>::FromError(
            Foundation::Error(Foundation::ErrorCode::NotFound, "菜谱不存在"));
    }
    if (StepResult != SQLITE_ROW) {
        return Foundation::Result<std::string>::FromError(
            StorageError(Database.NativeHandle()));
    }
    return ColumnText(Statement.Get(), 0);
}

Foundation::Result<void> StepDone(
    SqliteDatabase& Database,
    sqlite3_stmt* Statement) {
    if (sqlite3_step(Statement) != SQLITE_DONE) {
        return Foundation::Result<void>::FromError(StorageError(Database.NativeHandle()));
    }
    return Foundation::Result<void>();
}

}  // namespace Menu::Infrastructure::SqliteMealPlanRepositorySupport
