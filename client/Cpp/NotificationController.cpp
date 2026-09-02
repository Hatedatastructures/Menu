#include "NotificationController.hpp"

#include <utility>

#include <QGuiApplication>
#include <QTimer>

namespace Menu::Client {

NotificationController::NotificationController(QObject* Parent, BackendFactory Factory)
    : QObject(Parent), Backend(Factory ? Factory() : CreateNotificationBackend()) {
    RefreshPermission();
    if (auto* Application = qobject_cast<QGuiApplication*>(QCoreApplication::instance())) {
        connect(Application, &QGuiApplication::applicationStateChanged, this,
            [this](Qt::ApplicationState State) {
                if (State == Qt::ApplicationActive) {
                    RefreshPermission();
                }
            });
    }
}

NotificationController::~NotificationController() = default;

void NotificationController::ApplyState(const NotificationBackendState& State) {
    const bool Changed = SupportedValue != State.Supported ||
        PermissionGrantedValue != State.Granted || StatusValue != State.Status;
    SupportedValue = State.Supported;
    PermissionGrantedValue = State.Granted;
    StatusValue = State.Status;
    if (Changed) {
        emit StateChanged();
    }
}

bool NotificationController::Supported() const noexcept {
    return SupportedValue;
}

bool NotificationController::PermissionGranted() const noexcept {
    return PermissionGrantedValue;
}

bool NotificationController::Enabled() const noexcept {
    return EnabledValue;
}

QString NotificationController::Status() const {
    return StatusValue;
}

void NotificationController::SetEnabled(bool EnabledValueInput) {
    if (EnabledValue == EnabledValueInput) {
        return;
    }
    EnabledValue = EnabledValueInput;
    if (!EnabledValue) {
        CancelAllTimers();
    }
    emit StateChanged();
}

void NotificationController::SetStatus(QString StatusValueInput) {
    if (StatusValue == StatusValueInput) {
        return;
    }
    StatusValue = std::move(StatusValueInput);
    emit StateChanged();
}

void NotificationController::RefreshPermission() {
    if (Backend) {
        ApplyState(Backend->Refresh());
    }
}

void NotificationController::RequestPermission() {
    RefreshPermission();
    if (!SupportedValue || PermissionGrantedValue || !Backend) {
        return;
    }
    Backend->RequestPermission();
    SetStatus(QStringLiteral("请在系统弹窗中允许通知"));
    QTimer::singleShot(800, this, [this]() { RefreshPermission(); });
}

void NotificationController::ScheduleTimer(
    const QString& Id,
    int DelaySeconds,
    const QString& Title,
    const QString& Body) {
    if (!EnabledValue || Id.isEmpty() || DelaySeconds <= 0) {
        return;
    }
    RefreshPermission();
    if (!SupportedValue || !Backend) {
        return;
    }
    if (!PermissionGrantedValue) {
        SetStatus(QStringLiteral("通知权限、系统开关或通知渠道未开启"));
        return;
    }
    Backend->ScheduleTimer(Id, DelaySeconds, Title, Body);
}

void NotificationController::CancelTimer(const QString& Id) {
    if (Id.isEmpty()) {
        return;
    }
    if (Backend) {
        Backend->CancelTimer(Id);
    }
}

void NotificationController::CancelAllTimers() {
    if (Backend) {
        Backend->CancelAllTimers();
    }
}

void NotificationController::OpenSettings() {
    if (Backend) {
        Backend->OpenSettings();
    }
}

void NotificationController::ShowTestNotification() {
    if (!EnabledValue) {
        return;
    }
    if (SupportedValue && Backend) {
        Backend->ShowTestNotification();
    }
}

}  // namespace Menu::Client
