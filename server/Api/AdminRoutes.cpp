#include "Api/AdminRoutes.hpp"

#include "Api/ApiErrors.hpp"
#include "Api/ApiMiddleware.hpp"
#include "Api/RecipeDtos.hpp"

#include <string>

namespace Menu::Api::AdminRoutes {

std::optional<Transport::HttpResponse> Route(
    const ApiRouteContext& Context,
    const Transport::HttpRequest& Request,
    const Parsing::TargetParts& Target,
    std::string_view RequestId) {
    constexpr std::string_view RecipePrefix = "/api/v1/admin/recipes";
    constexpr std::string_view IngredientPrefix = "/api/v1/admin/ingredients";
    const bool IsRecipePath = Target.Path == RecipePrefix ||
        Target.Path.starts_with(std::string(RecipePrefix) + "/");
    const bool IsIngredientPath = Target.Path == IngredientPrefix ||
        Target.Path.starts_with(std::string(IngredientPrefix) + "/");

    if (Target.Path == "/api/v1/admin/stats") {
        const auto UserResult = Middleware::RequireAdmin(Context.Authentication, Request);
        if (!UserResult.HasValue()) {
            return Middleware::AdminError(UserResult.ErrorValue(), RequestId);
        }
        const auto RecipesResult = Context.AdminService.ListAllRecipes();
        const auto IngredientsResult = Context.AdminIngredientService.ListAll();
        if (!RecipesResult.HasValue() || !IngredientsResult.HasValue()) {
            return Middleware::InternalError(RequestId);
        }
        int PublishedCount = 0;
        int DraftCount = 0;
        int ArchivedCount = 0;
        for (const Domain::Recipe& Recipe : RecipesResult.Value()) {
            if (Recipe.Status == "published") {
                ++PublishedCount;
            } else if (Recipe.Status == "draft") {
                ++DraftCount;
            } else if (Recipe.Status == "archived") {
                ++ArchivedCount;
            }
        }
        boost::json::object Body;
        Body["publishedRecipes"] = PublishedCount;
        Body["draftRecipes"] = DraftCount;
        Body["archivedRecipes"] = ArchivedCount;
        Body["ingredientCount"] = IngredientsResult.Value().size();
        return Middleware::JsonResponse(200, std::move(Body), RequestId);
    }

    if (IsRecipePath) {
        const auto UserResult = Middleware::RequireAdmin(Context.Authentication, Request);
        if (!UserResult.HasValue()) {
            return Middleware::AdminError(UserResult.ErrorValue(), RequestId);
        }
        if (Request.Method == "GET" && Target.Path == RecipePrefix) {
            const auto Result = Context.AdminService.ListAllRecipes();
            if (!Result.HasValue()) {
                return Middleware::AdminError(Result.ErrorValue(), RequestId);
            }
            return Middleware::JsonResponse(
                200, RecipeDtos::ToArray(Result.Value()), RequestId);
        }
        if (Request.Method == "POST" && Target.Path == RecipePrefix) {
            const auto Payload = Parsing::ReadRecipePayload(Request.Body);
            if (!Payload.HasValue()) {
                return ApiErrors::Create(400, "invalid_body", "请求内容无效", RequestId);
            }
            const auto Result = Context.AdminService.CreateRecipe(Payload.Value());
            if (!Result.HasValue()) {
                return Middleware::AdminError(Result.ErrorValue(), RequestId);
            }
            return Middleware::JsonResponse(
                201, RecipeDtos::ToObject(Result.Value()), RequestId);
        }
        if (Target.Path.size() <= RecipePrefix.size() + 1U) {
            return ApiErrors::Create(404, "route_not_found", "接口不存在", RequestId);
        }
        const std::string RecipeId = Target.Path.substr(RecipePrefix.size() + 1U);
        if (!Parsing::IsSafeIdentifier(RecipeId)) {
            return ApiErrors::Create(400, "invalid_id", "菜谱 ID 无效", RequestId);
        }
        if (Request.Method == "PATCH") {
            const auto Payload = Parsing::ReadRecipePayload(Request.Body);
            if (!Payload.HasValue()) {
                return ApiErrors::Create(400, "invalid_body", "请求内容无效", RequestId);
            }
            const auto Result = Context.AdminService.UpdateRecipe(RecipeId, Payload.Value());
            if (!Result.HasValue()) {
                return Middleware::AdminError(Result.ErrorValue(), RequestId);
            }
            return Middleware::JsonResponse(
                200, RecipeDtos::ToObject(Result.Value()), RequestId);
        }
        if (Request.Method == "DELETE") {
            const auto Result = Context.AdminService.DeleteRecipe(RecipeId);
            if (!Result.HasValue()) {
                return Middleware::AdminError(Result.ErrorValue(), RequestId);
            }
            Transport::HttpResponse Response;
            Response.Status = 204;
            Response.ContentType = "application/json; charset=utf-8";
            Response.Headers.push_back(Transport::HttpHeader{
                "X-Request-Id", std::string(RequestId)});
            return Response;
        }
        return ApiErrors::Create(404, "route_not_found", "接口不存在", RequestId);
    }

    if (IsIngredientPath) {
        const auto UserResult = Middleware::RequireAdmin(Context.Authentication, Request);
        if (!UserResult.HasValue()) {
            return Middleware::AdminError(UserResult.ErrorValue(), RequestId);
        }
        if (Request.Method == "GET" && Target.Path == IngredientPrefix) {
            const auto Result = Context.AdminIngredientService.ListAll();
            if (!Result.HasValue()) {
                return Middleware::AdminError(Result.ErrorValue(), RequestId);
            }
            return Middleware::JsonResponse(
                200, RecipeDtos::ToArray(Result.Value()), RequestId);
        }
        if (Request.Method == "POST" && Target.Path == IngredientPrefix) {
            const auto Payload = Parsing::ReadIngredientPayload(Request.Body);
            if (!Payload.HasValue()) {
                return ApiErrors::Create(400, "invalid_body", "请求内容无效", RequestId);
            }
            const auto Result = Context.AdminIngredientService.Create(Payload.Value());
            if (!Result.HasValue()) {
                return Middleware::AdminError(Result.ErrorValue(), RequestId);
            }
            return Middleware::JsonResponse(
                201, RecipeDtos::ToObject(Result.Value()), RequestId);
        }
        if (Target.Path.size() <= IngredientPrefix.size() + 1U) {
            return ApiErrors::Create(404, "route_not_found", "接口不存在", RequestId);
        }
        const std::string IngredientId = Target.Path.substr(IngredientPrefix.size() + 1U);
        if (!Parsing::IsSafeIdentifier(IngredientId)) {
            return ApiErrors::Create(400, "invalid_id", "食材 ID 无效", RequestId);
        }
        if (Request.Method == "PATCH") {
            const auto Payload = Parsing::ReadIngredientPayload(Request.Body);
            if (!Payload.HasValue()) {
                return ApiErrors::Create(400, "invalid_body", "请求内容无效", RequestId);
            }
            const auto Result = Context.AdminIngredientService.Update(
                IngredientId, Payload.Value());
            if (!Result.HasValue()) {
                return Middleware::AdminError(Result.ErrorValue(), RequestId);
            }
            return Middleware::JsonResponse(
                200, RecipeDtos::ToObject(Result.Value()), RequestId);
        }
        if (Request.Method == "DELETE") {
            const auto Result = Context.AdminIngredientService.Delete(IngredientId);
            if (!Result.HasValue()) {
                return Middleware::AdminError(Result.ErrorValue(), RequestId);
            }
            Transport::HttpResponse Response;
            Response.Status = 204;
            Response.ContentType = "application/json; charset=utf-8";
            Response.Headers.push_back(Transport::HttpHeader{
                "X-Request-Id", std::string(RequestId)});
            return Response;
        }
        return ApiErrors::Create(404, "route_not_found", "接口不存在", RequestId);
    }

    return std::nullopt;
}

}  // namespace Menu::Api::AdminRoutes
