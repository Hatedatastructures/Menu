#pragma once

#include <Application/AuthService.hpp>

#include <boost/json/object.hpp>

namespace Menu::Api {

class AuthDtos final {
public:
    static boost::json::object ToObject(const Application::AuthResponse& Response);
};

}  // namespace Menu::Api
