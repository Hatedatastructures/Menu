#pragma once

#include <Transport/Core/HttpRequest.hpp>

#include <string>
#include <utility>
#include <vector>

namespace Menu::Transport {

struct HttpResponse {
    int Status = 200;
    std::string ContentType = "text/plain; charset=utf-8";
    std::string Body;
    std::vector<HttpHeader> Headers;
    bool KeepAlive = true;
};

}  // namespace Menu::Transport
