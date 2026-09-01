#include <ClientApi.hpp>

#include <QGuiApplication>
#include <QDebug>
#include <QQmlError>
#include <QQmlContext>
#include <QQmlApplicationEngine>
#include <QTimer>
#include <QWindow>

int main(int ArgumentCount, char* Arguments[]) {
    QGuiApplication Application(ArgumentCount, Arguments);
    const bool IsSmoke = qEnvironmentVariableIsSet("MENU_QML_SMOKE");
    if (IsSmoke) {
        qInfo().noquote() << "QmlSmoke: main-entered";
    }
    QString ServiceUrl = qEnvironmentVariable("MENU_API_BASE_URL");
    if (ServiceUrl.isEmpty()) {
#if defined(Q_OS_ANDROID)
        ServiceUrl = QStringLiteral("http://10.0.2.2:8080");
#else
        ServiceUrl = QStringLiteral("http://127.0.0.1:8080");
#endif
    }
    Menu::Client::ClientApi Api(ServiceUrl);
    if (IsSmoke) {
        qInfo().noquote() << "QmlSmoke: api-created";
    }
    QQmlApplicationEngine Engine;
    Engine.rootContext()->setContextProperty(QStringLiteral("MenuApi"), &Api);
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
