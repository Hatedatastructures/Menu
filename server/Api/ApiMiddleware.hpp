#pragma once

#include <Application/AuthService.hpp>
#include <Foundation/Result.hpp>
#include <Transport/Core/HttpRequest.hpp>
#include <Transport/Core/HttpResponse.hpp>

#include <boost/json/value.hpp>

#include <string>
#include <string_view>
#include <vector>

namespace Menu::Api::Middleware {

bool IsAllowedOrigin(
    const std::vector<std::string>& AllowedOrigins,
    std::string_view Origin);

std::string RequestId(const Transport::HttpRequest& Request);

Transport::HttpResponse PreflightResponse();

Transport::HttpResponse JsonResponse(
    int Status,
    boost::json::value Body,
    std::string_view Id);

Transport::HttpResponse InternalError(std::string_view Id);

Transport::HttpResponse AuthError(
    const Foundation::Error& ErrorValue,
    std::string_view Id);

Transport::HttpResponse AdminError(
    const Foundation::Error& ErrorValue,
    std::string_view Id);

Transport::HttpResponse WorkflowError(
    const Foundation::Error& ErrorValue,
    std::string_view Id);

Foundation::Result<Application::AuthUser> RequireUser(
    Application::AuthService* Authentication,
    const Transport::HttpRequest& Request);

Foundation::Result<Application::AuthUser> RequireAdmin(
    Application::AuthService* Authentication,
    const Transport::HttpRequest& Request);

}  // namespace Menu::Api::Middleware
