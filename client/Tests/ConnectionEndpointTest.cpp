#include <QtTest/QtTest>

#include <ConnectionEndpoint.hpp>

class ConnectionEndpointTest final : public QObject {
    Q_OBJECT

private slots:
    void ParsesIpv4Endpoint();
    void ParsesIpv6Endpoint();
    void RejectsUnsafeOrIncompleteEndpoint();
    void UsesDistinctStableScopeKeys();
};

void ConnectionEndpointTest::ParsesIpv4Endpoint() {
    const auto Result = Menu::Client::ParseEndpoint(
        QStringLiteral("http://192.168.1.20:8080/"));
    QVERIFY(Result.IsValid);
    QCOMPARE(Result.Endpoint.Scheme, QStringLiteral("http"));
    QCOMPARE(Result.Endpoint.Host, QStringLiteral("192.168.1.20"));
    QCOMPARE(Result.Endpoint.Port, quint16(8080));
    QCOMPARE(Result.Endpoint.ToUrl(), QStringLiteral("http://192.168.1.20:8080"));
}

void ConnectionEndpointTest::ParsesIpv6Endpoint() {
    const auto Result = Menu::Client::ParseEndpoint(
        QStringLiteral("https://[fe80::1]:9443/api"));
    QVERIFY(Result.IsValid);
    QCOMPARE(Result.Endpoint.Host, QStringLiteral("fe80::1"));
    QCOMPARE(Result.Endpoint.Port, quint16(9443));
    QCOMPARE(Result.Endpoint.Path, QStringLiteral("/api"));
    QCOMPARE(Result.Endpoint.ToUrl(), QStringLiteral("https://[fe80::1]:9443/api"));
}

void ConnectionEndpointTest::RejectsUnsafeOrIncompleteEndpoint() {
    for (const QString& Value : {
             QStringLiteral("192.168.1.20:8080"),
             QStringLiteral("ftp://192.168.1.20:8080"),
             QStringLiteral("http://:8080"),
             QStringLiteral("http://192.168.1.20:0"),
             QStringLiteral("http://192.168.1.20:65536"),
             QStringLiteral("http://user:password@192.168.1.20:8080"),
             QStringLiteral("http://192.168.1.20:8080?token=secret")}) {
        const auto Result = Menu::Client::ParseEndpoint(Value);
        QVERIFY2(!Result.IsValid, qPrintable(Value));
        QVERIFY(!Result.Error.isEmpty());
    }
}

void ConnectionEndpointTest::UsesDistinctStableScopeKeys() {
    const QString First = Menu::Client::EndpointScopeKey(
        QStringLiteral("http://127.0.0.1:8080"));
    const QString Same = Menu::Client::EndpointScopeKey(
        QStringLiteral("http://127.0.0.1:8080/"));
    const QString Other = Menu::Client::EndpointScopeKey(
        QStringLiteral("http://127.0.0.1:8081"));
    QVERIFY(!First.isEmpty());
    QCOMPARE(First, Same);
    QVERIFY(First != Other);
}

QTEST_GUILESS_MAIN(ConnectionEndpointTest)
#include "ConnectionEndpointTest.moc"
