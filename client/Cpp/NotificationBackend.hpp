#pragma once

#include <QString>

#include <memory>

namespace Menu::Client {

struct NotificationBackendState final {
    bool Supported = false;
    bool Granted = false;
    QString Status;
};

class NotificationBackend {
public:
    virtual ~NotificationBackend() = default;

    [[nodiscard]] virtual NotificationBackendState Refresh() = 0;
    virtual void RequestPermission() = 0;
    virtual void ScheduleTimer(
        const QString& Id,
        int DelaySeconds,
        const QString& Title,
        const QString& Body) = 0;
    virtual void CancelTimer(const QString& Id) = 0;
    virtual void CancelAllTimers() = 0;
    virtual void OpenSettings() = 0;
    virtual void ShowTestNotification() = 0;
};

[[nodiscard]] std::unique_ptr<NotificationBackend> CreateNotificationBackend();

}  // namespace Menu::Client
