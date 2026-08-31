#pragma once

#include <Foundation/Result.hpp>

#include <string>
#include <string_view>

namespace Menu::Infrastructure {

class PasswordHasher final {
public:
    static Foundation::Result<std::string> Hash(std::string_view Password);

    static Foundation::Result<bool> Verify(
        std::string_view Password,
        std::string_view EncodedHash);
};

}  // namespace Menu::Infrastructure
