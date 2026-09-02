#include "CookingTimerCoordinator.hpp"

#include <QCryptographicHash>

namespace Menu::Client {

CookingTimerCoordinator::CookingTimerCoordinator(
    NotificationController& Notifications,
    QObject* Parent)
    : QObject(Parent), NotificationsValue(Notifications) {}

QString CookingTimerCoordinator::TimerId(
    const QString& RecipeId,
    int StepOrder) {
    const QByteArray Key = RecipeId.toUtf8() + ':' + QByteArray::number(StepOrder);
    return QStringLiteral("cook-") + QString::fromLatin1(
        QCryptographicHash::hash(Key, QCryptographicHash::Sha256).toHex().left(24));
}

void CookingTimerCoordinator::StartTimer(
    const QString& RecipeId,
    int StepOrder,
    int DelaySeconds,
    const QString& Title,
    const QString& Body) {
    if (RecipeId.isEmpty() || StepOrder < 1 || DelaySeconds <= 0) {
        return;
    }
    const QString Id = TimerId(RecipeId, StepOrder);
    NotificationsValue.ScheduleTimer(Id, DelaySeconds, Title, Body);
}

void CookingTimerCoordinator::StopTimer(
    const QString& RecipeId,
    int StepOrder) {
    if (RecipeId.isEmpty() || StepOrder < 1) {
        return;
    }
    const QString Id = TimerId(RecipeId, StepOrder);
    NotificationsValue.CancelTimer(Id);
}

void CookingTimerCoordinator::StopAllTimers() {
    NotificationsValue.CancelAllTimers();
}

}  // namespace Menu::Client
