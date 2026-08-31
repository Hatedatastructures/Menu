#pragma once

#include <boost/asio/any_io_executor.hpp>
#include <boost/asio/thread_pool.hpp>

#include <atomic>
#include <cstddef>
#include <functional>

namespace Menu::Application {

class StorageExecutor final {
public:
    using Work = std::function<void()>;
    using Completion = std::function<void()>;

    explicit StorageExecutor(std::size_t ThreadCount = 1);
    ~StorageExecutor();

    StorageExecutor(const StorageExecutor&) = delete;
    StorageExecutor& operator=(const StorageExecutor&) = delete;

    bool Submit(
        Work WorkValue,
        Completion CompletionValue,
        boost::asio::any_io_executor CompletionExecutor);

    void Shutdown() noexcept;
    void Join() noexcept;

private:
    boost::asio::thread_pool WorkerPool;
    std::atomic<bool> Stopped = false;
};

}  // namespace Menu::Application
