#pragma once

#include <Application/AuthRepository.hpp>
#include <Infrastructure/SqliteDatabase.hpp>

namespace Menu::Infrastructure {

class SqliteAuthRepository final : public Application::AuthRepository {
public:
    explicit SqliteAuthRepository(SqliteDatabase& DatabaseValue)
        : Database(DatabaseValue) {}

    Foundation::Result<std::size_t> CountUsers() override;

    Foundation::Result<std::optional<Application::AuthUser>> FindUserByEmail(
        std::string_view Email) override;

    Foundation::Result<Application::AuthUser> CreateUser(Application::AuthUser User) override;

    Foundation::Result<void> StoreTokens(
        std::string_view UserId,
        std::string_view AccessTokenHash,
        std::int64_t AccessExpiresAt,
        std::string_view RefreshTokenHash,
        std::int64_t RefreshExpiresAt) override;

    Foundation::Result<std::optional<Application::AuthUser>> FindUserByAccessTokenHash(
        std::string_view TokenHash,
        std::int64_t Now) override;

    Foundation::Result<std::optional<Application::AuthUser>> FindUserByRefreshTokenHash(
        std::string_view TokenHash,
        std::int64_t Now) override;

    Foundation::Result<void> RotateRefreshToken(
        std::string_view OldRefreshTokenHash,
        std::string_view UserId,
        std::string_view AccessTokenHash,
        std::int64_t AccessExpiresAt,
        std::string_view NewRefreshTokenHash,
        std::int64_t NewRefreshExpiresAt,
        std::int64_t RevokedAt) override;

private:
    SqliteDatabase& Database;
};

}  // namespace Menu::Infrastructure
