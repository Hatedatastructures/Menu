#pragma once

#include <Application/AuthService.hpp>
#include <Application/AuthRepository.hpp>

#include <string_view>

namespace Menu::Infrastructure::SqliteAuthServiceTokenFlow {

Foundation::Result<Application::AuthResponse> IssueTokens(
    Application::AuthRepository& Repository,
    const Application::AuthUser& User);

Foundation::Result<Application::AuthResponse> Refresh(
    Application::AuthRepository& Repository,
    std::string_view RefreshToken);

}  // namespace Menu::Infrastructure::SqliteAuthServiceTokenFlow
