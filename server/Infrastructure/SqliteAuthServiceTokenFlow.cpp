#include "Infrastructure/SqliteAuthServiceTokenFlow.hpp"

#include "Infrastructure/PasswordHasher.hpp"
#include "Infrastructure/SqliteAuthServiceValidation.hpp"
#include "Infrastructure/TokenService.hpp"

#include <cstdint>

namespace Menu::Infrastructure::SqliteAuthServiceTokenFlow {
namespace {

using SqliteAuthServiceValidation::AuthenticationFailure;

}  // namespace

Foundation::Result<Application::AuthResponse> IssueTokens(
    Application::AuthRepository& Repository,
    const Application::AuthUser& User) {
    const auto AccessToken = TokenService::CreateToken();
    const auto RefreshToken = TokenService::CreateToken();
    if (!AccessToken.HasValue() || !RefreshToken.HasValue()) {
        return Foundation::Result<Application::AuthResponse>::FromError(
            Foundation::Error(Foundation::ErrorCode::StorageUnavailable, "令牌生成失败"));
    }
    const auto AccessHash = TokenService::HashToken(AccessToken.Value());
    const auto RefreshHash = TokenService::HashToken(RefreshToken.Value());
    if (!AccessHash.HasValue() || !RefreshHash.HasValue()) {
        return Foundation::Result<Application::AuthResponse>::FromError(
            Foundation::Error(Foundation::ErrorCode::StorageUnavailable, "令牌摘要失败"));
    }
    constexpr std::int64_t AccessLifetime = 15 * 60;
    constexpr std::int64_t RefreshLifetime = 30 * 24 * 60 * 60;
    const std::int64_t Now = TokenService::CurrentUnixSeconds();
    const auto StoreResult = Repository.StoreTokens(
        User.Id,
        AccessHash.Value(),
        Now + AccessLifetime,
        RefreshHash.Value(),
        Now + RefreshLifetime);
    if (!StoreResult.HasValue()) {
        return Foundation::Result<Application::AuthResponse>::FromError(
            StoreResult.ErrorValue());
    }
    return Application::AuthResponse{
        User, AccessToken.Value(), RefreshToken.Value(), static_cast<int>(AccessLifetime)};
}

Foundation::Result<Application::AuthResponse> Refresh(
    Application::AuthRepository& Repository,
    std::string_view RefreshToken) {
    if (RefreshToken.empty() || RefreshToken.size() > 256U) {
        return Foundation::Result<Application::AuthResponse>::FromError(AuthenticationFailure());
    }
    const auto TokenHash = TokenService::HashToken(RefreshToken);
    if (!TokenHash.HasValue()) {
        return Foundation::Result<Application::AuthResponse>::FromError(AuthenticationFailure());
    }
    const auto UserResult = Repository.FindUserByRefreshTokenHash(
        TokenHash.Value(), TokenService::CurrentUnixSeconds());
    if (!UserResult.HasValue()) {
        return Foundation::Result<Application::AuthResponse>::FromError(
            UserResult.ErrorValue());
    }
    if (!UserResult.Value().has_value()) {
        return Foundation::Result<Application::AuthResponse>::FromError(AuthenticationFailure());
    }
    const Application::AuthUser& User = UserResult.Value().value();
    const auto AccessToken = TokenService::CreateToken();
    const auto NewRefreshToken = TokenService::CreateToken();
    if (!AccessToken.HasValue() || !NewRefreshToken.HasValue()) {
        return Foundation::Result<Application::AuthResponse>::FromError(
            Foundation::Error(Foundation::ErrorCode::StorageUnavailable, "令牌生成失败"));
    }
    const auto AccessHash = TokenService::HashToken(AccessToken.Value());
    const auto NewRefreshHash = TokenService::HashToken(NewRefreshToken.Value());
    if (!AccessHash.HasValue() || !NewRefreshHash.HasValue()) {
        return Foundation::Result<Application::AuthResponse>::FromError(
            Foundation::Error(Foundation::ErrorCode::StorageUnavailable, "令牌摘要失败"));
    }
    constexpr std::int64_t AccessLifetime = 15 * 60;
    constexpr std::int64_t RefreshLifetime = 30 * 24 * 60 * 60;
    const std::int64_t Now = TokenService::CurrentUnixSeconds();
    const auto RotateResult = Repository.RotateRefreshToken(
        TokenHash.Value(),
        User.Id,
        AccessHash.Value(),
        Now + AccessLifetime,
        NewRefreshHash.Value(),
        Now + RefreshLifetime,
        Now);
    if (!RotateResult.HasValue()) {
        return Foundation::Result<Application::AuthResponse>::FromError(
            RotateResult.ErrorValue());
    }
    return Application::AuthResponse{
        User, AccessToken.Value(), NewRefreshToken.Value(), static_cast<int>(AccessLifetime)};
}

}  // namespace Menu::Infrastructure::SqliteAuthServiceTokenFlow
