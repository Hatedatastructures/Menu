#include "Infrastructure/SqliteRecipeRepositorySupport.hpp"

#include "Infrastructure/SqliteWorkflowSupport.hpp"

#include <boost/json/array.hpp>
#include <boost/json/serialize.hpp>

#include <array>
#include <string>
#include <string_view>
#include <utility>

namespace Menu::Infrastructure::SqliteRecipeRepositorySupport {
namespace {

using SqliteWorkflowSupport::BindDouble;
using SqliteWorkflowSupport::BindInteger;
using SqliteWorkflowSupport::BindText;
using SqliteWorkflowSupport::Prepare;
using SqliteWorkflowSupport::StatementGuard;
using SqliteWorkflowSupport::StorageError;

}  // namespace

std::string SerializeStrings(const std::vector<std::string>& Values) {
    boost::json::array Array;
    for (const std::string& Value : Values) {
        Array.emplace_back(Value);
    }
    return boost::json::serialize(Array);
}

Foundation::Result<void> StepDone(
    SqliteDatabase& Database,
    sqlite3_stmt* Statement) {
    if (sqlite3_step(Statement) != SQLITE_DONE) {
        return Foundation::Result<void>::FromError(StorageError(Database.NativeHandle()));
    }
    return Foundation::Result<void>();
}

Foundation::Result<void> InsertIngredientRows(
    SqliteDatabase& Database,
    const Domain::Recipe& RecipeValue) {
    for (const Domain::RecipeIngredient& IngredientValue : RecipeValue.Ingredients) {
        auto StatementResult = Prepare(
            Database,
            "INSERT INTO RecipeIngredients "
            "(RecipeId, IngredientId, Quantity, Unit, ServingFactor, Preparation, Required) "
            "VALUES (?, ?, ?, ?, ?, ?, ?);");
        if (!StatementResult.HasValue()) {
            return Foundation::Result<void>::FromError(StatementResult.ErrorValue());
        }
        StatementGuard Statement = std::move(StatementResult).Value();
        if (auto Result = BindText(Statement.Get(), 1, RecipeValue.Id); !Result.HasValue()) {
            return Result;
        }
        if (auto Result = BindText(Statement.Get(), 2, IngredientValue.IngredientId);
            !Result.HasValue()) {
            return Result;
        }
        if (auto Result = BindDouble(Statement.Get(), 3, IngredientValue.Quantity);
            !Result.HasValue()) {
            return Result;
        }
        if (auto Result = BindText(Statement.Get(), 4, IngredientValue.Unit); !Result.HasValue()) {
            return Result;
        }
        if (auto Result = BindDouble(Statement.Get(), 5, IngredientValue.ServingFactor);
            !Result.HasValue()) {
            return Result;
        }
        if (auto Result = BindText(Statement.Get(), 6, IngredientValue.Preparation);
            !Result.HasValue()) {
            return Result;
        }
        if (auto Result = BindInteger(
                Statement.Get(), 7, IngredientValue.Required ? 1 : 0);
            !Result.HasValue()) {
            return Result;
        }
        if (auto Result = StepDone(Database, Statement.Get()); !Result.HasValue()) {
            return Result;
        }
    }
    return Foundation::Result<void>();
}

Foundation::Result<void> InsertStepRows(
    SqliteDatabase& Database,
    const Domain::Recipe& RecipeValue) {
    for (const Domain::RecipeStep& StepValue : RecipeValue.Steps) {
        auto StatementResult = Prepare(
            Database,
            "INSERT INTO RecipeSteps "
            "(RecipeId, StepOrder, Title, Instruction, DurationSeconds, HasTimer) "
            "VALUES (?, ?, ?, ?, ?, ?);");
        if (!StatementResult.HasValue()) {
            return Foundation::Result<void>::FromError(StatementResult.ErrorValue());
        }
        StatementGuard Statement = std::move(StatementResult).Value();
        if (auto Result = BindText(Statement.Get(), 1, RecipeValue.Id); !Result.HasValue()) {
            return Result;
        }
        if (auto Result = BindInteger(Statement.Get(), 2, StepValue.StepOrder);
            !Result.HasValue()) {
            return Result;
        }
        if (auto Result = BindText(Statement.Get(), 3, StepValue.Title); !Result.HasValue()) {
            return Result;
        }
        if (auto Result = BindText(Statement.Get(), 4, StepValue.Instruction);
            !Result.HasValue()) {
            return Result;
        }
        if (auto Result = BindInteger(Statement.Get(), 5, StepValue.DurationSeconds);
            !Result.HasValue()) {
            return Result;
        }
        if (auto Result = BindInteger(Statement.Get(), 6, StepValue.HasTimer ? 1 : 0);
            !Result.HasValue()) {
            return Result;
        }
        if (auto Result = StepDone(Database, Statement.Get()); !Result.HasValue()) {
            return Result;
        }
    }
    return Foundation::Result<void>();
}

Foundation::Result<void> DeleteRecipeChildren(
    SqliteDatabase& Database,
    std::string_view RecipeId) {
    for (const char* Sql : {
             "DELETE FROM RecipeIngredients WHERE RecipeId = ?;",
             "DELETE FROM RecipeSteps WHERE RecipeId = ?;"}) {
        auto StatementResult = Prepare(Database, Sql);
        if (!StatementResult.HasValue()) {
            return Foundation::Result<void>::FromError(StatementResult.ErrorValue());
        }
        StatementGuard Statement = std::move(StatementResult).Value();
        if (auto Result = BindText(Statement.Get(), 1, RecipeId); !Result.HasValue()) {
            return Result;
        }
        if (auto Result = StepDone(Database, Statement.Get()); !Result.HasValue()) {
            return Result;
        }
    }
    return Foundation::Result<void>();
}

Foundation::Result<void> SaveRecipeBase(
    SqliteDatabase& Database,
    const Domain::Recipe& RecipeValue,
    bool Update) {
    const char* Sql = Update
                          ? "UPDATE Recipes SET Slug = ?, Name = ?, Cuisine = ?, Description = ?, "
                            "PrepMinutes = ?, CookMinutes = ?, Servings = ?, Difficulty = ?, "
                            "ImagePath = ?, Status = ?, AllergensJson = ?, CookwareJson = ?, "
                            "UpdatedAt = CURRENT_TIMESTAMP WHERE Id = ?;"
                          : "INSERT INTO Recipes "
                            "(Id, Slug, Name, Cuisine, Description, PrepMinutes, CookMinutes, "
                            "Servings, Difficulty, ImagePath, Status, AllergensJson, CookwareJson) "
                            "VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);";
    auto StatementResult = Prepare(Database, Sql);
    if (!StatementResult.HasValue()) {
        return Foundation::Result<void>::FromError(StatementResult.ErrorValue());
    }
    StatementGuard Statement = std::move(StatementResult).Value();
    int Index = 1;
    if (!Update) {
        if (auto Result = BindText(Statement.Get(), Index++, RecipeValue.Id);
            !Result.HasValue()) {
            return Result;
        }
    }
    if (auto Result = BindText(Statement.Get(), Index++, RecipeValue.Slug); !Result.HasValue()) {
        return Result;
    }
    if (auto Result = BindText(Statement.Get(), Index++, RecipeValue.Name); !Result.HasValue()) {
        return Result;
    }
    if (auto Result = BindText(Statement.Get(), Index++, RecipeValue.Cuisine);
        !Result.HasValue()) {
        return Result;
    }
    if (auto Result = BindText(Statement.Get(), Index++, RecipeValue.Description);
        !Result.HasValue()) {
        return Result;
    }
    if (auto Result = BindInteger(Statement.Get(), Index++, RecipeValue.PrepMinutes);
        !Result.HasValue()) {
        return Result;
    }
    if (auto Result = BindInteger(Statement.Get(), Index++, RecipeValue.CookMinutes);
        !Result.HasValue()) {
        return Result;
    }
    if (auto Result = BindInteger(Statement.Get(), Index++, RecipeValue.Servings);
        !Result.HasValue()) {
        return Result;
    }
    if (auto Result = BindInteger(Statement.Get(), Index++, RecipeValue.Difficulty);
        !Result.HasValue()) {
        return Result;
    }
    if (auto Result = BindText(Statement.Get(), Index++, RecipeValue.ImagePath);
        !Result.HasValue()) {
        return Result;
    }
    if (auto Result = BindText(Statement.Get(), Index++, RecipeValue.Status); !Result.HasValue()) {
        return Result;
    }
    if (auto Result = BindText(
            Statement.Get(), Index++, SerializeStrings(RecipeValue.Allergens));
        !Result.HasValue()) {
        return Result;
    }
    if (auto Result = BindText(
            Statement.Get(), Index++, SerializeStrings(RecipeValue.Cookware));
        !Result.HasValue()) {
        return Result;
    }
    if (Update) {
        if (auto Result = BindText(Statement.Get(), Index, RecipeValue.Id); !Result.HasValue()) {
            return Result;
        }
    }
    if (auto Result = StepDone(Database, Statement.Get()); !Result.HasValue()) {
        return Result;
    }
    if (Update && sqlite3_changes(Database.NativeHandle()) != 1) {
        return Foundation::Result<void>::FromError(
            Foundation::Error(Foundation::ErrorCode::NotFound, "菜谱不存在"));
    }
    return Foundation::Result<void>();
}

}  // namespace Menu::Infrastructure::SqliteRecipeRepositorySupport
