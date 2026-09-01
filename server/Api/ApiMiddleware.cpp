#include "Api/ApiMiddleware.hpp"

#include "Api/ApiErrors.hpp"

#include <boost/json/serialize.hpp>

#include <atomic>
#include <cstdint>
#include <ranges>
#include <string>

namespace Menu::Api::Middleware {

bool IsAllowedOrigin(
    const std::vector<std::string>& AllowedOrigins,
    std::string_view Origin) {
    return std::ranges::any_of(AllowedOrigins, [Origin](const std::string& AllowedOrigin) {
        return AllowedOrigin == Origin;
    });
}

std::string RequestId(const Transport::HttpRequest& Request) {
    const std::string Provided = Request.HeaderValue("X-Request-Id");
    if (!Provided.empty() && Provided.size() <= 64U &&
        std::ranges::all_of(Provided, [](char Character) {
            return Character >= 33 && Character <= 126;
        })) {
        return Provided;
    }
    static std::atomic<std::uint64_t> NextId = 1;
    return "menu-" + std::to_string(NextId.fetch_add(1, std::memory_order_relaxed));
}

Transport::HttpResponse PreflightResponse() {
    Transport::HttpResponse Response;
    Response.Status = 204;
    Response.ContentType = "text/plain; charset=utf-8";
    Response.Headers = {
        {"Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS"},
        {"Access-Control-Allow-Headers", "Content-Type, Authorization, X-Request-Id"},
        {"Access-Control-Max-Age", "600"},
        {"Vary", "Origin, Access-Control-Request-Method, Access-Control-Request-Headers"}};
    return Response;
}

Transport::HttpResponse JsonResponse(
    int Status,
    boost::json::value Body,
    std::string_view Id) {
    Transport::HttpResponse Response;
    Response.Status = Status;
    Response.ContentType = "application/json; charset=utf-8";
    Response.Body = boost::json::serialize(std::move(Body));
    Response.Headers.push_back(Transport::HttpHeader{"X-Request-Id", std::string(Id)});
    return Response;
}

Transport::HttpResponse InternalError(std::string_view Id) {
    return ApiErrors::Create(503, "storage_unavailable", "服务暂时不可用", Id);
}

Transport::HttpResponse AuthError(
    const Foundation::Error& ErrorValue,
    std::string_view Id) {
    int Status = 503;
    std::string_view Code = "auth_unavailable";
    std::string_view Message = "认证服务暂时不可用";
    if (ErrorValue.CodeValue() == Foundation::ErrorCode::InvalidArgument) {
        Status = 400;
        Code = "invalid_body";
        Message = "请求内容无效";
    } else if (ErrorValue.CodeValue() == Foundation::ErrorCode::Conflict) {
        Status = 409;
        Code = "email_already_registered";
        Message = "邮箱已注册";
    } else if (ErrorValue.CodeValue() == Foundation::ErrorCode::AuthenticationFailed) {
        Status = 401;
        Code = "authentication_failed";
        Message = "认证信息无效";
    }
    return ApiErrors::Create(Status, Code, Message, Id);
}

Transport::HttpResponse AdminError(
    const Foundation::Error& ErrorValue,
    std::string_view Id) {
    int Status = 503;
    std::string_view Code = "admin_unavailable";
    std::string_view Message = "管理服务暂时不可用";
    if (ErrorValue.CodeValue() == Foundation::ErrorCode::AuthenticationFailed) {
        Status = 401;
        Code = "authentication_required";
        Message = "请先登录";
    } else if (ErrorValue.CodeValue() == Foundation::ErrorCode::Forbidden) {
        Status = 403;
        Code = "admin_forbidden";
        Message = "需要管理员权限";
    } else if (ErrorValue.CodeValue() == Foundation::ErrorCode::InvalidArgument) {
        Status = 400;
        Code = "invalid_body";
        Message = "请求内容无效";
    } else if (ErrorValue.CodeValue() == Foundation::ErrorCode::Conflict) {
        Status = 409;
        Code = "recipe_conflict";
        Message = "菜谱已存在";
    } else if (ErrorValue.CodeValue() == Foundation::ErrorCode::NotFound) {
        Status = 404;
        Code = "recipe_not_found";
        Message = "菜谱不存在";
    }
    return ApiErrors::Create(Status, Code, Message, Id);
}

Transport::HttpResponse WorkflowError(
    const Foundation::Error& ErrorValue,
    std::string_view Id) {
    int Status = 503;
    std::string_view Code = "workflow_unavailable";
    std::string_view Message = "工作流服务暂时不可用";
    if (ErrorValue.CodeValue() == Foundation::ErrorCode::AuthenticationFailed) {
        Status = 401;
        Code = "authentication_required";
        Message = "请先登录";
    } else if (ErrorValue.CodeValue() == Foundation::ErrorCode::Forbidden) {
        Status = 403;
        Code = "workflow_forbidden";
        Message = "无权访问该工作流";
    } else if (ErrorValue.CodeValue() == Foundation::ErrorCode::InvalidArgument) {
        Status = 400;
        Code = "invalid_body";
        Message = "请求内容无效";
    } else if (ErrorValue.CodeValue() == Foundation::ErrorCode::NotFound) {
        Status = 404;
        Code = "resource_not_found";
        Message = "资源不存在";
    } else if (ErrorValue.CodeValue() == Foundation::ErrorCode::Conflict) {
        Status = 409;
        Code = "workflow_conflict";
        Message = "工作流状态冲突";
    }
    return ApiErrors::Create(Status, Code, Message, Id);
}

Foundation::Result<Application::AuthUser> RequireUser(
    Application::AuthService* Authentication,
    const Transport::HttpRequest& Request) {
    if (Authentication == nullptr) {
        return Foundation::Result<Application::AuthUser>::FromError(
            Foundation::Error(Foundation::ErrorCode::StorageUnavailable, "认证服务不可用"));
    }
    const std::string Authorization = Request.HeaderValue("Authorization");
    constexpr std::string_view Prefix = "Bearer ";
    if (!Authorization.starts_with(Prefix) || Authorization.size() == Prefix.size()) {
        return Foundation::Result<Application::AuthUser>::FromError(
            Foundation::Error(Foundation::ErrorCode::AuthenticationFailed, "需要认证"));
    }
    const auto UserResult = Authentication->Authenticate(
        std::string_view(Authorization).substr(Prefix.size()));
    if (!UserResult.HasValue()) {
        return UserResult;
    }
    return UserResult.Value();
}

Foundation::Result<Application::AuthUser> RequireAdmin(
    Application::AuthService* Authentication,
    const Transport::HttpRequest& Request) {
    const auto UserResult = RequireUser(Authentication, Request);
    if (!UserResult.HasValue()) {
        return UserResult;
    }
    if (!UserResult.Value().IsAdmin) {
        return Foundation::Result<Application::AuthUser>::FromError(
            Foundation::Error(Foundation::ErrorCode::Forbidden, "需要管理员权限"));
    }
    return UserResult.Value();
}

}  // namespace Menu::Api::Middleware
