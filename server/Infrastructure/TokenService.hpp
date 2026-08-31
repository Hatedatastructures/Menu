#pragma once

#include <Foundation/Result.hpp>

#include <cstddef>
#include <cstdint>
#include <string>
#include <string_view>

namespace Menu::Infrastructure {

class TokenService final {
public:
    static Foundation::Result<std::string> CreateToken(std::size_t ByteCount = 32U);

    static Foundation::Result<std::string> HashToken(std::string_view Token);

    [[nodiscard]] static std::int64_t CurrentUnixSeconds() noexcept;
};

}  // namespace Menu::Infrastructure
