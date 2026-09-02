#include <gtest/gtest.h>

#include <Runtime/ServerConfiguration.hpp>

#include <filesystem>
#include <fstream>

namespace {

void WriteText(const std::filesystem::path& Path, const std::string& Text) {
    std::ofstream Output(Path, std::ios::binary | std::ios::trunc);
    ASSERT_TRUE(Output.good());
    Output << Text;
}

}  // namespace

TEST(ServerConfiguration, LoadsExplicitValues) {
    const auto Path = std::filesystem::temp_directory_path() /
        "menu-server-configuration-test.json";
    WriteText(Path, R"json({
        "address":"0.0.0.0",
        "port":9090,
        "databasePath":"data/menu.db",
        "mediaDirectory":"data/media",
        "maxBodyBytes":2048,
        "maxTargetBytes":4096,
        "timeoutSeconds":30,
        "corsOrigins":["https://admin.example"]
    })json");

    const auto Result = Menu::Runtime::LoadServerConfiguration(Path);
    ASSERT_TRUE(Result.HasValue());
    EXPECT_EQ(Result.Value().Http.Address, "0.0.0.0");
    EXPECT_EQ(Result.Value().Http.Port, 9090);
    EXPECT_EQ(Result.Value().DatabasePath, std::filesystem::path("data/menu.db"));
    EXPECT_EQ(Result.Value().Http.TimeoutSeconds, 30);
    EXPECT_EQ(Result.Value().CorsOrigins.size(), 1U);
    std::error_code Error;
    std::filesystem::remove(Path, Error);
}

TEST(ServerConfiguration, RejectsInvalidValues) {
    const auto Path = std::filesystem::temp_directory_path() /
        "menu-server-configuration-invalid-test.json";
    WriteText(Path, R"json({"address":"127.0.0.1","port":0})json");
    const auto Result = Menu::Runtime::LoadServerConfiguration(Path);
    EXPECT_FALSE(Result.HasValue());
    std::error_code Error;
    std::filesystem::remove(Path, Error);
}

TEST(ServerConfiguration, RejectsMissingFile) {
    const auto Path = std::filesystem::temp_directory_path() /
        "menu-server-configuration-missing-test.json";
    std::error_code Error;
    std::filesystem::remove(Path, Error);
    const auto Result = Menu::Runtime::LoadServerConfiguration(Path);
    EXPECT_FALSE(Result.HasValue());
}
