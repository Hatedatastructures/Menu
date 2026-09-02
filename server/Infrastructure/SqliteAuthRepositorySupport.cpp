#include "Infrastructure/SqliteAuthRepositorySupport.hpp"

#include "Infrastructure/SqliteWorkflowSupport.hpp"
#include "Infrastructure/TokenService.hpp"

#include <string>
#include <string_view>
#include <utility>

namespace Menu::Infrastructure::SqliteAuthRepositorySupport {
namespace {

using SqliteWorkflowSupport::BindInteger;
using SqliteWorkflowSupport::BindText;
using SqliteWorkflowSupport::ColumnText;
using SqliteWorkflowSupport::Prepare;
using SqliteWorkflowSupport::StatementGuard;
using SqliteWorkflowSupport::StorageError;

Application::AuthUser ReadUser(sqlite3_stmt* Statement) {
    Application::AuthUser User;
    User.Id = ColumnText(Statement, 0);
    User.Email = ColumnText(Statement, 1);
    User.PasswordHash = ColumnText(Statement, 2);
    User.DisplayName = ColumnText(Statement, 3);
    User.IsAdmin = sqlite3_column_int(Statement, 4) != 0;
    return User;
}

}  // namespace

Foundation::Result<std::optional<Application::AuthUser>> FindUser(
    SqliteDatabase& Database,
    std::string_view Sql,
    std::string_view TokenHash,
    std::int64_t Now,
    bool HasExpiry) {
    auto StatementResult = Prepare(Database, Sql);
    if (!StatementResult.HasValue()) {
        return Foundation::Result<std::optional<Application::AuthUser>>::FromError(
            StatementResult.ErrorValue());
    }
    StatementGuard Statement = std::move(StatementResult).Value();
    if (auto Result = BindText(Statement.Get(), 1, TokenHash); !Result.HasValue()) {
        return Foundation::Result<std::optional<Application::AuthUser>>::FromError(
            Result.ErrorValue());
    }
    if (HasExpiry) {
        if (auto Result = BindInteger(Statement.Get(), 2, Now); !Result.HasValue()) {
            return Foundation::Result<std::optional<Application::AuthUser>>::FromError(
                Result.ErrorValue());
        }
    }
    const int StepResult = sqlite3_step(Statement.Get());
    if (StepResult == SQLITE_DONE) {
        return std::optional<Application::AuthUser>();
    }
    if (StepResult != SQLITE_ROW) {
        return Foundation::Result<std::optional<Application::AuthUser>>::FromError(
            StorageError(Database.NativeHandle()));
    }
    return std::optional<Application::AuthUser>(ReadUser(Statement.Get()));
}

Foundation::Result<void> InsertToken(
    SqliteDatabase& Database,
    std::string_view Table,
    std::string_view UserId,
    std::string_view TokenHash,
    std::int64_t ExpiresAt) {
    const std::string Sql = "INSERT INTO " + std::string(Table) +
                            " (Id, UserId, TokenHash, ExpiresAt) VALUES (?, ?, ?, ?);";
    auto StatementResult = Prepare(Database, Sql);
    if (!StatementResult.HasValue()) {
        return Foundation::Result<void>::FromError(StatementResult.ErrorValue());
    }
    StatementGuard Statement = std::move(StatementResult).Value();
    const auto Id = TokenService::CreateToken(16U);
    if (!Id.HasValue()) {
        return Foundation::Result<void>::FromError(Id.ErrorValue());
    }
    if (auto Result = BindText(Statement.Get(), 1, Id.Value()); !Result.HasValue()) {
        return Result;
    }
    if (auto Result = BindText(Statement.Get(), 2, UserId); !Result.HasValue()) {
        return Result;
    }
    if (auto Result = BindText(Statement.Get(), 3, TokenHash); !Result.HasValue()) {
        return Result;
    }
    if (auto Result = BindText(Statement.Get(), 4, std::to_string(ExpiresAt));
        !Result.HasValue()) {
        return Result;
    }
    if (sqlite3_step(Statement.Get()) != SQLITE_DONE) {
        return Foundation::Result<void>::FromError(StorageError(Database.NativeHandle()));
    }
    return Foundation::Result<void>();
}

Foundation::Result<void> RollbackWith(
    Foundation::Result<void> Result,
    SqliteDatabase& Database) {
    if (!Result.HasValue()) {
        Database.Rollback();
    }
    return Result;
}

}  // namespace Menu::Infrastructure::SqliteAuthRepositorySupport
