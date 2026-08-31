#include "Application/StorageExecutor.hpp"

#include <boost/asio/post.hpp>

#include <exception>
#include <memory>
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
    return Submit(
        std::move(WorkValue),
        std::move(CompletionValue),
        [CompletionExecutor = std::move(CompletionExecutor)](
            std::function<void()> Handler) mutable {
            boost::asio::post(CompletionExecutor, std::move(Handler));
        });
}

bool StorageExecutor::Submit(
    Work WorkValue,
    Completion CompletionValue,
    CompletionDispatcher CompletionDispatcherValue) {
    if (Stopped.load(std::memory_order_acquire) || !WorkValue || !CompletionValue ||
        !CompletionDispatcherValue) {
        return false;
    }

    auto CompletionState = std::make_shared<Completion>(std::move(CompletionValue));
    auto CompletionDelivered = std::make_shared<std::atomic<bool>>(false);
    const auto DeliverCompletion = [CompletionState, CompletionDelivered](
                                      StorageExecutionResult Result) {
        if (CompletionDelivered->exchange(true, std::memory_order_acq_rel)) {
            return;
        }
        try {
            (*CompletionState)(std::move(Result));
        } catch (...) {
        }
    };

    try {
        boost::asio::post(
            WorkerPool,
            [WorkValue = std::move(WorkValue),
             CompletionDispatcherValue = std::move(CompletionDispatcherValue),
             DeliverCompletion]() mutable {
                StorageExecutionResult Result;
                try {
                    WorkValue();
                } catch (...) {
                    Result.WorkFailure = std::current_exception();
                }

                auto CompletionTask = [DeliverCompletion, Result]() mutable {
                    DeliverCompletion(std::move(Result));
                };
                try {
                    CompletionDispatcherValue(std::move(CompletionTask));
                } catch (...) {
                    Result.CompletionFailure = std::current_exception();
                    DeliverCompletion(std::move(Result));
                }
            });
    } catch (...) {
        StorageExecutionResult Result;
        Result.CompletionFailure = std::current_exception();
        DeliverCompletion(std::move(Result));
        return false;
    }
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
