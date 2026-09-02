#include "ClientApi.hpp"

#include "ConnectionEndpoint.hpp"
#include "ConnectionSettingsStore.hpp"

#include <QNetworkReply>
#include <QNetworkRequest>
#include <QThreadPool>

namespace Menu::Client {
bool ClientApi::IsTestingConnection() const noexcept {
    return TestingConnection;
}

bool ClientApi::ConnectionReachable() const noexcept {
    return Reachable;
}

QString ClientApi::ConnectionStatus() const {
    return ConnectionState;
}

bool ClientApi::ApplyBaseUrl(const QString& BaseUrlValue) {
    const EndpointParseResult Parsed = ParseEndpoint(BaseUrlValue);
    if (!Parsed.IsValid) {
        SetError(Parsed.Error);
        ConnectionState = Parsed.Error;
        Reachable = false;
        emit ConnectionTestChanged();
        return false;
    }

    const QString NormalizedValue = Parsed.Endpoint.ToUrl();
    if (NormalizedValue == ServiceUrl) {
        return true;
    }

    if (ConnectionTestReply != nullptr) {
        ConnectionTestReply->abort();
        ConnectionTestReply = nullptr;
    }
    ++ConnectionTestGeneration;
    TestingConnection = false;
    Reachable = false;
    ConnectionState = QStringLiteral("未测试");

    const QString PreviousUserId = Session.UserId();
    const bool WasAuthenticated = IsAuthenticated();
    ResetPrivateState(PreviousUserId, false);
    Session.Clear();
    Transport.SetAccessToken({});
    ServiceUrl = NormalizedValue;
    Transport.SetBaseUrl(ServiceUrl);
    PersistConnectionSettings();
    if (WasAuthenticated) {
        emit AuthenticationChanged();
    }
    emit baseUrlChanged();
    emit ConnectionTestChanged();
    return true;
}

void ClientApi::TestBaseUrl(const QString& BaseUrlValue) {
    const QString Candidate = BaseUrlValue.trimmed().isEmpty()
        ? ServiceUrl
        : BaseUrlValue;
    const EndpointParseResult Parsed = ParseEndpoint(Candidate);
    if (!Parsed.IsValid) {
        ConnectionState = Parsed.Error;
        TestingConnection = false;
        Reachable = false;
        emit ConnectionTestChanged();
        return;
    }

    if (ConnectionTestReply != nullptr) {
        ConnectionTestReply->abort();
        ConnectionTestReply = nullptr;
    }
    const std::uint64_t Generation = ++ConnectionTestGeneration;
    TestingConnection = true;
    Reachable = false;
    ConnectionState = QStringLiteral("正在测试连接…");
    emit ConnectionTestChanged();

    QNetworkReply* Reply = Transport.GetAtBaseUrl(
        Parsed.Endpoint.ToUrl(), QStringLiteral("/healthz"));
    ConnectionTestReply = Reply;
    connect(Reply, &QNetworkReply::finished, this, [this, Reply, Generation]() {
        const int StatusCode = Reply->attribute(
            QNetworkRequest::HttpStatusCodeAttribute).toInt();
        const bool Success = Reply->error() == QNetworkReply::NoError &&
            StatusCode >= 200 && StatusCode < 300;
        Reply->deleteLater();
        if (Generation != ConnectionTestGeneration) {
            return;
        }
        ConnectionTestReply = nullptr;
        TestingConnection = false;
        Reachable = Success;
        ConnectionState = Success
            ? QStringLiteral("连接成功")
            : QStringLiteral("连接失败，请检查地址、端口和服务端状态");
        emit ConnectionTestChanged();
    });
}

void ClientApi::LoadConnectionSettings() {
    const EndpointParseResult Parsed = ParseEndpoint(
        ConnectionSettingsStore::LoadBaseUrl(CacheDirectory));
    if (Parsed.IsValid) {
        ServiceUrl = Parsed.Endpoint.ToUrl();
    }
}

void ClientApi::PersistConnectionSettings() const {
    const QString Directory = CacheDirectory;
    const QString BaseUrl = ServiceUrl;
    const auto Generation = ConnectionGeneration;
    const std::uint64_t GenerationValue =
        Generation->fetch_add(1, std::memory_order_acq_rel) + 1;
    QThreadPool::globalInstance()->start(
        [Directory, BaseUrl, Generation, GenerationValue]() {
            if (Generation->load(std::memory_order_acquire) == GenerationValue) {
                (void)ConnectionSettingsStore::SaveBaseUrl(Directory, BaseUrl);
            }
        });
}

}  // namespace Menu::Client
