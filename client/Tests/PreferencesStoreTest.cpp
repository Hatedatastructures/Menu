#include <QtTest/QtTest>

#include <PreferencesStore.hpp>

class PreferencesStoreTest final : public QObject {
    Q_OBJECT

private slots:
    void RoundTripsPreferences();
    void RejectsMalformedOrOutOfRangePreferences();
};

void PreferencesStoreTest::RoundTripsPreferences() {
    Menu::Client::PreferencesSnapshot Snapshot;
    Snapshot.ServingCount = 4;
    Snapshot.PreferredCuisines = {QStringLiteral("中餐"), QStringLiteral("日系")};
    Snapshot.Allergies = {QStringLiteral("花生")};
    Snapshot.AvailableMinutes = 60;
    Snapshot.PantryIngredientIds = {QStringLiteral("ingredient.rice")};
    Snapshot.Cookware = {QStringLiteral("炒锅")};
    Snapshot.RemindersEnabled = false;

    const auto Restored = Menu::Client::PreferencesStore::Decode(
        Menu::Client::PreferencesStore::Encode(Snapshot));
    QVERIFY(Restored.has_value());
    QCOMPARE(Restored->ServingCount, 4);
    QCOMPARE(Restored->PreferredCuisines, Snapshot.PreferredCuisines);
    QCOMPARE(Restored->Allergies, Snapshot.Allergies);
    QCOMPARE(Restored->AvailableMinutes, 60);
    QCOMPARE(Restored->PantryIngredientIds, Snapshot.PantryIngredientIds);
    QCOMPARE(Restored->Cookware, Snapshot.Cookware);
    QCOMPARE(Restored->RemindersEnabled, false);
}

void PreferencesStoreTest::RejectsMalformedOrOutOfRangePreferences() {
    QVERIFY(!Menu::Client::PreferencesStore::Decode(
                 QByteArrayLiteral("not-json"))
                 .has_value());
    QVERIFY(!Menu::Client::PreferencesStore::Decode(QByteArrayLiteral(
                 R"json({"schemaVersion":2,"kind":"preferences","servings":0,
                    "cuisines":[],"allergies":[],"availableMinutes":45,
                    "pantryIngredientIds":[],"cookware":[],"remindersEnabled":true})json"))
                 .has_value());
}

QTEST_GUILESS_MAIN(PreferencesStoreTest)
#include "PreferencesStoreTest.moc"
