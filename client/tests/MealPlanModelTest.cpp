#include <QtTest/QtTest>

#include <MealPlanModel.hpp>

class MealPlanModelTest final : public QObject {
    Q_OBJECT

private slots:
    void StoresPlansAndMergedIngredientDisplayFields();
};

void MealPlanModelTest::StoresPlansAndMergedIngredientDisplayFields() {
    Menu::Client::MealPlanModel Model;
    const QJsonDocument Document = QJsonDocument::fromJson(QByteArrayLiteral(R"json(
        [{
            "id": "plan.one",
            "planDate": "2026-09-01",
            "items": [{"recipeId":"recipe.tomato","recipeName":"番茄炒蛋","servings":2}],
            "combinedIngredients": [{"ingredientId":"ingredient.tomato","ingredientName":"番茄","category":"蔬菜","defaultUnit":"g","quantity":500,"unit":"g","required":true}]
        }]
    )json"));

    Model.SetPlans(Document.array());

    QCOMPARE(Model.Count(), 1);
    QCOMPARE(Model.rowCount(), 1);
    const QVariantMap Plan = Model.PlanAt(0);
    QCOMPARE(Plan.value(QStringLiteral("planDate")).toString(), QStringLiteral("2026-09-01"));
    const QVariantMap Ingredient = Plan.value(QStringLiteral("combinedIngredients"))
                                       .toList()
                                       .at(0)
                                       .toMap();
    QCOMPARE(Ingredient.value(QStringLiteral("ingredientName")).toString(), QStringLiteral("番茄"));
    QCOMPARE(Ingredient.value(QStringLiteral("category")).toString(), QStringLiteral("蔬菜"));
    QCOMPARE(Ingredient.value(QStringLiteral("quantity")).toDouble(), 500.0);
}

QTEST_GUILESS_MAIN(MealPlanModelTest)
#include "MealPlanModelTest.moc"
