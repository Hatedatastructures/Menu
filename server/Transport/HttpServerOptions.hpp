#pragma once

#include <boost/asio/any_io_executor.hpp>

#include <cstddef>
#include <cstdint>
#include <functional>
#include <string>

namespace Menu::Transport {

using HttpExecutorPost = std::function<void(
    boost::asio::any_io_executor,
    std::function<void()>)>;

struct HttpServerOptions {
    std::string Address = "127.0.0.1";
    std::uint16_t Port = 8080;
    std::size_t MaxBodyBytes = 1024U * 1024U;
    std::size_t MaxTargetBytes = 8U * 1024U;
    std::size_t MaxHeaderBytes = 64U * 1024U;
    int TimeoutSeconds = 10;
};

}  // namespace Menu::Transport
