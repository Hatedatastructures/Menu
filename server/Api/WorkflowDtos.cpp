#include "Api/WorkflowDtos.hpp"

#include <utility>

namespace Menu::Api {
namespace {

boost::json::object ToIngredientObject(const Domain::RecipeIngredient& Ingredient) {
    boost::json::object Result;
    Result["ingredientId"] = Ingredient.IngredientId;
    Result["ingredientName"] = Ingredient.IngredientName;
    Result["category"] = Ingredient.IngredientCategory;
    Result["defaultUnit"] = Ingredient.IngredientDefaultUnit;
    Result["isPantryStaple"] = Ingredient.IngredientIsPantryStaple;
    Result["quantity"] = Ingredient.Quantity;
    Result["unit"] = Ingredient.Unit;
    Result["required"] = Ingredient.Required;
    Result["preparation"] = Ingredient.Preparation;
    return Result;
}

boost::json::array ToIngredientArray(
    const std::vector<Domain::RecipeIngredient>& Ingredients) {
    boost::json::array Result;
    for (const Domain::RecipeIngredient& Ingredient : Ingredients) {
        Result.emplace_back(ToIngredientObject(Ingredient));
    }
    return Result;
}

boost::json::array ToTagArray(const std::vector<std::string>& Tags) {
    boost::json::array Result;
    for (const std::string& Tag : Tags) {
        Result.emplace_back(Tag);
    }
    return Result;
}

}  // namespace

boost::json::object WorkflowDtos::ToObject(const Domain::MealPlan& Plan) {
    boost::json::object Result;
    Result["id"] = Plan.Id;
    Result["planDate"] = Plan.PlanDate;
    boost::json::array Items;
    for (const Domain::MealPlanItem& Item : Plan.Items) {
        boost::json::object ItemObject;
        ItemObject["id"] = Item.Id;
        ItemObject["recipeId"] = Item.RecipeId;
        ItemObject["recipeName"] = Item.RecipeName;
        ItemObject["imagePath"] = Item.ImagePath;
        ItemObject["servings"] = Item.Servings;
        ItemObject["sortOrder"] = Item.SortOrder;
        ItemObject["ingredients"] = ToIngredientArray(Item.Ingredients);
        Items.emplace_back(std::move(ItemObject));
    }
    Result["items"] = std::move(Items);
    Result["combinedIngredients"] = ToIngredientArray(Plan.CombinedIngredients);
    return Result;
}

boost::json::object WorkflowDtos::ToObject(const Domain::CookingSession& Session) {
    boost::json::object Result;
    Result["id"] = Session.Id;
    Result["recipeId"] = Session.RecipeId;
    Result["currentStepOrder"] = Session.CurrentStepOrder;
    Result["state"] = Session.State;
    Result["startedAt"] = Session.StartedAt;
    Result["updatedAt"] = Session.UpdatedAt;
    return Result;
}

boost::json::object WorkflowDtos::ToObject(const Domain::Feedback& FeedbackValue) {
    boost::json::object Result;
    Result["id"] = FeedbackValue.Id;
    Result["recipeId"] = FeedbackValue.RecipeId;
    Result["outcome"] = FeedbackValue.Outcome;
    Result["tags"] = ToTagArray(FeedbackValue.Tags);
    Result["comment"] = FeedbackValue.Comment;
    Result["createdAt"] = FeedbackValue.CreatedAt;
    return Result;
}

boost::json::array WorkflowDtos::ToArray(const std::vector<Domain::MealPlan>& Plans) {
    boost::json::array Result;
    for (const Domain::MealPlan& Plan : Plans) {
        Result.emplace_back(ToObject(Plan));
    }
    return Result;
}

}  // namespace Menu::Api
