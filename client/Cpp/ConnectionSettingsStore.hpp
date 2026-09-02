#pragma once

#include <QString>

namespace Menu::Client {

class ConnectionSettingsStore final {
public:
    [[nodiscard]] static QString LoadBaseUrl(const QString& Directory);
    [[nodiscard]] static bool SaveBaseUrl(
        const QString& Directory,
        const QString& BaseUrl);
};

}  // namespace Menu::Client
