#include "ClientApi.hpp"

#include "ClientApiSupport.hpp"
#include "PreferencesStore.hpp"

#include <QDir>
#include <QFutureWatcher>
#include <QThreadPool>
#include <QtConcurrent/QtConcurrentRun>

namespace Menu::Client {

void ClientApi::LoadPreferences() {
    PreferencesLoading = true;
    const QString Path = PreferencesFilePath();
    const QString ExpectedUserId = Session.UserId();
    const std::uint64_t Generation = RequestGeneration;
    const auto PreferenceEpoch = PreferenceGeneration;
    const std::uint64_t PreferenceEpochValue =
        PreferenceEpoch->load(std::memory_order_acquire);
    auto* Watcher = new QFutureWatcher<QByteArray>(this);
    connect(Watcher, &QFutureWatcher<QByteArray>::finished, this,
        [this, Watcher, ExpectedUserId, Generation, PreferenceEpoch, PreferenceEpochValue]() {
            const QByteArray Data = Watcher->result();
            Watcher->deleteLater();
            if (Generation != RequestGeneration || ExpectedUserId != Session.UserId()) {
                if (ExpectedUserId == Session.UserId()) {
                    PreferencesLoading = false;
                }
                return;
            }
            if (PreferenceEpoch->load(std::memory_order_acquire) != PreferenceEpochValue) {
                PreferencesLoading = false;
                return;
            }
            if (!Data.isEmpty() && ApplyPreferences(Data)) {
                emit PreferencesChanged();
            }
            PreferencesLoading = false;
        });
    Watcher->setFuture(QtConcurrent::run([Path]() {
        return Support::ReadCacheFile(Path);
    }));
}

bool ClientApi::ApplyPreferences(const QByteArray& Data) {
    const auto Snapshot = PreferencesStore::Decode(Data);
    if (!Snapshot.has_value()) {
        return false;
    }
    return Preferences.Apply(*Snapshot);
}

void ClientApi::PersistPreferences() const {
    const PreferencesSnapshot Snapshot = Preferences.Snapshot();
    const QString Path = PreferencesFilePath();
    const QByteArray Data = PreferencesStore::Encode(Snapshot);
    const auto Epoch = CacheEpoch;
    const std::uint64_t EpochValue = Epoch->load(std::memory_order_acquire);
    const auto Generation = PreferenceGeneration;
    const std::uint64_t GenerationValue =
        Generation->fetch_add(1, std::memory_order_acq_rel) + 1;
    QThreadPool::globalInstance()->start(
        [Path, Data, Epoch, EpochValue, Generation, GenerationValue]() {
            QMutexLocker Locker(&Support::CacheIoMutex());
            if (Epoch->load(std::memory_order_acquire) == EpochValue &&
                Generation->load(std::memory_order_acquire) == GenerationValue) {
                (void)Support::WriteCacheFile(Path, Data);
            }
        });
}

QString ClientApi::PreferencesFilePath() const {
    const QString FileName = Session.UserId().isEmpty()
        ? QStringLiteral("preferences.json")
        : QStringLiteral("preferences.") + Support::UserScopeHash(Session.UserId()) +
              QStringLiteral(".json");
    return QDir(CacheScopeDirectory()).filePath(FileName);
}

}  // namespace Menu::Client
