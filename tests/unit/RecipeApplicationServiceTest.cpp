#include <gtest/gtest.h>

#include <Application/RecipeApplicationService.hpp>
#include <TestFixtures/ApplicationFixtures.hpp>

#include <algorithm>

TEST(RecipeApplicationServiceTest, ReturnsOnlyThreeRecipesThatFitUserConstraints) {
    auto Service = Menu::Tests::ApplicationFixtures::WithSeedRecipes();
    const Menu::Domain::RecommendationRequest Request{
        2, 35, {"中餐"}, {"花生"}, {"ingredient.rice"}, {"炒锅"}, {}};

    const auto Result = Service.RecommendTonight(Request);

    ASSERT_TRUE(Result.HasValue());
    ASSERT_EQ(Result.Value().size(), 2U);
    for (const auto& Recommendation : Result.Value()) {
        EXPECT_EQ(Recommendation.RecipeValue.Cuisine, "中餐");
        EXPECT_EQ(Recommendation.RecipeValue.Allergens.end(),
                  std::find(Recommendation.RecipeValue.Allergens.begin(),
                            Recommendation.RecipeValue.Allergens.end(), "花生"));
    }
}

TEST(RecipeApplicationServiceTest, ListsPublishedRecipesThroughRepositoryPort) {
    auto Service = Menu::Tests::ApplicationFixtures::WithSeedRecipes();

    const auto Result = Service.ListPublishedRecipes();

    ASSERT_TRUE(Result.HasValue());
    EXPECT_EQ(Result.Value().size(), 4U);
}
