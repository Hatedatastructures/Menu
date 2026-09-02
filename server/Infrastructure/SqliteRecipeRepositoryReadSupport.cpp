#include "Infrastructure/SqliteRecipeRepositorySupport.hpp"

#include "Infrastructure/SqliteWorkflowSupport.hpp"

#include <boost/json/parse.hpp>

#include <string_view>
#include <utility>

namespace Menu::Infrastructure::SqliteRecipeRepositorySupport {
namespace {

using SqliteWorkflowSupport::BindText;
using SqliteWorkflowSupport::ColumnText;
using SqliteWorkflowSupport::Prepare;
using SqliteWorkflowSupport::StatementGuard;
using SqliteWorkflowSupport::StorageError;

}  // namespace

Foundation::Result<std::vector<std::string>> ParseStringArray(
    sqlite3* Handle,
    std::string_view JsonText) {
    boost::system::error_code Error;
    const boost::json::value Parsed = boost::json::parse(JsonText, Error);
    if (Error || !Parsed.is_array()) {
        return Foundation::Result<std::vector<std::string>>::FromError(StorageError(Handle));
    }
    std::vector<std::string> Values;
    for (const boost::json::value& Value : Parsed.as_array()) {
        if (!Value.is_string()) {
            return Foundation::Result<std::vector<std::string>>::FromError(
                StorageError(Handle));
        }
        Values.emplace_back(Value.as_string().c_str());
    }
    return Values;
}

Foundation::Result<Domain::Recipe> ReadRecipeRow(
    SqliteDatabase& Database,
    sqlite3_stmt* Statement) {
    Domain::Recipe RecipeValue;
    RecipeValue.Id = ColumnText(Statement, 0);
    RecipeValue.Slug = ColumnText(Statement, 1);
    RecipeValue.Name = ColumnText(Statement, 2);
    RecipeValue.Cuisine = ColumnText(Statement, 3);
    RecipeValue.Description = ColumnText(Statement, 4);
    RecipeValue.PrepMinutes = sqlite3_column_int(Statement, 5);
    RecipeValue.CookMinutes = sqlite3_column_int(Statement, 6);
    RecipeValue.Servings = sqlite3_column_int(Statement, 7);
    RecipeValue.Difficulty = sqlite3_column_int(Statement, 8);
    RecipeValue.ImagePath = ColumnText(Statement, 9);
    RecipeValue.Status = ColumnText(Statement, 10);

    const auto Allergens = ParseStringArray(Database.NativeHandle(), ColumnText(Statement, 11));
    if (!Allergens.HasValue()) {
        return Foundation::Result<Domain::Recipe>::FromError(Allergens.ErrorValue());
    }
    RecipeValue.Allergens = Allergens.Value();
    const auto Cookware = ParseStringArray(Database.NativeHandle(), ColumnText(Statement, 12));
    if (!Cookware.HasValue()) {
        return Foundation::Result<Domain::Recipe>::FromError(Cookware.ErrorValue());
    }
    RecipeValue.Cookware = Cookware.Value();

    auto IngredientStatementResult = Prepare(
        Database,
        "SELECT RecipeIngredients.IngredientId, RecipeIngredients.Quantity, "
        "RecipeIngredients.Unit, RecipeIngredients.Required, RecipeIngredients.ServingFactor, "
        "RecipeIngredients.Preparation, Ingredients.Name, Ingredients.Category, "
        "Ingredients.DefaultUnit, Ingredients.IsPantryStaple "
        "FROM RecipeIngredients JOIN Ingredients ON Ingredients.Id = RecipeIngredients.IngredientId "
        "WHERE RecipeIngredients.RecipeId = ? ORDER BY RecipeIngredients.IngredientId;");
    if (!IngredientStatementResult.HasValue()) {
        return Foundation::Result<Domain::Recipe>::FromError(
            IngredientStatementResult.ErrorValue());
    }
    StatementGuard IngredientStatement = std::move(IngredientStatementResult).Value();
    if (auto Result = BindText(IngredientStatement.Get(), 1, RecipeValue.Id);
        !Result.HasValue()) {
        return Foundation::Result<Domain::Recipe>::FromError(Result.ErrorValue());
    }
    while (true) {
        const int StepResult = sqlite3_step(IngredientStatement.Get());
        if (StepResult == SQLITE_DONE) {
            break;
        }
        if (StepResult != SQLITE_ROW) {
            return Foundation::Result<Domain::Recipe>::FromError(
                StorageError(Database.NativeHandle()));
        }
        RecipeValue.Ingredients.emplace_back(
            ColumnText(IngredientStatement.Get(), 0),
            sqlite3_column_double(IngredientStatement.Get(), 1),
            ColumnText(IngredientStatement.Get(), 2),
            sqlite3_column_int(IngredientStatement.Get(), 3) != 0,
            sqlite3_column_double(IngredientStatement.Get(), 4),
            ColumnText(IngredientStatement.Get(), 5));
        Domain::RecipeIngredient& IngredientValue = RecipeValue.Ingredients.back();
        IngredientValue.IngredientName = ColumnText(IngredientStatement.Get(), 6);
        IngredientValue.IngredientCategory = ColumnText(IngredientStatement.Get(), 7);
        IngredientValue.IngredientDefaultUnit = ColumnText(IngredientStatement.Get(), 8);
        IngredientValue.IngredientIsPantryStaple =
            sqlite3_column_int(IngredientStatement.Get(), 9) != 0;
    }

    auto StepStatementResult = Prepare(
        Database,
        "SELECT StepOrder, Title, Instruction, DurationSeconds, HasTimer "
        "FROM RecipeSteps WHERE RecipeId = ? ORDER BY StepOrder;");
    if (!StepStatementResult.HasValue()) {
        return Foundation::Result<Domain::Recipe>::FromError(StepStatementResult.ErrorValue());
    }
    StatementGuard StepStatement = std::move(StepStatementResult).Value();
    if (auto Result = BindText(StepStatement.Get(), 1, RecipeValue.Id); !Result.HasValue()) {
        return Foundation::Result<Domain::Recipe>::FromError(Result.ErrorValue());
    }
    while (true) {
        const int StepResult = sqlite3_step(StepStatement.Get());
        if (StepResult == SQLITE_DONE) {
            break;
        }
        if (StepResult != SQLITE_ROW) {
            return Foundation::Result<Domain::Recipe>::FromError(
                StorageError(Database.NativeHandle()));
        }
        RecipeValue.Steps.push_back(Domain::RecipeStep{
            sqlite3_column_int(StepStatement.Get(), 0),
            ColumnText(StepStatement.Get(), 1),
            ColumnText(StepStatement.Get(), 2),
            sqlite3_column_int(StepStatement.Get(), 3),
            sqlite3_column_int(StepStatement.Get(), 4) != 0});
    }
    return RecipeValue;
}

Foundation::Result<Domain::Recipe> ReadPublishedRecipe(
    SqliteDatabase& Database,
    std::string_view Id) {
    auto StatementResult = Prepare(
        Database,
        "SELECT Id, Slug, Name, Cuisine, Description, PrepMinutes, CookMinutes, Servings, "
        "Difficulty, ImagePath, Status, AllergensJson, CookwareJson FROM Recipes "
        "WHERE Id = ? AND Status = 'published';");
    if (!StatementResult.HasValue()) {
        return Foundation::Result<Domain::Recipe>::FromError(StatementResult.ErrorValue());
    }
    StatementGuard Statement = std::move(StatementResult).Value();
    if (auto Result = BindText(Statement.Get(), 1, Id); !Result.HasValue()) {
        return Foundation::Result<Domain::Recipe>::FromError(Result.ErrorValue());
    }
    const int StepResult = sqlite3_step(Statement.Get());
    if (StepResult == SQLITE_DONE) {
        return Foundation::Result<Domain::Recipe>::FromError(
            Foundation::Error(Foundation::ErrorCode::NotFound, "菜谱不存在"));
    }
    if (StepResult != SQLITE_ROW) {
        return Foundation::Result<Domain::Recipe>::FromError(
            StorageError(Database.NativeHandle()));
    }
    return ReadRecipeRow(Database, Statement.Get());
}

}  // namespace Menu::Infrastructure::SqliteRecipeRepositorySupport
