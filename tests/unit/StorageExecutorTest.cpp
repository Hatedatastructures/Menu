#include <gtest/gtest.h>

#include <Application/StorageExecutor.hpp>

#include <boost/asio/executor_work_guard.hpp>
#include <boost/asio/io_context.hpp>
#include <boost/asio/post.hpp>

#include <atomic>
#include <chrono>
#include <functional>
#include <future>
#include <stdexcept>
#include <thread>

TEST(StorageExecutorTest, CompletesValueOwnedWorkAndRejectsAfterShutdown) {
    boost::asio::io_context IoContext;
    auto WorkGuard = boost::asio::make_work_guard(IoContext);
    Menu::Application::StorageExecutor Executor(2);
    std::promise<void> CompletionPromise;
    std::atomic<int> CompletionCount = 0;

    for (int Index = 0; Index < 20; ++Index) {
        const bool Accepted = Executor.Submit(
            [Index]() {
                const int OwnedValue = Index;
                static_cast<void>(OwnedValue);
            },
            [&CompletionCount, &CompletionPromise](
                Menu::Application::StorageExecutionResult Result) {
                if (!Result.Succeeded()) {
                    return;
                }
                if (CompletionCount.fetch_add(1) + 1 == 20) {
                    CompletionPromise.set_value();
                }
            },
            IoContext.get_executor());
        ASSERT_TRUE(Accepted);
    }

    std::thread IoThread([&IoContext]() {
        IoContext.run();
    });
    EXPECT_EQ(CompletionPromise.get_future().wait_for(std::chrono::seconds(5)),
              std::future_status::ready);
    Executor.Shutdown();
    EXPECT_FALSE(Executor.Submit(
        []() {},
        [](Menu::Application::StorageExecutionResult) {},
        IoContext.get_executor()));
    WorkGuard.reset();
    IoContext.stop();
    IoThread.join();
    EXPECT_EQ(CompletionCount.load(), 20);
}

TEST(StorageExecutorTest, PropagatesWorkExceptionToCompletion) {
    boost::asio::io_context IoContext;
    auto WorkGuard = boost::asio::make_work_guard(IoContext);
    Menu::Application::StorageExecutor Executor(1);
    std::promise<Menu::Application::StorageExecutionResult> CompletionPromise;
    auto CompletionFuture = CompletionPromise.get_future();

    const bool Accepted = Executor.Submit(
        []() {
            throw std::runtime_error("work failed");
        },
        [&CompletionPromise](Menu::Application::StorageExecutionResult Result) {
            CompletionPromise.set_value(std::move(Result));
        },
        IoContext.get_executor());
    ASSERT_TRUE(Accepted);

    std::thread IoThread([&IoContext]() {
        IoContext.run();
    });
    ASSERT_EQ(CompletionFuture.wait_for(std::chrono::seconds(5)),
              std::future_status::ready);
    const auto Result = CompletionFuture.get();
    ASSERT_FALSE(Result.Succeeded());
    ASSERT_NE(Result.WorkFailure, nullptr);
    EXPECT_EQ(Result.CompletionFailure, nullptr);

    Executor.Shutdown();
    WorkGuard.reset();
    IoContext.stop();
    IoThread.join();
}

TEST(StorageExecutorTest, ReportsCompletionDispatchFailureWhenExecutorIsStopped) {
    boost::asio::io_context IoContext;
    IoContext.stop();
    Menu::Application::StorageExecutor Executor(1);
    std::promise<Menu::Application::StorageExecutionResult> CompletionPromise;
    auto CompletionFuture = CompletionPromise.get_future();

    const bool Accepted = Executor.Submit(
        []() {},
        [&CompletionPromise](Menu::Application::StorageExecutionResult Result) {
            CompletionPromise.set_value(std::move(Result));
        },
        [&IoContext](std::function<void()> Handler) {
            if (IoContext.stopped()) {
                throw std::runtime_error("completion executor stopped");
            }
            boost::asio::post(IoContext, std::move(Handler));
        });
    ASSERT_TRUE(Accepted);

    ASSERT_EQ(CompletionFuture.wait_for(std::chrono::seconds(5)),
              std::future_status::ready);
    const auto Result = CompletionFuture.get();
    ASSERT_FALSE(Result.Succeeded());
    EXPECT_EQ(Result.WorkFailure, nullptr);
    ASSERT_NE(Result.CompletionFailure, nullptr);
    Executor.Shutdown();
}
