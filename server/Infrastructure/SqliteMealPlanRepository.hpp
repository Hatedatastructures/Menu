#pragma once

#include <Application/MealPlanRepository.hpp>
#include <Infrastructure/SqliteDatabase.hpp>

namespace Menu::Infrastructure {

class SqliteMealPlanRepository final : public Application::MealPlanRepository {
public:
    explicit SqliteMealPlanRepository(SqliteDatabase& DatabaseValue)
        : Database(DatabaseValue) {}

    Foundation::Result<Domain::MealPlan> Save(
        std::string_view UserId,
        std::string_view PlanDate,
        std::vector<Domain::MealPlanItem> Items) override;

    Foundation::Result<std::vector<Domain::MealPlan>> List(
        std::string_view UserId,
        std::string_view FromDate,
        std::string_view ToDate) override;

private:
    Foundation::Result<Domain::MealPlan> ReadPlan(std::string_view PlanId);
    Foundation::Result<std::string> FindRecipeName(std::string_view RecipeId);

    SqliteDatabase& Database;
};

}  // namespace Menu::Infrastructure
