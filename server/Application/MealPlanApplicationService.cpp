#include "Application/MealPlanApplicationService.hpp"

#include <Domain/PlanRules.hpp>

#include <cctype>
#include <utility>

namespace Menu::Application {
namespace {

Foundation::Error InvalidPlan(std::string Message) {
    return Foundation::Error(Foundation::ErrorCode::InvalidArgument, std::move(Message));
}

bool IsValidDate(std::string_view Value) {
    if (Value.size() != 10U || Value[4] != '-' || Value[7] != '-') {
        return false;
    }
    for (std::size_t Index = 0; Index < Value.size(); ++Index) {
        if (Index == 4U || Index == 7U) {
            continue;
        }
        if (std::isdigit(static_cast<unsigned char>(Value[Index])) == 0) {
            return false;
        }
    }
    const int Month = (Value[5] - '0') * 10 + Value[6] - '0';
    const int Day = (Value[8] - '0') * 10 + Value[9] - '0';
    return Month >= 1 && Month <= 12 && Day >= 1 && Day <= 31;
}

}  // namespace

MealPlanApplicationService::MealPlanApplicationService(
    std::unique_ptr<MealPlanRepository> RepositoryValue)
    : Repository(std::move(RepositoryValue)) {}

Foundation::Result<Domain::MealPlan> MealPlanApplicationService::Save(
    std::string_view UserId,
    std::string_view PlanDate,
    std::vector<Domain::MealPlanItem> Items) {
    if (!Repository || UserId.empty() || !IsValidDate(PlanDate) || Items.size() > 14U) {
        return Foundation::Result<Domain::MealPlan>::FromError(
            InvalidPlan("计划信息无效"));
    }
    for (const Domain::MealPlanItem& Item : Items) {
        if (Item.RecipeId.empty() || Item.RecipeId.size() > 128U || Item.Servings < 1 ||
            Item.Servings > 24 || Item.SortOrder < 0) {
            return Foundation::Result<Domain::MealPlan>::FromError(
                InvalidPlan("计划菜谱项无效"));
        }
    }
    auto Result = Repository->Save(UserId, PlanDate, std::move(Items));
    if (!Result.HasValue()) {
        return Result;
    }
    Result.Value().CombinedIngredients =
        Domain::PlanRules::MergeIngredients(Result.Value().Items);
    return Result;
}

Foundation::Result<std::vector<Domain::MealPlan>> MealPlanApplicationService::List(
    std::string_view UserId,
    std::string_view FromDate,
    std::string_view ToDate) {
    if (!Repository || UserId.empty() || !IsValidDate(FromDate) || !IsValidDate(ToDate) ||
        FromDate > ToDate) {
        return Foundation::Result<std::vector<Domain::MealPlan>>::FromError(
            InvalidPlan("计划日期范围无效"));
    }
    auto Result = Repository->List(UserId, FromDate, ToDate);
    if (!Result.HasValue()) {
        return Result;
    }
    for (Domain::MealPlan& Plan : Result.Value()) {
        Plan.CombinedIngredients = Domain::PlanRules::MergeIngredients(Plan.Items);
    }
    return Result;
}

}  // namespace Menu::Application
