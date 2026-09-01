#pragma once

#include <Domain/MealPlan.hpp>

#include <vector>

namespace Menu::Domain {

class PlanRules final {
public:
    static std::vector<RecipeIngredient> MergeIngredients(
        const std::vector<MealPlanItem>& Items);
};

}  // namespace Menu::Domain
