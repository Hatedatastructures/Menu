#include "Api/RecipeDtos.hpp"

#include <utility>

namespace Menu::Api {
namespace {

boost::json::array ToStringArray(const std::vector<std::string>& Values) {
    boost::json::array Result;
    for (const std::string& Value : Values) {
        Result.emplace_back(Value);
    }
    return Result;
}

boost::json::array ToIngredientArray(
    const std::vector<Domain::RecipeIngredient>& Ingredients) {
    boost::json::array Result;
    for (const Domain::RecipeIngredient& IngredientValue : Ingredients) {
        boost::json::object Item;
        Item["ingredientId"] = IngredientValue.IngredientId;
        Item["quantity"] = IngredientValue.Quantity;
        Item["unit"] = IngredientValue.Unit;
        Item["required"] = IngredientValue.Required;
        Item["servingFactor"] = IngredientValue.ServingFactor;
        Item["preparation"] = IngredientValue.Preparation;
        Item["ingredientName"] = IngredientValue.IngredientName;
        Item["category"] = IngredientValue.IngredientCategory;
        Item["defaultUnit"] = IngredientValue.IngredientDefaultUnit;
        Item["isPantryStaple"] = IngredientValue.IngredientIsPantryStaple;
        Result.emplace_back(std::move(Item));
    }
    return Result;
}

boost::json::array ToStepArray(const std::vector<Domain::RecipeStep>& Steps) {
    boost::json::array Result;
    for (const Domain::RecipeStep& StepValue : Steps) {
        boost::json::object Item;
        Item["stepOrder"] = StepValue.StepOrder;
        Item["title"] = StepValue.Title;
        Item["instruction"] = StepValue.Instruction;
        Item["durationSeconds"] = StepValue.DurationSeconds;
        Item["hasTimer"] = StepValue.HasTimer;
        Result.emplace_back(std::move(Item));
    }
    return Result;
}

}  // namespace

boost::json::object RecipeDtos::ToObject(const Domain::Recipe& RecipeValue) {
    boost::json::object Result;
    Result["id"] = RecipeValue.Id;
    Result["slug"] = RecipeValue.Slug;
    Result["name"] = RecipeValue.Name;
    Result["cuisine"] = RecipeValue.Cuisine;
    Result["description"] = RecipeValue.Description;
    Result["prepMinutes"] = RecipeValue.PrepMinutes;
    Result["cookMinutes"] = RecipeValue.CookMinutes;
    Result["totalMinutes"] = RecipeValue.PrepMinutes + RecipeValue.CookMinutes;
    Result["servings"] = RecipeValue.Servings;
    Result["difficulty"] = RecipeValue.Difficulty;
    Result["imagePath"] = RecipeValue.ImagePath;
    Result["status"] = RecipeValue.Status;
    Result["allergens"] = ToStringArray(RecipeValue.Allergens);
    Result["cookware"] = ToStringArray(RecipeValue.Cookware);
    Result["ingredients"] = ToIngredientArray(RecipeValue.Ingredients);
    Result["steps"] = ToStepArray(RecipeValue.Steps);
    return Result;
}

boost::json::object RecipeDtos::ToObject(const Domain::Ingredient& IngredientValue) {
    boost::json::object Result;
    Result["id"] = IngredientValue.Id;
    Result["name"] = IngredientValue.Name;
    Result["aliases"] = ToStringArray(IngredientValue.Aliases);
    Result["category"] = IngredientValue.Category;
    Result["defaultUnit"] = IngredientValue.DefaultUnit;
    Result["isPantryStaple"] = IngredientValue.IsPantryStaple;
    Result["substituteGroup"] = IngredientValue.SubstituteGroup;
    Result["storeSkuMapping"] = IngredientValue.StoreSkuMapping;
    return Result;
}

boost::json::object RecipeDtos::ToObject(
    const Domain::Recommendation& RecommendationValue) {
    boost::json::object Result;
    Result["recipe"] = ToObject(RecommendationValue.RecipeValue);
    Result["availableIngredientCount"] =
        RecommendationValue.Availability.AvailableCount;
    Result["missingIngredientCount"] = RecommendationValue.Availability.MissingCount;
    Result["score"] = RecommendationValue.Score;
    return Result;
}

boost::json::array RecipeDtos::ToArray(const std::vector<Domain::Recipe>& Recipes) {
    boost::json::array Result;
    for (const Domain::Recipe& RecipeValue : Recipes) {
        Result.emplace_back(ToObject(RecipeValue));
    }
    return Result;
}

boost::json::array RecipeDtos::ToArray(
    const std::vector<Domain::Ingredient>& Ingredients) {
    boost::json::array Result;
    for (const Domain::Ingredient& IngredientValue : Ingredients) {
        Result.emplace_back(ToObject(IngredientValue));
    }
    return Result;
}

boost::json::array RecipeDtos::ToArray(
    const std::vector<Domain::Recommendation>& Recommendations) {
    boost::json::array Result;
    for (const Domain::Recommendation& RecommendationValue : Recommendations) {
        Result.emplace_back(ToObject(RecommendationValue));
    }
    return Result;
}

}  // namespace Menu::Api
