#include <QtTest/QtTest>

#include <NotificationController.hpp>

#include <memory>

namespace {

class FakeNotificationBackend final : public Menu::Client::NotificationBackend {
public:
    Menu::Client::NotificationBackendState State{true, true, QStringLiteral("fake-ready")};
    int RequestCount = 0;
    int ScheduleCount = 0;
    int CancelCount = 0;
    int CancelAllCount = 0;
    int SettingsCount = 0;
    int TestCount = 0;

    Menu::Client::NotificationBackendState Refresh() override { return State; }
    void RequestPermission() override { ++RequestCount; }
    void ScheduleTimer(const QString&, int, const QString&, const QString&) override { ++ScheduleCount; }
    void CancelTimer(const QString&) override { ++CancelCount; }
    void CancelAllTimers() override { ++CancelAllCount; }
    void OpenSettings() override { ++SettingsCount; }
    void ShowTestNotification() override { ++TestCount; }
};

}  // namespace

class NotificationControllerTest final : public QObject {
    Q_OBJECT

private slots:
    void StartsInDeterministicState();
    void NoopOperationsAreSafeOffAndroid();
    void DelegatesOperationsToInjectedBackend();
};

void NotificationControllerTest::StartsInDeterministicState() {
    Menu::Client::NotificationController Controller;
#if defined(Q_OS_ANDROID)
    QVERIFY(Controller.Supported());
#else
    QVERIFY(!Controller.Supported());
    QVERIFY(!Controller.PermissionGranted());
#endif
    QVERIFY(!Controller.Status().isEmpty());
}

void NotificationControllerTest::NoopOperationsAreSafeOffAndroid() {
    Menu::Client::NotificationController Controller;
    Controller.RefreshPermission();
    Controller.RequestPermission();
    Controller.ScheduleTimer(QStringLiteral("fixture"), 1,
                             QStringLiteral("测试"), QStringLiteral("正文"));
    Controller.CancelTimer(QStringLiteral("fixture"));
#if !defined(Q_OS_ANDROID)
    QVERIFY(!Controller.PermissionGranted());
    QVERIFY(Controller.Status().contains(QStringLiteral("平台")));
#endif
}

void NotificationControllerTest::DelegatesOperationsToInjectedBackend() {
    auto Backend = std::make_unique<FakeNotificationBackend>();
    FakeNotificationBackend* BackendPointer = Backend.get();
    auto Holder = std::make_shared<std::unique_ptr<FakeNotificationBackend>>(std::move(Backend));
    Menu::Client::NotificationController Controller(nullptr,
        [Holder]() mutable { return std::unique_ptr<Menu::Client::NotificationBackend>(std::move(*Holder)); });

    Controller.RefreshPermission();
    QVERIFY(Controller.Supported());
    QVERIFY(Controller.PermissionGranted());
    Controller.RequestPermission();
    Controller.ScheduleTimer(QStringLiteral("timer"), 10,
        QStringLiteral("title"), QStringLiteral("body"));
    Controller.CancelTimer(QStringLiteral("timer"));
    Controller.CancelAllTimers();
    Controller.OpenSettings();
    Controller.ShowTestNotification();

    QCOMPARE(BackendPointer->RequestCount, 0);
    QCOMPARE(BackendPointer->ScheduleCount, 1);
    QCOMPARE(BackendPointer->CancelCount, 1);
    QCOMPARE(BackendPointer->CancelAllCount, 1);
    QCOMPARE(BackendPointer->SettingsCount, 1);
    QCOMPARE(BackendPointer->TestCount, 1);
}

QTEST_GUILESS_MAIN(NotificationControllerTest)
#include "NotificationControllerTest.moc"
