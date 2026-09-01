#pragma once

#include "Api/ApiParsing.hpp"
#include "Api/ApiRouteContext.hpp"

#include <Transport/Core/HttpResponse.hpp>

#include <optional>
#include <string_view>

namespace Menu::Api::RecipeRoutes {

std::optional<Transport::HttpResponse> Route(
    const ApiRouteContext& Context,
    const Transport::HttpRequest& Request,
    const Parsing::TargetParts& Target,
    std::string_view RequestId);

}  // namespace Menu::Api::RecipeRoutes
