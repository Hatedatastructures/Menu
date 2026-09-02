#include "Api/ApiRouter.hpp"

#include "Api/ApiErrors.hpp"
#include "Api/ApiMiddleware.hpp"
#include "Api/ApiParsing.hpp"
#include "Api/ApiRouteContext.hpp"
#include "Api/AdminRoutes.hpp"
#include "Api/AuthRoutes.hpp"
#include "Api/RecipeRoutes.hpp"
#include "Api/WorkflowRoutes.hpp"

#include <boost/asio/system_executor.hpp>
#include <boost/asio/post.hpp>
#include <boost/json/object.hpp>

#include <atomic>
#include <memory>
#include <optional>
#include <utility>

namespace Menu::Api {
namespace {

struct RouteState {
    std::optional<Transport::HttpResponse> Response;
};

}  // namespace

ApiRouter::ApiRouter(ApiDependencies DependenciesValue)
    : Service(DependenciesValue.Service),
      Storage(DependenciesValue.Storage),
      Authentication(DependenciesValue.Authentication),
      AdminService(DependenciesValue.AdminService),
      AdminIngredientService(DependenciesValue.AdminIngredientService),
      MealPlans(DependenciesValue.MealPlans),
      CookingSessions(DependenciesValue.CookingSessions),
      Feedbacks(DependenciesValue.Feedbacks),
      CorsOrigins(std::move(DependenciesValue.CorsOrigins)),
      CompletionDispatcher([](std::function<void()> Handler) {
          boost::asio::post(boost::asio::system_executor(), std::move(Handler));
      }) {}

void ApiRouter::Handle(
    Transport::HttpRequest RequestValue,
    Transport::HttpResponseCallback Complete) {
    if (!Complete) {
        return;
    }
    const std::string Id = Middleware::RequestId(RequestValue);
    const std::string Origin = RequestValue.HeaderValue("Origin");
    const bool IsPreflight = RequestValue.Method == "OPTIONS";
    if (!Origin.empty() && !Middleware::IsAllowedOrigin(CorsOrigins, Origin)) {
        Complete(ApiErrors::Create(403, "cors_denied", "跨域来源未被允许", Id));
        return;
    }
    const Transport::HttpResponseCallback Finish =
        [Complete = std::move(Complete), Origin, IsPreflight](
            Transport::HttpResponse Response) mutable {
            if (!Origin.empty()) {
                Response.Headers.push_back(
                    Transport::HttpHeader{"Access-Control-Allow-Origin", Origin});
                if (!IsPreflight) {
                    Response.Headers.push_back(Transport::HttpHeader{"Vary", "Origin"});
                }
            }
            Complete(std::move(Response));
        };
    const Parsing::TargetParts Target = Parsing::SplitTarget(RequestValue.Target);
    if (RequestValue.Method == "OPTIONS") {
        Finish(Middleware::PreflightResponse());
        return;
    }
    if (RequestValue.Method == "GET" && Target.Path == "/healthz") {
        boost::json::object Body;
        Body["status"] = "ok";
        Finish(Middleware::JsonResponse(200, std::move(Body), Id));
        return;
    }
    if (RequestValue.Method == "GET" && Target.Path == "/readyz") {
        boost::json::object Body;
        Body["status"] = Ready.load(std::memory_order_acquire) ? "ready" : "starting";
        Finish(Middleware::JsonResponse(
            Ready.load(std::memory_order_acquire) ? 200 : 503,
            std::move(Body), Id));
        return;
    }

    auto State = std::make_shared<RouteState>();
    const bool Accepted = Storage.Submit(
        [this,
         RequestValue = std::move(RequestValue),
         Id,
         State]() mutable {
            State->Response = Route(RequestValue, Id);
        },
        [State, Id, Finish](Application::StorageExecutionResult Execution) mutable {
            if (!Execution.Succeeded() || !State->Response.has_value()) {
                Finish(Middleware::InternalError(Id));
                return;
            }
            Finish(std::move(State->Response).value());
        },
        CompletionDispatcher);
    if (!Accepted) {
        Finish(Middleware::InternalError(Id));
    }
}

void ApiRouter::SetReady(bool ReadyValue) noexcept {
    Ready.store(ReadyValue, std::memory_order_release);
}

void ApiRouter::SetCompletionDispatcher(
    Application::StorageExecutor::CompletionDispatcher DispatcherValue) {
    if (DispatcherValue) {
        CompletionDispatcher = std::move(DispatcherValue);
    }
}

Transport::HttpResponse ApiRouter::Route(
    const Transport::HttpRequest& Request,
    const std::string& Id) {
    const Parsing::TargetParts Target = Parsing::SplitTarget(Request.Target);
    const ApiRouteContext Context{
        Service, Authentication, AdminService, AdminIngredientService,
        MealPlans, CookingSessions, Feedbacks};
    if (const auto Response = AuthRoutes::Route(Context, Request, Target, Id);
        Response.has_value()) {
        return Response.value();
    }
    if (const auto Response = AdminRoutes::Route(Context, Request, Target, Id);
        Response.has_value()) {
        return Response.value();
    }
    if (const auto Response = RecipeRoutes::Route(Context, Request, Target, Id);
        Response.has_value()) {
        return Response.value();
    }
    if (const auto Response = WorkflowRoutes::Route(Context, Request, Target, Id);
        Response.has_value()) {
        return Response.value();
    }
    return ApiErrors::Create(404, "route_not_found", "接口不存在", Id);
}

}  // namespace Menu::Api
