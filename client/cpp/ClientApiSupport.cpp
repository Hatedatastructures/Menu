#include "ClientApiSupport.hpp"

#include <QCryptographicHash>
#include <QDir>
#include <QFile>
#include <QFileInfo>
#include <QSaveFile>

namespace Menu::Client::Support {

QString NormalizeBaseUrl(QString Value) {
    while (Value.endsWith('/')) {
        Value.chop(1);
    }
    return Value;
}

QString UserScopeHash(const QString& UserId) {
    return QString::fromLatin1(QCryptographicHash::hash(
        UserId.toUtf8(), QCryptographicHash::Sha256).toHex());
}

QByteArray ReadCacheFile(const QString& Path) {
    QFile File(Path);
    if (!File.open(QIODevice::ReadOnly)) {
        return {};
    }
    return File.readAll();
}

bool WriteCacheFile(const QString& Path, const QByteArray& Data) {
    const QFileInfo Information(Path);
    if (!QDir().mkpath(Information.absolutePath())) {
        return false;
    }
    QSaveFile File(Path);
    if (!File.open(QIODevice::WriteOnly)) {
        return false;
    }
    if (File.write(Data) != Data.size()) {
        return false;
    }
    return File.commit();
}

}  // namespace Menu::Client::Support
