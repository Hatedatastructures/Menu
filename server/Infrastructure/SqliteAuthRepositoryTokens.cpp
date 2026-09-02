#include "Infrastructure/SqliteAuthRepository.hpp"

#include "Infrastructure/SqliteAuthRepositorySupport.hpp"
#include "Infrastructure/SqliteWorkflowSupport.hpp"

#include <string_view>
#include <utility>

namespace Menu::Infrastructure {
namespace {

using SqliteAuthRepositorySupport::FindUser;
using SqliteAuthRepositorySupport::InsertToken;
using SqliteAuthRepositorySupport::RollbackWith;
using SqliteWorkflowSupport::BindText;
using SqliteWorkflowSupport::Prepare;
using SqliteWorkflowSupport::StatementGuard;

}  // namespace

Foundation::Result<void> SqliteAuthRepository::StoreTokens(
    std::string_view UserId,
    std::string_view AccessTokenHash,
    std::int64_t AccessExpiresAt,
    std::string_view RefreshTokenHash,
    std::int64_t RefreshExpiresAt) {
    const auto BeginResult = Database.BeginTransaction();
    if (!BeginResult.HasValue()) {
        return BeginResult;
    }
    if (auto Result = InsertToken(
            Database, "AccessTokens", UserId, AccessTokenHash, AccessExpiresAt);
        !Result.HasValue()) {
        return RollbackWith(Result, Database);
    }
    if (auto Result = InsertToken(
            Database, "RefreshTokens", UserId, RefreshTokenHash, RefreshExpiresAt);
        !Result.HasValue()) {
        return RollbackWith(Result, Database);
    }
    const auto CommitResult = Database.Commit();
    if (!CommitResult.HasValue()) {
        Database.Rollback();
        return CommitResult;
    }
    return Foundation::Result<void>();
}

Foundation::Result<std::optional<Application::AuthUser>>
SqliteAuthRepository::FindUserByAccessTokenHash(
    std::string_view TokenHash,
    std::int64_t Now) {
    return FindUser(
        Database,
        "SELECT Users.Id, Users.Email, Users.PasswordHash, Users.DisplayName, Users.IsAdmin "
        "FROM AccessTokens JOIN Users ON Users.Id = AccessTokens.UserId "
        "WHERE AccessTokens.TokenHash = ? AND CAST(AccessTokens.ExpiresAt AS INTEGER) > ? "
        "AND AccessTokens.RevokedAt IS NULL LIMIT 1;",
        TokenHash,
        Now,
        true);
}

Foundation::Result<std::optional<Application::AuthUser>>
SqliteAuthRepository::FindUserByRefreshTokenHash(
    std::string_view TokenHash,
    std::int64_t Now) {
    return FindUser(
        Database,
        "SELECT Users.Id, Users.Email, Users.PasswordHash, Users.DisplayName, Users.IsAdmin "
        "FROM RefreshTokens JOIN Users ON Users.Id = RefreshTokens.UserId "
        "WHERE RefreshTokens.TokenHash = ? AND CAST(RefreshTokens.ExpiresAt AS INTEGER) > ? "
        "AND RefreshTokens.RevokedAt IS NULL LIMIT 1;",
        TokenHash,
        Now,
        true);
}

Foundation::Result<void> SqliteAuthRepository::RotateRefreshToken(
    std::string_view OldRefreshTokenHash,
    std::string_view UserId,
    std::string_view AccessTokenHash,
    std::int64_t AccessExpiresAt,
    std::string_view NewRefreshTokenHash,
    std::int64_t NewRefreshExpiresAt,
    std::int64_t RevokedAt) {
    const auto BeginResult = Database.BeginTransaction();
    if (!BeginResult.HasValue()) {
        return BeginResult;
    }
    auto StatementResult = Prepare(
        Database,
        "UPDATE RefreshTokens SET RevokedAt = ? WHERE TokenHash = ? AND UserId = ? "
        "AND RevokedAt IS NULL;");
    if (!StatementResult.HasValue()) {
        Database.Rollback();
        return Foundation::Result<void>::FromError(StatementResult.ErrorValue());
    }
    StatementGuard Statement = std::move(StatementResult).Value();
    if (auto Result = BindText(Statement.Get(), 1, std::to_string(RevokedAt));
        !Result.HasValue()) {
        Database.Rollback();
        return Result;
    }
    if (auto Result = BindText(Statement.Get(), 2, OldRefreshTokenHash); !Result.HasValue()) {
        Database.Rollback();
        return Result;
    }
    if (auto Result = BindText(Statement.Get(), 3, UserId); !Result.HasValue()) {
        Database.Rollback();
        return Result;
    }
    if (sqlite3_step(Statement.Get()) != SQLITE_DONE ||
        sqlite3_changes(Database.NativeHandle()) != 1) {
        Database.Rollback();
        return Foundation::Result<void>::FromError(
            Foundation::Error(Foundation::ErrorCode::AuthenticationFailed, "刷新令牌无效"));
    }
    if (auto Result = InsertToken(
            Database, "AccessTokens", UserId, AccessTokenHash, AccessExpiresAt);
        !Result.HasValue()) {
        return RollbackWith(Result, Database);
    }
    if (auto Result = InsertToken(
            Database, "RefreshTokens", UserId, NewRefreshTokenHash, NewRefreshExpiresAt);
        !Result.HasValue()) {
        return RollbackWith(Result, Database);
    }
    const auto CommitResult = Database.Commit();
    if (!CommitResult.HasValue()) {
        Database.Rollback();
        return CommitResult;
    }
    return Foundation::Result<void>();
}

}  // namespace Menu::Infrastructure
