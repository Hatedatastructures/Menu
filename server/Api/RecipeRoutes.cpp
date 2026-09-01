#include "Api/RecipeRoutes.hpp"

#include "Api/ApiErrors.hpp"
#include "Api/ApiMiddleware.hpp"
#include "Api/RecipeDtos.hpp"

#include <string>

namespace Menu::Api::RecipeRoutes {

std::optional<Transport::HttpResponse> Route(
    const ApiRouteContext& Context,
    const Transport::HttpRequest& Request,
    const Parsing::TargetParts& Target,
    std::string_view RequestId) {
    if (Request.Method == "GET" && Target.Path == "/api/v1/recipes") {
        const auto OptionsResult = Parsing::ReadRecipeListOptions(Target.Query);
        if (!OptionsResult.HasValue()) {
            return ApiErrors::Create(400, "invalid_query", "查询参数无效", RequestId);
        }
        const auto RecipesResult = Context.Service.ListPublishedRecipes();
        if (!RecipesResult.HasValue()) {
            return Middleware::InternalError(RequestId);
        }
        std::vector<Domain::Recipe> Recipes;
        for (const Domain::Recipe& Recipe : RecipesResult.Value()) {
            const Parsing::RecipeListOptions& Options = OptionsResult.Value();
            if (!Options.Cuisine.empty() && Recipe.Cuisine != Options.Cuisine) {
                continue;
            }
            if (Options.HasMaxMinutes &&
                Recipe.PrepMinutes + Recipe.CookMinutes > Options.MaxMinutes) {
                continue;
            }
            if (Options.Difficulty != 0 && Recipe.Difficulty != Options.Difficulty) {
                continue;
            }
            Recipes.push_back(Recipe);
            if (Recipes.size() == Options.Limit) {
                break;
            }
        }
        return Middleware::JsonResponse(
            200, RecipeDtos::ToArray(Recipes), RequestId);
    }

    if (Request.Method == "GET" && Target.Path == "/api/v1/ingredients") {
        if (!Target.Query.empty()) {
            return ApiErrors::Create(400, "invalid_query", "查询参数无效", RequestId);
        }
        const auto IngredientsResult = Context.Service.ListIngredients();
        if (!IngredientsResult.HasValue()) {
            return Middleware::InternalError(RequestId);
        }
        return Middleware::JsonResponse(
            200, RecipeDtos::ToArray(IngredientsResult.Value()), RequestId);
    }

    constexpr std::string_view RecipePrefix = "/api/v1/recipes/";
    if (Request.Method == "GET" && Target.Path.starts_with(RecipePrefix)) {
        const std::string RecipeId = Target.Path.substr(RecipePrefix.size());
        if (!Parsing::IsSafeIdentifier(RecipeId)) {
            return ApiErrors::Create(400, "invalid_id", "菜谱 ID 无效", RequestId);
        }
        const auto RecipeResult = Context.Service.FindPublishedRecipe(RecipeId);
        if (!RecipeResult.HasValue()) {
            return Middleware::InternalError(RequestId);
        }
        if (!RecipeResult.Value().has_value()) {
            return ApiErrors::Create(404, "recipe_not_found", "菜谱不存在", RequestId);
        }
        return Middleware::JsonResponse(
            200, RecipeDtos::ToObject(RecipeResult.Value().value()), RequestId);
    }

    if (Request.Method == "POST" &&
        Target.Path == "/api/v1/recommendations/tonight") {
        const auto RequestResult = Parsing::ReadRecommendationRequest(Request.Body);
        if (!RequestResult.HasValue()) {
            return ApiErrors::Create(400, "invalid_body", "请求内容无效", RequestId);
        }
        const auto Recommendations = Context.Service.RecommendTonight(RequestResult.Value());
        if (!Recommendations.HasValue()) {
            if (Recommendations.ErrorValue().CodeValue() ==
                Foundation::ErrorCode::InvalidArgument) {
                return ApiErrors::Create(400, "invalid_body", "请求内容无效", RequestId);
            }
            return Middleware::InternalError(RequestId);
        }
        return Middleware::JsonResponse(
            200, RecipeDtos::ToArray(Recommendations.Value()), RequestId);
    }

    return std::nullopt;
}

}  // namespace Menu::Api::RecipeRoutes
