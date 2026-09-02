#include <gtest/gtest.h>

#include <Composition/ServerApplication.hpp>

#include <chrono>
#include <filesystem>
#include <string>
#include <utility>

namespace {

TEST(ServerApplicationTest, CreatesStartsAndStopsWithSharedComposition) {
    const auto Suffix = std::to_string(
        std::chrono::steady_clock::now().time_since_epoch().count());
    const std::filesystem::path Root =
        std::filesystem::current_path() / ("MenuCompositionTest-" + Suffix);

    Menu::Runtime::ServerConfiguration Configuration;
    Configuration.Http.Address = "127.0.0.1";
    Configuration.Http.Port = 0;
    Configuration.DatabasePath = Root / "Data" / "menu.db";
    Configuration.MediaDirectory = Root / "Media";

    auto Result = Menu::Server::ServerApplication::Create(std::move(Configuration));
    ASSERT_TRUE(Result.HasValue()) << Result.ErrorValue().MessageValue();
    auto Application = std::move(Result).Value();
    ASSERT_TRUE(Application->Start().HasValue());
    EXPECT_GT(Application->LocalPort(), 0U);
    Application->Stop();
    Application.reset();

    std::error_code Error;
    std::filesystem::remove_all(Root, Error);
    EXPECT_FALSE(Error);
}

}  // namespace
