#include "Infrastructure/SqliteMealPlanRepository.hpp"

#include "Infrastructure/SqliteWorkflowSupport.hpp"
#include "Infrastructure/TokenService.hpp"

#include <optional>
#include <utility>

namespace Menu::Infrastructure {
namespace {

using SqliteWorkflowSupport::BindText;
using SqliteWorkflowSupport::ColumnText;
using SqliteWorkflowSupport::Prepare;
using SqliteWorkflowSupport::StatementGuard;
using SqliteWorkflowSupport::StorageError;

Foundation::Error NotFound(std::string Message) {
    return Foundation::Error(Foundation::ErrorCode::NotFound, std::move(Message));
}

Foundation::Result<void> StepDone(SqliteDatabase& Database, sqlite3_stmt* Statement) {
    if (sqlite3_step(Statement) != SQLITE_DONE) {
        return Foundation::Result<void>::FromError(StorageError(Database.NativeHandle()));
    }
    return Foundation::Result<void>();
}

}  // namespace

Foundation::Result<std::string> SqliteMealPlanRepository::FindRecipeName(
    std::string_view RecipeId) {
    auto StatementResult = Prepare(
        Database, "SELECT Name FROM Recipes WHERE Id = ? LIMIT 1;");
    if (!StatementResult.HasValue()) {
        return Foundation::Result<std::string>::FromError(StatementResult.ErrorValue());
    }
    StatementGuard Statement = std::move(StatementResult).Value();
    const auto BindResult = BindText(Statement.Get(), 1, RecipeId);
    if (!BindResult.HasValue()) {
        return Foundation::Result<std::string>::FromError(BindResult.ErrorValue());
    }
    const int StepResult = sqlite3_step(Statement.Get());
    if (StepResult == SQLITE_DONE) {
        return Foundation::Result<std::string>::FromError(
            NotFound("菜谱不存在"));
    }
    if (StepResult != SQLITE_ROW) {
        return Foundation::Result<std::string>::FromError(StorageError(Database.NativeHandle()));
    }
    return ColumnText(Statement.Get(), 0);
}

Foundation::Result<Domain::MealPlan> SqliteMealPlanRepository::Save(
    std::string_view UserId,
    std::string_view PlanDate,
    std::vector<Domain::MealPlanItem> Items) {
    const auto BeginResult = Database.BeginTransaction();
    if (!BeginResult.HasValue()) {
        return Foundation::Result<Domain::MealPlan>::FromError(BeginResult.ErrorValue());
    }

    std::optional<std::string> ExistingPlanId;
    auto ExistingResult = Prepare(
        Database,
        "SELECT Id FROM MealPlans WHERE UserId = ? AND PlanDate = ? LIMIT 1;");
    if (!ExistingResult.HasValue()) {
        Database.Rollback();
        return Foundation::Result<Domain::MealPlan>::FromError(ExistingResult.ErrorValue());
    }
    StatementGuard ExistingStatement = std::move(ExistingResult).Value();
    const auto UserBind = BindText(ExistingStatement.Get(), 1, UserId);
    const auto DateBind = BindText(ExistingStatement.Get(), 2, PlanDate);
    if (!UserBind.HasValue() || !DateBind.HasValue()) {
        Database.Rollback();
        return Foundation::Result<Domain::MealPlan>::FromError(
            !UserBind.HasValue() ? UserBind.ErrorValue() : DateBind.ErrorValue());
    }
    const int ExistingStep = sqlite3_step(ExistingStatement.Get());
    if (ExistingStep == SQLITE_ROW) {
        ExistingPlanId = ColumnText(ExistingStatement.Get(), 0);
    } else if (ExistingStep != SQLITE_DONE) {
        Database.Rollback();
        return Foundation::Result<Domain::MealPlan>::FromError(
            StorageError(Database.NativeHandle()));
    }

    std::string PlanId;
    if (ExistingPlanId.has_value()) {
        PlanId = ExistingPlanId.value();
        auto DeleteResult = Prepare(
            Database, "DELETE FROM MealPlanItems WHERE MealPlanId = ?;");
        if (!DeleteResult.HasValue()) {
            Database.Rollback();
            return Foundation::Result<Domain::MealPlan>::FromError(DeleteResult.ErrorValue());
        }
        StatementGuard DeleteStatement = std::move(DeleteResult).Value();
        const auto DeleteBind = BindText(DeleteStatement.Get(), 1, PlanId);
        if (!DeleteBind.HasValue()) {
            Database.Rollback();
            return Foundation::Result<Domain::MealPlan>::FromError(DeleteBind.ErrorValue());
        }
        const auto DeleteStep = StepDone(Database, DeleteStatement.Get());
        if (!DeleteStep.HasValue()) {
            Database.Rollback();
            return Foundation::Result<Domain::MealPlan>::FromError(DeleteStep.ErrorValue());
        }
    } else {
        const auto NewPlanId = TokenService::CreateToken(16U);
        if (!NewPlanId.HasValue()) {
            Database.Rollback();
            return Foundation::Result<Domain::MealPlan>::FromError(NewPlanId.ErrorValue());
        }
        PlanId = "plan." + NewPlanId.Value();
        auto InsertResult = Prepare(
            Database,
            "INSERT INTO MealPlans (Id, UserId, PlanDate) VALUES (?, ?, ?);");
        if (!InsertResult.HasValue()) {
            Database.Rollback();
            return Foundation::Result<Domain::MealPlan>::FromError(InsertResult.ErrorValue());
        }
        StatementGuard InsertStatement = std::move(InsertResult).Value();
        const auto IdBind = BindText(InsertStatement.Get(), 1, PlanId);
        const auto InsertUserBind = BindText(InsertStatement.Get(), 2, UserId);
        const auto InsertDateBind = BindText(InsertStatement.Get(), 3, PlanDate);
        if (!IdBind.HasValue() || !InsertUserBind.HasValue() || !InsertDateBind.HasValue()) {
            Database.Rollback();
            const Foundation::Error ErrorValue = !IdBind.HasValue()
                ? IdBind.ErrorValue()
                : !InsertUserBind.HasValue() ? InsertUserBind.ErrorValue()
                                             : InsertDateBind.ErrorValue();
            return Foundation::Result<Domain::MealPlan>::FromError(ErrorValue);
        }
        const auto InsertStep = StepDone(Database, InsertStatement.Get());
        if (!InsertStep.HasValue()) {
            Database.Rollback();
            return Foundation::Result<Domain::MealPlan>::FromError(InsertStep.ErrorValue());
        }
    }

    auto InsertItemResult = Prepare(
        Database,
        "INSERT INTO MealPlanItems "
        "(Id, MealPlanId, RecipeId, Servings, SortOrder) VALUES (?, ?, ?, ?, ?);");
    if (!InsertItemResult.HasValue()) {
        Database.Rollback();
        return Foundation::Result<Domain::MealPlan>::FromError(InsertItemResult.ErrorValue());
    }
    for (const Domain::MealPlanItem& Item : Items) {
        const auto RecipeResult = FindRecipeName(Item.RecipeId);
        if (!RecipeResult.HasValue()) {
            Database.Rollback();
            return Foundation::Result<Domain::MealPlan>::FromError(RecipeResult.ErrorValue());
        }
        const auto ItemIdResult = TokenService::CreateToken(16U);
        if (!ItemIdResult.HasValue()) {
            Database.Rollback();
            return Foundation::Result<Domain::MealPlan>::FromError(ItemIdResult.ErrorValue());
        }
        StatementGuard ItemStatement = std::move(InsertItemResult).Value();
        const auto IdBind = BindText(ItemStatement.Get(), 1, "plan-item." + ItemIdResult.Value());
        const auto PlanBind = BindText(ItemStatement.Get(), 2, PlanId);
        const auto RecipeBind = BindText(ItemStatement.Get(), 3, Item.RecipeId);
        const auto ServingsBind = SqliteWorkflowSupport::BindInteger(
            ItemStatement.Get(), 4, Item.Servings);
        const auto SortBind = SqliteWorkflowSupport::BindInteger(
            ItemStatement.Get(), 5, Item.SortOrder);
        if (!IdBind.HasValue() || !PlanBind.HasValue() || !RecipeBind.HasValue() ||
            !ServingsBind.HasValue() || !SortBind.HasValue()) {
            Database.Rollback();
            const Foundation::Error ErrorValue = !IdBind.HasValue()
                ? IdBind.ErrorValue()
                : !PlanBind.HasValue() ? PlanBind.ErrorValue()
                : !RecipeBind.HasValue() ? RecipeBind.ErrorValue()
                : !ServingsBind.HasValue() ? ServingsBind.ErrorValue()
                                           : SortBind.ErrorValue();
            return Foundation::Result<Domain::MealPlan>::FromError(ErrorValue);
        }
        const auto ItemStep = StepDone(Database, ItemStatement.Get());
        if (!ItemStep.HasValue()) {
            Database.Rollback();
            return Foundation::Result<Domain::MealPlan>::FromError(ItemStep.ErrorValue());
        }
        InsertItemResult = Prepare(
            Database,
            "INSERT INTO MealPlanItems "
            "(Id, MealPlanId, RecipeId, Servings, SortOrder) VALUES (?, ?, ?, ?, ?);");
        if (!InsertItemResult.HasValue()) {
            Database.Rollback();
            return Foundation::Result<Domain::MealPlan>::FromError(InsertItemResult.ErrorValue());
        }
    }

    const auto CommitResult = Database.Commit();
    if (!CommitResult.HasValue()) {
        Database.Rollback();
        return Foundation::Result<Domain::MealPlan>::FromError(CommitResult.ErrorValue());
    }
    return ReadPlan(PlanId);
}

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
    const auto UserBind = BindText(Statement.Get(), 1, UserId);
    const auto FromBind = BindText(Statement.Get(), 2, FromDate);
    const auto ToBind = BindText(Statement.Get(), 3, ToDate);
    if (!UserBind.HasValue() || !FromBind.HasValue() || !ToBind.HasValue()) {
        return Foundation::Result<std::vector<Domain::MealPlan>>::FromError(
            !UserBind.HasValue() ? UserBind.ErrorValue()
            : !FromBind.HasValue() ? FromBind.ErrorValue() : ToBind.ErrorValue());
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
    const auto PlanBind = BindText(PlanStatement.Get(), 1, PlanId);
    if (!PlanBind.HasValue()) {
        return Foundation::Result<Domain::MealPlan>::FromError(PlanBind.ErrorValue());
    }
    const int PlanStep = sqlite3_step(PlanStatement.Get());
    if (PlanStep == SQLITE_DONE) {
        return Foundation::Result<Domain::MealPlan>::FromError(
            NotFound("计划不存在"));
    }
    if (PlanStep != SQLITE_ROW) {
        return Foundation::Result<Domain::MealPlan>::FromError(StorageError(Database.NativeHandle()));
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
    const auto ItemBind = BindText(ItemStatement.Get(), 1, PlanId);
    if (!ItemBind.HasValue()) {
        return Foundation::Result<Domain::MealPlan>::FromError(ItemBind.ErrorValue());
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
        const auto Ingredients = SqliteWorkflowSupport::ReadRecipeIngredients(
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
