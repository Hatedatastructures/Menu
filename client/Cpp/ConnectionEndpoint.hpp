#pragma once

#include <QString>
#include <QtGlobal>

namespace Menu::Client {

struct ConnectionEndpoint final {
    QString Scheme = QStringLiteral("http");
    QString Host;
    quint16 Port = 8080;
    QString Path = QStringLiteral("/");

    [[nodiscard]] QString ToUrl() const;
};

struct EndpointParseResult final {
    bool IsValid = false;
    ConnectionEndpoint Endpoint;
    QString Error;
};

[[nodiscard]] EndpointParseResult ParseEndpoint(const QString& Value);
[[nodiscard]] QString EndpointScopeKey(const QString& Value);
[[nodiscard]] QString DefaultBaseUrl();

}  // namespace Menu::Client
