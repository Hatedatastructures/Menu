#include "Infrastructure/SeedDataRows.hpp"

#include "Infrastructure/SqliteWorkflowSupport.hpp"

#include <boost/json/array.hpp>
#include <boost/json/serialize.hpp>

#include <array>
#include <string>
#include <string_view>
#include <utility>

namespace Menu::Infrastructure::SeedDataRows {
namespace {

using SqliteWorkflowSupport::BindDouble;
using SqliteWorkflowSupport::BindInt;
using SqliteWorkflowSupport::BindText;
using SqliteWorkflowSupport::Prepare;
using SqliteWorkflowSupport::StatementGuard;
using SqliteWorkflowSupport::StorageError;

Foundation::Result<void> StepDone(
    SqliteDatabase& Database,
    sqlite3_stmt* Statement) {
    if (sqlite3_step(Statement) != SQLITE_DONE) {
        return Foundation::Result<void>::FromError(StorageError(Database.NativeHandle()));
    }
    return Foundation::Result<void>();
}

std::string SerializeStrings(const std::vector<std::string>& Values) {
    boost::json::array Array;
    for (const std::string& Value : Values) {
        Array.emplace_back(Value);
    }
    return boost::json::serialize(Array);
}

Foundation::Result<void> InsertRecipeIngredients(
    SqliteDatabase& Database,
    const Domain::Recipe& RecipeValue) {
    for (const Domain::RecipeIngredient& IngredientValue : RecipeValue.Ingredients) {
        auto StatementResult = Prepare(
            Database,
            "INSERT OR IGNORE INTO RecipeIngredients "
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
        if (auto Result = BindInt(Statement.Get(), 7, IngredientValue.Required ? 1 : 0);
            !Result.HasValue()) {
            return Result;
        }
        if (auto Result = StepDone(Database, Statement.Get()); !Result.HasValue()) {
            return Result;
        }
    }
    return Foundation::Result<void>();
}

Foundation::Result<void> InsertRecipeSteps(
    SqliteDatabase& Database,
    const Domain::Recipe& RecipeValue) {
    for (const Domain::RecipeStep& StepValue : RecipeValue.Steps) {
        auto StatementResult = Prepare(
            Database,
            "INSERT OR IGNORE INTO RecipeSteps "
            "(RecipeId, StepOrder, Title, Instruction, DurationSeconds, HasTimer) "
            "VALUES (?, ?, ?, ?, ?, ?);");
        if (!StatementResult.HasValue()) {
            return Foundation::Result<void>::FromError(StatementResult.ErrorValue());
        }
        StatementGuard Statement = std::move(StatementResult).Value();
        if (auto Result = BindText(Statement.Get(), 1, RecipeValue.Id); !Result.HasValue()) {
            return Result;
        }
        if (auto Result = BindInt(Statement.Get(), 2, StepValue.StepOrder); !Result.HasValue()) {
            return Result;
        }
        if (auto Result = BindText(Statement.Get(), 3, StepValue.Title); !Result.HasValue()) {
            return Result;
        }
        if (auto Result = BindText(Statement.Get(), 4, StepValue.Instruction);
            !Result.HasValue()) {
            return Result;
        }
        if (auto Result = BindInt(Statement.Get(), 5, StepValue.DurationSeconds);
            !Result.HasValue()) {
            return Result;
        }
        if (auto Result = BindInt(Statement.Get(), 6, StepValue.HasTimer ? 1 : 0);
            !Result.HasValue()) {
            return Result;
        }
        if (auto Result = StepDone(Database, Statement.Get()); !Result.HasValue()) {
            return Result;
        }
    }
    return Foundation::Result<void>();
}

}  // namespace

Foundation::Result<void> InsertRecipe(
    SqliteDatabase& Database,
    const Domain::Recipe& RecipeValue) {
    auto StatementResult = Prepare(
        Database,
        "INSERT OR IGNORE INTO Recipes "
        "(Id, Slug, Name, Cuisine, Description, PrepMinutes, CookMinutes, Servings, Difficulty, "
        "ImagePath, Status, AllergensJson, CookwareJson) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);");
    if (!StatementResult.HasValue()) {
        return Foundation::Result<void>::FromError(StatementResult.ErrorValue());
    }
    StatementGuard Statement = std::move(StatementResult).Value();
    const std::array<std::string_view, 5> TextValues = {
        RecipeValue.Id, RecipeValue.Slug, RecipeValue.Name,
        RecipeValue.Cuisine, RecipeValue.Description};
    for (std::size_t Index = 0; Index < TextValues.size(); ++Index) {
        if (auto Result = BindText(
                Statement.Get(), static_cast<int>(Index + 1), TextValues[Index]);
            !Result.HasValue()) {
            return Result;
        }
    }
    const std::array<int, 4> IntValues = {
        RecipeValue.PrepMinutes, RecipeValue.CookMinutes,
        RecipeValue.Servings, RecipeValue.Difficulty};
    for (std::size_t Index = 0; Index < IntValues.size(); ++Index) {
        if (auto Result = BindInt(
                Statement.Get(), static_cast<int>(Index + 6), IntValues[Index]);
            !Result.HasValue()) {
            return Result;
        }
    }
    if (auto Result = BindText(Statement.Get(), 10, RecipeValue.ImagePath); !Result.HasValue()) {
        return Result;
    }
    if (auto Result = BindText(Statement.Get(), 11, RecipeValue.Status); !Result.HasValue()) {
        return Result;
    }
    if (auto Result = BindText(Statement.Get(), 12, SerializeStrings(RecipeValue.Allergens));
        !Result.HasValue()) {
        return Result;
    }
    if (auto Result = BindText(Statement.Get(), 13, SerializeStrings(RecipeValue.Cookware));
        !Result.HasValue()) {
        return Result;
    }
    if (auto Result = StepDone(Database, Statement.Get()); !Result.HasValue()) {
        return Result;
    }
    if (auto Result = InsertRecipeIngredients(Database, RecipeValue); !Result.HasValue()) {
        return Result;
    }
    return InsertRecipeSteps(Database, RecipeValue);
}

}  // namespace Menu::Infrastructure::SeedDataRows
