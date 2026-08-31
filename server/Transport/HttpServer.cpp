#include "Transport/HttpServer.hpp"

#include "Transport/HttpSession.hpp"

#include <boost/asio/ip/address.hpp>
#include <boost/asio/strand.hpp>
#include <boost/asio/post.hpp>

#include <atomic>
#include <utility>

namespace Menu::Transport {

struct HttpServer::State {
    State(boost::asio::io_context& IoContextValue, HttpServerOptions OptionsValue)
        : IoContext(IoContextValue),
          Acceptor(boost::asio::make_strand(IoContextValue)),
          Options(std::move(OptionsValue)) {}

    boost::asio::io_context& IoContext;
    boost::asio::ip::tcp::acceptor Acceptor;
    HttpServerOptions Options;
    HttpHandler Handler;
    HttpExecutorPost Post;
    std::atomic<bool> Accepting = false;
    std::uint16_t BoundPort = 0;
};

HttpServer::HttpServer(
    boost::asio::io_context& IoContextValue,
    HttpServerOptions OptionsValue,
    HttpExecutorPost PostValue)
    : ServerState(std::make_shared<State>(IoContextValue, std::move(OptionsValue))) {
    ServerState->Post = std::move(PostValue);
}

HttpServer::~HttpServer() {
    Stop();
}

Foundation::Result<void> HttpServer::Start(HttpHandler HandlerValue) {
    if (!HandlerValue) {
        return Foundation::Result<void>::FromError(
            Foundation::Error(Foundation::ErrorCode::InvalidArgument, "HTTP handler 不能为空"));
    }
    if (ServerState->Accepting.load(std::memory_order_acquire)) {
        return Foundation::Result<void>::FromError(
            Foundation::Error(Foundation::ErrorCode::InvalidArgument, "HTTP 服务已启动"));
    }

    boost::system::error_code Error;
    const auto Address = boost::asio::ip::make_address(ServerState->Options.Address, Error);
    if (Error) {
        return Foundation::Result<void>::FromError(
            Foundation::Error(Foundation::ErrorCode::InvalidArgument, "监听地址无效"));
    }
    const boost::asio::ip::tcp::endpoint Endpoint(Address, ServerState->Options.Port);
    ServerState->Acceptor.open(Endpoint.protocol(), Error);
    if (Error) {
        return Foundation::Result<void>::FromError(
            Foundation::Error(Foundation::ErrorCode::TransportUnavailable, "打开监听器失败"));
    }
    ServerState->Acceptor.set_option(
        boost::asio::socket_base::reuse_address(true), Error);
    if (Error) {
        ServerState->Acceptor.close();
        return Foundation::Result<void>::FromError(
            Foundation::Error(Foundation::ErrorCode::TransportUnavailable, "设置监听器失败"));
    }
    ServerState->Acceptor.bind(Endpoint, Error);
    if (Error) {
        ServerState->Acceptor.close();
        return Foundation::Result<void>::FromError(
            Foundation::Error(Foundation::ErrorCode::TransportUnavailable, "绑定监听地址失败"));
    }
    ServerState->Acceptor.listen(boost::asio::socket_base::max_listen_connections, Error);
    if (Error) {
        ServerState->Acceptor.close();
        return Foundation::Result<void>::FromError(
            Foundation::Error(Foundation::ErrorCode::TransportUnavailable, "启动监听失败"));
    }

    ServerState->Handler = std::move(HandlerValue);
    ServerState->BoundPort = ServerState->Acceptor.local_endpoint(Error).port();
    if (Error) {
        ServerState->Acceptor.close();
        return Foundation::Result<void>::FromError(
            Foundation::Error(Foundation::ErrorCode::TransportUnavailable, "读取监听端口失败"));
    }
    ServerState->Accepting.store(true, std::memory_order_release);
    AcceptNext(ServerState);
    return Foundation::Result<void>();
}

void HttpServer::Stop() noexcept {
    if (!ServerState || !ServerState->Accepting.load(std::memory_order_acquire)) {
        return;
    }
    ServerState->Accepting.store(false, std::memory_order_release);
    boost::system::error_code Error;
    ServerState->Acceptor.close(Error);
}

std::uint16_t HttpServer::LocalPort() const noexcept {
    return ServerState == nullptr ? 0 : ServerState->BoundPort;
}

void HttpServer::AcceptNext(const std::shared_ptr<State>& StateValue) {
    if (!StateValue->Accepting.load(std::memory_order_acquire)) {
        return;
    }
    StateValue->Acceptor.async_accept(
        boost::asio::make_strand(StateValue->IoContext),
        [StateValue](boost::system::error_code Error,
                     boost::asio::ip::tcp::socket SocketValue) {
            if (!Error && StateValue->Accepting.load(std::memory_order_acquire)) {
                std::make_shared<HttpSession>(
                    std::move(SocketValue),
                    StateValue->Options,
                    StateValue->Handler,
                    StateValue->Post)->Run();
            }
            if (StateValue->Accepting.load(std::memory_order_acquire)) {
                AcceptNext(StateValue);
            }
        });
}

}  // namespace Menu::Transport
