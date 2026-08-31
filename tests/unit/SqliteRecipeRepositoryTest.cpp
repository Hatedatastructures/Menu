#include <gtest/gtest.h>

#include <Infrastructure/MigrationRunner.hpp>
#include <Infrastructure/SeedData.hpp>
#include <Infrastructure/SqliteRecipeRepository.hpp>
#include <TestFixtures/DatabaseFixtures.hpp>

TEST(SqliteRecipeRepositoryTest, ReturnsPublishedRecipesAndTheirSteps) {
    auto Database = Menu::Tests::DatabaseFixtures::OpenTemporary();
    ASSERT_TRUE(Menu::Infrastructure::MigrationRunner::Apply(Database).HasValue());
    ASSERT_TRUE(Menu::Infrastructure::SeedData::InsertIfEmpty(Database).HasValue());
    Menu::Infrastructure::SqliteRecipeRepository Repository(Database);

    const auto Recipes = Repository.ListPublished();

    ASSERT_TRUE(Recipes.HasValue());
    ASSERT_FALSE(Recipes.Value().empty());
    EXPECT_EQ(Recipes.Value()[0].Status, "published");
    EXPECT_FALSE(Recipes.Value()[0].Ingredients.empty());
    EXPECT_FALSE(Recipes.Value()[0].Steps.empty());
}

TEST(SqliteRecipeRepositoryTest, FindsKnownRecipeAndReturnsEmptyForUnknownId) {
    auto Database = Menu::Tests::DatabaseFixtures::OpenTemporary();
    ASSERT_TRUE(Menu::Infrastructure::MigrationRunner::Apply(Database).HasValue());
    ASSERT_TRUE(Menu::Infrastructure::SeedData::InsertIfEmpty(Database).HasValue());
    Menu::Infrastructure::SqliteRecipeRepository Repository(Database);

    const auto Recipes = Repository.ListPublished();
    ASSERT_TRUE(Recipes.HasValue());
    ASSERT_FALSE(Recipes.Value().empty());

    const auto Found = Repository.FindPublishedById(Recipes.Value()[0].Id);
    const auto Missing = Repository.FindPublishedById("recipe.does-not-exist");
    ASSERT_TRUE(Found.HasValue());
    ASSERT_TRUE(Missing.HasValue());
    ASSERT_TRUE(Found.Value().has_value());
    EXPECT_FALSE(Missing.Value().has_value());
    EXPECT_EQ(Found.Value()->Id, Recipes.Value()[0].Id);
}
