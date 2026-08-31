#pragma once

#include <string>
#include <string_view>
#include <utility>
#include <vector>

namespace Menu::Transport {

struct HttpHeader {
    std::string Name;
    std::string Value;
};

struct HttpRequest {
    std::string Method;
    std::string Target;
    std::string Body;
    std::vector<HttpHeader> Headers;
    bool KeepAlive = true;

    [[nodiscard]] std::string HeaderValue(std::string_view NameValue) const {
        for (const HttpHeader& Header : Headers) {
            if (Header.Name == NameValue) {
                return Header.Value;
            }
        }
        return {};
    }
};

}  // namespace Menu::Transport
