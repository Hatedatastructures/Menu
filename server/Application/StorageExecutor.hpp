#pragma once

#include <boost/asio/any_io_executor.hpp>
#include <boost/asio/thread_pool.hpp>

#include <atomic>
#include <cstddef>
#include <exception>
#include <functional>
#include <memory>

namespace Menu::Application {

struct StorageExecutionResult {
    std::exception_ptr WorkFailure;
    std::exception_ptr CompletionFailure;

    [[nodiscard]] bool Succeeded() const noexcept {
        return WorkFailure == nullptr && CompletionFailure == nullptr;
    }
};

class StorageExecutor final {
public:
    using Work = std::function<void()>;
    using Completion = std::function<void(StorageExecutionResult)>;
    using CompletionDispatcher = std::function<void(std::function<void()>)>;

    explicit StorageExecutor(std::size_t ThreadCount = 1);
    ~StorageExecutor();

    StorageExecutor(const StorageExecutor&) = delete;
    StorageExecutor& operator=(const StorageExecutor&) = delete;

    bool Submit(
        Work WorkValue,
        Completion CompletionValue,
        boost::asio::any_io_executor CompletionExecutor);

    bool Submit(
        Work WorkValue,
        Completion CompletionValue,
        CompletionDispatcher CompletionDispatcherValue);

    void Shutdown() noexcept;
    void Join() noexcept;

private:
    boost::asio::thread_pool WorkerPool;
    std::atomic<bool> Stopped = false;
};

}  // namespace Menu::Application
