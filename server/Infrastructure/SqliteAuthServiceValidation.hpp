#pragma once

#include <Foundation/Result.hpp>

#include <string>
#include <string_view>

namespace Menu::Infrastructure::SqliteAuthServiceValidation {

Foundation::Error InvalidInput(std::string Message);
Foundation::Error AuthenticationFailure();

Foundation::Result<std::string> NormalizeEmail(std::string_view Email);
Foundation::Result<void> ValidatePassword(std::string_view Password);
bool IsValidDisplayName(std::string_view DisplayName);

}  // namespace Menu::Infrastructure::SqliteAuthServiceValidation
