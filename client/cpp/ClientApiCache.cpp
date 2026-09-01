#include "ClientApi.hpp"

#include "ClientApiSupport.hpp"

#include <QDir>
#include <QFile>
#include <QFutureWatcher>
#include <QJsonArray>
#include <QJsonDocument>
#include <QJsonObject>
#include <QJsonParseError>
#include <QJsonValue>
#include <QNetworkReply>
#include <QNetworkRequest>
#include <QThreadPool>
#include <QtConcurrent/QtConcurrentRun>

namespace Menu::Client {

void ClientApi::ClearCache() {
    const QString PrivatePlanPath = CacheFilePath(DataKind::Plans);
    ++RequestGeneration;
    const auto Epoch = CacheEpoch;
    const std::uint64_t EpochValue =
        Epoch->fetch_add(1, std::memory_order_acq_rel) + 1;
    const QList<QNetworkReply*> Replies = ActiveReplies;
    for (QNetworkReply* Reply : Replies) {
        if (Reply != nullptr) {
            Reply->abort();
        }
    }
    ActiveReplies.clear();
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
    const QString Directory = CacheDirectory;
    const QString PreferencesPath = PreferencesFilePath();
    QThreadPool::globalInstance()->start(
        [Directory, PrivatePlanPath, PreferencesPath, Epoch, EpochValue]() {
            if (Epoch->load(std::memory_order_acquire) != EpochValue) {
                return;
            }
            (void)QFile::remove(QDir(Directory).filePath(
                QStringLiteral("recipes.json")));
            (void)QFile::remove(QDir(Directory).filePath(
                QStringLiteral("ingredients.json")));
            (void)QFile::remove(QDir(Directory).filePath(
                QStringLiteral("tonight.json")));
            (void)QFile::remove(PrivatePlanPath);
            (void)QFile::remove(QDir(Directory).filePath(
                QStringLiteral("plans.json")));
            (void)QFile::remove(PreferencesPath);
            (void)QFile::remove(QDir(Directory).filePath(
                QStringLiteral("preferences.json")));
        });
    SetError(QStringLiteral("离线缓存已清除"));
}

void ClientApi::LoadPreferences() {
    PreferencesLoading = true;
    const QString Path = PreferencesFilePath();
    const QString ExpectedUserId = UserIdValue;
    const std::uint64_t Generation = RequestGeneration;
    auto* Watcher = new QFutureWatcher<QByteArray>(this);
    connect(Watcher, &QFutureWatcher<QByteArray>::finished, this,
        [this, Watcher, ExpectedUserId, Generation]() {
            const QByteArray Data = Watcher->result();
            Watcher->deleteLater();
            if (Generation != RequestGeneration || ExpectedUserId != UserIdValue) {
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
    QJsonParseError ParseError;
    const QJsonDocument Document = QJsonDocument::fromJson(Data, &ParseError);
    if (ParseError.error != QJsonParseError::NoError || !Document.isObject()) {
        return false;
    }
    const QJsonObject Object = Document.object();
    if (Object.value(QStringLiteral("schemaVersion")).toInt(-1) !=
            Support::CacheSchemaVersion ||
        Object.value(QStringLiteral("kind")).toString() !=
            QStringLiteral("preferences")) {
        return false;
    }
    const int NextServingCount = Object.value(QStringLiteral("servings")).toInt(-1);
    const QJsonValue CuisinesValue = Object.value(QStringLiteral("cuisines"));
    const QJsonValue AllergiesValueJson = Object.value(QStringLiteral("allergies"));
    const QJsonValue AvailableMinutesValueJson =
        Object.value(QStringLiteral("availableMinutes"));
    const QJsonValue PantryIngredientIdsValueJson =
        Object.value(QStringLiteral("pantryIngredientIds"));
    const QJsonValue CookwareValueJson = Object.value(QStringLiteral("cookware"));
    const QJsonValue RemindersValue = Object.value(QStringLiteral("remindersEnabled"));
    if (NextServingCount < 1 || NextServingCount > 24 ||
        !CuisinesValue.isArray() || !AllergiesValueJson.isArray() ||
        !AvailableMinutesValueJson.isDouble() ||
        AvailableMinutesValueJson.toInt(-1) < 10 ||
        AvailableMinutesValueJson.toInt(-1) > 240 ||
        !PantryIngredientIdsValueJson.isArray() || !CookwareValueJson.isArray() ||
        !RemindersValue.isBool()) {
        return false;
    }
    const auto ReadStringList = [](const QJsonArray& Values) {
        QStringList Result;
        for (const QJsonValue& Value : Values) {
            if (Value.isString() && !Value.toString().isEmpty() &&
                Value.toString().size() <= 32) {
                Result.push_back(Value.toString());
            }
        }
        return Result;
    };
    ServingCountValue = NextServingCount;
    PreferredCuisinesValue = ReadStringList(CuisinesValue.toArray());
    AllergiesValue = ReadStringList(AllergiesValueJson.toArray());
    AvailableMinutesValue = AvailableMinutesValueJson.toInt();
    PantryIngredientIdsValue = ReadStringList(PantryIngredientIdsValueJson.toArray());
    CookwareValue = ReadStringList(CookwareValueJson.toArray());
    RemindersEnabledValue = RemindersValue.toBool();
    return true;
}

void ClientApi::PersistPreferences() const {
    QJsonArray Cuisines;
    for (const QString& Value : PreferredCuisinesValue) {
        Cuisines.append(Value);
    }
    QJsonArray Cookware;
    for (const QString& Value : CookwareValue) {
        Cookware.append(Value);
    }
    QJsonArray Allergies;
    for (const QString& Value : AllergiesValue) {
        Allergies.append(Value);
    }
    QJsonArray PantryIngredientIds;
    for (const QString& Value : PantryIngredientIdsValue) {
        PantryIngredientIds.append(Value);
    }
    QJsonObject Object;
    Object.insert(QStringLiteral("schemaVersion"), Support::CacheSchemaVersion);
    Object.insert(QStringLiteral("kind"), QStringLiteral("preferences"));
    Object.insert(QStringLiteral("servings"), ServingCountValue);
    Object.insert(QStringLiteral("cuisines"), Cuisines);
    Object.insert(QStringLiteral("allergies"), Allergies);
    Object.insert(QStringLiteral("availableMinutes"), AvailableMinutesValue);
    Object.insert(QStringLiteral("pantryIngredientIds"), PantryIngredientIds);
    Object.insert(QStringLiteral("cookware"), Cookware);
    Object.insert(QStringLiteral("remindersEnabled"), RemindersEnabledValue);
    const QString Path = PreferencesFilePath();
    const QByteArray Data = QJsonDocument(Object).toJson(QJsonDocument::Compact);
    const auto Epoch = CacheEpoch;
    const std::uint64_t EpochValue = Epoch->load(std::memory_order_acquire);
    const auto Generation = PreferenceGeneration;
    const std::uint64_t GenerationValue =
        Generation->fetch_add(1, std::memory_order_acq_rel) + 1;
    QThreadPool::globalInstance()->start(
        [Path, Data, Epoch, EpochValue, Generation, GenerationValue]() {
        if (Epoch->load(std::memory_order_acquire) == EpochValue &&
            Generation->load(std::memory_order_acquire) == GenerationValue) {
            (void)Support::WriteCacheFile(Path, Data);
        }
    });
}

QString ClientApi::PreferencesFilePath() const {
    const QString FileName = UserIdValue.isEmpty()
        ? QStringLiteral("preferences.json")
        : QStringLiteral("preferences.") + Support::UserScopeHash(UserIdValue) +
              QStringLiteral(".json");
    return QDir(CacheDirectory).filePath(FileName);
}

void ClientApi::ReadCache(DataKind Kind, std::uint64_t RequestGenerationValue) {
    const QString Path = CacheFilePath(Kind);
    auto* Watcher = new QFutureWatcher<QByteArray>(this);
    connect(Watcher, &QFutureWatcher<QByteArray>::finished, this,
        [this, Kind, Watcher, RequestGenerationValue]() {
            const QByteArray Data = Watcher->result();
            Watcher->deleteLater();
            if (RequestGenerationValue != RequestGeneration) {
                return;
            }
            if (CancelRequested) {
                FinishDataOperation();
                return;
            }
            if (!Data.isEmpty() && ApplyData(Kind, Data, true)) {
                FinishDataOperation();
                return;
            }
            if (RecipeValues.rowCount() == 0 && IngredientValues.rowCount() == 0 &&
                RecommendationValues.rowCount() == 0 && PlanValues.rowCount() == 0) {
                SetError(QStringLiteral("网络不可用，且没有可用的离线数据"));
            }
            FinishDataOperation();
        });
    Watcher->setFuture(QtConcurrent::run([Path]() {
        return Support::ReadCacheFile(Path);
    }));
}

void ClientApi::HandleReply(
    DataKind Kind,
    QNetworkReply* Reply,
    std::uint64_t RequestGenerationValue) {
    const int StatusCode = Reply->attribute(
        QNetworkRequest::HttpStatusCodeAttribute).toInt();
    const bool IsSuccessful = Reply->error() == QNetworkReply::NoError &&
        StatusCode >= 200 && StatusCode < 300;
    const QByteArray Data = Reply->isOpen() ? Reply->readAll() : QByteArray();
    ActiveReplies.removeAll(Reply);
    Reply->deleteLater();

    if (RequestGenerationValue != RequestGeneration) {
        return;
    }

    if (CancelRequested) {
        FinishDataOperation();
        return;
    }

    if (IsSuccessful && ApplyData(Kind, Data)) {
        PersistCache(Kind, Data);
        FinishDataOperation();
        return;
    }

    NetworkFailure = true;
    SetOffline(true);
    if (StatusCode == 401 && !AccessTokenValue.isEmpty()) {
        Logout();
        SetError(QStringLiteral("登录已过期，请重新登录"));
        return;
    }
    ReadCache(Kind, RequestGenerationValue);
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
                          Support::UserScopeHash(UserIdValue) +
                          QStringLiteral(".json");
    return QDir(CacheDirectory).filePath(FileName);
}

}  // namespace Menu::Client
