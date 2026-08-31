#include "Transport/HttpSession.hpp"

#include <boost/asio/dispatch.hpp>
#include <boost/asio/post.hpp>
#include <boost/beast/http/read.hpp>

#include <algorithm>
#include <string>
#include <utility>

namespace Menu::Transport {
HttpSession::HttpSession(
    boost::asio::ip::tcp::socket SocketValue,
    HttpServerOptions OptionsValue,
    HttpHandler HandlerValue,
    HttpExecutorPost PostValue)
    : Socket(std::move(SocketValue)),
      Timer(Socket.get_executor()),
      Options(std::move(OptionsValue)),
      Handler(std::move(HandlerValue)),
      Post(std::move(PostValue)) {
    if (!Post) {
        Post = [](boost::asio::any_io_executor Executor,
                  std::function<void()> HandlerValue) {
            boost::asio::post(Executor, std::move(HandlerValue));
        };
    }
}

void HttpSession::Run() {
    const auto Self = shared_from_this();
    boost::asio::dispatch(Socket.get_executor(), [Self]() {
        Self->ReadRequest();
    });
}

void HttpSession::ReadRequest() {
    if (Closing.load(std::memory_order_acquire)) {
        return;
    }
    Parser.emplace();
    Parser->body_limit(Options.MaxBodyBytes);
    Parser->header_limit(Options.MaxHeaderBytes);
    SetDeadline();
    const auto Self = shared_from_this();
    boost::beast::http::async_read(
        Socket,
        Buffer,
        *Parser,
        [Self](boost::system::error_code Error, std::size_t) {
            Self->Timer.cancel();
            if (Error == boost::beast::http::error::end_of_stream) {
                Self->Close();
                return;
            }
            if (Error == boost::beast::http::error::body_limit) {
                Self->WriteError(413, "请求体过大");
                return;
            }
            if (Error) {
                Self->WriteError(400, "请求格式无效");
                return;
            }
            if (Self->Parser->get().target().size() > Self->Options.MaxTargetBytes) {
                Self->WriteError(414, "请求目标过长");
                return;
            }
            Self->DispatchRequest();
        });
}

void HttpSession::DispatchRequest() {
    if (Closing.load(std::memory_order_acquire) || !Parser.has_value()) {
        return;
    }
    const HttpRequest RequestValue = ConvertRequest(Parser->get());
    const auto Self = shared_from_this();
    const auto Delivered = std::make_shared<std::atomic<bool>>(false);
    const HttpResponseCallback Complete =
        [Self, Delivered, RequestKeepAlive = RequestValue.KeepAlive](
            HttpResponse ResponseValue) {
            if (Delivered->exchange(true, std::memory_order_acq_rel)) {
                return;
            }
            ResponseValue.KeepAlive = ResponseValue.KeepAlive && RequestKeepAlive;
            try {
                Self->Post(
                    Self->Socket.get_executor(),
                    [Self, ResponseValue = std::move(ResponseValue)]() mutable {
                        Self->WriteResponse(std::move(ResponseValue));
                    });
            } catch (...) {
                Self->RequestClose();
            }
        };
    try {
        Handler(RequestValue, Complete);
    } catch (...) {
        Complete(HttpResponse{500, "text/plain; charset=utf-8", "服务器内部错误", {}, false});
    }
}

void HttpSession::WriteResponse(HttpResponse ResponseValue) {
    if (Closing.load(std::memory_order_acquire)) {
        return;
    }
    const int SafeStatus =
        ResponseValue.Status >= 100 && ResponseValue.Status <= 599
            ? ResponseValue.Status
            : 500;
    Response = BeastResponse{
        static_cast<boost::beast::http::status>(SafeStatus), 11};
    Response.set(boost::beast::http::field::server, "Menu");
    Response.set(boost::beast::http::field::content_type, ResponseValue.ContentType);
    for (const HttpHeader& Header : ResponseValue.Headers) {
        Response.set(Header.Name, Header.Value);
    }
    Response.keep_alive(ResponseValue.KeepAlive);
    Response.body() = std::move(ResponseValue.Body);
    Response.prepare_payload();
    SetDeadline();
    const auto Self = shared_from_this();
    boost::beast::http::async_write(
        Socket,
        Response,
        [Self](boost::system::error_code Error, std::size_t) {
            Self->Timer.cancel();
            if (Error || !Self->Response.keep_alive()) {
                Self->Close();
                return;
            }
            Self->ReadRequest();
        });
}

void HttpSession::WriteError(int Status, std::string Body) {
    HttpResponse ResponseValue;
    ResponseValue.Status = Status;
    ResponseValue.Body = std::move(Body);
    ResponseValue.KeepAlive = false;
    WriteResponse(std::move(ResponseValue));
}

void HttpSession::SetDeadline() {
    Timer.expires_after(std::chrono::seconds(std::max(1, Options.TimeoutSeconds)));
    const auto Self = shared_from_this();
    Timer.async_wait([Self](boost::system::error_code Error) {
        if (!Error) {
            Self->Close();
        }
    });
}

void HttpSession::RequestClose() {
    if (Closing.exchange(true, std::memory_order_acq_rel)) {
        return;
    }
    const auto Self = shared_from_this();
    try {
        Post(Socket.get_executor(), [Self]() {
            Self->CloseOnExecutor();
        });
    } catch (...) {
        // The executor is stopped. The shared session lifecycle releases the
        // socket without touching Asio objects from this thread.
    }
}

void HttpSession::Close() {
    if (Closing.exchange(true, std::memory_order_acq_rel)) {
        return;
    }
    CloseOnExecutor();
}

void HttpSession::CloseOnExecutor() {
    boost::system::error_code Error;
    Timer.cancel();
    Socket.shutdown(boost::asio::ip::tcp::socket::shutdown_both, Error);
    Socket.close(Error);
}

HttpRequest HttpSession::ConvertRequest(
    const boost::beast::http::request<boost::beast::http::string_body>& RequestValue) {
    const auto CopyString = [](boost::beast::string_view Value) {
        return std::string(Value.data(), Value.size());
    };
    HttpRequest Result;
    Result.Method = CopyString(RequestValue.method_string());
    Result.Target = CopyString(RequestValue.target());
    Result.Body = RequestValue.body();
    Result.KeepAlive = RequestValue.keep_alive();
    for (const auto& Field : RequestValue) {
        Result.Headers.push_back(
            HttpHeader{CopyString(Field.name_string()), CopyString(Field.value())});
    }
    return Result;
}

}  // namespace Menu::Transport
