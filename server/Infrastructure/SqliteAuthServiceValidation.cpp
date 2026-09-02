#include "Infrastructure/SqliteAuthServiceValidation.hpp"

#include <algorithm>
#include <cctype>
#include <ranges>
#include <utility>

namespace Menu::Infrastructure::SqliteAuthServiceValidation {
namespace {

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

}  // namespace

Foundation::Error InvalidInput(std::string Message) {
    return Foundation::Error(Foundation::ErrorCode::InvalidArgument, std::move(Message));
}

Foundation::Error AuthenticationFailure() {
    return Foundation::Error(
        Foundation::ErrorCode::AuthenticationFailed, "邮箱或密码错误");
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
    if (DisplayName.empty() || DisplayName.size() > 64U) {
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

}  // namespace Menu::Infrastructure::SqliteAuthServiceValidation
