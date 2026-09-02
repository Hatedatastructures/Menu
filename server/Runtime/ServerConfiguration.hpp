#pragma once

#include <Foundation/Result.hpp>
#include <Transport/HttpServerOptions.hpp>

#include <filesystem>
#include <string>
#include <vector>

namespace Menu::Runtime {

struct ServerConfiguration final {
    Transport::HttpServerOptions Http;
    std::filesystem::path DatabasePath = "server/Data/menu.db";
    std::filesystem::path MediaDirectory = "assets/media";
    std::vector<std::string> CorsOrigins;
};

[[nodiscard]] Foundation::Result<ServerConfiguration> LoadServerConfiguration(
    const std::filesystem::path& ConfigPath);

}  // namespace Menu::Runtime
