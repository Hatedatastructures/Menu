#include "Api/AuthRoutes.hpp"

#include "Api/ApiErrors.hpp"
#include "Api/ApiMiddleware.hpp"
#include "Api/AuthDtos.hpp"

#include <boost/json/object.hpp>

#include <string>

namespace Menu::Api::AuthRoutes {

std::optional<Transport::HttpResponse> Route(
    const ApiRouteContext& Context,
    const Transport::HttpRequest& Request,
    const Parsing::TargetParts& Target,
    std::string_view RequestId) {
    if (Request.Method != "POST" ||
        (Target.Path != "/api/v1/auth/register" &&
         Target.Path != "/api/v1/auth/login" &&
         Target.Path != "/api/v1/auth/refresh")) {
        return std::nullopt;
    }
    const std::string Kind = Target.Path.ends_with("/register")
        ? "register"
        : Target.Path.ends_with("/login") ? "login" : "refresh";
    const auto Payload = Parsing::ReadAuthPayload(Request.Body, Kind);
    if (!Payload.HasValue()) {
        return ApiErrors::Create(400, "invalid_body", "请求内容无效", RequestId);
    }
    Foundation::Result<Application::AuthResponse> AuthResult =
        Foundation::Result<Application::AuthResponse>::FromError(
            Foundation::Error(Foundation::ErrorCode::AuthenticationFailed, "认证失败"));
    if (Kind == "register") {
        AuthResult = Context.Authentication.Register(
            Payload.Value().Email,
            Payload.Value().Password,
            Payload.Value().DisplayName);
    } else if (Kind == "login") {
        AuthResult = Context.Authentication.Login(
            Payload.Value().Email, Payload.Value().Password);
    } else {
        AuthResult = Context.Authentication.Refresh(Payload.Value().RefreshToken);
    }
    if (!AuthResult.HasValue()) {
        return Middleware::AuthError(AuthResult.ErrorValue(), RequestId);
    }
    return Middleware::JsonResponse(
        Kind == "register" ? 201 : 200,
        AuthDtos::ToObject(AuthResult.Value()),
        RequestId);
}

}  // namespace Menu::Api::AuthRoutes
