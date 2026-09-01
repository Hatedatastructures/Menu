#include <gtest/gtest.h>

#include <Domain/PlanRules.hpp>

TEST(PlanRulesTest, MergesRepeatedIngredientsAndPreservesDisplayFields) {
    Menu::Domain::MealPlanItem FirstItem;
    FirstItem.RecipeId = "recipe.one";
    FirstItem.Servings = 2;
    FirstItem.Ingredients.push_back(Menu::Domain::RecipeIngredient{
        "ingredient.tomato", 200.0, "g", true, 1.0, "切块"});
    FirstItem.Ingredients.back().IngredientName = "番茄";
    FirstItem.Ingredients.back().IngredientCategory = "蔬菜";
    FirstItem.Ingredients.back().IngredientDefaultUnit = "g";
    FirstItem.Ingredients.back().IngredientIsPantryStaple = false;

    Menu::Domain::MealPlanItem SecondItem;
    SecondItem.RecipeId = "recipe.two";
    SecondItem.Servings = 2;
    SecondItem.Ingredients.push_back(Menu::Domain::RecipeIngredient{
        "ingredient.tomato", 300.0, "g", true, 1.0, "切片"});
    SecondItem.Ingredients.back().IngredientName = "番茄";
    SecondItem.Ingredients.back().IngredientCategory = "蔬菜";
    SecondItem.Ingredients.back().IngredientDefaultUnit = "g";
    SecondItem.Ingredients.back().IngredientIsPantryStaple = false;

    const auto Merged = Menu::Domain::PlanRules::MergeIngredients(
        std::vector<Menu::Domain::MealPlanItem>{FirstItem, SecondItem});

    ASSERT_EQ(Merged.size(), 1U);
    EXPECT_EQ(Merged[0].IngredientId, "ingredient.tomato");
    EXPECT_DOUBLE_EQ(Merged[0].Quantity, 500.0);
    EXPECT_EQ(Merged[0].IngredientName, "番茄");
    EXPECT_EQ(Merged[0].IngredientCategory, "蔬菜");
    EXPECT_EQ(Merged[0].IngredientDefaultUnit, "g");
}

TEST(PlanRulesTest, KeepsDifferentUnitsSeparate) {
    Menu::Domain::MealPlanItem PlanItem;
    PlanItem.Ingredients = {
        Menu::Domain::RecipeIngredient{"ingredient.oil", 10.0, "ml"},
        Menu::Domain::RecipeIngredient{"ingredient.oil", 1.0, "勺"}};

    const auto Merged = Menu::Domain::PlanRules::MergeIngredients(
        std::vector<Menu::Domain::MealPlanItem>{PlanItem});

    ASSERT_EQ(Merged.size(), 2U);
    EXPECT_EQ(Merged[0].Unit, "ml");
    EXPECT_EQ(Merged[1].Unit, "勺");
}
