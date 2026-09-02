#include "Infrastructure/SeedData.hpp"

#include "Infrastructure/SeedCatalog.hpp"
#include "Infrastructure/SeedDataRows.hpp"

namespace Menu::Infrastructure {

Foundation::Result<void> SeedData::InsertIfEmpty(SqliteDatabase& Database) {
    const auto IngredientCount = Database.ScalarInt("SELECT COUNT(*) FROM Ingredients;");
    const auto RecipeCount = Database.ScalarInt("SELECT COUNT(*) FROM Recipes;");
    if (!IngredientCount.HasValue()) {
        return Foundation::Result<void>::FromError(IngredientCount.ErrorValue());
    }
    if (!RecipeCount.HasValue()) {
        return Foundation::Result<void>::FromError(RecipeCount.ErrorValue());
    }
    if (IngredientCount.Value() > 0 && RecipeCount.Value() > 0) {
        return Foundation::Result<void>();
    }

    const auto BeginResult = Database.BeginTransaction();
    if (!BeginResult.HasValue()) {
        return BeginResult;
    }
    for (const Domain::Ingredient& IngredientValue : SeedCatalog::Ingredients()) {
        if (auto Result = SeedDataRows::InsertIngredient(Database, IngredientValue);
            !Result.HasValue()) {
            Database.Rollback();
            return Result;
        }
    }
    for (const Domain::Recipe& RecipeValue : SeedCatalog::Recipes()) {
        if (auto Result = SeedDataRows::InsertRecipe(Database, RecipeValue);
            !Result.HasValue()) {
            Database.Rollback();
            return Result;
        }
    }
    const auto CommitResult = Database.Commit();
    if (!CommitResult.HasValue()) {
        Database.Rollback();
        return CommitResult;
    }
    return Foundation::Result<void>();
}

}  // namespace Menu::Infrastructure
