#include "Infrastructure/SqliteAuthService.hpp"

#include "Infrastructure/PasswordHasher.hpp"
#include "Infrastructure/TokenService.hpp"

#include <algorithm>
#include <cctype>
#include <optional>
#include <ranges>
#include <string>
#include <string_view>
#include <utility>

namespace Menu::Infrastructure {
namespace {

Foundation::Error InvalidInput(std::string Message) {
    return Foundation::Error(Foundation::ErrorCode::InvalidArgument, std::move(Message));
}

Foundation::Error AuthenticationFailure() {
    return Foundation::Error(
        Foundation::ErrorCode::AuthenticationFailed, "邮箱或密码错误");
}

std::string TrimLowerAscii(std::string_view Value) {
    std::size_t First = 0;
    while (First < Value.size() &&
           std::isspace(static_cast<unsigned char>(Value[First])) != 0) {
        ++First;
    }
    std::size_t Last = Value.size();
    while (Last > First &&
           std::isspace(static_cast<unsigned char>(Value[Last - 1])) != 0) {
        --Last;
    }
    std::string Result(Value.substr(First, Last - First));
    for (char& Character : Result) {
        Character = static_cast<char>(
            std::tolower(static_cast<unsigned char>(Character)));
    }
    return Result;
}

Foundation::Result<std::string> NormalizeEmail(std::string_view Email) {
    const std::string Normalized = TrimLowerAscii(Email);
    const std::size_t At = Normalized.find('@');
    if (Normalized.size() < 5U || Normalized.size() > 320U || At == std::string::npos ||
        At == 0U || At + 1U >= Normalized.size() ||
        Normalized.find('@', At + 1U) != std::string::npos ||
        std::ranges::any_of(Normalized, [](char Character) {
            return std::isspace(static_cast<unsigned char>(Character)) != 0 ||
                   Character == '\r' || Character == '\n';
        })) {
        return Foundation::Result<std::string>::FromError(
            InvalidInput("邮箱格式无效"));
    }
    return Normalized;
}

Foundation::Result<void> ValidatePassword(std::string_view Password) {
    if (Password.size() < 12U || Password.size() > 128U) {
        return Foundation::Result<void>::FromError(
            InvalidInput("密码长度必须在 12 到 128 个字符之间"));
    }
    return Foundation::Result<void>();
}

bool IsValidDisplayName(std::string_view DisplayName) {
    if (DisplayName.empty() || DisplayName.size() > 80U) {
        return false;
    }
    bool HasVisibleCharacter = false;
    for (unsigned char Character : DisplayName) {
        if (Character < 0x20U || Character == 0x7FU) {
            return false;
        }
        if (std::isspace(Character) == 0) {
            HasVisibleCharacter = true;
        }
    }
    return HasVisibleCharacter;
}

Foundation::Result<Application::AuthResponse> CreateTokenResponse(
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
    return CreateTokenResponse(*Repository, CreatedUser.Value());
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
        return Foundation::Result<Application::AuthResponse>::FromError(AuthenticationFailure());
    }
    const auto UserResult = Repository->FindUserByEmail(NormalizedEmail.Value());
    if (!UserResult.HasValue()) {
        return Foundation::Result<Application::AuthResponse>::FromError(
            UserResult.ErrorValue());
    }
    if (!UserResult.Value().has_value()) {
        return Foundation::Result<Application::AuthResponse>::FromError(AuthenticationFailure());
    }
    const Application::AuthUser& User = UserResult.Value().value();
    const auto VerifyResult = PasswordHasher::Verify(Password, User.PasswordHash);
    if (!VerifyResult.HasValue() || !VerifyResult.Value()) {
        return Foundation::Result<Application::AuthResponse>::FromError(AuthenticationFailure());
    }
    return CreateTokenResponse(*Repository, User);
}

Foundation::Result<Application::AuthResponse> SqliteAuthService::Refresh(
    std::string_view RefreshToken) {
    if (!Repository || RefreshToken.empty() || RefreshToken.size() > 256U) {
        return Foundation::Result<Application::AuthResponse>::FromError(AuthenticationFailure());
    }
    const auto TokenHash = TokenService::HashToken(RefreshToken);
    if (!TokenHash.HasValue()) {
        return Foundation::Result<Application::AuthResponse>::FromError(AuthenticationFailure());
    }
    const auto UserResult = Repository->FindUserByRefreshTokenHash(
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
    const auto RotateResult = Repository->RotateRefreshToken(
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

Foundation::Result<Application::AuthUser> SqliteAuthService::Authenticate(
    std::string_view AccessToken) {
    if (!Repository || AccessToken.empty() || AccessToken.size() > 256U) {
        return Foundation::Result<Application::AuthUser>::FromError(AuthenticationFailure());
    }
    const auto TokenHash = TokenService::HashToken(AccessToken);
    if (!TokenHash.HasValue()) {
        return Foundation::Result<Application::AuthUser>::FromError(AuthenticationFailure());
    }
    const auto UserResult = Repository->FindUserByAccessTokenHash(
        TokenHash.Value(), TokenService::CurrentUnixSeconds());
    if (!UserResult.HasValue()) {
        return Foundation::Result<Application::AuthUser>::FromError(UserResult.ErrorValue());
    }
    if (!UserResult.Value().has_value()) {
        return Foundation::Result<Application::AuthUser>::FromError(AuthenticationFailure());
    }
    return UserResult.Value().value();
}

}  // namespace Menu::Infrastructure
