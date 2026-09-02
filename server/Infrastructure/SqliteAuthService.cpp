#include "Infrastructure/SqliteAuthService.hpp"

#include "Infrastructure/PasswordHasher.hpp"
#include "Infrastructure/SqliteAuthServiceTokenFlow.hpp"
#include "Infrastructure/SqliteAuthServiceValidation.hpp"
#include "Infrastructure/TokenService.hpp"

#include <utility>

namespace Menu::Infrastructure {
namespace {

using SqliteAuthServiceTokenFlow::IssueTokens;
using SqliteAuthServiceValidation::AuthenticationFailure;
using SqliteAuthServiceValidation::InvalidInput;
using SqliteAuthServiceValidation::IsValidDisplayName;
using SqliteAuthServiceValidation::NormalizeEmail;
using SqliteAuthServiceValidation::ValidatePassword;

}  // namespace

SqliteAuthService::SqliteAuthService(
    std::unique_ptr<Application::AuthRepository> RepositoryValue)
    : Repository(std::move(RepositoryValue)) {}

Foundation::Result<Application::AuthResponse> SqliteAuthService::Register(
    std::string_view Email,
    std::string_view Password,
    std::string_view DisplayName) {
    if (!Repository || !IsValidDisplayName(DisplayName)) {
        return Foundation::Result<Application::AuthResponse>::FromError(
            InvalidInput("注册信息无效"));
    }
    const auto NormalizedEmail = NormalizeEmail(Email);
    const auto PasswordResult = ValidatePassword(Password);
    if (!NormalizedEmail.HasValue() || !PasswordResult.HasValue()) {
        return Foundation::Result<Application::AuthResponse>::FromError(
            InvalidInput("注册信息无效"));
    }
    const auto PasswordHash = PasswordHasher::Hash(Password);
    const auto UserId = TokenService::CreateToken(16U);
    if (!PasswordHash.HasValue() || !UserId.HasValue()) {
        return Foundation::Result<Application::AuthResponse>::FromError(
            Foundation::Error(Foundation::ErrorCode::StorageUnavailable, "注册安全组件不可用"));
    }
    Application::AuthUser User{
        "user." + UserId.Value(),
        NormalizedEmail.Value(),
        std::string(DisplayName),
        PasswordHash.Value(),
        false};
    const auto CreatedUser = Repository->CreateUser(std::move(User));
    if (!CreatedUser.HasValue()) {
        return Foundation::Result<Application::AuthResponse>::FromError(
            CreatedUser.ErrorValue());
    }
    return IssueTokens(*Repository, CreatedUser.Value());
}

Foundation::Result<Application::AuthResponse> SqliteAuthService::Login(
    std::string_view Email,
    std::string_view Password) {
    if (!Repository) {
        return Foundation::Result<Application::AuthResponse>::FromError(
            Foundation::Error(Foundation::ErrorCode::StorageUnavailable, "认证仓储不可用"));
    }
    if (Password.empty() || Password.size() > 128U) {
        return Foundation::Result<Application::AuthResponse>::FromError(
            AuthenticationFailure());
    }
    const auto NormalizedEmail = NormalizeEmail(Email);
    if (!NormalizedEmail.HasValue()) {
        return Foundation::Result<Application::AuthResponse>::FromError(
            AuthenticationFailure());
    }
    const auto UserResult = Repository->FindUserByEmail(NormalizedEmail.Value());
    if (!UserResult.HasValue()) {
        return Foundation::Result<Application::AuthResponse>::FromError(
            UserResult.ErrorValue());
    }
    if (!UserResult.Value().has_value()) {
        return Foundation::Result<Application::AuthResponse>::FromError(
            AuthenticationFailure());
    }
    const Application::AuthUser& User = UserResult.Value().value();
    const auto VerifyResult = PasswordHasher::Verify(Password, User.PasswordHash);
    if (!VerifyResult.HasValue() || !VerifyResult.Value()) {
        return Foundation::Result<Application::AuthResponse>::FromError(
            AuthenticationFailure());
    }
    return IssueTokens(*Repository, User);
}

Foundation::Result<Application::AuthResponse> SqliteAuthService::Refresh(
    std::string_view RefreshToken) {
    if (!Repository) {
        return Foundation::Result<Application::AuthResponse>::FromError(
            AuthenticationFailure());
    }
    return SqliteAuthServiceTokenFlow::Refresh(*Repository, RefreshToken);
}

Foundation::Result<Application::AuthUser> SqliteAuthService::Authenticate(
    std::string_view AccessToken) {
    if (!Repository || AccessToken.empty() || AccessToken.size() > 256U) {
        return Foundation::Result<Application::AuthUser>::FromError(
            AuthenticationFailure());
    }
    const auto TokenHash = TokenService::HashToken(AccessToken);
    if (!TokenHash.HasValue()) {
        return Foundation::Result<Application::AuthUser>::FromError(
            AuthenticationFailure());
    }
    const auto UserResult = Repository->FindUserByAccessTokenHash(
        TokenHash.Value(), TokenService::CurrentUnixSeconds());
    if (!UserResult.HasValue()) {
        return Foundation::Result<Application::AuthUser>::FromError(
            UserResult.ErrorValue());
    }
    if (!UserResult.Value().has_value()) {
        return Foundation::Result<Application::AuthUser>::FromError(
            AuthenticationFailure());
    }
    return UserResult.Value().value();
}

}  // namespace Menu::Infrastructure
