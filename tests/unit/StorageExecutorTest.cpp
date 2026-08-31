#include <gtest/gtest.h>

#include <Application/StorageExecutor.hpp>

#include <boost/asio/executor_work_guard.hpp>
#include <boost/asio/io_context.hpp>

#include <atomic>
#include <chrono>
#include <future>
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
            [&CompletionCount, &CompletionPromise]() {
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
    EXPECT_FALSE(Executor.Submit([]() {}, []() {}, IoContext.get_executor()));
    WorkGuard.reset();
    IoContext.stop();
    IoThread.join();
    EXPECT_EQ(CompletionCount.load(), 20);
}
