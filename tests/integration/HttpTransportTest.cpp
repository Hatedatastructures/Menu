#include <gtest/gtest.h>

#include <Transport/Core/HttpHandler.hpp>
#include <Transport/HttpServer.hpp>

#include <boost/asio/connect.hpp>
#include <boost/asio/io_context.hpp>
#include <boost/asio/ip/tcp.hpp>
#include <boost/asio/write.hpp>
#include <boost/beast/core/flat_buffer.hpp>
#include <boost/beast/http.hpp>

#include <chrono>
#include <atomic>
#include <future>
#include <string>
#include <stdexcept>
#include <thread>
#include <utility>
#include <vector>

TEST(HttpTransportTest, HandlesRealHttpRequestWithInjectedHandler) {
    boost::asio::io_context IoContext;
    Menu::Transport::HttpServerOptions Options;
    Options.Address = "127.0.0.1";
    Options.Port = 0;
    Menu::Transport::HttpServer Server(IoContext, Options);
    const auto StartResult = Server.Start(
        [](Menu::Transport::HttpRequest Request,
           Menu::Transport::HttpResponseCallback Complete) {
            Menu::Transport::HttpResponse Response;
            Response.Status = 200;
            Response.Body = "received:" + Request.Target;
            Complete(std::move(Response));
        });
    ASSERT_TRUE(StartResult.HasValue());

    std::thread IoThread([&IoContext]() {
        IoContext.run();
    });

    boost::asio::ip::tcp::socket Socket(IoContext);
    Socket.connect({boost::asio::ip::make_address("127.0.0.1"), Server.LocalPort()});
    boost::beast::http::request<boost::beast::http::string_body> Request(
        boost::beast::http::verb::get, "/test", 11);
    Request.set(boost::beast::http::field::host, "localhost");
    boost::beast::http::write(Socket, Request);

    boost::beast::flat_buffer Buffer;
    boost::beast::http::response<boost::beast::http::string_body> Response;
    boost::beast::http::read(Socket, Buffer, Response);

    EXPECT_EQ(Response.result(), boost::beast::http::status::ok);
    EXPECT_EQ(Response.body(), "received:/test");
    Socket.shutdown(boost::asio::ip::tcp::socket::shutdown_both);
    Server.Stop();
    IoContext.stop();
    IoThread.join();
}

TEST(HttpTransportTest, StopsWhileAcceptCallbacksAndConnectionsAreConcurrent) {
    boost::asio::io_context IoContext;
    Menu::Transport::HttpServerOptions Options;
    Options.Address = "127.0.0.1";
    Options.Port = 0;
    Menu::Transport::HttpServer Server(IoContext, Options);
    std::atomic<int> RequestCount = 0;
    const auto StartResult = Server.Start(
        [&RequestCount](Menu::Transport::HttpRequest,
                        Menu::Transport::HttpResponseCallback Complete) {
            RequestCount.fetch_add(1, std::memory_order_relaxed);
            Complete(Menu::Transport::HttpResponse{});
        });
    ASSERT_TRUE(StartResult.HasValue());

    std::thread IoThread([&IoContext]() {
        IoContext.run();
    });
    const std::uint16_t Port = Server.LocalPort();
    std::vector<std::thread> ConnectorThreads;
    for (int Index = 0; Index < 16; ++Index) {
        ConnectorThreads.emplace_back([Port]() {
            boost::asio::io_context ClientContext;
            boost::asio::ip::tcp::socket Socket(ClientContext);
            boost::system::error_code Error;
            Socket.connect(
                {boost::asio::ip::make_address("127.0.0.1"), Port}, Error);
            Socket.close(Error);
        });
    }
    std::thread StopThread([&Server]() {
        for (int Index = 0; Index < 100; ++Index) {
            Server.Stop();
            std::this_thread::yield();
        }
    });

    StopThread.join();
    for (std::thread& ConnectorThread : ConnectorThreads) {
        ConnectorThread.join();
    }
    Server.Stop();
    IoContext.stop();
    IoThread.join();
    EXPECT_NE(Server.LocalPort(), 0U);
}

