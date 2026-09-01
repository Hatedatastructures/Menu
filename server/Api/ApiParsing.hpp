#pragma once

#include <Domain/MealPlan.hpp>
#include <Domain/Ingredient.hpp>
#include <Domain/Recommendation.hpp>
#include <Domain/Recipe.hpp>
#include <Foundation/Result.hpp>

#include <boost/json/object.hpp>

#include <map>
#include <optional>
#include <string>
#include <string_view>
#include <utility>
#include <vector>

namespace Menu::Api::Parsing {

struct TargetParts {
    std::string Path;
    std::string Query;
};

struct RecipeListOptions {
    std::string Cuisine;
    int MaxMinutes = 0;
    int Difficulty = 0;
    std::size_t Limit = 20;
    bool HasMaxMinutes = false;
};

struct AuthPayload {
    std::string Email;
    std::string Password;
    std::string DisplayName;
    std::string RefreshToken;
};

struct MealPlanPayload {
    std::string PlanDate;
    std::vector<Domain::MealPlanItem> Items;
};

struct CookingSessionUpdatePayload {
    int CurrentStepOrder = 1;
    std::string State;
};

struct FeedbackPayload {
    std::string RecipeId;
    std::string Outcome;
    std::vector<std::string> Tags;
    std::string Comment;
};

TargetParts SplitTarget(std::string_view Target);

Foundation::Result<std::map<std::string, std::string>> ParseQuery(
    std::string_view Query);

Foundation::Result<int> ParseInteger(
    std::string_view Value,
    int Minimum,
    int Maximum);

Foundation::Result<std::string> ReadRequiredString(
    const boost::json::object& Object,
    std::string_view Key,
    std::size_t MaximumLength);

Foundation::Result<int> ReadJsonInteger(
    const boost::json::object& Object,
    std::string_view Key,
    int Minimum,
    int Maximum);

Foundation::Result<std::vector<std::string>> ReadJsonStringArray(
    const boost::json::object& Object,
    std::string_view Key);

Foundation::Result<Domain::Recipe> ReadRecipePayload(std::string_view Body);

Foundation::Result<Domain::Ingredient> ReadIngredientPayload(std::string_view Body);

Foundation::Result<AuthPayload> ReadAuthPayload(
    std::string_view Body,
    std::string_view Kind);

Foundation::Result<RecipeListOptions> ReadRecipeListOptions(
    std::string_view Query);

Foundation::Result<Domain::RecommendationRequest> ReadRecommendationRequest(
    std::string_view Body);

Foundation::Result<MealPlanPayload> ReadMealPlanPayload(std::string_view Body);

Foundation::Result<CookingSessionUpdatePayload> ReadCookingSessionUpdatePayload(
    std::string_view Body);

Foundation::Result<FeedbackPayload> ReadFeedbackPayload(std::string_view Body);

Foundation::Result<std::pair<std::string, std::string>> ReadPlanQuery(
    std::string_view Query);

bool IsSafeIdentifier(std::string_view Value);

}  // namespace Menu::Api::Parsing
