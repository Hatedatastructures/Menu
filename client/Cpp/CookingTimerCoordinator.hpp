#pragma once

#include "NotificationController.hpp"

#include <QObject>
#include <QString>

namespace Menu::Client {

class CookingTimerCoordinator final : public QObject {
    Q_OBJECT

public:
    explicit CookingTimerCoordinator(
        NotificationController& Notifications,
        QObject* Parent = nullptr);

    [[nodiscard]] static QString TimerId(
        const QString& RecipeId,
        int StepOrder);

    Q_INVOKABLE void StartTimer(
        const QString& RecipeId,
        int StepOrder,
        int DelaySeconds,
        const QString& Title,
        const QString& Body);
    Q_INVOKABLE void StopTimer(
        const QString& RecipeId,
        int StepOrder);
    Q_INVOKABLE void StopAllTimers();

private:
    NotificationController& NotificationsValue;
};

}  // namespace Menu::Client
