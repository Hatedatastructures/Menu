#pragma once

#include <string>
#include <utility>

namespace Menu::Foundation {

enum class ErrorCode {
    Unknown,
    InvalidArgument,
    NotFound,
    StorageUnavailable,
    NotReady,
};

class Error {
public:
    Error() = default;

    Error(ErrorCode CodeValue, std::string MessageValue)
        : Code(CodeValue), Message(std::move(MessageValue)) {}

    [[nodiscard]] ErrorCode CodeValue() const noexcept {
        return Code;
    }

    [[nodiscard]] const std::string& MessageValue() const noexcept {
        return Message;
    }

    friend bool operator==(const Error& Left, const Error& Right) noexcept {
        return Left.Code == Right.Code && Left.Message == Right.Message;
    }

private:
    ErrorCode Code = ErrorCode::Unknown;
    std::string Message;
};

}  // namespace Menu::Foundation
