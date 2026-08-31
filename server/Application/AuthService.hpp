#pragma once

#include <Application/AuthRepository.hpp>

#include <memory>
#include <string>
#include <string_view>

namespace Menu::Application {

struct AuthResponse {
    AuthUser User;
    std::string AccessToken;
    std::string RefreshToken;
    int AccessTokenExpiresInSeconds = 0;
};

class AuthService {
public:
    virtual ~AuthService() = default;

    virtual Foundation::Result<AuthResponse> Register(
        std::string_view Email,
        std::string_view Password,
        std::string_view DisplayName) = 0;

    virtual Foundation::Result<AuthResponse> Login(
        std::string_view Email,
        std::string_view Password) = 0;

    virtual Foundation::Result<AuthResponse> Refresh(
        std::string_view RefreshToken) = 0;

    virtual Foundation::Result<AuthUser> Authenticate(
        std::string_view AccessToken) = 0;
};

}  // namespace Menu::Application
