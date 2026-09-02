#include "Infrastructure/SqliteRecipeRepository.hpp"

#include "Infrastructure/SqliteRecipeRepositorySupport.hpp"
#include "Infrastructure/SqliteWorkflowSupport.hpp"

#include <string_view>
#include <utility>

namespace Menu::Infrastructure {
namespace {

using SqliteRecipeRepositorySupport::ReadPublishedRecipe;
using SqliteRecipeRepositorySupport::ReadRecipeRow;
using SqliteWorkflowSupport::BindText;
using SqliteWorkflowSupport::ColumnText;
using SqliteWorkflowSupport::Prepare;
using SqliteWorkflowSupport::StatementGuard;
using SqliteWorkflowSupport::StorageError;

constexpr std::string_view RecipeColumns =
    "Id, Slug, Name, Cuisine, Description, PrepMinutes, CookMinutes, Servings, "
    "Difficulty, ImagePath, Status, AllergensJson, CookwareJson";

Foundation::Result<std::vector<Domain::Recipe>> ReadRecipes(
    SqliteDatabase& Database,
    std::string_view Query) {
    auto StatementResult = Prepare(Database, Query);
    if (!StatementResult.HasValue()) {
        return Foundation::Result<std::vector<Domain::Recipe>>::FromError(
            StatementResult.ErrorValue());
    }
    StatementGuard Statement = std::move(StatementResult).Value();
    std::vector<Domain::Recipe> Recipes;
    while (true) {
        const int StepResult = sqlite3_step(Statement.Get());
        if (StepResult == SQLITE_DONE) {
            break;
        }
        if (StepResult != SQLITE_ROW) {
            return Foundation::Result<std::vector<Domain::Recipe>>::FromError(
                StorageError(Database.NativeHandle()));
        }
        auto RecipeResult = ReadRecipeRow(Database, Statement.Get());
        if (!RecipeResult.HasValue()) {
            return Foundation::Result<std::vector<Domain::Recipe>>::FromError(
                RecipeResult.ErrorValue());
        }
        Recipes.push_back(std::move(RecipeResult).Value());
    }
    return Recipes;
}

}  // namespace

Foundation::Result<std::vector<Domain::Recipe>> SqliteRecipeRepository::ListPublished() {
    return ReadRecipes(
        Database,
        std::string("SELECT ") + std::string(RecipeColumns) +
            " FROM Recipes WHERE Status = 'published' ORDER BY Id;");
}

Foundation::Result<std::optional<Domain::Recipe>>
SqliteRecipeRepository::FindPublishedById(std::string_view Id) {
    auto RecipeResult = ReadPublishedRecipe(Database, Id);
    if (!RecipeResult.HasValue()) {
        if (RecipeResult.ErrorValue().CodeValue() == Foundation::ErrorCode::NotFound) {
            return std::optional<Domain::Recipe>();
        }
        return Foundation::Result<std::optional<Domain::Recipe>>::FromError(
            RecipeResult.ErrorValue());
    }
    return std::optional<Domain::Recipe>(std::move(RecipeResult).Value());
}

Foundation::Result<std::vector<Domain::Recipe>> SqliteRecipeRepository::ListAll() {
    return ReadRecipes(
        Database,
        std::string("SELECT ") + std::string(RecipeColumns) +
            " FROM Recipes ORDER BY UpdatedAt DESC, Id;");
}

}  // namespace Menu::Infrastructure
