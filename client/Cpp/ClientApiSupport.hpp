#pragma once

#include <QByteArray>
#include <QMutex>
#include <QString>

namespace Menu::Client::Support {

inline constexpr int CacheSchemaVersion = 2;
inline constexpr int RequestTimeoutMilliseconds = 8000;

QString NormalizeBaseUrl(QString Value);
QString UserScopeHash(const QString& UserId);
QByteArray ReadCacheFile(const QString& Path);
bool WriteCacheFile(const QString& Path, const QByteArray& Data);
QMutex& CacheIoMutex();

}  // namespace Menu::Client::Support
