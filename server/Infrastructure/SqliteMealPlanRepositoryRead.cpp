#include "Infrastructure/SqliteMealPlanRepository.hpp"

#include "Infrastructure/SqliteWorkflowSupport.hpp"

#include <string>
#include <string_view>
#include <utility>

namespace Menu::Infrastructure {
namespace {

using SqliteWorkflowSupport::BindText;
using SqliteWorkflowSupport::ColumnText;
using SqliteWorkflowSupport::Prepare;
using SqliteWorkflowSupport::ReadRecipeIngredients;
using SqliteWorkflowSupport::StatementGuard;
using SqliteWorkflowSupport::StorageError;

Foundation::Error PlanNotFound() {
    return Foundation::Error(Foundation::ErrorCode::NotFound, "计划不存在");
}

}  // namespace

Foundation::Result<std::vector<Domain::MealPlan>> SqliteMealPlanRepository::List(
    std::string_view UserId,
    std::string_view FromDate,
    std::string_view ToDate) {
    auto StatementResult = Prepare(
        Database,
        "SELECT Id FROM MealPlans WHERE UserId = ? AND PlanDate >= ? AND PlanDate <= "
        "? ORDER BY PlanDate, Id;");
    if (!StatementResult.HasValue()) {
        return Foundation::Result<std::vector<Domain::MealPlan>>::FromError(
            StatementResult.ErrorValue());
    }
    StatementGuard Statement = std::move(StatementResult).Value();
    if (auto Result = BindText(Statement.Get(), 1, UserId); !Result.HasValue()) {
        return Foundation::Result<std::vector<Domain::MealPlan>>::FromError(Result.ErrorValue());
    }
    if (auto Result = BindText(Statement.Get(), 2, FromDate); !Result.HasValue()) {
        return Foundation::Result<std::vector<Domain::MealPlan>>::FromError(Result.ErrorValue());
    }
    if (auto Result = BindText(Statement.Get(), 3, ToDate); !Result.HasValue()) {
        return Foundation::Result<std::vector<Domain::MealPlan>>::FromError(Result.ErrorValue());
    }
    std::vector<Domain::MealPlan> Plans;
    while (true) {
        const int StepResult = sqlite3_step(Statement.Get());
        if (StepResult == SQLITE_DONE) {
            break;
        }
        if (StepResult != SQLITE_ROW) {
            return Foundation::Result<std::vector<Domain::MealPlan>>::FromError(
                StorageError(Database.NativeHandle()));
        }
        const auto PlanResult = ReadPlan(ColumnText(Statement.Get(), 0));
        if (!PlanResult.HasValue()) {
            return Foundation::Result<std::vector<Domain::MealPlan>>::FromError(
                PlanResult.ErrorValue());
        }
        Plans.push_back(PlanResult.Value());
    }
    return Plans;
}

Foundation::Result<Domain::MealPlan> SqliteMealPlanRepository::ReadPlan(
    std::string_view PlanId) {
    auto PlanResult = Prepare(
        Database, "SELECT UserId, PlanDate FROM MealPlans WHERE Id = ? LIMIT 1;");
    if (!PlanResult.HasValue()) {
        return Foundation::Result<Domain::MealPlan>::FromError(PlanResult.ErrorValue());
    }
    StatementGuard PlanStatement = std::move(PlanResult).Value();
    if (auto Result = BindText(PlanStatement.Get(), 1, PlanId); !Result.HasValue()) {
        return Foundation::Result<Domain::MealPlan>::FromError(Result.ErrorValue());
    }
    const int PlanStep = sqlite3_step(PlanStatement.Get());
    if (PlanStep == SQLITE_DONE) {
        return Foundation::Result<Domain::MealPlan>::FromError(PlanNotFound());
    }
    if (PlanStep != SQLITE_ROW) {
        return Foundation::Result<Domain::MealPlan>::FromError(
            StorageError(Database.NativeHandle()));
    }
    Domain::MealPlan Plan;
    Plan.Id = std::string(PlanId);
    Plan.UserId = ColumnText(PlanStatement.Get(), 0);
    Plan.PlanDate = ColumnText(PlanStatement.Get(), 1);

    auto ItemResult = Prepare(
        Database,
        "SELECT mpi.Id, mpi.RecipeId, r.Name, r.ImagePath, mpi.Servings, "
        "mpi.SortOrder, r.Servings FROM MealPlanItems mpi "
        "JOIN Recipes r ON r.Id = mpi.RecipeId WHERE mpi.MealPlanId = ? "
        "ORDER BY mpi.SortOrder, mpi.Id;");
    if (!ItemResult.HasValue()) {
        return Foundation::Result<Domain::MealPlan>::FromError(ItemResult.ErrorValue());
    }
    StatementGuard ItemStatement = std::move(ItemResult).Value();
    if (auto Result = BindText(ItemStatement.Get(), 1, PlanId); !Result.HasValue()) {
        return Foundation::Result<Domain::MealPlan>::FromError(Result.ErrorValue());
    }
    while (true) {
        const int StepResult = sqlite3_step(ItemStatement.Get());
        if (StepResult == SQLITE_DONE) {
            break;
        }
        if (StepResult != SQLITE_ROW) {
            return Foundation::Result<Domain::MealPlan>::FromError(
                StorageError(Database.NativeHandle()));
        }
        Domain::MealPlanItem Item;
        Item.Id = ColumnText(ItemStatement.Get(), 0);
        Item.RecipeId = ColumnText(ItemStatement.Get(), 1);
        Item.RecipeName = ColumnText(ItemStatement.Get(), 2);
        Item.ImagePath = ColumnText(ItemStatement.Get(), 3);
        Item.Servings = sqlite3_column_int(ItemStatement.Get(), 4);
        Item.SortOrder = sqlite3_column_int(ItemStatement.Get(), 5);
        const int RecipeServings = sqlite3_column_int(ItemStatement.Get(), 6);
        const auto Ingredients = ReadRecipeIngredients(
            Database, Item.RecipeId, RecipeServings, Item.Servings);
        if (!Ingredients.HasValue()) {
            return Foundation::Result<Domain::MealPlan>::FromError(Ingredients.ErrorValue());
        }
        Item.Ingredients = Ingredients.Value();
        Plan.Items.push_back(std::move(Item));
    }
    return Plan;
}

}  // namespace Menu::Infrastructure
