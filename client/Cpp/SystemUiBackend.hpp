#pragma once

#include <memory>

namespace Menu::Client {

class SystemUiBackend {
public:
    virtual ~SystemUiBackend() = default;
    virtual void ApplyTheme(bool Dark) = 0;
};

[[nodiscard]] std::unique_ptr<SystemUiBackend> CreateSystemUiBackend();

}  // namespace Menu::Client
