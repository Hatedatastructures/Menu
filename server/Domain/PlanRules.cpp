#include "Domain/PlanRules.hpp"

#include <utility>

namespace Menu::Domain {

std::vector<RecipeIngredient> PlanRules::MergeIngredients(
    const std::vector<MealPlanItem>& Items) {
    std::vector<RecipeIngredient> Merged;
    for (const MealPlanItem& Item : Items) {
        for (const RecipeIngredient& Ingredient : Item.Ingredients) {
            auto Existing = Merged.end();
            for (auto Candidate = Merged.begin(); Candidate != Merged.end(); ++Candidate) {
                if (Candidate->IngredientId == Ingredient.IngredientId &&
                    Candidate->Unit == Ingredient.Unit) {
                    Existing = Candidate;
                    break;
                }
            }
            if (Existing == Merged.end()) {
                RecipeIngredient Copy = Ingredient;
                Copy.ServingFactor = 1.0;
                Merged.push_back(std::move(Copy));
            } else {
                Existing->Quantity += Ingredient.Quantity;
                Existing->Required = Existing->Required || Ingredient.Required;
                if (Existing->IngredientName.empty()) {
                    Existing->IngredientName = Ingredient.IngredientName;
                }
                if (Existing->IngredientCategory.empty()) {
                    Existing->IngredientCategory = Ingredient.IngredientCategory;
                }
                if (Existing->IngredientDefaultUnit.empty()) {
                    Existing->IngredientDefaultUnit = Ingredient.IngredientDefaultUnit;
                }
                Existing->IngredientIsPantryStaple =
                    Existing->IngredientIsPantryStaple || Ingredient.IngredientIsPantryStaple;
            }
        }
    }
    return Merged;
}

}  // namespace Menu::Domain
