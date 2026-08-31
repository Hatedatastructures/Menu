#pragma once

#include <Foundation/Result.hpp>
#include <Transport/Core/HttpHandler.hpp>
#include <Transport/HttpServerOptions.hpp>

#include <boost/asio/ip/tcp.hpp>
#include <boost/asio/steady_timer.hpp>
#include <boost/beast/core/flat_buffer.hpp>
#include <boost/beast/http/parser.hpp>
#include <boost/beast/http/string_body.hpp>
#include <boost/beast/http/write.hpp>

#include <atomic>
#include <memory>
#include <optional>

namespace Menu::Transport {

class HttpSession final : public std::enable_shared_from_this<HttpSession> {
public:
    HttpSession(
        boost::asio::ip::tcp::socket SocketValue,
        HttpServerOptions OptionsValue,
        HttpHandler HandlerValue,
        HttpExecutorPost PostValue = {});

    void Run();

private:
    using RequestParser = boost::beast::http::request_parser<
        boost::beast::http::string_body>;
    using BeastResponse = boost::beast::http::response<
        boost::beast::http::string_body>;

    void ReadRequest();
    void DispatchRequest();
    void WriteResponse(HttpResponse ResponseValue);
    void WriteError(int Status, std::string Body);
    void SetDeadline();
    void RequestClose();
    void Close();
    void CloseOnExecutor();

    static HttpRequest ConvertRequest(
        const boost::beast::http::request<boost::beast::http::string_body>& RequestValue);

    boost::asio::ip::tcp::socket Socket;
    boost::asio::steady_timer Timer;
    boost::beast::flat_buffer Buffer;
    std::optional<RequestParser> Parser;
    HttpServerOptions Options;
    HttpHandler Handler;
    HttpExecutorPost Post;
    BeastResponse Response;
    std::atomic<bool> Closing = false;
};

}  // namespace Menu::Transport
