#include <QtTest/QtTest>

#include <RecipeModel.hpp>

class RecipeModelTest final : public QObject {
    Q_OBJECT

private slots:
    void StoresRecipeAndIngredientDisplayFields();
};

void RecipeModelTest::StoresRecipeAndIngredientDisplayFields() {
    Menu::Client::RecipeModel Model;
    const QJsonDocument Document = QJsonDocument::fromJson(QByteArrayLiteral(R"json(
        [{
            "id": "recipe.tomato",
            "name": "番茄炒蛋",
            "cuisine": "中餐",
            "imagePath": "assets/media/tomato-egg.png",
            "totalMinutes": 20,
            "servings": 2,
            "difficulty": 1,
            "ingredients": [{
                "ingredientId": "ingredient.tomato",
                "ingredientName": "番茄",
                "category": "蔬菜",
                "defaultUnit": "g",
                "isPantryStaple": false,
                "quantity": 200,
                "unit": "g",
                "required": true
            }]
        }]
    )json"));
    const QJsonArray Recipes = Document.array();

    Model.SetRecipes(Recipes);

    QCOMPARE(Model.Count(), 1);
    QCOMPARE(Model.rowCount(), 1);
    const QModelIndex Index = Model.index(0, 0);
    QCOMPARE(Model.data(Index, Menu::Client::RecipeModel::NameRole).toString(), "番茄炒蛋");
    QCOMPARE(Model.data(Index, Menu::Client::RecipeModel::CuisineRole).toString(), "中餐");
    QCOMPARE(Model.data(Index, Menu::Client::RecipeModel::MediaUrlRole).toString(),
        "qrc:/qt/qml/MenuClient/assets/media/tomato-egg.png");
    const QVariantList Ingredients =
        Model.data(Index, Menu::Client::RecipeModel::IngredientsRole).toList();
    QCOMPARE(Ingredients.size(), 1);
    const QVariantMap Ingredient = Ingredients.at(0).toMap();
    QCOMPARE(Ingredient.value("ingredientName").toString(), "番茄");
    QCOMPARE(Ingredient.value("category").toString(), "蔬菜");
    QCOMPARE(Ingredient.value("defaultUnit").toString(), "g");
    QCOMPARE(Ingredient.value("isPantryStaple").toBool(), false);
}

QTEST_GUILESS_MAIN(RecipeModelTest)
#include "RecipeModelTest.moc"
