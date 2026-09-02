#include "Infrastructure/SqliteAuthRepository.hpp"

#include "Infrastructure/SqliteAuthRepositorySupport.hpp"
#include "Infrastructure/SqliteWorkflowSupport.hpp"

#include <utility>

namespace Menu::Infrastructure {
namespace {

using SqliteAuthRepositorySupport::FindUser;
using SqliteWorkflowSupport::BindInteger;
using SqliteWorkflowSupport::BindText;
using SqliteWorkflowSupport::Prepare;
using SqliteWorkflowSupport::StatementGuard;
using SqliteWorkflowSupport::StorageError;

}  // namespace

Foundation::Result<std::size_t> SqliteAuthRepository::CountUsers() {
    const auto Count = Database.ScalarInt("SELECT COUNT(*) FROM Users;");
    if (!Count.HasValue()) {
        return Foundation::Result<std::size_t>::FromError(Count.ErrorValue());
    }
    return static_cast<std::size_t>(Count.Value());
}

Foundation::Result<std::optional<Application::AuthUser>>
SqliteAuthRepository::FindUserByEmail(std::string_view Email) {
    return FindUser(
        Database,
        "SELECT Id, Email, PasswordHash, DisplayName, IsAdmin FROM Users WHERE Email = ? LIMIT 1;",
        Email,
        0,
        false);
}

Foundation::Result<Application::AuthUser>
SqliteAuthRepository::CreateUser(Application::AuthUser User) {
    const auto BeginResult = Database.BeginTransaction();
    if (!BeginResult.HasValue()) {
        return Foundation::Result<Application::AuthUser>::FromError(BeginResult.ErrorValue());
    }
    const auto Count = CountUsers();
    if (!Count.HasValue()) {
        Database.Rollback();
        return Foundation::Result<Application::AuthUser>::FromError(Count.ErrorValue());
    }
    User.IsAdmin = Count.Value() == 0U;
    auto StatementResult = Prepare(
        Database,
        "INSERT INTO Users (Id, Email, PasswordHash, DisplayName, IsAdmin) VALUES (?, ?, ?, ?, ?);");
    if (!StatementResult.HasValue()) {
        Database.Rollback();
        return Foundation::Result<Application::AuthUser>::FromError(StatementResult.ErrorValue());
    }
    StatementGuard Statement = std::move(StatementResult).Value();
    if (auto Result = BindText(Statement.Get(), 1, User.Id); !Result.HasValue()) {
        Database.Rollback();
        return Foundation::Result<Application::AuthUser>::FromError(Result.ErrorValue());
    }
    if (auto Result = BindText(Statement.Get(), 2, User.Email); !Result.HasValue()) {
        Database.Rollback();
        return Foundation::Result<Application::AuthUser>::FromError(Result.ErrorValue());
    }
    if (auto Result = BindText(Statement.Get(), 3, User.PasswordHash); !Result.HasValue()) {
        Database.Rollback();
        return Foundation::Result<Application::AuthUser>::FromError(Result.ErrorValue());
    }
    if (auto Result = BindText(Statement.Get(), 4, User.DisplayName); !Result.HasValue()) {
        Database.Rollback();
        return Foundation::Result<Application::AuthUser>::FromError(Result.ErrorValue());
    }
    if (auto Result = BindInteger(Statement.Get(), 5, User.IsAdmin ? 1 : 0);
        !Result.HasValue()) {
        Database.Rollback();
        return Foundation::Result<Application::AuthUser>::FromError(Result.ErrorValue());
    }
    const int StepResult = sqlite3_step(Statement.Get());
    if (StepResult != SQLITE_DONE) {
        Database.Rollback();
        if (StepResult == SQLITE_CONSTRAINT) {
            return Foundation::Result<Application::AuthUser>::FromError(
                Foundation::Error(Foundation::ErrorCode::Conflict, "邮箱已注册"));
        }
        return Foundation::Result<Application::AuthUser>::FromError(
            StorageError(Database.NativeHandle()));
    }
    const auto CommitResult = Database.Commit();
    if (!CommitResult.HasValue()) {
        Database.Rollback();
        return Foundation::Result<Application::AuthUser>::FromError(CommitResult.ErrorValue());
    }
    return User;
}

}  // namespace Menu::Infrastructure
