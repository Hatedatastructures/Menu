#include "Infrastructure/SqliteMealPlanRepository.hpp"

#include "Infrastructure/SqliteMealPlanRepositorySupport.hpp"
#include "Infrastructure/SqliteWorkflowSupport.hpp"
#include "Infrastructure/TokenService.hpp"

#include <string>
#include <string_view>
#include <utility>

namespace Menu::Infrastructure {
namespace {

using SqliteMealPlanRepositorySupport::FindRecipeName;
using SqliteMealPlanRepositorySupport::StepDone;
using SqliteWorkflowSupport::BindInteger;
using SqliteWorkflowSupport::BindText;
using SqliteWorkflowSupport::Prepare;
using SqliteWorkflowSupport::StatementGuard;

}  // namespace

Foundation::Result<Domain::MealPlan> SqliteMealPlanRepository::Save(
    std::string_view UserId,
    std::string_view PlanDate,
    std::vector<Domain::MealPlanItem> Items) {
    const auto BeginResult = Database.BeginTransaction();
    if (!BeginResult.HasValue()) {
        return Foundation::Result<Domain::MealPlan>::FromError(BeginResult.ErrorValue());
    }

    std::string PlanId;
    auto ExistingResult = Prepare(
        Database,
        "SELECT Id FROM MealPlans WHERE UserId = ? AND PlanDate = ? LIMIT 1;");
    if (!ExistingResult.HasValue()) {
        Database.Rollback();
        return Foundation::Result<Domain::MealPlan>::FromError(ExistingResult.ErrorValue());
    }
    StatementGuard ExistingStatement = std::move(ExistingResult).Value();
    if (auto Result = BindText(ExistingStatement.Get(), 1, UserId); !Result.HasValue()) {
        Database.Rollback();
        return Foundation::Result<Domain::MealPlan>::FromError(Result.ErrorValue());
    }
    if (auto Result = BindText(ExistingStatement.Get(), 2, PlanDate); !Result.HasValue()) {
        Database.Rollback();
        return Foundation::Result<Domain::MealPlan>::FromError(Result.ErrorValue());
    }
    const int ExistingStep = sqlite3_step(ExistingStatement.Get());
    if (ExistingStep == SQLITE_ROW) {
        PlanId = SqliteWorkflowSupport::ColumnText(ExistingStatement.Get(), 0);
        auto DeleteResult = Prepare(
            Database, "DELETE FROM MealPlanItems WHERE MealPlanId = ?;");
        if (!DeleteResult.HasValue()) {
            Database.Rollback();
            return Foundation::Result<Domain::MealPlan>::FromError(DeleteResult.ErrorValue());
        }
        StatementGuard DeleteStatement = std::move(DeleteResult).Value();
        if (auto Result = BindText(DeleteStatement.Get(), 1, PlanId); !Result.HasValue()) {
            Database.Rollback();
            return Foundation::Result<Domain::MealPlan>::FromError(Result.ErrorValue());
        }
        if (auto Result = StepDone(Database, DeleteStatement.Get()); !Result.HasValue()) {
            Database.Rollback();
            return Foundation::Result<Domain::MealPlan>::FromError(Result.ErrorValue());
        }
    } else if (ExistingStep == SQLITE_DONE) {
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
        if (auto Result = BindText(InsertStatement.Get(), 1, PlanId); !Result.HasValue()) {
            Database.Rollback();
            return Foundation::Result<Domain::MealPlan>::FromError(Result.ErrorValue());
        }
        if (auto Result = BindText(InsertStatement.Get(), 2, UserId); !Result.HasValue()) {
            Database.Rollback();
            return Foundation::Result<Domain::MealPlan>::FromError(Result.ErrorValue());
        }
        if (auto Result = BindText(InsertStatement.Get(), 3, PlanDate); !Result.HasValue()) {
            Database.Rollback();
            return Foundation::Result<Domain::MealPlan>::FromError(Result.ErrorValue());
        }
        if (auto Result = StepDone(Database, InsertStatement.Get()); !Result.HasValue()) {
            Database.Rollback();
            return Foundation::Result<Domain::MealPlan>::FromError(Result.ErrorValue());
        }
    } else {
        Database.Rollback();
        return Foundation::Result<Domain::MealPlan>::FromError(
            SqliteWorkflowSupport::StorageError(Database.NativeHandle()));
    }

    for (const Domain::MealPlanItem& Item : Items) {
        const auto RecipeResult =
            SqliteMealPlanRepositorySupport::FindRecipeName(Database, Item.RecipeId);
        if (!RecipeResult.HasValue()) {
            Database.Rollback();
            return Foundation::Result<Domain::MealPlan>::FromError(RecipeResult.ErrorValue());
        }
        const auto ItemIdResult = TokenService::CreateToken(16U);
        if (!ItemIdResult.HasValue()) {
            Database.Rollback();
            return Foundation::Result<Domain::MealPlan>::FromError(ItemIdResult.ErrorValue());
        }
        auto InsertItemResult = Prepare(
            Database,
            "INSERT INTO MealPlanItems "
            "(Id, MealPlanId, RecipeId, Servings, SortOrder) VALUES (?, ?, ?, ?, ?);");
        if (!InsertItemResult.HasValue()) {
            Database.Rollback();
            return Foundation::Result<Domain::MealPlan>::FromError(InsertItemResult.ErrorValue());
        }
        StatementGuard ItemStatement = std::move(InsertItemResult).Value();
        if (auto Result = BindText(
                ItemStatement.Get(), 1, "plan-item." + ItemIdResult.Value());
            !Result.HasValue()) {
            Database.Rollback();
            return Foundation::Result<Domain::MealPlan>::FromError(Result.ErrorValue());
        }
        if (auto Result = BindText(ItemStatement.Get(), 2, PlanId); !Result.HasValue()) {
            Database.Rollback();
            return Foundation::Result<Domain::MealPlan>::FromError(Result.ErrorValue());
        }
        if (auto Result = BindText(ItemStatement.Get(), 3, Item.RecipeId); !Result.HasValue()) {
            Database.Rollback();
            return Foundation::Result<Domain::MealPlan>::FromError(Result.ErrorValue());
        }
        if (auto Result = BindInteger(ItemStatement.Get(), 4, Item.Servings);
            !Result.HasValue()) {
            Database.Rollback();
            return Foundation::Result<Domain::MealPlan>::FromError(Result.ErrorValue());
        }
        if (auto Result = BindInteger(ItemStatement.Get(), 5, Item.SortOrder);
            !Result.HasValue()) {
            Database.Rollback();
            return Foundation::Result<Domain::MealPlan>::FromError(Result.ErrorValue());
        }
        if (auto Result = StepDone(Database, ItemStatement.Get()); !Result.HasValue()) {
            Database.Rollback();
            return Foundation::Result<Domain::MealPlan>::FromError(Result.ErrorValue());
        }
    }

    const auto CommitResult = Database.Commit();
    if (!CommitResult.HasValue()) {
        Database.Rollback();
        return Foundation::Result<Domain::MealPlan>::FromError(CommitResult.ErrorValue());
    }
    return ReadPlan(PlanId);
}

}  // namespace Menu::Infrastructure
