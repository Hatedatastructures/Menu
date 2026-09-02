#pragma once

#include <memory>

namespace Menu::Client {

class ScreenAwakeBackend {
public:
    virtual ~ScreenAwakeBackend() = default;
    virtual void Apply(bool Enabled) = 0;
};

[[nodiscard]] std::unique_ptr<ScreenAwakeBackend> CreateScreenAwakeBackend();

}  // namespace Menu::Client
