#include "Api/ApiParsing.hpp"

#include "Api/ApiParsingSupport.hpp"

#include <boost/json/parse.hpp>

#include <set>
#include <string>

namespace Menu::Api::Parsing {

Foundation::Result<AuthPayload> ReadAuthPayload(
    std::string_view Body,
    std::string_view Kind) {
    if (Body.empty() || Body.size() > 64U * 1024U) {
        return Foundation::Result<AuthPayload>::FromError(
            Support::InvalidRequest("认证请求体无效"));
    }
    boost::system::error_code ParseError;
    const boost::json::value Parsed = boost::json::parse(Body, ParseError);
    if (ParseError || !Parsed.is_object()) {
        return Foundation::Result<AuthPayload>::FromError(
            Support::InvalidRequest("认证请求 JSON 无效"));
    }
    const std::set<std::string> AllowedKeys = Kind == "register"
        ? std::set<std::string>{"email", "password", "displayName"}
        : Kind == "login" ? std::set<std::string>{"email", "password"}
                          : std::set<std::string>{"refreshToken"};
    for (const auto& Entry : Parsed.as_object()) {
        if (!AllowedKeys.contains(std::string(Entry.key().data(), Entry.key().size()))) {
            return Foundation::Result<AuthPayload>::FromError(
                Support::InvalidRequest("认证请求包含未知字段"));
        }
    }

    AuthPayload Payload;
    if (Kind == "refresh") {
        const auto RefreshToken = ReadRequiredString(
            Parsed.as_object(), "refreshToken", 256U);
        if (!RefreshToken.HasValue()) {
            return Foundation::Result<AuthPayload>::FromError(RefreshToken.ErrorValue());
        }
        Payload.RefreshToken = RefreshToken.Value();
        return Payload;
    }
    const auto Email = ReadRequiredString(Parsed.as_object(), "email", 320U);
    const std::size_t PasswordMaximumLength = Kind == "login" ? 64U * 1024U : 128U;
    const auto Password = ReadRequiredString(
        Parsed.as_object(), "password", PasswordMaximumLength);
    if (!Email.HasValue() || !Password.HasValue()) {
        return Foundation::Result<AuthPayload>::FromError(
            Support::InvalidRequest("认证请求字段无效"));
    }
    Payload.Email = Email.Value();
    Payload.Password = Password.Value();
    if (Kind == "register") {
        const auto DisplayName = ReadRequiredString(
            Parsed.as_object(), "displayName", 80U);
        if (!DisplayName.HasValue()) {
            return Foundation::Result<AuthPayload>::FromError(DisplayName.ErrorValue());
        }
        Payload.DisplayName = DisplayName.Value();
    }
    return Payload;
}

}  // namespace Menu::Api::Parsing
