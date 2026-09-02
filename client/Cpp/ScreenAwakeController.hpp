#pragma once

#include "ScreenAwakeBackend.hpp"

#include <QObject>

#include <functional>
#include <memory>

namespace Menu::Client {

class ScreenAwakeController final : public QObject {
    Q_OBJECT
    Q_PROPERTY(bool Enabled READ Enabled NOTIFY EnabledChanged)

public:
    using BackendFactory = std::function<std::unique_ptr<ScreenAwakeBackend>()>;

    explicit ScreenAwakeController(
        QObject* Parent = nullptr,
        BackendFactory Factory = {});
    ~ScreenAwakeController() override;

    [[nodiscard]] bool Enabled() const noexcept;

    Q_INVOKABLE void SetEnabled(bool EnabledValue);

signals:
    void EnabledChanged();

private:
    bool EnabledValue = false;
    std::unique_ptr<ScreenAwakeBackend> Backend;
};

}  // namespace Menu::Client
