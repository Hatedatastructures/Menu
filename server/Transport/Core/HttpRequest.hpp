#pragma once

#include <algorithm>
#include <cctype>
#include <ranges>
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
        const auto EqualName = [](std::string_view Left, std::string_view Right) {
            return Left.size() == Right.size() &&
                   std::ranges::equal(Left, Right, [](char LeftCharacter, char RightCharacter) {
                       return std::tolower(static_cast<unsigned char>(LeftCharacter)) ==
                              std::tolower(static_cast<unsigned char>(RightCharacter));
                   });
        };
        for (const HttpHeader& Header : Headers) {
            if (EqualName(Header.Name, NameValue)) {
                return Header.Value;
            }
        }
        return {};
    }
};

}  // namespace Menu::Transport
