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
    const QString PreviousUserId = Session.UserId();
    ResetPrivateState(PreviousUserId, true);
    Session.Clear();
    Transport.SetAccessToken({});
    emit AuthenticationChanged();
}

void ClientApi::HandleAuthenticationReply(
    QNetworkReply* Reply,
    bool IsRegistration,
    bool IsRefresh,
    std::uint64_t RequestGenerationValue,
    std::uint64_t AuthenticationGenerationValue) {
    const int StatusCode = Reply->attribute(
        QNetworkRequest::HttpStatusCodeAttribute).toInt();
    const bool IsSuccessful = Reply->error() == QNetworkReply::NoError &&
        StatusCode >= 200 && StatusCode < 300;
    const QByteArray Data = Reply->isOpen() ? Reply->readAll() : QByteArray();
    Reply->deleteLater();
    if (RequestGenerationValue != RequestGeneration ||
        AuthenticationGenerationValue != AuthenticationGeneration) {
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
    // Go gin-server 返回 vo.Result 包装: { code, msg, data: {...}, requestId, timeStamp }
    // 需要从 data 字段中提取实际的认证数据
    const QJsonObject Data = Object.value(QStringLiteral("data")).toObject();
    const QString AccessToken = Data.value(QStringLiteral("accessToken")).toString();
    const QString RefreshToken = Data.value(QStringLiteral("refreshToken")).toString();
    const QVariantMap User = Data.value(QStringLiteral("user")).toObject().toVariantMap();
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
    ResetPrivateState(Session.UserId(), false);
    Session.Set(
        std::move(AccessToken),
        std::move(RefreshToken),
        std::move(UserDisplayName),
        std::move(UserId));
    Transport.SetAccessToken(Session.AccessToken());
    emit AuthenticationChanged();
    LoadPreferences();
}

void ClientApi::ResetPrivateState(
    const QString& UserIdToRemove,
    bool RemovePlanCache) {
    ++RequestGeneration;
    ++AuthenticationGeneration;
    const auto Epoch = CacheEpoch;
    const std::uint64_t EpochValue =
        Epoch->fetch_add(1, std::memory_order_acq_rel) + 1;
    Transport.AbortAll();
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
        emit cookingSessionChanged();
    }

    if (!RemovePlanCache || UserIdToRemove.isEmpty()) {
        return;
    }
    const QString PrivatePlanPath = CacheFilePath(DataKind::Plans);
    const QString LegacyPlanPath = QDir(CacheDirectory).filePath(
        QStringLiteral("plans.json"));
    QThreadPool::globalInstance()->start(
        [PrivatePlanPath, LegacyPlanPath, Epoch, EpochValue]() {
            QMutexLocker Locker(&Support::CacheIoMutex());
            if (Epoch->load(std::memory_order_acquire) != EpochValue) {
                return;
            }
            (void)QFile::remove(PrivatePlanPath);
            (void)QFile::remove(LegacyPlanPath);
        });
}

void ClientApi::ResetPreferencesToDefaults() {
    const bool Changed = Preferences.ServingCount() != 2 ||
        Preferences.PreferredCuisines() !=
            QStringList{QStringLiteral("中餐"), QStringLiteral("西餐"), QStringLiteral("日系")} ||
        !Preferences.Allergies().isEmpty() || Preferences.AvailableMinutes() != 45 ||
        !Preferences.PantryIngredientIds().isEmpty() ||
        Preferences.Cookware() != QStringList{
            QStringLiteral("炒锅"), QStringLiteral("平底锅"), QStringLiteral("烤箱")} ||
        !Preferences.RemindersEnabled();
    Preferences.Reset();
    if (Changed) {
        emit PreferencesChanged();
    }
}

}  // namespace Menu::Client
