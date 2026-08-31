#include <gtest/gtest.h>

#include <Domain/RecipeRules.hpp>
#include <TestFixtures/RecipeFixtures.hpp>

TEST(RecipeRulesTest, ScalesRequiredAndOptionalIngredientsByServingRatio) {
    const Menu::Domain::Recipe Recipe =
        Menu::Tests::RecipeFixtures::SingleRecipeWithIngredients(
            2,
            {Menu::Domain::RecipeIngredient{"ingredient.rice", 200.0, "g", true},
             Menu::Domain::RecipeIngredient{"ingredient.salt", 2.0, "g", false}});

    const auto Result = Menu::Domain::RecipeRules::ScaleRecipeIngredients(Recipe, 4);

    ASSERT_TRUE(Result.HasValue());
    ASSERT_EQ(Result.Value().size(), 2U);
    EXPECT_DOUBLE_EQ(Result.Value()[0].Quantity, 400.0);
    EXPECT_TRUE(Result.Value()[0].Required);
    EXPECT_DOUBLE_EQ(Result.Value()[1].Quantity, 4.0);
    EXPECT_FALSE(Result.Value()[1].Required);
}

TEST(RecipeRulesTest, RejectsInvalidServingCountAndDuplicateStepOrder) {
    Menu::Domain::Recipe Recipe =
        Menu::Tests::RecipeFixtures::SingleRecipeWithIngredients(
            0,
            {Menu::Domain::RecipeIngredient{"ingredient.rice", 200.0, "g", true}});

    const auto ServingResult = Menu::Domain::RecipeRules::ValidateRecipe(Recipe);
    ASSERT_FALSE(ServingResult.HasValue());

    Recipe.Servings = 2;
    Recipe.Steps[1].StepOrder = Recipe.Steps[0].StepOrder;
    const auto StepResult = Menu::Domain::RecipeRules::ValidateRecipe(Recipe);
    EXPECT_FALSE(StepResult.HasValue());
}

TEST(RecipeRulesTest, NormalizesRegisteredIngredientAlias) {
    const std::vector<Menu::Domain::Ingredient> Ingredients = {
        Menu::Domain::Ingredient{
            "ingredient.tomato", "番茄", {"西红柿", "蕃茄"}, "蔬菜", "g", true, "", ""},
    };

    const auto Result = Menu::Domain::RecipeRules::NormalizeIngredientAlias(
        "  西红柿  ", Ingredients);

    ASSERT_TRUE(Result.HasValue());
    EXPECT_EQ(Result.Value(), "ingredient.tomato");
}

TEST(RecipeRulesTest, CountsRequiredIngredientsAgainstPantryIds) {
    const Menu::Domain::Recipe Recipe =
        Menu::Tests::RecipeFixtures::SingleRecipeWithIngredients(
            2,
            {Menu::Domain::RecipeIngredient{"ingredient.rice", 200.0, "g", true},
             Menu::Domain::RecipeIngredient{"ingredient.salt", 2.0, "g", false},
             Menu::Domain::RecipeIngredient{"ingredient.egg", 2.0, "个", true}});

    const auto Availability = Menu::Domain::RecipeRules::BuildIngredientAvailability(
        Recipe, {"ingredient.rice"});

    EXPECT_EQ(Availability.AvailableCount, 1U);
    EXPECT_EQ(Availability.MissingCount, 1U);
}
