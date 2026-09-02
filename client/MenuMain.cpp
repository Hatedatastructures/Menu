#include <ClientApi.hpp>
#include <CookingTimerCoordinator.hpp>
#include <DisplayMetrics.hpp>
#include <NotificationController.hpp>
#include <ScreenAwakeController.hpp>
#include <ThemeController.hpp>

#include <QGuiApplication>
#include <QDebug>
#include <QQmlError>
#include <QQmlContext>
#include <QQmlApplicationEngine>
#include <QQuickStyle>
#include <QTimer>
#include <QWindow>

int main(int ArgumentCount, char* Arguments[]) {
    QGuiApplication Application(ArgumentCount, Arguments);
    Application.setOrganizationName(QStringLiteral("Menu"));
    Application.setApplicationName(QStringLiteral("CookFlow"));
    QQuickStyle::setStyle(QStringLiteral("Basic"));
    const bool IsSmoke = qEnvironmentVariableIsSet("MENU_QML_SMOKE");
    if (IsSmoke) {
        qInfo().noquote() << "QmlSmoke: main-entered";
    }
    const QString ServiceUrl = qEnvironmentVariable("MENU_API_BASE_URL");
    Menu::Client::ClientApi Api(ServiceUrl);
    Menu::Client::DisplayMetrics Display(&Application);
    Menu::Client::NotificationController Notifications(&Application);
    Menu::Client::ScreenAwakeController ScreenAwake(&Application);
    Menu::Client::CookingTimerCoordinator CookingTimers(Notifications, &Application);
    Menu::Client::ThemeController Theme(&Application);
    Notifications.SetEnabled(Api.RemindersEnabled());
    QObject::connect(&Api, &Menu::Client::ClientApi::PreferencesChanged,
        &Notifications, [&Api, &Notifications]() {
            Notifications.SetEnabled(Api.RemindersEnabled());
        });
    if (IsSmoke) {
        qInfo().noquote() << "QmlSmoke: api-created";
    }
    QQmlApplicationEngine Engine;
    Engine.rootContext()->setContextProperty(QStringLiteral("MenuApi"), &Api);
    Engine.rootContext()->setContextProperty(QStringLiteral("MenuDisplay"), &Display);
    Engine.rootContext()->setContextProperty(
        QStringLiteral("MenuNotifications"), &Notifications);
    Engine.rootContext()->setContextProperty(
        QStringLiteral("MenuScreenAwake"), &ScreenAwake);
    Engine.rootContext()->setContextProperty(QStringLiteral("MenuTimers"), &CookingTimers);
    Engine.rootContext()->setContextProperty(QStringLiteral("MenuTheme"), &Theme);
    if (IsSmoke) {
        qInfo().noquote() << "QmlSmoke: engine-created";
    }
    QObject::connect(&Engine, &QQmlApplicationEngine::warnings,
        [](const QList<QQmlError>& Errors) {
            for (const QQmlError& Error : Errors) {
                qWarning().noquote() << "QmlWarning:" << Error.toString();
            }
        });
    if (IsSmoke) {
        qInfo().noquote() << "QmlSmoke: loading";
    }
    Engine.loadFromModule("MenuClient", "Main");
    if (IsSmoke) {
        qInfo().noquote() << "QmlSmoke: loaded" << Engine.rootObjects().size();
    }
    if (Engine.rootObjects().isEmpty()) {
        return 1;
    }

    QObject* RootObject = Engine.rootObjects().constFirst();
    if (auto* Window = qobject_cast<QWindow*>(RootObject)) {
        Display.AttachWindow(Window);
        Theme.RefreshSystemBars();
        QTimer::singleShot(500, &Theme, &Menu::Client::ThemeController::RefreshSystemBars);
    }
    if (IsSmoke) {
        qInfo().noquote() << "QmlSmoke: root" << RootObject->metaObject()->className();
        if (auto* Window = qobject_cast<QWindow*>(RootObject)) {
            qInfo().noquote() << "QmlSmoke: window" << Window->size()
                              << "visible=" << Window->isVisible();
        }
    }

    if (IsSmoke) {
        if (auto* Window = qobject_cast<QWindow*>(Engine.rootObjects().constFirst())) {
            Window->hide();
            qInfo().noquote() << "QmlSmoke: window-hidden";
        }
    }

    QTimer SmokeTimer;
    if (IsSmoke) {
        SmokeTimer.setSingleShot(true);
        SmokeTimer.setInterval(1500);
        QObject::connect(&SmokeTimer, &QTimer::timeout, []() {
            qInfo().noquote() << "QmlSmoke: timer-fired";
            QCoreApplication::quit();
        });
        SmokeTimer.start();
        qInfo().noquote() << "QmlSmoke: timer-started";
    }
    return Application.exec();
}
