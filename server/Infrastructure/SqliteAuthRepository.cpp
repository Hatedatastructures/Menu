#include "Infrastructure/SqliteAuthRepository.hpp"

#include "Infrastructure/TokenService.hpp"

#include <string>
#include <string_view>
#include <utility>

namespace Menu::Infrastructure {
namespace {

class StatementGuard final {
public:
    explicit StatementGuard(sqlite3_stmt* StatementValue)
        : Statement(StatementValue) {}

    StatementGuard(const StatementGuard&) = delete;
    StatementGuard& operator=(const StatementGuard&) = delete;

    StatementGuard(StatementGuard&& Other) noexcept
        : Statement(std::exchange(Other.Statement, nullptr)) {}

    StatementGuard& operator=(StatementGuard&& Other) noexcept {
        if (this != &Other) {
            Reset();
            Statement = std::exchange(Other.Statement, nullptr);
        }
        return *this;
    }

    ~StatementGuard() {
        Reset();
    }

    [[nodiscard]] sqlite3_stmt* Get() const noexcept {
        return Statement;
    }

private:
    void Reset() noexcept {
        if (Statement != nullptr) {
            sqlite3_finalize(Statement);
            Statement = nullptr;
        }
    }

    sqlite3_stmt* Statement = nullptr;
};

Foundation::Error StorageError(sqlite3* Handle) {
    const char* Message = Handle == nullptr ? nullptr : sqlite3_errmsg(Handle);
    return Foundation::Error(
        Foundation::ErrorCode::StorageUnavailable,
        Message == nullptr ? "SQLite 查询失败" : std::string(Message));
}

Foundation::Result<StatementGuard> Prepare(SqliteDatabase& Database, std::string_view Sql) {
    sqlite3_stmt* RawStatement = nullptr;
    const std::string SqlText(Sql);
    if (sqlite3_prepare_v2(
            Database.NativeHandle(), SqlText.c_str(), static_cast<int>(SqlText.size()),
            &RawStatement, nullptr) != SQLITE_OK) {
        return Foundation::Result<StatementGuard>::FromError(
            StorageError(Database.NativeHandle()));
    }
    return StatementGuard(RawStatement);
}

Foundation::Result<void> BindText(sqlite3_stmt* Statement, int Index, std::string_view Value) {
    const std::string Text(Value);
    if (sqlite3_bind_text(Statement, Index, Text.c_str(), -1, SQLITE_TRANSIENT) != SQLITE_OK) {
        return Foundation::Result<void>::FromError(
            Foundation::Error(Foundation::ErrorCode::StorageUnavailable, "SQLite 绑定失败"));
    }
    return Foundation::Result<void>();
}

Foundation::Result<void> BindInteger(sqlite3_stmt* Statement, int Index, std::int64_t Value) {
    if (sqlite3_bind_int64(Statement, Index, Value) != SQLITE_OK) {
        return Foundation::Result<void>::FromError(
            Foundation::Error(Foundation::ErrorCode::StorageUnavailable, "SQLite 绑定失败"));
    }
    return Foundation::Result<void>();
}

std::string ColumnText(sqlite3_stmt* Statement, int Column) {
    const unsigned char* TextValue = sqlite3_column_text(Statement, Column);
    return TextValue == nullptr ? std::string() : reinterpret_cast<const char*>(TextValue);
}

Application::AuthUser ReadUser(sqlite3_stmt* Statement) {
    Application::AuthUser User;
    User.Id = ColumnText(Statement, 0);
    User.Email = ColumnText(Statement, 1);
    User.PasswordHash = ColumnText(Statement, 2);
    User.DisplayName = ColumnText(Statement, 3);
    User.IsAdmin = sqlite3_column_int(Statement, 4) != 0;
    return User;
}

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
    const auto TokenBind = BindText(Statement.Get(), 1, TokenHash);
    if (!TokenBind.HasValue()) {
        return Foundation::Result<std::optional<Application::AuthUser>>::FromError(
            TokenBind.ErrorValue());
    }
    if (HasExpiry) {
        const auto NowBind = BindInteger(Statement.Get(), 2, Now);
        if (!NowBind.HasValue()) {
            return Foundation::Result<std::optional<Application::AuthUser>>::FromError(
                NowBind.ErrorValue());
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
    const auto IdBind = BindText(Statement.Get(), 1, Id.Value());
    const auto UserBind = BindText(Statement.Get(), 2, UserId);
    const auto HashBind = BindText(Statement.Get(), 3, TokenHash);
    const auto ExpiryBind = BindText(Statement.Get(), 4, std::to_string(ExpiresAt));
    if (!IdBind.HasValue()) {
        return IdBind;
    }
    if (!UserBind.HasValue()) {
        return UserBind;
    }
    if (!HashBind.HasValue()) {
        return HashBind;
    }
    if (!ExpiryBind.HasValue()) {
        return ExpiryBind;
    }
    if (sqlite3_step(Statement.Get()) != SQLITE_DONE) {
        return Foundation::Result<void>::FromError(StorageError(Database.NativeHandle()));
    }
    return Foundation::Result<void>();
}

Foundation::Result<void> RollbackWith(Foundation::Result<void> Result, SqliteDatabase& Database) {
    if (!Result.HasValue()) {
        Database.Rollback();
    }
    return Result;
}

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
    const auto IdBind = BindText(Statement.Get(), 1, User.Id);
    const auto EmailBind = BindText(Statement.Get(), 2, User.Email);
    const auto PasswordBind = BindText(Statement.Get(), 3, User.PasswordHash);
    const auto DisplayBind = BindText(Statement.Get(), 4, User.DisplayName);
    const auto AdminBind = BindInteger(Statement.Get(), 5, User.IsAdmin ? 1 : 0);
    if (!IdBind.HasValue() || !EmailBind.HasValue() || !PasswordBind.HasValue() ||
        !DisplayBind.HasValue() || !AdminBind.HasValue()) {
        Database.Rollback();
        const Foundation::Error ErrorValue = !IdBind.HasValue()
                                                 ? IdBind.ErrorValue()
                                                 : !EmailBind.HasValue()
                                                       ? EmailBind.ErrorValue()
                                                       : !PasswordBind.HasValue()
                                                             ? PasswordBind.ErrorValue()
                                                             : !DisplayBind.HasValue()
                                                                   ? DisplayBind.ErrorValue()
                                                                   : AdminBind.ErrorValue();
        return Foundation::Result<Application::AuthUser>::FromError(ErrorValue);
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
    const auto AccessResult = InsertToken(
        Database, "AccessTokens", UserId, AccessTokenHash, AccessExpiresAt);
    if (!AccessResult.HasValue()) {
        return RollbackWith(AccessResult, Database);
    }
    const auto RefreshResult = InsertToken(
        Database, "RefreshTokens", UserId, RefreshTokenHash, RefreshExpiresAt);
    if (!RefreshResult.HasValue()) {
        return RollbackWith(RefreshResult, Database);
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
    const auto RevokedBind = BindText(Statement.Get(), 1, std::to_string(RevokedAt));
    const auto OldHashBind = BindText(Statement.Get(), 2, OldRefreshTokenHash);
    const auto UserBind = BindText(Statement.Get(), 3, UserId);
    if (!RevokedBind.HasValue() || !OldHashBind.HasValue() || !UserBind.HasValue()) {
        Database.Rollback();
        const Foundation::Error ErrorValue = !RevokedBind.HasValue()
                                                 ? RevokedBind.ErrorValue()
                                                 : !OldHashBind.HasValue()
                                                       ? OldHashBind.ErrorValue()
                                                       : UserBind.ErrorValue();
        return Foundation::Result<void>::FromError(ErrorValue);
    }
    if (sqlite3_step(Statement.Get()) != SQLITE_DONE || sqlite3_changes(Database.NativeHandle()) != 1) {
        Database.Rollback();
        return Foundation::Result<void>::FromError(
            Foundation::Error(Foundation::ErrorCode::AuthenticationFailed, "刷新令牌无效"));
    }
    const auto AccessResult = InsertToken(
        Database, "AccessTokens", UserId, AccessTokenHash, AccessExpiresAt);
    if (!AccessResult.HasValue()) {
        return RollbackWith(AccessResult, Database);
    }
    const auto RefreshResult = InsertToken(
        Database, "RefreshTokens", UserId, NewRefreshTokenHash, NewRefreshExpiresAt);
    if (!RefreshResult.HasValue()) {
        return RollbackWith(RefreshResult, Database);
    }
    const auto CommitResult = Database.Commit();
    if (!CommitResult.HasValue()) {
        Database.Rollback();
        return CommitResult;
    }
    return Foundation::Result<void>();
}

}  // namespace Menu::Infrastructure
