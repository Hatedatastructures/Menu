#include "ClientApi.hpp"

#include "ClientApiSupport.hpp"

#include <QDir>
#include <QFile>
#include <QJsonArray>
#include <QJsonDocument>
#include <QJsonObject>
#include <QJsonParseError>
#include <QNetworkReply>
#include <QNetworkRequest>
#include <QThreadPool>

#include <utility>

namespace Menu::Client {

void ClientApi::Logout() {
    const QString PreviousUserId = UserIdValue;
    ResetPrivateState(PreviousUserId, true);
    AccessTokenValue.clear();
    RefreshTokenValue.clear();
    DisplayNameValue.clear();
    UserIdValue.clear();
    emit AuthenticationChanged();
}

void ClientApi::HandleAuthenticationReply(
    QNetworkReply* Reply,
    bool IsRegistration,
    bool IsRefresh,
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
    if (!IsSuccessful) {
        if (IsRefresh) {
            Logout();
            SetError(QStringLiteral("登录已过期，请重新登录"));
        } else {
            SetError(IsRegistration
                    ? QStringLiteral("创建账户失败，请检查邮箱和密码")
                    : QStringLiteral("登录失败，请检查邮箱和密码"));
        }
        return;
    }
    QJsonParseError ParseError;
    const QJsonDocument Document = QJsonDocument::fromJson(Data, &ParseError);
    if (ParseError.error != QJsonParseError::NoError || !Document.isObject()) {
        SetError(QStringLiteral("登录响应无效"));
        return;
    }
    const QJsonObject Object = Document.object();
    const QString AccessToken = Object.value(QStringLiteral("accessToken")).toString();
    const QString RefreshToken = Object.value(QStringLiteral("refreshToken")).toString();
    const QVariantMap User = Object.value(QStringLiteral("user")).toObject().toVariantMap();
    const QString UserDisplayName = User.value(QStringLiteral("displayName")).toString();
    const QString UserId = User.value(QStringLiteral("id")).toString();
    if (AccessToken.isEmpty() || RefreshToken.isEmpty() || UserDisplayName.isEmpty() ||
        UserId.isEmpty()) {
        SetError(QStringLiteral("登录响应缺少账户信息"));
        return;
    }
    SetAuthentication(AccessToken, RefreshToken, UserDisplayName, UserId);
    LoadData();
    LoadPlans();
}

void ClientApi::SetAuthentication(
    QString AccessToken,
    QString RefreshToken,
    QString UserDisplayName,
    QString UserId) {
    ResetPrivateState(UserIdValue, false);
    AccessTokenValue = std::move(AccessToken);
    RefreshTokenValue = std::move(RefreshToken);
    DisplayNameValue = std::move(UserDisplayName);
    UserIdValue = std::move(UserId);
    emit AuthenticationChanged();
    LoadPreferences();
}

void ClientApi::ResetPrivateState(
    const QString& UserIdToRemove,
    bool RemovePlanCache) {
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
    AuthenticationReply = nullptr;
    CookingSessionReply = nullptr;
    FeedbackReply = nullptr;
    PendingOperations = 0;
    DataLoadPending = false;
    PlanLoadPending = false;
    CancelRequested = false;
    NetworkFailure = false;
    PreferencesLoading = false;
    SetLoading(false);

    const bool HadActiveRecipe = !ActiveRecipeValue.isEmpty();
    const bool HadCookingSession = !CurrentCookingSessionValue.isEmpty();
    ActiveRecipeValue.clear();
    CurrentCookingSessionValue.clear();
    PlanValues.SetPlans(QJsonArray());
    ResetPreferencesToDefaults();
    if (HadActiveRecipe) {
        emit ActiveRecipeChanged();
    }
    if (HadCookingSession) {
        emit CookingSessionChanged();
    }

    if (!RemovePlanCache || UserIdToRemove.isEmpty()) {
        return;
    }
    const QString PrivatePlanPath = QDir(CacheDirectory).filePath(
        QStringLiteral("plans.") + Support::UserScopeHash(UserIdToRemove) +
        QStringLiteral(".json"));
    const QString LegacyPlanPath = QDir(CacheDirectory).filePath(
        QStringLiteral("plans.json"));
    QThreadPool::globalInstance()->start(
        [PrivatePlanPath, LegacyPlanPath, Epoch, EpochValue]() {
            if (Epoch->load(std::memory_order_acquire) != EpochValue) {
                return;
            }
            (void)QFile::remove(PrivatePlanPath);
            (void)QFile::remove(LegacyPlanPath);
        });
}

void ClientApi::ResetPreferencesToDefaults() {
    const bool Changed = ServingCountValue != 2 ||
        PreferredCuisinesValue !=
            QStringList{QStringLiteral("中餐"), QStringLiteral("西餐"), QStringLiteral("日系")} ||
        !AllergiesValue.isEmpty() || AvailableMinutesValue != 45 ||
        !PantryIngredientIdsValue.isEmpty() ||
        CookwareValue != QStringList{
            QStringLiteral("炒锅"), QStringLiteral("平底锅"), QStringLiteral("烤箱")} ||
        !RemindersEnabledValue;
    ServingCountValue = 2;
    PreferredCuisinesValue = {
        QStringLiteral("中餐"), QStringLiteral("西餐"), QStringLiteral("日系")};
    AllergiesValue.clear();
    AvailableMinutesValue = 45;
    PantryIngredientIdsValue.clear();
    CookwareValue = {
        QStringLiteral("炒锅"), QStringLiteral("平底锅"), QStringLiteral("烤箱")};
    RemindersEnabledValue = true;
    if (Changed) {
        emit PreferencesChanged();
    }
}

}  // namespace Menu::Client