TEST(HttpTransportTest, RejectsMalformedRequestAndClosesConnection) {
    boost::asio::io_context IoContext;
    Menu::Transport::HttpServer Server(IoContext, Menu::Transport::HttpServerOptions{});
    ASSERT_TRUE(Server.Start(
                       [](Menu::Transport::HttpRequest,
                          Menu::Transport::HttpResponseCallback Complete) {
                           Complete(Menu::Transport::HttpResponse{});
                       })
                    .HasValue());
    std::thread IoThread([&IoContext]() {
        IoContext.run();
    });

    boost::asio::ip::tcp::socket Socket(IoContext);
    Socket.connect({boost::asio::ip::make_address("127.0.0.1"), Server.LocalPort()});
    const std::string MalformedRequest = "not-http\r\n\r\n";
    boost::asio::write(Socket, boost::asio::buffer(MalformedRequest));
    boost::beast::flat_buffer Buffer;
    boost::beast::http::response<boost::beast::http::string_body> Response;
    boost::beast::http::read(Socket, Buffer, Response);

    EXPECT_EQ(static_cast<int>(Response.result()), 400);
    EXPECT_FALSE(Response.keep_alive());
    Socket.shutdown(boost::asio::ip::tcp::socket::shutdown_both);
    Server.Stop();
    IoContext.stop();
    IoThread.join();
}

TEST(HttpTransportTest, RejectsBodyAboveConfiguredLimit) {
    boost::asio::io_context IoContext;
    Menu::Transport::HttpServerOptions Options;
    Options.Port = 0;
    Options.MaxBodyBytes = 4;
    Menu::Transport::HttpServer Server(IoContext, Options);
    ASSERT_TRUE(Server.Start(
                       [](Menu::Transport::HttpRequest,
                          Menu::Transport::HttpResponseCallback Complete) {
                           Complete(Menu::Transport::HttpResponse{});
                       })
                    .HasValue());
    std::thread IoThread([&IoContext]() {
        IoContext.run();
    });

    boost::asio::ip::tcp::socket Socket(IoContext);
    Socket.connect({boost::asio::ip::make_address("127.0.0.1"), Server.LocalPort()});
    boost::beast::http::request<boost::beast::http::string_body> Request(
        boost::beast::http::verb::post, "/", 11);
    Request.set(boost::beast::http::field::host, "localhost");
    Request.body() = "12345";
    Request.prepare_payload();
    boost::beast::http::write(Socket, Request);
    boost::beast::flat_buffer Buffer;
    boost::beast::http::response<boost::beast::http::string_body> Response;
    boost::beast::http::read(Socket, Buffer, Response);

    EXPECT_EQ(static_cast<int>(Response.result()), 413);
    EXPECT_FALSE(Response.keep_alive());
    Socket.shutdown(boost::asio::ip::tcp::socket::shutdown_both);
    Server.Stop();
    IoContext.stop();
    IoThread.join();
}

TEST(HttpTransportTest, RejectsTargetAboveConfiguredLimit) {
    boost::asio::io_context IoContext;
    Menu::Transport::HttpServerOptions Options;
    Options.Port = 0;
    Options.MaxTargetBytes = 4;
    Menu::Transport::HttpServer Server(IoContext, Options);
    ASSERT_TRUE(Server.Start(
                       [](Menu::Transport::HttpRequest,
                          Menu::Transport::HttpResponseCallback Complete) {
                           Complete(Menu::Transport::HttpResponse{});
                       })
                    .HasValue());
    std::thread IoThread([&IoContext]() {
        IoContext.run();
    });

    boost::asio::ip::tcp::socket Socket(IoContext);
    Socket.connect({boost::asio::ip::make_address("127.0.0.1"), Server.LocalPort()});
    boost::beast::http::request<boost::beast::http::string_body> Request(
        boost::beast::http::verb::get, "/long", 11);
    Request.set(boost::beast::http::field::host, "localhost");
    boost::beast::http::write(Socket, Request);
    boost::beast::flat_buffer Buffer;
    boost::beast::http::response<boost::beast::http::string_body> Response;
    boost::beast::http::read(Socket, Buffer, Response);

    EXPECT_EQ(static_cast<int>(Response.result()), 414);
    EXPECT_FALSE(Response.keep_alive());
    Socket.shutdown(boost::asio::ip::tcp::socket::shutdown_both);
    Server.Stop();
    IoContext.stop();
    IoThread.join();
}

