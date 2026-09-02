#include "ConnectionEndpoint.hpp"

#include <QCryptographicHash>
#include <QUrl>

#include <utility>

namespace Menu::Client {
namespace {

EndpointParseResult Invalid(QString Error) {
    EndpointParseResult Result;
    Result.Error = std::move(Error);
    return Result;
}

}  // namespace

QString ConnectionEndpoint::ToUrl() const {
    QUrl Url;
    Url.setScheme(Scheme.toLower());
    Url.setHost(Host);
    Url.setPort(static_cast<int>(Port));
    const QString NormalizedPath = Path.isEmpty() ? QStringLiteral("/") : Path;
    Url.setPath(NormalizedPath == QStringLiteral("/") ? QString() : NormalizedPath);
    QString Value = Url.toString(QUrl::FullyEncoded);
    while (Value.endsWith('/')) {
        Value.chop(1);
    }
    return Value;
}

EndpointParseResult ParseEndpoint(const QString& Value) {
    const QString Input = Value.trimmed();
    if (Input.isEmpty() || Input.size() > 512) {
        return Invalid(QStringLiteral("服务器地址不能为空且不能超过 512 个字符"));
    }

    const QUrl Url(Input);
    const QString Scheme = Url.scheme().toLower();
    if (!Url.isValid() || (Scheme != QStringLiteral("http") &&
                          Scheme != QStringLiteral("https"))) {
        return Invalid(QStringLiteral("仅支持 http 或 https 地址"));
    }
    if (Url.host().isEmpty() || !Url.userName().isEmpty() ||
        !Url.password().isEmpty() || !Url.query().isEmpty() ||
        !Url.fragment().isEmpty()) {
        return Invalid(QStringLiteral("服务器地址必须只包含协议、主机、端口和路径"));
    }

    const int ExplicitPort = Url.port(-1);
    const int Port = ExplicitPort < 0
        ? (Scheme == QStringLiteral("https") ? 443 : 80)
        : ExplicitPort;
    if (Port < 1 || Port > 65535) {
        return Invalid(QStringLiteral("端口必须在 1 到 65535 之间"));
    }

    ConnectionEndpoint Endpoint;
    Endpoint.Scheme = Scheme;
    Endpoint.Host = Url.host(QUrl::FullyDecoded);
    Endpoint.Port = static_cast<quint16>(Port);
    Endpoint.Path = Url.path(QUrl::FullyEncoded);
    if (Endpoint.Path.isEmpty()) {
        Endpoint.Path = QStringLiteral("/");
    }
    while (Endpoint.Path.size() > 1 && Endpoint.Path.endsWith('/')) {
        Endpoint.Path.chop(1);
    }

    EndpointParseResult Result;
    Result.IsValid = true;
    Result.Endpoint = std::move(Endpoint);
    return Result;
}

QString EndpointScopeKey(const QString& Value) {
    const EndpointParseResult Result = ParseEndpoint(Value);
    const QString Canonical = Result.IsValid
        ? Result.Endpoint.ToUrl()
        : Value.trimmed().toLower();
    return QString::fromLatin1(QCryptographicHash::hash(
        Canonical.toUtf8(), QCryptographicHash::Sha256).toHex());
}

QString DefaultBaseUrl() {
#if defined(Q_OS_ANDROID)
    return QStringLiteral("http://10.0.2.2:8080");
#else
    return QStringLiteral("http://127.0.0.1:8080");
#endif
}

}  // namespace Menu::Client
