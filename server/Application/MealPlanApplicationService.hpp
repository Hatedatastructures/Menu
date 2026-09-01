#pragma once

#include <Application/MealPlanRepository.hpp>

#include <memory>
#include <string_view>

namespace Menu::Application {

class MealPlanApplicationService final {
public:
    explicit MealPlanApplicationService(std::unique_ptr<MealPlanRepository> RepositoryValue);

    Foundation::Result<Domain::MealPlan> Save(
        std::string_view UserId,
        std::string_view PlanDate,
        std::vector<Domain::MealPlanItem> Items);

    Foundation::Result<std::vector<Domain::MealPlan>> List(
        std::string_view UserId,
        std::string_view FromDate,
        std::string_view ToDate);

private:
    std::unique_ptr<MealPlanRepository> Repository;
};

}  // namespace Menu::Application
