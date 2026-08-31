#pragma once

#include <Foundation/Result.hpp>

#include <cstdint>
#include <optional>
#include <string>
#include <string_view>

namespace Menu::Application {

struct AuthUser {
    std::string Id;
    std::string Email;
    std::string DisplayName;
    std::string PasswordHash;
    bool IsAdmin = false;
};

class AuthRepository {
public:
    virtual ~AuthRepository() = default;

    virtual Foundation::Result<std::size_t> CountUsers() = 0;

    virtual Foundation::Result<std::optional<AuthUser>> FindUserByEmail(
        std::string_view Email) = 0;

    virtual Foundation::Result<AuthUser> CreateUser(AuthUser User) = 0;

    virtual Foundation::Result<void> StoreTokens(
        std::string_view UserId,
        std::string_view AccessTokenHash,
        std::int64_t AccessExpiresAt,
        std::string_view RefreshTokenHash,
        std::int64_t RefreshExpiresAt) = 0;

    virtual Foundation::Result<std::optional<AuthUser>> FindUserByAccessTokenHash(
        std::string_view TokenHash,
        std::int64_t Now) = 0;

    virtual Foundation::Result<std::optional<AuthUser>> FindUserByRefreshTokenHash(
        std::string_view TokenHash,
        std::int64_t Now) = 0;

    virtual Foundation::Result<void> RotateRefreshToken(
        std::string_view OldRefreshTokenHash,
        std::string_view UserId,
        std::string_view AccessTokenHash,
        std::int64_t AccessExpiresAt,
        std::string_view NewRefreshTokenHash,
        std::int64_t NewRefreshExpiresAt,
        std::int64_t RevokedAt) = 0;
};

}  // namespace Menu::Application
