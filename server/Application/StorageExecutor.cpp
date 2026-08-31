#include "Application/StorageExecutor.hpp"

#include <boost/asio/post.hpp>

#include <utility>

namespace Menu::Application {

StorageExecutor::StorageExecutor(std::size_t ThreadCount)
    : WorkerPool(ThreadCount == 0U ? 1U : ThreadCount) {}

StorageExecutor::~StorageExecutor() {
    Shutdown();
}

bool StorageExecutor::Submit(
    Work WorkValue,
    Completion CompletionValue,
    boost::asio::any_io_executor CompletionExecutor) {
    if (Stopped.load(std::memory_order_acquire) || !WorkValue || !CompletionValue) {
        return false;
    }

    boost::asio::post(
        WorkerPool,
        [WorkValue = std::move(WorkValue),
         CompletionValue = std::move(CompletionValue),
         CompletionExecutor = std::move(CompletionExecutor)]() mutable {
            try {
                WorkValue();
            } catch (...) {
            }
            boost::asio::post(CompletionExecutor, std::move(CompletionValue));
        });
    return true;
}

void StorageExecutor::Shutdown() noexcept {
    bool Expected = false;
    if (Stopped.compare_exchange_strong(Expected, true, std::memory_order_acq_rel)) {
        WorkerPool.join();
    }
}

void StorageExecutor::Join() noexcept {
    Shutdown();
}

}  // namespace Menu::Application
