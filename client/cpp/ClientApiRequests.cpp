#include "ClientApi.hpp"

#include "ClientApiSupport.hpp"

#include <QDate>
#include <QJsonArray>
#include <QJsonDocument>
#include <QJsonObject>
#include <QNetworkAccessManager>
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
        Item.insert(QStringLiteral("servings"), ServingCountValue);
        Item.insert(QStringLiteral("sortOrder"), SortOrder++);
        Items.append(Item);
    }
    QJsonObject Body;
    Body.insert(QStringLiteral("planDate"), PlanDate);
    Body.insert(QStringLiteral("items"), Items);
    const std::uint64_t Generation = RequestGeneration;
    QNetworkReply* Reply = NetworkManager->post(
        CreateRequest(QStringLiteral("/api/v1/plans")),
        QJsonDocument(Body).toJson(QJsonDocument::Compact));
    ConfigureReply(Reply);
    connect(Reply, &QNetworkReply::finished, this, [this, Reply, Generation]() {
        const int StatusCode = Reply->attribute(
            QNetworkRequest::HttpStatusCodeAttribute).toInt();
        const bool IsSuccessful = Reply->error() == QNetworkReply::NoError &&
            StatusCode >= 200 && StatusCode < 300;
        ActiveReplies.removeAll(Reply);
        Reply->deleteLater();
        if (Generation != RequestGeneration) {
            return;
        }
        if (StatusCode == 401 && !AccessTokenValue.isEmpty()) {
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
    QNetworkReply* Reply = NetworkManager->post(
        CreateRequest(QStringLiteral("/api/v1/auth/login")),
        QJsonDocument(Body).toJson(QJsonDocument::Compact));
    ConfigureReply(Reply);
    AuthenticationReply = Reply;
    connect(Reply, &QNetworkReply::finished, this, [this, Reply, Generation]() {
        HandleAuthenticationReply(Reply, false, false, Generation);
    });
}

void ClientApi::RefreshSession() {
    if (RefreshTokenValue.isEmpty()) {
        SetError(QStringLiteral("没有可用的登录会话，请重新登录"));
        return;
    }
    QJsonObject Body;
    Body.insert(QStringLiteral("refreshToken"), RefreshTokenValue);
    const std::uint64_t Generation = RequestGeneration;
    QNetworkReply* Reply = NetworkManager->post(
        CreateRequest(QStringLiteral("/api/v1/auth/refresh")),
        QJsonDocument(Body).toJson(QJsonDocument::Compact));
    ConfigureReply(Reply);
    AuthenticationReply = Reply;
    connect(Reply, &QNetworkReply::finished, this, [this, Reply, Generation]() {
        HandleAuthenticationReply(Reply, false, true, Generation);
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
    QNetworkReply* Reply = NetworkManager->post(
        CreateRequest(QStringLiteral("/api/v1/auth/register")),
        QJsonDocument(Body).toJson(QJsonDocument::Compact));
    ConfigureReply(Reply);
    AuthenticationReply = Reply;
    connect(Reply, &QNetworkReply::finished, this, [this, Reply, Generation]() {
        HandleAuthenticationReply(Reply, true, false, Generation);
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
    QNetworkRequest Request = CreateRequest(Resource);
    QNetworkReply* Reply = nullptr;
    if (Kind == DataKind::Recommendations) {
        Request.setHeader(QNetworkRequest::ContentTypeHeader, "application/json");
        Reply = NetworkManager->post(Request, RecommendationRequestBody());
    } else {
        Reply = NetworkManager->get(Request);
    }
    ConfigureReply(Reply);
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
    Body.insert(QStringLiteral("servings"), ServingCountValue);
    Body.insert(QStringLiteral("availableMinutes"), AvailableMinutesValue);
    Body.insert(QStringLiteral("cuisines"), ToJsonArray(PreferredCuisinesValue));
    Body.insert(QStringLiteral("allergies"), ToJsonArray(AllergiesValue));
    Body.insert(QStringLiteral("pantryIngredientIds"),
        ToJsonArray(PantryIngredientIdsValue));
    Body.insert(QStringLiteral("cookware"), ToJsonArray(CookwareValue));
    Body.insert(QStringLiteral("recentRecipeIds"), QJsonArray());
    return QJsonDocument(Body).toJson(QJsonDocument::Compact);
}

QNetworkRequest ClientApi::CreateRequest(const QString& Resource) const {
    QNetworkRequest Request(QUrl(ServiceUrl + Resource));
    Request.setRawHeader("Accept", "application/json");
    Request.setHeader(QNetworkRequest::ContentTypeHeader, "application/json");
    if (!AccessTokenValue.isEmpty()) {
        Request.setRawHeader(
            "Authorization", QByteArrayLiteral("Bearer ") + AccessTokenValue.toUtf8());
    }
    return Request;
}

void ClientApi::ConfigureReply(QNetworkReply* Reply) {
    ActiveReplies.push_back(Reply);
    auto* TimeoutTimer = new QTimer(Reply);
    TimeoutTimer->setSingleShot(true);
    TimeoutTimer->setInterval(Support::RequestTimeoutMilliseconds);
    connect(TimeoutTimer, &QTimer::timeout, Reply, &QNetworkReply::abort);
    connect(Reply, &QNetworkReply::finished, TimeoutTimer, &QTimer::stop);
    TimeoutTimer->start();
}

}  // namespace Menu::Client
