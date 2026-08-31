#pragma once

#include <Application/AuthService.hpp>

#include <memory>

namespace Menu::Infrastructure {

class SqliteAuthService final : public Application::AuthService {
public:
    explicit SqliteAuthService(std::unique_ptr<Application::AuthRepository> RepositoryValue);

    Foundation::Result<Application::AuthResponse> Register(
        std::string_view Email,
        std::string_view Password,
        std::string_view DisplayName) override;

    Foundation::Result<Application::AuthResponse> Login(
        std::string_view Email,
        std::string_view Password) override;

    Foundation::Result<Application::AuthResponse> Refresh(
        std::string_view RefreshToken) override;

    Foundation::Result<Application::AuthUser> Authenticate(
        std::string_view AccessToken) override;

private:
    Foundation::Result<Application::AuthResponse> IssueTokens(
        Application::AuthUser User);

    std::unique_ptr<Application::AuthRepository> Repository;
};

}  // namespace Menu::Infrastructure
