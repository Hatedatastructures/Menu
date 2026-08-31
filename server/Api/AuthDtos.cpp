#include "Api/AuthDtos.hpp"

#include <utility>

namespace Menu::Api {

boost::json::object AuthDtos::ToObject(const Application::AuthResponse& Response) {
    boost::json::object User;
    User["id"] = Response.User.Id;
    User["email"] = Response.User.Email;
    User["displayName"] = Response.User.DisplayName;
    User["isAdmin"] = Response.User.IsAdmin;

    boost::json::object Result;
    Result["user"] = std::move(User);
    Result["accessToken"] = Response.AccessToken;
    Result["refreshToken"] = Response.RefreshToken;
    Result["accessTokenExpiresInSeconds"] = Response.AccessTokenExpiresInSeconds;
    return Result;
}

}  // namespace Menu::Api
