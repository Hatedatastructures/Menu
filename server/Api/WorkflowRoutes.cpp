#include "Api/WorkflowRoutes.hpp"

#include "Api/ApiErrors.hpp"
#include "Api/ApiMiddleware.hpp"
#include "Api/WorkflowDtos.hpp"

#include <boost/json/parse.hpp>

#include <string>

namespace Menu::Api::WorkflowRoutes {

std::optional<Transport::HttpResponse> Route(
    const ApiRouteContext& Context,
    const Transport::HttpRequest& Request,
    const Parsing::TargetParts& Target,
    std::string_view RequestId) {
    if (Target.Path == "/api/v1/plans") {
        if (Context.MealPlans == nullptr) {
            return ApiErrors::Create(503, "workflow_unavailable", "计划服务暂时不可用", RequestId);
        }
        const auto UserResult = Middleware::RequireUser(Context.Authentication, Request);
        if (!UserResult.HasValue()) {
            return Middleware::WorkflowError(UserResult.ErrorValue(), RequestId);
        }
        if (Request.Method == "GET") {
            const auto QueryResult = Parsing::ReadPlanQuery(Target.Query);
            if (!QueryResult.HasValue()) {
                return ApiErrors::Create(400, "invalid_query", "查询参数无效", RequestId);
            }
            const auto PlansResult = Context.MealPlans->List(
                UserResult.Value().Id,
                QueryResult.Value().first,
                QueryResult.Value().second);
            if (!PlansResult.HasValue()) {
                return Middleware::WorkflowError(PlansResult.ErrorValue(), RequestId);
            }
            return Middleware::JsonResponse(
                200, WorkflowDtos::ToArray(PlansResult.Value()), RequestId);
        }
        if (Request.Method == "POST") {
            if (!Target.Query.empty()) {
                return ApiErrors::Create(400, "invalid_query", "查询参数无效", RequestId);
            }
            const auto Payload = Parsing::ReadMealPlanPayload(Request.Body);
            if (!Payload.HasValue()) {
                return ApiErrors::Create(400, "invalid_body", "请求内容无效", RequestId);
            }
            const auto PlanResult = Context.MealPlans->Save(
                UserResult.Value().Id, Payload.Value().PlanDate, Payload.Value().Items);
            if (!PlanResult.HasValue()) {
                return Middleware::WorkflowError(PlanResult.ErrorValue(), RequestId);
            }
            return Middleware::JsonResponse(
                201, WorkflowDtos::ToObject(PlanResult.Value()), RequestId);
        }
        return ApiErrors::Create(404, "route_not_found", "接口不存在", RequestId);
    }

    constexpr std::string_view CookingSessionPrefix = "/api/v1/cooking-sessions";
    if (Target.Path == CookingSessionPrefix ||
        Target.Path.starts_with(std::string(CookingSessionPrefix) + "/")) {
        if (Context.CookingSessions == nullptr) {
            return ApiErrors::Create(503, "workflow_unavailable", "做饭服务暂时不可用", RequestId);
        }
        const auto UserResult = Middleware::RequireUser(Context.Authentication, Request);
        if (!UserResult.HasValue()) {
            return Middleware::WorkflowError(UserResult.ErrorValue(), RequestId);
        }
        if (Request.Method == "POST" && Target.Path == CookingSessionPrefix) {
            boost::system::error_code ParseError;
            const boost::json::value Parsed = boost::json::parse(Request.Body, ParseError);
            if (ParseError || !Parsed.is_object()) {
                return ApiErrors::Create(400, "invalid_body", "请求内容无效", RequestId);
            }
            const auto RecipeId = Parsing::ReadRequiredString(
                Parsed.as_object(), "recipeId", 128U);
            if (!RecipeId.HasValue() || Parsed.as_object().size() != 1U) {
                return ApiErrors::Create(400, "invalid_body", "请求内容无效", RequestId);
            }
            const auto SessionResult = Context.CookingSessions->Create(
                UserResult.Value().Id, RecipeId.Value());
            if (!SessionResult.HasValue()) {
                return Middleware::WorkflowError(SessionResult.ErrorValue(), RequestId);
            }
            return Middleware::JsonResponse(
                201, WorkflowDtos::ToObject(SessionResult.Value()), RequestId);
        }
        constexpr std::size_t SessionIdOffset = CookingSessionPrefix.size() + 1U;
        if (Request.Method == "PATCH" && Target.Path.size() > SessionIdOffset) {
            const std::string SessionId = Target.Path.substr(SessionIdOffset);
            if (!Parsing::IsSafeIdentifier(SessionId)) {
                return ApiErrors::Create(400, "invalid_id", "会话 ID 无效", RequestId);
            }
            const auto Payload = Parsing::ReadCookingSessionUpdatePayload(Request.Body);
            if (!Payload.HasValue()) {
                return ApiErrors::Create(400, "invalid_body", "请求内容无效", RequestId);
            }
            const auto SessionResult = Context.CookingSessions->Update(
                UserResult.Value().Id,
                SessionId,
                Payload.Value().CurrentStepOrder,
                Payload.Value().State);
            if (!SessionResult.HasValue()) {
                return Middleware::WorkflowError(SessionResult.ErrorValue(), RequestId);
            }
            return Middleware::JsonResponse(
                200, WorkflowDtos::ToObject(SessionResult.Value()), RequestId);
        }
        return ApiErrors::Create(404, "route_not_found", "接口不存在", RequestId);
    }

    if (Request.Method == "POST" && Target.Path == "/api/v1/feedback") {
        if (Context.Feedbacks == nullptr) {
            return ApiErrors::Create(503, "workflow_unavailable", "反馈服务暂时不可用", RequestId);
        }
        const auto UserResult = Middleware::RequireUser(Context.Authentication, Request);
        if (!UserResult.HasValue()) {
            return Middleware::WorkflowError(UserResult.ErrorValue(), RequestId);
        }
        const auto Payload = Parsing::ReadFeedbackPayload(Request.Body);
        if (!Payload.HasValue()) {
            return ApiErrors::Create(400, "invalid_body", "请求内容无效", RequestId);
        }
        const auto FeedbackResult = Context.Feedbacks->Create(
            UserResult.Value().Id,
            Payload.Value().RecipeId,
            Payload.Value().Outcome,
            Payload.Value().Tags,
            Payload.Value().Comment);
        if (!FeedbackResult.HasValue()) {
            return Middleware::WorkflowError(FeedbackResult.ErrorValue(), RequestId);
        }
        return Middleware::JsonResponse(
            201, WorkflowDtos::ToObject(FeedbackResult.Value()), RequestId);
    }

    return std::nullopt;
}

}  // namespace Menu::Api::WorkflowRoutes