TEST(HttpTransportTest, SerializesKeepAlivePipelineResponses) {
    boost::asio::io_context IoContext;
    Menu::Transport::HttpServer Server(IoContext, Menu::Transport::HttpServerOptions{});
    ASSERT_TRUE(Server.Start(
                       [](Menu::Transport::HttpRequest Request,
                          Menu::Transport::HttpResponseCallback Complete) {
                           Menu::Transport::HttpResponse Response;
                           Response.Body = Request.Target;
                           Complete(std::move(Response));
                       })
                    .HasValue());
    std::thread IoThread([&IoContext]() {
        IoContext.run();
    });

    boost::asio::ip::tcp::socket Socket(IoContext);
    Socket.connect({boost::asio::ip::make_address("127.0.0.1"), Server.LocalPort()});
    const std::string Requests =
        "GET /one HTTP/1.1\r\nHost: localhost\r\nConnection: keep-alive\r\n\r\n"
        "GET /two HTTP/1.1\r\nHost: localhost\r\nConnection: close\r\n\r\n";
    boost::asio::write(Socket, boost::asio::buffer(Requests));
    boost::beast::flat_buffer Buffer;
    boost::beast::http::response<boost::beast::http::string_body> FirstResponse;
    boost::beast::http::read(Socket, Buffer, FirstResponse);
    boost::beast::http::response<boost::beast::http::string_body> SecondResponse;
    boost::beast::http::read(Socket, Buffer, SecondResponse);

    EXPECT_EQ(FirstResponse.body(), "/one");
    EXPECT_TRUE(FirstResponse.keep_alive());
    EXPECT_EQ(SecondResponse.body(), "/two");
    EXPECT_FALSE(SecondResponse.keep_alive());
    Socket.shutdown(boost::asio::ip::tcp::socket::shutdown_both);
    Server.Stop();
    IoContext.stop();
    IoThread.join();
}

TEST(HttpTransportTest, ReportsResponseDispatchFailureWithoutCrossThreadSocketAccess) {
    boost::asio::io_context IoContext;
    Menu::Transport::HttpServerOptions Options;
    Options.Port = 0;
    std::atomic<int> PostCount = 0;
    std::promise<void> FailurePromise;
    auto FailureFuture = FailurePromise.get_future();
    Menu::Transport::HttpExecutorPost ThrowingPost =
        [&PostCount, &FailurePromise](boost::asio::any_io_executor,
                                      std::function<void()>) {
            if (PostCount.fetch_add(1, std::memory_order_relaxed) + 1 == 2) {
                FailurePromise.set_value();
            }
            throw std::runtime_error("completion executor stopped");
        };
    Menu::Transport::HttpServer Server(IoContext, Options, ThrowingPost);
    ASSERT_TRUE(Server.Start(
                       [](Menu::Transport::HttpRequest,
                          Menu::Transport::HttpResponseCallback Complete) {
                           Complete(Menu::Transport::HttpResponse{});
                       })
                    .HasValue());
    std::thread IoThread([&IoContext]() {
        IoContext.run();
    });

    boost::asio::ip::tcp::socket Socket(IoContext);
    Socket.connect({boost::asio::ip::make_address("127.0.0.1"), Server.LocalPort()});
    const std::string Request = "GET / HTTP/1.1\r\nHost: localhost\r\n\r\n";
    boost::asio::write(Socket, boost::asio::buffer(Request));
    EXPECT_EQ(FailureFuture.wait_for(std::chrono::seconds(5)),
              std::future_status::ready);

    Server.Stop();
    Socket.close();
    IoContext.stop();
    IoThread.join();
    EXPECT_EQ(PostCount.load(std::memory_order_relaxed), 2);
}
