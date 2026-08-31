#include <gtest/gtest.h>

#include <Domain/Recommendation.hpp>
#include <TestFixtures/RecipeFixtures.hpp>

TEST(RecommendationRulesTest, FiltersAllergiesTimeCuisineAndCookware) {
    Menu::Domain::Recipe SafeRecipe =
        Menu::Tests::RecipeFixtures::SingleRecipeWithIngredients(
            2,
            {Menu::Domain::RecipeIngredient{"ingredient.rice", 200.0, "g", true}});
    SafeRecipe.Id = "recipe.safe";
    SafeRecipe.Cuisine = "中餐";

    Menu::Domain::Recipe AllergyRecipe = SafeRecipe;
    AllergyRecipe.Id = "recipe.allergy";
    AllergyRecipe.Allergens = {"花生"};

    Menu::Domain::Recipe SlowRecipe = SafeRecipe;
    SlowRecipe.Id = "recipe.slow";
    SlowRecipe.CookMinutes = 60;

    const Menu::Domain::RecommendationRequest Request{
        2, 35, {"中餐"}, {"花生"}, {"ingredient.rice"}, {"炒锅"}, {}};
    const auto Result = Menu::Domain::RuleBasedRecommendation::Rank(
        Request, {SafeRecipe, AllergyRecipe, SlowRecipe});

    ASSERT_TRUE(Result.HasValue());
    ASSERT_EQ(Result.Value().size(), 1U);
    EXPECT_EQ(Result.Value()[0].RecipeValue.Id, "recipe.safe");
}

TEST(RecommendationRulesTest, ReturnsAtMostThreeStableRecommendations) {
    std::vector<Menu::Domain::Recipe> Recipes;
    for (int Index = 0; Index < 5; ++Index) {
        auto Recipe = Menu::Tests::RecipeFixtures::SingleRecipeWithIngredients(
            2,
            {Menu::Domain::RecipeIngredient{"ingredient.rice", 200.0, "g", true}});
        Recipe.Id = "recipe." + std::to_string(Index);
        Recipes.push_back(std::move(Recipe));
    }

    const Menu::Domain::RecommendationRequest Request{
        2, 35, {"中餐"}, {}, {"ingredient.rice"}, {"炒锅"}, {}};
    const auto Result = Menu::Domain::RuleBasedRecommendation::Rank(Request, Recipes);

    ASSERT_TRUE(Result.HasValue());
    ASSERT_EQ(Result.Value().size(), 3U);
    EXPECT_EQ(Result.Value()[0].RecipeValue.Id, "recipe.0");
    EXPECT_EQ(Result.Value()[1].RecipeValue.Id, "recipe.1");
    EXPECT_EQ(Result.Value()[2].RecipeValue.Id, "recipe.2");
}
