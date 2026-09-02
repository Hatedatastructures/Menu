#include <QtTest/QtTest>

#include <CookingTimerCoordinator.hpp>

class CookingTimerCoordinatorTest final : public QObject {
    Q_OBJECT

private slots:
    void GeneratesStableDistinctIds();
    void IgnoresInvalidTimerRequests();
};

void CookingTimerCoordinatorTest::GeneratesStableDistinctIds() {
    const QString First = Menu::Client::CookingTimerCoordinator::TimerId(
        QStringLiteral("recipe.tomato"), 1);
    QCOMPARE(First, Menu::Client::CookingTimerCoordinator::TimerId(
                        QStringLiteral("recipe.tomato"), 1));
    QVERIFY(!First.isEmpty());
    QVERIFY(First != Menu::Client::CookingTimerCoordinator::TimerId(
        QStringLiteral("recipe.tomato"), 2));
    QVERIFY(First != Menu::Client::CookingTimerCoordinator::TimerId(
        QStringLiteral("recipe.tofu"), 1));
}

void CookingTimerCoordinatorTest::IgnoresInvalidTimerRequests() {
    Menu::Client::NotificationController Notifications;
    Menu::Client::CookingTimerCoordinator Coordinator(Notifications);
    Coordinator.StartTimer({}, 0, 0, {}, {});
    Coordinator.StopTimer({}, 0);
    Coordinator.StopAllTimers();
    QVERIFY(true);
}

QTEST_GUILESS_MAIN(CookingTimerCoordinatorTest)
#include "CookingTimerCoordinatorTest.moc"
