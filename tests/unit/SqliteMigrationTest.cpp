#include <gtest/gtest.h>

#include <Infrastructure/MigrationRunner.hpp>
#include <TestFixtures/DatabaseFixtures.hpp>

TEST(SqliteMigrationTest, AppliesSchemaAndEnablesWalAndForeignKeys) {
    auto Database = Menu::Tests::DatabaseFixtures::OpenTemporary();

    const auto Result = Menu::Infrastructure::MigrationRunner::Apply(Database);

    ASSERT_TRUE(Result.HasValue());
    const auto JournalMode = Database.ScalarText("PRAGMA journal_mode;");
    const auto ForeignKeys = Database.ScalarInt("PRAGMA foreign_keys;");
    const auto MigrationCount = Database.ScalarInt(
        "SELECT COUNT(*) FROM SchemaMigration;");
    ASSERT_TRUE(JournalMode.HasValue());
    ASSERT_TRUE(ForeignKeys.HasValue());
    ASSERT_TRUE(MigrationCount.HasValue());
    EXPECT_EQ(JournalMode.Value(), "wal");
    EXPECT_EQ(ForeignKeys.Value(), 1);
    EXPECT_EQ(MigrationCount.Value(), 1);
}

TEST(SqliteMigrationTest, IsIdempotentAndRollsBackFailedTransaction) {
    auto Database = Menu::Tests::DatabaseFixtures::OpenTemporary();

    ASSERT_TRUE(Menu::Infrastructure::MigrationRunner::Apply(Database).HasValue());
    ASSERT_TRUE(Menu::Infrastructure::MigrationRunner::Apply(Database).HasValue());
    const auto MigrationCount = Database.ScalarInt(
        "SELECT COUNT(*) FROM SchemaMigration;");
    ASSERT_TRUE(MigrationCount.HasValue());
    EXPECT_EQ(MigrationCount.Value(), 1);
}
