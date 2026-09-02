#include "ClientTransport.hpp"

#include "ClientApiSupport.hpp"

#include <QNetworkAccessManager>
#include <QNetworkReply>
#include <QNetworkRequest>
#include <QTimer>
#include <QUrl>

#include <utility>

namespace Menu::Client {

ClientTransport::ClientTransport(QObject* Parent)
    : QObject(Parent), Network(new QNetworkAccessManager(this)) {}

ClientTransport::~ClientTransport() {
    AbortAll();
}

void ClientTransport::SetBaseUrl(QString BaseUrlInput) {
    BaseUrlValue = std::move(BaseUrlInput);
}

void ClientTransport::SetAccessToken(QString AccessTokenInput) {
    AccessTokenValue = std::move(AccessTokenInput);
}

QNetworkReply* ClientTransport::Get(const QString& Resource) {
    return Send("GET", BaseUrlValue, Resource, {}, true);
}

QNetworkReply* ClientTransport::GetAtBaseUrl(
    const QString& BaseUrlInput,
    const QString& Resource) {
    return Send("GET", BaseUrlInput, Resource, {}, false);
}

QNetworkReply* ClientTransport::Post(
    const QString& Resource,
    const QByteArray& Body) {
    return Send("POST", BaseUrlValue, Resource, Body, true);
}

QNetworkReply* ClientTransport::Custom(
    const QByteArray& Method,
    const QString& Resource,
    const QByteArray& Body) {
    return Send(Method, BaseUrlValue, Resource, Body, true);
}

QNetworkReply* ClientTransport::Send(
    const QByteArray& Method,
    const QString& BaseUrlInput,
    const QString& Resource,
    const QByteArray& Body,
    bool IncludeAccessToken) {
    QNetworkRequest Request(QUrl(BaseUrlInput + Resource));
    Request.setRawHeader("Accept", QByteArrayLiteral("application/json"));
    Request.setHeader(QNetworkRequest::ContentTypeHeader, "application/json");
    if (IncludeAccessToken && !AccessTokenValue.isEmpty()) {
        Request.setRawHeader(
            "Authorization", QByteArrayLiteral("Bearer ") + AccessTokenValue.toUtf8());
    }

    QNetworkReply* Reply = nullptr;
    if (Method == QByteArrayLiteral("GET")) {
        Reply = Network->get(Request);
    } else if (Method == QByteArrayLiteral("POST")) {
        Reply = Network->post(Request, Body);
    } else {
        Reply = Network->sendCustomRequest(Request, Method, Body);
    }
    Track(Reply);
    return Reply;
}

void ClientTransport::Track(QNetworkReply* Reply) {
    if (Reply == nullptr) {
        return;
    }
    ActiveReplies.push_back(Reply);
    auto* TimeoutTimer = new QTimer(Reply);
    TimeoutTimer->setSingleShot(true);
    TimeoutTimer->setInterval(Support::RequestTimeoutMilliseconds);
    connect(TimeoutTimer, &QTimer::timeout, Reply, &QNetworkReply::abort);
    connect(Reply, &QNetworkReply::finished, TimeoutTimer, &QTimer::stop);
    connect(Reply, &QNetworkReply::finished, this, [this, Reply]() {
        ActiveReplies.removeAll(Reply);
    });
    connect(Reply, &QObject::destroyed, this, [this, Reply]() {
        ActiveReplies.removeAll(Reply);
    });
    TimeoutTimer->start();
}

void ClientTransport::AbortAll() noexcept {
    const QList<QNetworkReply*> Replies = ActiveReplies;
    for (QNetworkReply* Reply : Replies) {
        if (Reply != nullptr) {
            Reply->abort();
        }
    }
    ActiveReplies.clear();
}

}  // namespace Menu::Client
