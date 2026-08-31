#pragma once

#include <Transport/Core/HttpRequest.hpp>
#include <Transport/Core/HttpResponse.hpp>

#include <functional>

namespace Menu::Transport {

using HttpResponseCallback = std::function<void(HttpResponse)>;
using HttpHandler = std::function<void(HttpRequest, HttpResponseCallback)>;

}  // namespace Menu::Transport
