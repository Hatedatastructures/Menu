#include "ClientApi.hpp"

#include <QJsonArray>
#include <QJsonDocument>
#include <QJsonObject>
#include <QJsonParseError>
#include <QNetworkReply>
#include <QNetworkRequest>

namespace Menu::Client {

void ClientApi::OpenRecipe(int Index) {
    QVariantMap NextRecipe = RecommendationValues.RecommendationAt(Index);
    if (NextRecipe.isEmpty()) {
        NextRecipe = RecipeValues.RecipeAt(Index);
    }
    if (NextRecipe == ActiveRecipeValue) {
        return;
    }
    ActiveRecipeValue = NextRecipe;
    emit ActiveRecipeChanged();
}

void ClientApi::CloseRecipe() {
    if (ActiveRecipeValue.isEmpty()) {
        return;
    }
    ActiveRecipeValue.clear();
    emit ActiveRecipeChanged();
}

void ClientApi::StartCookingSession() {
    if (!IsAuthenticated()) {
        SetError(QStringLiteral("当前仅在本机保存做饭页面，登录后可同步进度"));
        return;
    }
    const QString RecipeId = ActiveRecipeValue.value(QStringLiteral("id")).toString();
    if (RecipeId.isEmpty()) {
        SetError(QStringLiteral("请先选择一道菜"));
        return;
    }
    QJsonObject Body;
    Body.insert(QStringLiteral("recipeId"), RecipeId);
    const std::uint64_t Generation = RequestGeneration;
    QNetworkReply* Reply = NetworkManager->post(
        CreateRequest(QStringLiteral("/api/v1/cooking-sessions")),
        QJsonDocument(Body).toJson(QJsonDocument::Compact));
    ConfigureReply(Reply);
    CookingSessionReply = Reply;
    connect(Reply, &QNetworkReply::finished, this, [this, Reply, Generation]() {
        HandleCookingSessionReply(Reply, Generation);
    });
}

void ClientApi::UpdateCookingSession(int CurrentStepOrder, const QString& State) {
    if (!IsAuthenticated() ||
        CurrentCookingSessionValue.value(QStringLiteral("id")).toString().isEmpty()) {
        return;
    }
    const QString SessionId = CurrentCookingSessionValue.value(
        QStringLiteral("id")).toString();
    QJsonObject Body;
    Body.insert(QStringLiteral("currentStepOrder"), CurrentStepOrder);
    Body.insert(QStringLiteral("state"), State);
    const std::uint64_t Generation = RequestGeneration;
    QNetworkReply* Reply = NetworkManager->sendCustomRequest(
        CreateRequest(QStringLiteral("/api/v1/cooking-sessions/") + SessionId),
        QByteArrayLiteral("PATCH"),
        QJsonDocument(Body).toJson(QJsonDocument::Compact));
    ConfigureReply(Reply);
    CookingSessionReply = Reply;
    connect(Reply, &QNetworkReply::finished, this, [this, Reply, Generation]() {
        HandleCookingSessionReply(Reply, Generation);
    });
}

void ClientApi::SubmitFeedback(
    const QString& Outcome,
    const QVariantList& Tags,
    const QString& Comment) {
    if (!IsAuthenticated()) {
        SetError(QStringLiteral("登录后才能保存反馈"));
        return;
    }
    QJsonArray TagValues;
    for (const QVariant& Value : Tags) {
        TagValues.append(Value.toString());
    }
    QJsonObject Body;
    Body.insert(
        QStringLiteral("recipeId"),
        ActiveRecipeValue.value(QStringLiteral("id")).toString());
    Body.insert(QStringLiteral("outcome"), Outcome);
    Body.insert(QStringLiteral("tags"), TagValues);
    Body.insert(QStringLiteral("comment"), Comment);
    const std::uint64_t Generation = RequestGeneration;
    QNetworkReply* Reply = NetworkManager->post(
        CreateRequest(QStringLiteral("/api/v1/feedback")),
        QJsonDocument(Body).toJson(QJsonDocument::Compact));
    ConfigureReply(Reply);
    FeedbackReply = Reply;
    connect(Reply, &QNetworkReply::finished, this, [this, Reply, Generation]() {
        HandleFeedbackReply(Reply, Generation);
    });
}

void ClientApi::HandleCookingSessionReply(
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
    if (StatusCode == 401 && !AccessTokenValue.isEmpty()) {
        Logout();
        SetError(QStringLiteral("登录已过期，请重新登录"));
        return;
    }
    if (!IsSuccessful) {
        SetError(QStringLiteral("做饭进度同步失败，当前页面仍可继续"));
        return;
    }
    QJsonParseError ParseError;
    const QJsonDocument Document = QJsonDocument::fromJson(Data, &ParseError);
    if (ParseError.error != QJsonParseError::NoError || !Document.isObject()) {
        SetError(QStringLiteral("做饭进度响应无效"));
        return;
    }
    CurrentCookingSessionValue = Document.object().toVariantMap();
    emit CookingSessionChanged();
}

void ClientApi::HandleFeedbackReply(
    QNetworkReply* Reply,
    std::uint64_t RequestGenerationValue) {
    const int StatusCode = Reply->attribute(
        QNetworkRequest::HttpStatusCodeAttribute).toInt();
    const bool IsSuccessful = Reply->error() == QNetworkReply::NoError &&
        StatusCode >= 200 && StatusCode < 300;
    ActiveReplies.removeAll(Reply);
    Reply->deleteLater();
    if (RequestGenerationValue != RequestGeneration) {
        return;
    }
    if (StatusCode == 401 && !AccessTokenValue.isEmpty()) {
        Logout();
        SetError(QStringLiteral("登录已过期，请重新登录"));
        return;
    }
    if (!IsSuccessful) {
        SetError(QStringLiteral("反馈保存失败，请稍后重试"));
    }
}

}  // namespace Menu::Client
