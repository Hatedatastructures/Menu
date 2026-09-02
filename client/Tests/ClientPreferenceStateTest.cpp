#include <QtTest/QtTest>

#include <ClientPreferenceState.hpp>

class ClientPreferenceStateTest final : public QObject {
    Q_OBJECT

private slots:
    void ClampsNumericValuesAndNormalizesLists();
    void RoundTripsSnapshot();
};

void ClientPreferenceStateTest::ClampsNumericValuesAndNormalizesLists() {
    Menu::Client::ClientPreferenceState State;
    QVERIFY(State.SetServingCount(99));
    QCOMPARE(State.ServingCount(), 24);
    QVERIFY(State.SetAvailableMinutes(1));
    QCOMPARE(State.AvailableMinutes(), 10);
    QVERIFY(State.SetAllergies({QStringLiteral(" 虾 "), QStringLiteral("虾"), QStringLiteral(""),
                                QStringLiteral("花生")}));
    QCOMPARE(State.Allergies(), QStringList({QStringLiteral("虾"), QStringLiteral("花生")}));
}

void ClientPreferenceStateTest::RoundTripsSnapshot() {
    Menu::Client::ClientPreferenceState State;
    QVERIFY(State.SetServingCount(4));
    QVERIFY(State.SetAvailableMinutes(60));
    QVERIFY(State.SetPantryIngredientIds({QStringLiteral("ingredient.rice")}));
    QVERIFY(State.SetRemindersEnabled(false));

    const auto Snapshot = State.Snapshot();
    Menu::Client::ClientPreferenceState Restored;
    QVERIFY(Restored.Apply(Snapshot));
    QCOMPARE(Restored.ServingCount(), 4);
    QCOMPARE(Restored.AvailableMinutes(), 60);
    QCOMPARE(Restored.PantryIngredientIds(), QStringList({QStringLiteral("ingredient.rice")}));
    QVERIFY(!Restored.RemindersEnabled());
    QVERIFY(!Restored.Apply(Snapshot));
}

QTEST_GUILESS_MAIN(ClientPreferenceStateTest)
#include "ClientPreferenceStateTest.moc"
