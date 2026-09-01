#pragma once

#include <Domain/MealPlan.hpp>
#include <Foundation/Result.hpp>

#include <string_view>
#include <vector>

namespace Menu::Application {

class MealPlanRepository {
public:
    virtual ~MealPlanRepository() = default;

    virtual Foundation::Result<Domain::MealPlan> Save(
        std::string_view UserId,
        std::string_view PlanDate,
        std::vector<Domain::MealPlanItem> Items) = 0;

    virtual Foundation::Result<std::vector<Domain::MealPlan>> List(
        std::string_view UserId,
        std::string_view FromDate,
        std::string_view ToDate) = 0;
};

}  // namespace Menu::Application
