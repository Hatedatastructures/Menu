#pragma once

#include <string>
#include <utility>
#include <vector>

namespace Menu::Domain {

struct RecipeIngredient {
    std::string IngredientId;
    double Quantity = 0.0;
    std::string Unit;
    bool Required = true;
    double ServingFactor = 1.0;
    std::string Preparation;

    RecipeIngredient() = default;

    RecipeIngredient(
        std::string IngredientIdValue,
        double QuantityValue,
        std::string UnitValue,
        bool RequiredValue = true,
        double ServingFactorValue = 1.0,
        std::string PreparationValue = {})
        : IngredientId(std::move(IngredientIdValue)),
          Quantity(QuantityValue),
          Unit(std::move(UnitValue)),
          Required(RequiredValue),
          ServingFactor(ServingFactorValue),
          Preparation(std::move(PreparationValue)) {}
};

struct RecipeStep {
    int StepOrder = 0;
    std::string Title;
    std::string Instruction;
    int DurationSeconds = 0;
    bool HasTimer = false;
};

struct Recipe {
    std::string Id;
    std::string Slug;
    std::string Name;
    std::string Cuisine;
    std::string Description;
    int PrepMinutes = 0;
    int CookMinutes = 0;
    int Servings = 1;
    int Difficulty = 1;
    std::string ImagePath;
    std::string Status;
    std::vector<RecipeIngredient> Ingredients;
    std::vector<RecipeStep> Steps;
    std::vector<std::string> Allergens;
    std::vector<std::string> Cookware;
};

}  // namespace Menu::Domain
