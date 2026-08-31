#pragma once

#include "Foundation/Error.hpp"

#include <cassert>
#include <optional>
#include <utility>
#include <variant>

namespace Menu::Foundation {

template <typename ValueType>
class Result {
public:
    Result(const ValueType& Value)
        : Storage(Value) {}

    Result(ValueType&& Value)
        : Storage(std::move(Value)) {}

    Result(const Error& Failure)
        : Storage(Failure) {}

    Result(Error&& Failure)
        : Storage(std::move(Failure)) {}

    [[nodiscard]] bool HasValue() const noexcept {
        return std::holds_alternative<ValueType>(Storage);
    }

    [[nodiscard]] const ValueType& Value() const& {
        assert(HasValue());
        return std::get<ValueType>(Storage);
    }

    [[nodiscard]] ValueType& Value() & {
        assert(HasValue());
        return std::get<ValueType>(Storage);
    }

    [[nodiscard]] ValueType&& Value() && {
        assert(HasValue());
        return std::get<ValueType>(std::move(Storage));
    }

    [[nodiscard]] const Error& ErrorValue() const& {
        assert(!HasValue());
        return std::get<Error>(Storage);
    }

    [[nodiscard]] Error&& ErrorValue() && {
        assert(!HasValue());
        return std::get<Error>(std::move(Storage));
    }

    static Result FromError(Error Failure) {
        return Result(std::move(Failure));
    }

private:
    std::variant<ValueType, Error> Storage;
};

template <>
class Result<void> {
public:
    Result() = default;

    Result(const Error& Failure)
        : FailureValue(Failure) {}

    Result(Error&& Failure)
        : FailureValue(std::move(Failure)) {}

    [[nodiscard]] bool HasValue() const noexcept {
        return !FailureValue.has_value();
    }

    [[nodiscard]] const Error& ErrorValue() const& {
        assert(!HasValue());
        return FailureValue.value();
    }

    static Result FromError(Error Failure) {
        return Result(std::move(Failure));
    }

private:
    std::optional<Error> FailureValue;
};

}  // namespace Menu::Foundation
