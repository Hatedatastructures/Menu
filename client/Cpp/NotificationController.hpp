#pragma once

#include "NotificationBackend.hpp"

#include <QObject>

#include <memory>
#include <functional>

namespace Menu::Client {

class NotificationController final : public QObject {
    Q_OBJECT
    Q_PROPERTY(bool Supported READ Supported CONSTANT)
    Q_PROPERTY(bool PermissionGranted READ PermissionGranted NOTIFY StateChanged)
    Q_PROPERTY(bool Enabled READ Enabled WRITE SetEnabled NOTIFY StateChanged)
    Q_PROPERTY(QString Status READ Status NOTIFY StateChanged)

public:
    using BackendFactory = std::function<std::unique_ptr<NotificationBackend>()>;

    explicit NotificationController(
        QObject* Parent = nullptr,
        BackendFactory Factory = {});
    ~NotificationController() override;

    [[nodiscard]] bool Supported() const noexcept;
    [[nodiscard]] bool PermissionGranted() const noexcept;
    [[nodiscard]] bool Enabled() const noexcept;
    [[nodiscard]] QString Status() const;

    void SetEnabled(bool EnabledValue);

    Q_INVOKABLE void RefreshPermission();
    Q_INVOKABLE void RequestPermission();
    Q_INVOKABLE void ScheduleTimer(
        const QString& Id,
        int DelaySeconds,
        const QString& Title,
        const QString& Body);
    Q_INVOKABLE void CancelTimer(const QString& Id);
    Q_INVOKABLE void CancelAllTimers();
    Q_INVOKABLE void OpenSettings();
    Q_INVOKABLE void ShowTestNotification();

signals:
    void StateChanged();

private:
    void ApplyState(const NotificationBackendState& State);
    void SetStatus(QString StatusValue);

    bool SupportedValue = false;
    bool PermissionGrantedValue = false;
    bool EnabledValue = true;
    QString StatusValue;
    std::unique_ptr<NotificationBackend> Backend;
};

}  // namespace Menu::Client
