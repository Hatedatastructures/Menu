#pragma once

#include <QObject>
#include <QList>
#include <QString>

class QNetworkAccessManager;
class QNetworkReply;

namespace Menu::Client {

class ClientTransport final : public QObject {
    Q_OBJECT

public:
    explicit ClientTransport(QObject* Parent = nullptr);
    ~ClientTransport() override;

    ClientTransport(const ClientTransport&) = delete;
    ClientTransport& operator=(const ClientTransport&) = delete;

    void SetBaseUrl(QString BaseUrlValue);
    void SetAccessToken(QString AccessTokenValue);

    [[nodiscard]] QNetworkReply* Get(const QString& Resource);
    [[nodiscard]] QNetworkReply* GetAtBaseUrl(
        const QString& BaseUrlValue,
        const QString& Resource);
    [[nodiscard]] QNetworkReply* Post(
        const QString& Resource,
        const QByteArray& Body);
    [[nodiscard]] QNetworkReply* Custom(
        const QByteArray& Method,
        const QString& Resource,
        const QByteArray& Body);

    void AbortAll() noexcept;

private:
    [[nodiscard]] QNetworkReply* Send(
        const QByteArray& Method,
        const QString& BaseUrlValue,
        const QString& Resource,
        const QByteArray& Body,
        bool IncludeAccessToken);
    void Track(QNetworkReply* Reply);

    QString BaseUrlValue;
    QString AccessTokenValue;
    QNetworkAccessManager* Network = nullptr;
    QList<QNetworkReply*> ActiveReplies;
};

}  // namespace Menu::Client
