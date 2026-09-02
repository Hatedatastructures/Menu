#include "NotificationBackend.hpp"

#if defined(Q_OS_ANDROID)
#include <QJniObject>
#include <QtCore/qcoreapplication_platform.h>
#endif

namespace Menu::Client {
namespace {

#if defined(Q_OS_ANDROID)

class AndroidNotificationBackend final : public NotificationBackend {
public:
    NotificationBackendState Refresh() override {
        const auto Context = QNativeInterface::QAndroidApplication::context();
        const bool Granted = Context.object() != nullptr &&
            QJniObject::callStaticMethod<jboolean>(
                "com/menu/cookflow/NotificationBridge",
                "hasNotificationPermission",
                "(Landroid/content/Context;)Z",
                Context.object<jobject>());
        return {true, Granted, Granted
            ? QStringLiteral("通知已开启")
            : QStringLiteral("通知权限、系统开关或通知渠道未开启")};
    }

    void RequestPermission() override {
        const auto Context = QNativeInterface::QAndroidApplication::context();
        if (Context.object() != nullptr) {
            QJniObject::callStaticMethod<void>(
                "com/menu/cookflow/NotificationBridge",
                "requestNotificationPermission",
                "(Landroid/content/Context;)V",
                Context.object<jobject>());
        }
    }

    void ScheduleTimer(
        const QString& Id,
        int DelaySeconds,
        const QString& Title,
        const QString& Body) override {
        const auto Context = QNativeInterface::QAndroidApplication::context();
        if (Context.object() == nullptr) {
            return;
        }
        const QJniObject JniId = QJniObject::fromString(Id);
        const QJniObject JniTitle = QJniObject::fromString(Title);
        const QJniObject JniBody = QJniObject::fromString(Body);
        QJniObject::callStaticMethod<void>(
            "com/menu/cookflow/NotificationBridge",
            "scheduleTimer",
            "(Landroid/content/Context;Ljava/lang/String;JLjava/lang/String;Ljava/lang/String;)V",
            Context.object<jobject>(), JniId.object<jstring>(),
            static_cast<jlong>(DelaySeconds), JniTitle.object<jstring>(),
            JniBody.object<jstring>());
    }

    void CancelTimer(const QString& Id) override {
        const auto Context = QNativeInterface::QAndroidApplication::context();
        if (Context.object() == nullptr) {
            return;
        }
        const QJniObject JniId = QJniObject::fromString(Id);
        QJniObject::callStaticMethod<void>(
            "com/menu/cookflow/NotificationBridge",
            "cancelTimer",
            "(Landroid/content/Context;Ljava/lang/String;)V",
            Context.object<jobject>(), JniId.object<jstring>());
    }

    void CancelAllTimers() override {
        const auto Context = QNativeInterface::QAndroidApplication::context();
        if (Context.object() != nullptr) {
            QJniObject::callStaticMethod<void>(
                "com/menu/cookflow/NotificationBridge",
                "cancelAllTimers",
                "(Landroid/content/Context;)V",
                Context.object<jobject>());
        }
    }

    void OpenSettings() override {
        const auto Context = QNativeInterface::QAndroidApplication::context();
        if (Context.object() != nullptr) {
            QJniObject::callStaticMethod<void>(
                "com/menu/cookflow/NotificationBridge",
                "openNotificationSettings",
                "(Landroid/content/Context;)V",
                Context.object<jobject>());
        }
    }

    void ShowTestNotification() override {
        const auto Context = QNativeInterface::QAndroidApplication::context();
        if (Context.object() != nullptr) {
            QJniObject::callStaticMethod<void>(
                "com/menu/cookflow/NotificationBridge",
                "showTestNotification",
                "(Landroid/content/Context;)V",
                Context.object<jobject>());
        }
    }
};

#else

class DesktopNotificationBackend final : public NotificationBackend {
public:
    NotificationBackendState Refresh() override {
        return {false, false, QStringLiteral("当前平台不支持系统通知")};
    }

    void RequestPermission() override {}
    void ScheduleTimer(
        const QString&, int, const QString&, const QString&) override {}
    void CancelTimer(const QString&) override {}
    void CancelAllTimers() override {}
    void OpenSettings() override {}
    void ShowTestNotification() override {}
};

#endif

}  // namespace

std::unique_ptr<NotificationBackend> CreateNotificationBackend() {
#if defined(Q_OS_ANDROID)
    return std::make_unique<AndroidNotificationBackend>();
#else
    return std::make_unique<DesktopNotificationBackend>();
#endif
}

}  // namespace Menu::Client
