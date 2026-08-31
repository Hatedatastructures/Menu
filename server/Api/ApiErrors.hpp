#pragma once

#include <Transport/Core/HttpResponse.hpp>

#include <string_view>

namespace Menu::Api {

class ApiErrors final {
public:
    static Transport::HttpResponse Create(
        int Status,
        std::string_view Code,
        std::string_view Message,
        std::string_view RequestId);
};

}  // namespace Menu::Api
