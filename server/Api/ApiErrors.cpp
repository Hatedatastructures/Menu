#include "Api/ApiErrors.hpp"

#include <boost/json/object.hpp>
#include <boost/json/serialize.hpp>

#include <string>

namespace Menu::Api {

Transport::HttpResponse ApiErrors::Create(
    int Status,
    std::string_view Code,
    std::string_view Message,
    std::string_view RequestId) {
    boost::json::object ErrorObject;
    ErrorObject["code"] = Code;
    ErrorObject["message"] = Message;
    ErrorObject["requestId"] = RequestId;
    boost::json::object Body;
    Body["error"] = std::move(ErrorObject);

    Transport::HttpResponse Response;
    Response.Status = Status;
    Response.ContentType = "application/json; charset=utf-8";
    Response.Body = boost::json::serialize(Body);
    Response.Headers.push_back(Transport::HttpHeader{
        "X-Request-Id", std::string(RequestId)});
    Response.KeepAlive = true;
    return Response;
}

}  // namespace Menu::Api
