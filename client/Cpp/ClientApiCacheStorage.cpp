#include "ClientApi.hpp"

#include "ClientApiSupport.hpp"
#include "ConnectionEndpoint.hpp"

#include <QDir>
#include <QFile>
#include <QJsonArray>
#include <QJsonDocument>
#include <QJsonObject>
#include <QJsonParseError>
#include <QThreadPool>

namespace Menu::Client {

void ClientApi::ClearCache() {
    const QString PrivatePlanPath = CacheFilePath(DataKind::Plans);
    ++RequestGeneration;
    const auto Epoch = CacheEpoch;
    const std::uint64_t EpochValue =
        Epoch->fetch_add(1, std::memory_order_acq_rel) + 1;
    Transport.AbortAll();
    PreferenceGeneration->fetch_add(1, std::memory_order_acq_rel);
    PendingOperations = 0;
    DataLoadPending = false;
    PlanLoadPending = false;
    CancelRequested = false;
    SetLoading(false);
    const QJsonArray EmptyValues;
    RecipeValues.SetRecipes(EmptyValues);
    IngredientValues.SetIngredients(EmptyValues);
    RecommendationValues.SetRecommendations(EmptyValues);
    PlanValues.SetPlans(EmptyValues);
    const QString Directory = CacheScopeDirectory();
    const QString PreferencesPath = PreferencesFilePath();
    QThreadPool::globalInstance()->start(
        [Directory, PrivatePlanPath, PreferencesPath, Epoch, EpochValue]() {
            QMutexLocker Locker(&Support::CacheIoMutex());
            if (Epoch->load(std::memory_order_acquire) != EpochValue) {
                return;
            }
            (void)QFile::remove(QDir(Directory).filePath(QStringLiteral("recipes.json")));
            (void)QFile::remove(QDir(Directory).filePath(QStringLiteral("ingredients.json")));
            (void)QFile::remove(QDir(Directory).filePath(QStringLiteral("tonight.json")));
            (void)QFile::remove(PrivatePlanPath);
            (void)QFile::remove(QDir(Directory).filePath(QStringLiteral("plans.json")));
            (void)QFile::remove(PreferencesPath);
            const QStringList PreferenceFiles = QDir(Directory).entryList(
                QStringList{QStringLiteral("preferences*.json")}, QDir::Files);
            for (const QString& File : PreferenceFiles) {
                (void)QFile::remove(QDir(Directory).filePath(File));
            }
        });
    SetError(QStringLiteral("离线缓存已清除"));
}

bool ClientApi::ApplyData(DataKind Kind, const QByteArray& Data, bool IsCached) {
    QJsonParseError ParseError;
    const QJsonDocument Document = QJsonDocument::fromJson(Data, &ParseError);
    if (ParseError.error != QJsonParseError::NoError) {
        return false;
    }
    QJsonArray Values;
    if (!IsCached && Document.isArray()) {
        Values = Document.array();
    } else if (IsCached && Document.isObject()) {
        const QJsonObject Object = Document.object();
        const QString ExpectedKind = Kind == DataKind::Recipes
            ? QStringLiteral("recipes")
            : Kind == DataKind::Ingredients
                  ? QStringLiteral("ingredients")
                  : Kind == DataKind::Recommendations
                        ? QStringLiteral("recommendations")
                        : QStringLiteral("plans");
        if (Object.value(QStringLiteral("schemaVersion")).toInt(-1) !=
                Support::CacheSchemaVersion ||
            Object.value(QStringLiteral("kind")).toString() != ExpectedKind ||
            !Object.value(QStringLiteral("data")).isArray()) {
            return false;
        }
        Values = Object.value(QStringLiteral("data")).toArray();
    } else {
        return false;
    }
    if (Kind == DataKind::Recipes) {
        RecipeValues.SetRecipes(Values);
    } else if (Kind == DataKind::Ingredients) {
        IngredientValues.SetIngredients(Values);
    } else if (Kind == DataKind::Recommendations) {
        RecommendationValues.SetRecommendations(Values);
    } else {
        PlanValues.SetPlans(Values);
    }
    return true;
}

void ClientApi::PersistCache(DataKind Kind, const QByteArray& Data) {
    const QString Path = CacheFilePath(Kind);
    const QString CacheKind = Kind == DataKind::Recipes
        ? QStringLiteral("recipes")
        : Kind == DataKind::Ingredients
              ? QStringLiteral("ingredients")
              : Kind == DataKind::Recommendations
                    ? QStringLiteral("recommendations")
                    : QStringLiteral("plans");
    const auto Epoch = CacheEpoch;
    const std::uint64_t EpochValue = Epoch->load(std::memory_order_acquire);
    QThreadPool::globalInstance()->start([Path, Data, CacheKind, Epoch, EpochValue]() {
        QMutexLocker Locker(&Support::CacheIoMutex());
        QJsonParseError ParseError;
        const QJsonDocument Document = QJsonDocument::fromJson(Data, &ParseError);
        if (ParseError.error != QJsonParseError::NoError || !Document.isArray()) {
            return;
        }
        QJsonObject Envelope;
        Envelope.insert(QStringLiteral("schemaVersion"), Support::CacheSchemaVersion);
        Envelope.insert(QStringLiteral("kind"), CacheKind);
        Envelope.insert(QStringLiteral("data"), Document.array());
        if (Epoch->load(std::memory_order_acquire) == EpochValue) {
            (void)Support::WriteCacheFile(
                Path, QJsonDocument(Envelope).toJson(QJsonDocument::Compact));
        }
    });
}

QString ClientApi::CacheFilePath(DataKind Kind) const {
    const QString FileName = Kind == DataKind::Recipes
        ? QStringLiteral("recipes.json")
        : Kind == DataKind::Ingredients
              ? QStringLiteral("ingredients.json")
              : Kind == DataKind::Recommendations
                    ? QStringLiteral("tonight.json")
                    : QStringLiteral("plans.") +
                          Support::UserScopeHash(Session.UserId()) +
                          QStringLiteral(".json");
    return QDir(CacheScopeDirectory()).filePath(FileName);
}

QString ClientApi::CacheScopeDirectory() const {
    return QDir(CacheDirectory).filePath(
        QStringLiteral("server-") + EndpointScopeKey(ServiceUrl));
}

}  // namespace Menu::Client
