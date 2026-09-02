#include "ConnectionSettingsStore.hpp"

#include "ClientApiSupport.hpp"

#include <QDir>
#include <QJsonDocument>
#include <QJsonObject>
#include <QJsonParseError>

namespace Menu::Client {
namespace {

QString SettingsPath(const QString& Directory) {
    return QDir(Directory).filePath(QStringLiteral("connection.json"));
}

}  // namespace

QString ConnectionSettingsStore::LoadBaseUrl(const QString& Directory) {
    const QByteArray Data = Support::ReadCacheFile(SettingsPath(Directory));
    if (Data.isEmpty()) {
        return {};
    }
    QJsonParseError ParseError;
    const QJsonDocument Document = QJsonDocument::fromJson(Data, &ParseError);
    if (ParseError.error != QJsonParseError::NoError || !Document.isObject()) {
        return {};
    }
    const QJsonObject Object = Document.object();
    if (Object.value(QStringLiteral("schemaVersion")).toInt(-1) !=
            Support::CacheSchemaVersion ||
        Object.value(QStringLiteral("kind")).toString() !=
            QStringLiteral("connection")) {
        return {};
    }
    return Object.value(QStringLiteral("baseUrl")).toString();
}

bool ConnectionSettingsStore::SaveBaseUrl(
    const QString& Directory,
    const QString& BaseUrl) {
    QJsonObject Object;
    Object.insert(QStringLiteral("schemaVersion"), Support::CacheSchemaVersion);
    Object.insert(QStringLiteral("kind"), QStringLiteral("connection"));
    Object.insert(QStringLiteral("baseUrl"), BaseUrl);
    return Support::WriteCacheFile(
        SettingsPath(Directory),
        QJsonDocument(Object).toJson(QJsonDocument::Compact));
}

}  // namespace Menu::Client
