#pragma once

#include <Foundation/Result.hpp>
#include <Transport/Core/HttpHandler.hpp>
#include <Transport/HttpServerOptions.hpp>

#include <boost/asio/io_context.hpp>
#include <boost/asio/ip/tcp.hpp>

#include <cstdint>
#include <memory>

namespace Menu::Transport {

class HttpServer final {
public:
    HttpServer(
        boost::asio::io_context& IoContextValue,
        HttpServerOptions OptionsValue,
        HttpExecutorPost PostValue = {});
    ~HttpServer();

    HttpServer(const HttpServer&) = delete;
    HttpServer& operator=(const HttpServer&) = delete;

    Foundation::Result<void> Start(HttpHandler HandlerValue);
    void Stop() noexcept;

    [[nodiscard]] std::uint16_t LocalPort() const noexcept;

private:
    struct State;

    static void AcceptNext(const std::shared_ptr<State>& StateValue);

    std::shared_ptr<State> ServerState;
};

}  // namespace Menu::Transport
