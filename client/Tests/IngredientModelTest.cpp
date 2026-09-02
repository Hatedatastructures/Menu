#include <QtTest/QtTest>

#include <IngredientModel.hpp>

class IngredientModelTest final : public QObject {
    Q_OBJECT

private slots:
    void ExposesDisplayFieldsForFilteredViews();
};

void IngredientModelTest::ExposesDisplayFieldsForFilteredViews() {
    Menu::Client::IngredientModel Model;
    const QJsonDocument Document = QJsonDocument::fromJson(QByteArrayLiteral(R"json(
        [{"id":"ingredient.tomato","name":"番茄","category":"蔬菜",
          "defaultUnit":"g","isPantryStaple":false}]
    )json"));
    Model.SetIngredients(Document.array());

    QCOMPARE(Model.Count(), 1);
    QCOMPARE(Model.IngredientAt(0).value(QStringLiteral("name")).toString(),
             QStringLiteral("番茄"));
    QCOMPARE(Model.IngredientAt(-1), QVariantMap());
}

QTEST_GUILESS_MAIN(IngredientModelTest)
#include "IngredientModelTest.moc"
