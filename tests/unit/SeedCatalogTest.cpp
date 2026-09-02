#include <Infrastructure/SeedCatalog.hpp>

#include <gtest/gtest.h>

#include <set>

TEST(SeedCatalogTest, ContainsStableIngredientsAndPublishedRecipes) {
    const auto Ingredients = Menu::Infrastructure::SeedCatalog::Ingredients();
    const auto Recipes = Menu::Infrastructure::SeedCatalog::Recipes();

    ASSERT_EQ(Ingredients.size(), 20U);
    ASSERT_EQ(Recipes.size(), 12U);

    std::set<std::string> IngredientIds;
    for (const auto& Ingredient : Ingredients) {
        EXPECT_TRUE(IngredientIds.insert(Ingredient.Id).second);
        EXPECT_FALSE(Ingredient.Name.empty());
        EXPECT_FALSE(Ingredient.DefaultUnit.empty());
    }

    std::set<std::string> RecipeIds;
    for (const auto& Recipe : Recipes) {
        EXPECT_TRUE(RecipeIds.insert(Recipe.Id).second);
        EXPECT_EQ(Recipe.Status, "published");
        EXPECT_EQ(Recipe.Steps.size(), 3U);
        EXPECT_FALSE(Recipe.Ingredients.empty());
        for (const auto& Step : Recipe.Steps) {
            if (Step.DurationSeconds <= 0) {
                EXPECT_FALSE(Step.HasTimer);
            }
        }
    }
}
