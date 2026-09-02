#include "ClientApi.hpp"

#include "ClientApiSupport.hpp"

#include <QDate>
#include <QJsonArray>
#include <QJsonDocument>
#include <QJsonObject>
#include <QNetworkReply>
#include <QNetworkRequest>
#include <QTimer>
#include <QUrl>

namespace Menu::Client {

void ClientApi::LoadData() {
    if (Loading) {
        DataLoadPending = true;
        return;
    }
    SetError({});
    SetOffline(false);
    NetworkFailure = false;
    CancelRequested = false;
    SetLoading(true);
    PendingOperations = 3;
    RequestData(DataKind::Recipes);
    RequestData(DataKind::Ingredients);
    RequestData(DataKind::Recommendations);
}

void ClientApi::LoadPlans() {
    if (!IsAuthenticated()) {
        SetError(QStringLiteral("登录后才能查看和编辑本周计划"));
        return;
    }
    if (Loading) {
        PlanLoadPending = true;
        return;
    }
    SetError({});
    SetOffline(false);
    NetworkFailure = false;
    CancelRequested = false;
    SetLoading(true);
    PendingOperations = 1;
    RequestData(DataKind::Plans);
}

void ClientApi::SavePlan(const QString& PlanDate, const QVariantList& RecipeIds) {
    if (!IsAuthenticated()) {
        SetError(QStringLiteral("登录后才能保存计划"));
        return;
    }
    QJsonArray Items;
    int SortOrder = 0;
    for (const QVariant& Value : RecipeIds) {
        const QString RecipeId = Value.toString();
        if (RecipeId.isEmpty()) {
            continue;
        }
        QJsonObject Item;
        Item.insert(QStringLiteral("recipeId"), RecipeId);
        Item.insert(QStringLiteral("servings"), Preferences.ServingCount());
        Item.insert(QStringLiteral("sortOrder"), SortOrder++);
        Items.append(Item);
    }
    QJsonObject Body;
    Body.insert(QStringLiteral("planDate"), PlanDate);
    Body.insert(QStringLiteral("items"), Items);
    const std::uint64_t Generation = RequestGeneration;
    QNetworkReply* Reply = Transport.Post(
        QStringLiteral("/api/v1/plans"), QJsonDocument(Body).toJson(QJsonDocument::Compact));
    connect(Reply, &QNetworkReply::finished, this, [this, Reply, Generation]() {
        const int StatusCode = Reply->attribute(
            QNetworkRequest::HttpStatusCodeAttribute).toInt();
        const bool IsSuccessful = Reply->error() == QNetworkReply::NoError &&
            StatusCode >= 200 && StatusCode < 300;
        Reply->deleteLater();
        if (Generation != RequestGeneration) {
            return;
        }
        if (StatusCode == 401 && Session.IsAuthenticated()) {
            Logout();
            SetError(QStringLiteral("登录已过期，请重新登录"));
            return;
        }
        if (!IsSuccessful) {
            SetError(QStringLiteral("计划保存失败，请检查网络后重试"));
            return;
        }
        LoadPlans();
    });
}

void ClientApi::Login(const QString& Email, const QString& Password) {
    QJsonObject Body;
    Body.insert(QStringLiteral("email"), Email);
    Body.insert(QStringLiteral("password"), Password);
    const std::uint64_t Generation = RequestGeneration;
    const std::uint64_t AuthenticationGenerationValue = ++AuthenticationGeneration;
    QNetworkReply* Reply = Transport.Post(
        QStringLiteral("/api/v1/auth/login"), QJsonDocument(Body).toJson(QJsonDocument::Compact));
    AuthenticationReply = Reply;
    connect(Reply, &QNetworkReply::finished, this,
        [this, Reply, Generation, AuthenticationGenerationValue]() {
        HandleAuthenticationReply(
            Reply, false, false, Generation, AuthenticationGenerationValue);
    });
}

void ClientApi::RefreshSession() {
    if (!Session.HasRefreshToken()) {
        SetError(QStringLiteral("没有可用的登录会话，请重新登录"));
        return;
    }
    QJsonObject Body;
    Body.insert(QStringLiteral("refreshToken"), Session.RefreshToken());
    const std::uint64_t Generation = RequestGeneration;
    const std::uint64_t AuthenticationGenerationValue = ++AuthenticationGeneration;
    QNetworkReply* Reply = Transport.Post(
        QStringLiteral("/api/v1/auth/refresh"), QJsonDocument(Body).toJson(QJsonDocument::Compact));
    AuthenticationReply = Reply;
    connect(Reply, &QNetworkReply::finished, this,
        [this, Reply, Generation, AuthenticationGenerationValue]() {
        HandleAuthenticationReply(
            Reply, false, true, Generation, AuthenticationGenerationValue);
    });
}

void ClientApi::Register(
    const QString& Email,
    const QString& Password,
    const QString& DisplayNameValue) {
    QJsonObject Body;
    Body.insert(QStringLiteral("email"), Email);
    Body.insert(QStringLiteral("password"), Password);
    Body.insert(QStringLiteral("displayName"), DisplayNameValue);
    const std::uint64_t Generation = RequestGeneration;
    const std::uint64_t AuthenticationGenerationValue = ++AuthenticationGeneration;
    QNetworkReply* Reply = Transport.Post(
        QStringLiteral("/api/v1/auth/register"), QJsonDocument(Body).toJson(QJsonDocument::Compact));
    AuthenticationReply = Reply;
    connect(Reply, &QNetworkReply::finished, this,
        [this, Reply, Generation, AuthenticationGenerationValue]() {
        HandleAuthenticationReply(
            Reply, true, false, Generation, AuthenticationGenerationValue);
    });
}

void ClientApi::RequestData(DataKind Kind) {
    const QString Resource = Kind == DataKind::Recipes
        ? QStringLiteral("/api/v1/recipes?limit=20")
        : Kind == DataKind::Ingredients
              ? QStringLiteral("/api/v1/ingredients")
              : Kind == DataKind::Recommendations
                    ? QStringLiteral("/api/v1/recommendations/tonight")
                    : QStringLiteral("/api/v1/plans?from=") +
                          QDate::currentDate().toString(Qt::ISODate) +
                          QStringLiteral("&to=") +
                          QDate::currentDate().addDays(6).toString(Qt::ISODate);
    QNetworkReply* Reply = nullptr;
    if (Kind == DataKind::Recommendations) {
        Reply = Transport.Post(Resource, RecommendationRequestBody());
    } else {
        Reply = Transport.Get(Resource);
    }
    const std::uint64_t Generation = RequestGeneration;
    connect(Reply, &QNetworkReply::finished, this, [this, Kind, Reply, Generation]() {
        HandleReply(Kind, Reply, Generation);
    });
}

QByteArray ClientApi::RecommendationRequestBody() const {
    const auto ToJsonArray = [](const QStringList& Values) {
        QJsonArray Result;
        for (const QString& Value : Values) {
            Result.append(Value);
        }
        return Result;
    };
    QJsonObject Body;
    Body.insert(QStringLiteral("servings"), Preferences.ServingCount());
    Body.insert(QStringLiteral("availableMinutes"), Preferences.AvailableMinutes());
    Body.insert(QStringLiteral("cuisines"), ToJsonArray(Preferences.PreferredCuisines()));
    Body.insert(QStringLiteral("allergies"), ToJsonArray(Preferences.Allergies()));
    Body.insert(QStringLiteral("pantryIngredientIds"),
        ToJsonArray(Preferences.PantryIngredientIds()));
    Body.insert(QStringLiteral("cookware"), ToJsonArray(Preferences.Cookware()));
    Body.insert(QStringLiteral("recentRecipeIds"), QJsonArray());
    return QJsonDocument(Body).toJson(QJsonDocument::Compact);
}

}  // namespace Menu::Client
