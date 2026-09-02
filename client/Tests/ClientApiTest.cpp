#include <QtTest/QtTest>

#include <ClientApi.hpp>
#include <ConnectionEndpoint.hpp>

#include <QFile>
#include <QDir>
#include <QJsonDocument>
#include <QJsonObject>
#include <QSignalSpy>
#include <QStringList>
#include <QTcpServer>
#include <QTcpSocket>
#include <QTemporaryDir>
#include <QThread>
#include <QThreadPool>

namespace {

QString ScopedCacheDirectory(const QString& CachePath, const QString& BaseUrl) {
    return QDir(CachePath).filePath(
        QStringLiteral("server-") + Menu::Client::EndpointScopeKey(BaseUrl));
}

QString ScopedCachePath(
    const QString& CachePath,
    const QString& BaseUrl,
    const QString& FileName) {
    return QDir(ScopedCacheDirectory(CachePath, BaseUrl)).filePath(FileName);
}

bool PreferencesMatch(
    const QString& Path,
    int ServingCount,
    const QStringList& Cuisines,
    const QStringList& Allergies,
    int AvailableMinutes,
    const QStringList& PantryIngredientIds,
    bool RemindersEnabled) {
    QFile File(Path);
    if (!File.open(QIODevice::ReadOnly)) {
        return false;
    }
    const QJsonObject Object = QJsonDocument::fromJson(File.readAll()).object();
    const QJsonArray CuisineValues = Object.value(QStringLiteral("cuisines")).toArray();
    const QJsonArray AllergyValues = Object.value(QStringLiteral("allergies")).toArray();
    const QJsonArray PantryValues =
        Object.value(QStringLiteral("pantryIngredientIds")).toArray();
    if (Object.value(QStringLiteral("schemaVersion")).toInt() != 2 ||
        Object.value(QStringLiteral("kind")).toString() !=
            QStringLiteral("preferences") ||
        Object.value(QStringLiteral("servings")).toInt() != ServingCount ||
        Object.value(QStringLiteral("availableMinutes")).toInt() != AvailableMinutes ||
        Object.value(QStringLiteral("remindersEnabled")).toBool() != RemindersEnabled ||
        CuisineValues.size() != Cuisines.size() ||
        AllergyValues.size() != Allergies.size() ||
        PantryValues.size() != PantryIngredientIds.size()) {
        return false;
    }
    for (int Index = 0; Index < Cuisines.size(); ++Index) {
        if (CuisineValues.at(Index).toString() != Cuisines.at(Index)) {
            return false;
        }
    }
    for (int Index = 0; Index < Allergies.size(); ++Index) {
        if (AllergyValues.at(Index).toString() != Allergies.at(Index)) {
            return false;
        }
    }
    for (int Index = 0; Index < PantryIngredientIds.size(); ++Index) {
        if (PantryValues.at(Index).toString() != PantryIngredientIds.at(Index)) {
            return false;
        }
    }
    return true;
}

QByteArray HttpResponse(const QByteArray& Body, int StatusCode = 200) {
    const QByteArray StatusText = StatusCode == 401 ? QByteArrayLiteral("Unauthorized")
                                                    : QByteArrayLiteral("OK");
    return QByteArrayLiteral("HTTP/1.1 ") + QByteArray::number(StatusCode) +
        QByteArrayLiteral(" ") + StatusText +
        QByteArrayLiteral("\r\nContent-Type: application/json\r\n") +
        QByteArrayLiteral("Content-Length: ") + QByteArray::number(Body.size()) +
        QByteArrayLiteral("\r\nConnection: close\r\n\r\n") + Body;
}

class FixtureServer final : public QObject {
    Q_OBJECT

public:
    explicit FixtureServer(QObject* Parent = nullptr)
        : QObject(Parent) {
        connect(&Server, &QTcpServer::newConnection, this, &FixtureServer::SendResponse);
    }

    bool Listen() {
        return Server.listen(QHostAddress::LocalHost);
    }

    [[nodiscard]] quint16 Port() const {
        return Server.serverPort();
    }

    void Stop() {
        Server.close();
    }

    [[nodiscard]] QByteArray RecommendationRequest() const {
        return RecommendationRequestValue;
    }

    void SetPlansUnauthorized() {
        PlansUnauthorized = true;
    }

private slots:
    void SendResponse() {
        while (Server.hasPendingConnections()) {
            QTcpSocket* Socket = Server.nextPendingConnection();
            connect(Socket, &QTcpSocket::readyRead, this, [this, Socket]() {
                const QByteArray Request = Socket->readAll();
                if (Request.startsWith("POST /api/v1/recommendations/tonight")) {
                    RecommendationRequestValue = Request;
                }
                const bool Unauthorized = PlansUnauthorized &&
                    Request.startsWith("GET /api/v1/plans");
                const QByteArray Body = Request.startsWith(
                                            "POST /api/v1/auth/login")
                    ? QByteArrayLiteral(
                          R"json({"user":{"id":"user.fixture","email":"cook@example.com","displayName":"厨房","isAdmin":false},"accessToken":"access-fixture","refreshToken":"refresh-fixture","accessTokenExpiresInSeconds":900})json")
                    : Request.startsWith("POST /api/v1/auth/refresh")
                          ? QByteArrayLiteral(
                                R"json({"user":{"id":"user.fixture","email":"cook@example.com","displayName":"厨房","isAdmin":false},"accessToken":"access-refreshed","refreshToken":"refresh-refreshed","accessTokenExpiresInSeconds":900})json")
                    : Request.startsWith("GET /api/v1/plans")
                          ? QByteArrayLiteral(
                                R"json([{"id":"plan.one","planDate":"2026-09-01","items":[{"recipeId":"recipe.tomato","recipeName":"番茄炒蛋","servings":2}],"combinedIngredients":[{"ingredientId":"ingredient.tomato","ingredientName":"番茄","category":"蔬菜","defaultUnit":"g","quantity":500,"unit":"g","required":true}]}])json")
                    : Request.startsWith(
                                            "POST /api/v1/recommendations/tonight")
                    ? QByteArrayLiteral(
                          R"json([{"recipe":{"id":"recipe.tomato","name":"番茄炒蛋","cuisine":"中餐","imagePath":"assets/media/tomato-egg.png","totalMinutes":20,"servings":2,"difficulty":1,"ingredients":[{"ingredientId":"ingredient.tomato","ingredientName":"番茄","category":"蔬菜","defaultUnit":"g","isPantryStaple":false,"quantity":200,"unit":"g","required":true}],"steps":[]},"availableIngredientCount":1,"missingIngredientCount":0,"score":10}])json")
                    : Request.startsWith("GET /api/v1/ingredients")
                          ? QByteArrayLiteral(
                                R"json([{"id":"ingredient.tomato","name":"番茄","category":"蔬菜","defaultUnit":"g","isPantryStaple":false,"aliases":["西红柿"],"substituteGroup":""}])json")
                          : QByteArrayLiteral(
                                R"json([{"id":"recipe.tomato","name":"番茄炒蛋","cuisine":"中餐","totalMinutes":20,"servings":2,"difficulty":1,"ingredients":[{"ingredientId":"ingredient.tomato","ingredientName":"番茄","category":"蔬菜","defaultUnit":"g","isPantryStaple":false,"quantity":200,"unit":"g","required":true}],"steps":[]}])json");
                Socket->write(HttpResponse(
                    Unauthorized ? QByteArrayLiteral("{\"error\":{\"code\":\"authentication_required\"}}") : Body,
                    Unauthorized ? 401 : 200));
                Socket->disconnectFromHost();
            });
        }
    }

private:
    QTcpServer Server;
    bool PlansUnauthorized = false;
    QByteArray RecommendationRequestValue;
};

class ThreadPoolLimitGuard final {
public:
    explicit ThreadPoolLimitGuard(QThreadPool* PoolValue)
        : Pool(PoolValue), OriginalMaximum(PoolValue->maxThreadCount()) {
        Pool->setMaxThreadCount(1);
    }

    ~ThreadPoolLimitGuard() {
        Pool->setMaxThreadCount(OriginalMaximum);
    }

private:
    QThreadPool* Pool;
    int OriginalMaximum;
};

class ThreadPoolDrainGuard final {
public:
    ~ThreadPoolDrainGuard() {
        QThreadPool::globalInstance()->waitForDone(5000);
    }
};

}  // namespace

class ClientApiTest final : public QObject {
    Q_OBJECT

private slots:
    void cleanup();
    void LoadsRemoteDataAndFallsBackToCache();
    void LogsInAndLoadsAuthenticatedPlans();
    void RefreshesAuthenticatedSession();
    void QueuesPublicRefreshBehindPlanLoad();
    void SendsConfiguredRecommendationPreferences();
    void PersistsPreferencesWithVersionedCache();
    void PersistsAndValidatesBaseUrl();
    void TestsBaseUrlConnectivity();
    void ClearsPrivateStateOnUnauthorizedResponse();
    void ClearsPrivateStateOnLogout();
    void ClearCacheBlocksStaleAsyncWrites();
    void CancelsPendingLoadsWithoutLateModelWrites();
    void DestroysWithInFlightRequests();
};

void ClientApiTest::cleanup() {
    QThreadPool::globalInstance()->waitForDone(5000);
}

void ClientApiTest::LoadsRemoteDataAndFallsBackToCache() {
    QTemporaryDir CacheDirectory;
    QVERIFY(CacheDirectory.isValid());
    ThreadPoolDrainGuard DrainGuard;

    FixtureServer Fixture;
    QVERIFY(Fixture.Listen());

    const QString BaseUrl = QStringLiteral("http://127.0.0.1:%1").arg(Fixture.Port());
    Menu::Client::ClientApi RemoteApi(BaseUrl, CacheDirectory.path());
    QSignalSpy RemoteReady(&RemoteApi, &Menu::Client::ClientApi::DataReady);
    RemoteApi.LoadData();
    QTRY_COMPARE_WITH_TIMEOUT(RemoteApi.Recipes()->rowCount(), 1, 5000);
    QTRY_COMPARE_WITH_TIMEOUT(RemoteApi.Ingredients()->rowCount(), 1, 5000);
    QTRY_COMPARE_WITH_TIMEOUT(RemoteApi.Recommendations()->rowCount(), 1, 5000);
    QVERIFY(!RemoteApi.IsOffline());
    QCOMPARE(
        RemoteApi.Recipes()
            ->data(RemoteApi.Recipes()->index(0), Menu::Client::RecipeModel::IngredientsRole)
            .toList()
            .at(0)
            .toMap()
            .value(QStringLiteral("ingredientName"))
            .toString(),
        QStringLiteral("番茄"));
    QVERIFY(RemoteReady.count() >= 1);
    QCOMPARE(
        RemoteApi.Recommendations()
            ->data(RemoteApi.Recommendations()->index(0),
                Menu::Client::RecommendationModel::NameRole)
            .toString(),
        QStringLiteral("番茄炒蛋"));
    QCOMPARE(
        RemoteApi.Recommendations()
            ->data(RemoteApi.Recommendations()->index(0),
                Menu::Client::RecommendationModel::MissingIngredientCountRole)
            .toInt(),
        0);

    QTRY_VERIFY_WITH_TIMEOUT(
        QFile::exists(ScopedCachePath(CacheDirectory.path(), BaseUrl, QStringLiteral("recipes.json"))), 5000);
    QTRY_VERIFY_WITH_TIMEOUT(
        QFile::exists(ScopedCachePath(CacheDirectory.path(), BaseUrl, QStringLiteral("ingredients.json"))), 5000);
    QTRY_VERIFY_WITH_TIMEOUT(
        QFile::exists(ScopedCachePath(CacheDirectory.path(), BaseUrl, QStringLiteral("tonight.json"))), 5000);
    QFile RecipeCache(ScopedCachePath(CacheDirectory.path(), BaseUrl, QStringLiteral("recipes.json")));
    QVERIFY(RecipeCache.open(QIODevice::ReadOnly));
    const QJsonObject RecipeEnvelope =
        QJsonDocument::fromJson(RecipeCache.readAll()).object();
    QCOMPARE(RecipeEnvelope.value(QStringLiteral("schemaVersion")).toInt(), 2);
    QCOMPARE(RecipeEnvelope.value(QStringLiteral("kind")).toString(),
        QStringLiteral("recipes"));
    QVERIFY(RecipeEnvelope.value(QStringLiteral("data")).isArray());

    Fixture.Stop();
    Menu::Client::ClientApi OfflineApi(
        BaseUrl,
        CacheDirectory.path());
    OfflineApi.LoadData();
    QTRY_COMPARE_WITH_TIMEOUT(OfflineApi.Recipes()->rowCount(), 1, 5000);
    QTRY_COMPARE_WITH_TIMEOUT(OfflineApi.Ingredients()->rowCount(), 1, 5000);
    QTRY_COMPARE_WITH_TIMEOUT(OfflineApi.Recommendations()->rowCount(), 1, 5000);
    QVERIFY(OfflineApi.IsOffline());
    QCOMPARE(OfflineApi.Recipes()->data(
                 OfflineApi.Recipes()->index(0), Menu::Client::RecipeModel::NameRole)
                 .toString(),
        QStringLiteral("番茄炒蛋"));

    QTemporaryDir InvalidCacheDirectory;
    QVERIFY(InvalidCacheDirectory.isValid());
    const QString InvalidBaseUrl = QStringLiteral("http://127.0.0.1:1");
    QVERIFY(QDir().mkpath(ScopedCacheDirectory(
        InvalidCacheDirectory.path(), InvalidBaseUrl)));
    QFile InvalidCache(ScopedCachePath(
        InvalidCacheDirectory.path(), InvalidBaseUrl, QStringLiteral("recipes.json")));
    QVERIFY(InvalidCache.open(QIODevice::WriteOnly));
    InvalidCache.write(QByteArrayLiteral(
        R"json({"schemaVersion":1,"kind":"recipes","data":[]})json"));
    InvalidCache.close();
    Menu::Client::ClientApi InvalidApi(
        QStringLiteral("http://127.0.0.1:1"), InvalidCacheDirectory.path());
    InvalidApi.LoadData();
    QTRY_VERIFY_WITH_TIMEOUT(!InvalidApi.IsLoading(), 5000);
    QCOMPARE(InvalidApi.Recipes()->rowCount(), 0);
    QVERIFY(InvalidApi.ErrorMessage().contains(QStringLiteral("没有可用的离线数据")));
}

void ClientApiTest::LogsInAndLoadsAuthenticatedPlans() {
    QTemporaryDir CacheDirectory;
    QVERIFY(CacheDirectory.isValid());
    ThreadPoolDrainGuard DrainGuard;
    FixtureServer Fixture;
    QVERIFY(Fixture.Listen());

    const QString BaseUrl = QStringLiteral("http://127.0.0.1:%1").arg(Fixture.Port());
    Menu::Client::ClientApi Api(BaseUrl, CacheDirectory.path());
    Api.Login(QStringLiteral("cook@example.com"), QStringLiteral("Menu-Cook-Password-2026"));
    QTRY_VERIFY_WITH_TIMEOUT(Api.IsAuthenticated(), 5000);
    QCOMPARE(Api.DisplayName(), QStringLiteral("厨房"));
    QTRY_COMPARE_WITH_TIMEOUT(Api.Recommendations()->rowCount(), 1, 5000);
    QTRY_COMPARE_WITH_TIMEOUT(Api.Plans()->rowCount(), 1, 5000);
    const QVariantMap Plan = Api.Plans()->PlanAt(0);
    QCOMPARE(Plan.value(QStringLiteral("planDate")).toString(), QStringLiteral("2026-09-01"));
    const QVariantMap Ingredient = Plan.value(QStringLiteral("combinedIngredients"))
                                       .toList()
                                       .at(0)
                                       .toMap();
    QCOMPARE(Ingredient.value(QStringLiteral("ingredientName")).toString(), QStringLiteral("番茄"));
    QTRY_VERIFY_WITH_TIMEOUT(
        !QDir(ScopedCacheDirectory(CacheDirectory.path(), BaseUrl)).entryList(
                QStringList{QStringLiteral("plans.*.json")}, QDir::Files).isEmpty(),
        5000);
}

void ClientApiTest::QueuesPublicRefreshBehindPlanLoad() {
    QTemporaryDir CacheDirectory;
    QVERIFY(CacheDirectory.isValid());
    ThreadPoolDrainGuard DrainGuard;
    FixtureServer Fixture;
    QVERIFY(Fixture.Listen());
    const QString BaseUrl = QStringLiteral("http://127.0.0.1:%1").arg(Fixture.Port());
    Menu::Client::ClientApi Api(BaseUrl, CacheDirectory.path());
    Api.Login(QStringLiteral("cook@example.com"), QStringLiteral("Menu-Cook-Password-2026"));
    QTRY_VERIFY_WITH_TIMEOUT(Api.IsAuthenticated(), 5000);
    Api.LoadPlans();
    Api.LoadData();

    QTRY_COMPARE_WITH_TIMEOUT(Api.Recommendations()->rowCount(), 1, 5000);
    QTRY_COMPARE_WITH_TIMEOUT(Api.Plans()->rowCount(), 1, 5000);
}

void ClientApiTest::RefreshesAuthenticatedSession() {
    QTemporaryDir CacheDirectory;
    QVERIFY(CacheDirectory.isValid());
    ThreadPoolDrainGuard DrainGuard;
    FixtureServer Fixture;
    QVERIFY(Fixture.Listen());

    const QString BaseUrl = QStringLiteral("http://127.0.0.1:%1").arg(Fixture.Port());
    Menu::Client::ClientApi Api(BaseUrl, CacheDirectory.path());
    Api.Login(QStringLiteral("cook@example.com"), QStringLiteral("Menu-Cook-Password-2026"));
    QTRY_VERIFY_WITH_TIMEOUT(Api.IsAuthenticated(), 5000);
    QVERIFY(QMetaObject::invokeMethod(&Api, "RefreshSession"));
    QTRY_COMPARE_WITH_TIMEOUT(Api.Recommendations()->rowCount(), 1, 5000);
    QTRY_COMPARE_WITH_TIMEOUT(Api.Plans()->rowCount(), 1, 5000);
}

void ClientApiTest::SendsConfiguredRecommendationPreferences() {
    QTemporaryDir CacheDirectory;
    QVERIFY(CacheDirectory.isValid());
    FixtureServer Fixture;
    QVERIFY(Fixture.Listen());

    Menu::Client::ClientApi Api(
        QStringLiteral("http://127.0.0.1:%1").arg(Fixture.Port()),
        CacheDirectory.path());
    Api.SetAvailableMinutes(30);
    Api.SetAllergies(QStringList{QStringLiteral("花生")});
    Api.SetPantryIngredientIds(QStringList{QStringLiteral("ingredient.rice")});
    Api.LoadData();
    QTRY_COMPARE_WITH_TIMEOUT(Api.Recommendations()->rowCount(), 1, 5000);
    QTRY_VERIFY_WITH_TIMEOUT(
        Fixture.RecommendationRequest().contains("\"availableMinutes\":30"), 5000);
    QTRY_VERIFY_WITH_TIMEOUT(
        Fixture.RecommendationRequest().contains("\"allergies\":[\"花生\"]"), 5000);
    QTRY_VERIFY_WITH_TIMEOUT(
        Fixture.RecommendationRequest().contains(
            "\"pantryIngredientIds\":[\"ingredient.rice\"]"),
        5000);
}

void ClientApiTest::PersistsPreferencesWithVersionedCache() {
    QTemporaryDir CacheDirectory;
    QVERIFY(CacheDirectory.isValid());
    ThreadPoolDrainGuard DrainGuard;
    const QString BaseUrl = QStringLiteral("http://127.0.0.1:1");
    {
        Menu::Client::ClientApi FirstApi(
            BaseUrl, CacheDirectory.path());
        FirstApi.SetServingCount(4);
        FirstApi.SetPreferredCuisines(QStringList{QStringLiteral("日系")});
        QVERIFY(FirstApi.setProperty(
            "Allergies", QStringList{QStringLiteral("花生")}));
        QVERIFY(FirstApi.setProperty("AvailableMinutes", 30));
        QVERIFY(FirstApi.setProperty(
            "PantryIngredientIds", QStringList{QStringLiteral("ingredient.rice")}));
        FirstApi.SetRemindersEnabled(false);
        QThreadPool::globalInstance()->waitForDone(5000);
        const bool Matches = PreferencesMatch(
            ScopedCachePath(CacheDirectory.path(), BaseUrl, QStringLiteral("preferences.json")),
            4, QStringList{QStringLiteral("日系")},
            QStringList{QStringLiteral("花生")}, 30,
            QStringList{QStringLiteral("ingredient.rice")}, false);
        QFile ActualFile(ScopedCachePath(
            CacheDirectory.path(), BaseUrl, QStringLiteral("preferences.json")));
        QVERIFY(ActualFile.open(QIODevice::ReadOnly));
        const QByteArray Actual = ActualFile.readAll();
        QVERIFY2(Matches, Actual.constData());
    }

    Menu::Client::ClientApi SecondApi(
        BaseUrl, CacheDirectory.path());
    QTRY_COMPARE_WITH_TIMEOUT(SecondApi.ServingCount(), 4, 5000);
    QTRY_COMPARE_WITH_TIMEOUT(
        SecondApi.PreferredCuisines(), QStringList{QStringLiteral("日系")}, 5000);
    QCOMPARE(SecondApi.property("Allergies").toStringList(),
        QStringList{QStringLiteral("花生")});
    QCOMPARE(SecondApi.property("AvailableMinutes").toInt(), 30);
    QCOMPARE(SecondApi.property("PantryIngredientIds").toStringList(),
        QStringList{QStringLiteral("ingredient.rice")});
    QCOMPARE(SecondApi.RemindersEnabled(), false);
}

void ClientApiTest::PersistsAndValidatesBaseUrl() {
    QTemporaryDir CacheDirectory;
    QVERIFY(CacheDirectory.isValid());
    ThreadPoolDrainGuard DrainGuard;
    const QString PreviousBaseUrl = QStringLiteral("http://127.0.0.1:8080");
    const QString NextBaseUrl = QStringLiteral("http://192.168.1.20:8081");
    QVERIFY(QDir().mkpath(ScopedCacheDirectory(
        CacheDirectory.path(), PreviousBaseUrl)));
    Menu::Client::ClientApi Api(PreviousBaseUrl, CacheDirectory.path());
    QFile ExistingCache(ScopedCachePath(
        CacheDirectory.path(), PreviousBaseUrl, QStringLiteral("recipes.json")));
    QVERIFY(ExistingCache.open(QIODevice::WriteOnly));
    QVERIFY(ExistingCache.write("old-server-cache") > 0);
    ExistingCache.close();

    QVERIFY(Api.ApplyBaseUrl(NextBaseUrl));
    QCOMPARE(Api.BaseUrl(), NextBaseUrl);
    QVERIFY(!Api.ApplyBaseUrl(QStringLiteral("ftp://192.168.1.20:8081")));
    QCOMPARE(Api.BaseUrl(), NextBaseUrl);
    QVERIFY(Api.ErrorMessage().contains(QStringLiteral("http")));
    QTRY_VERIFY_WITH_TIMEOUT(
        QFile::exists(ScopedCachePath(
            CacheDirectory.path(), PreviousBaseUrl, QStringLiteral("recipes.json"))), 5000);
    QVERIFY(!QFile::exists(ScopedCachePath(
        CacheDirectory.path(), NextBaseUrl, QStringLiteral("recipes.json"))));

    const QString ConnectionPath = CacheDirectory.filePath(
        QStringLiteral("connection.json"));
    QTRY_VERIFY_WITH_TIMEOUT(QFile::exists(ConnectionPath), 5000);
    Menu::Client::ClientApi Restored({}, CacheDirectory.path());
    QCOMPARE(Restored.BaseUrl(), NextBaseUrl);
}

void ClientApiTest::TestsBaseUrlConnectivity() {
    QTemporaryDir CacheDirectory;
    QVERIFY(CacheDirectory.isValid());
    ThreadPoolDrainGuard DrainGuard;
    FixtureServer Fixture;
    QVERIFY(Fixture.Listen());
    Menu::Client::ClientApi Api(
        QStringLiteral("http://127.0.0.1:1"), CacheDirectory.path());
    QSignalSpy TestChanged(&Api, &Menu::Client::ClientApi::ConnectionTestChanged);

    Api.TestBaseUrl(QStringLiteral("http://127.0.0.1:%1").arg(Fixture.Port()));
    QTRY_VERIFY_WITH_TIMEOUT(Api.ConnectionReachable(), 5000);
    QVERIFY(TestChanged.count() > 0);
    QVERIFY(Api.ConnectionStatus().contains(QStringLiteral("成功")));

    Api.TestBaseUrl(QStringLiteral("http://127.0.0.1:1"));
    QTRY_VERIFY_WITH_TIMEOUT(!Api.IsTestingConnection(), 5000);
    QVERIFY(!Api.ConnectionReachable());
    QVERIFY(Api.ConnectionStatus().contains(QStringLiteral("失败")));
}

void ClientApiTest::ClearsPrivateStateOnUnauthorizedResponse() {
    QTemporaryDir CacheDirectory;
    QVERIFY(CacheDirectory.isValid());
    ThreadPoolDrainGuard DrainGuard;
    FixtureServer Fixture;
    QVERIFY(Fixture.Listen());
    const QString BaseUrl = QStringLiteral("http://127.0.0.1:%1").arg(Fixture.Port());
    Menu::Client::ClientApi Api(BaseUrl, CacheDirectory.path());
    Api.Login(QStringLiteral("cook@example.com"), QStringLiteral("Menu-Cook-Password-2026"));
    QTRY_VERIFY_WITH_TIMEOUT(Api.IsAuthenticated(), 5000);
    QTRY_COMPARE_WITH_TIMEOUT(Api.Plans()->rowCount(), 1, 5000);
    QVERIFY(!QDir(ScopedCacheDirectory(CacheDirectory.path(), BaseUrl)).entryList(
        QStringList{QStringLiteral("plans.*.json")}, QDir::Files).isEmpty());

    Fixture.SetPlansUnauthorized();
    Api.LoadPlans();
    QTRY_VERIFY_WITH_TIMEOUT(!Api.IsAuthenticated(), 5000);
    QCOMPARE(Api.Plans()->rowCount(), 0);
    QTRY_VERIFY_WITH_TIMEOUT(QDir(ScopedCacheDirectory(CacheDirectory.path(), BaseUrl)).entryList(
        QStringList{QStringLiteral("plans.*.json")}, QDir::Files).isEmpty(), 5000);
    QVERIFY(Api.ErrorMessage().contains(QStringLiteral("登录已过期")));
}

void ClientApiTest::ClearsPrivateStateOnLogout() {
    QTemporaryDir CacheDirectory;
    QVERIFY(CacheDirectory.isValid());
    ThreadPoolDrainGuard DrainGuard;
    FixtureServer Fixture;
    QVERIFY(Fixture.Listen());

    const QString BaseUrl = QStringLiteral("http://127.0.0.1:%1").arg(Fixture.Port());
    Menu::Client::ClientApi Api(BaseUrl, CacheDirectory.path());
    Api.Login(QStringLiteral("cook@example.com"), QStringLiteral("Menu-Cook-Password-2026"));
    QTRY_VERIFY_WITH_TIMEOUT(Api.IsAuthenticated(), 5000);
    QTRY_COMPARE_WITH_TIMEOUT(Api.Plans()->rowCount(), 1, 5000);

    Api.LoadData();
    QTRY_COMPARE_WITH_TIMEOUT(Api.Recommendations()->rowCount(), 1, 5000);
    Api.OpenRecipe(0);
    QVERIFY(!Api.ActiveRecipe().isEmpty());
    Api.SetServingCount(5);

    Api.Logout();
    QVERIFY(!Api.IsAuthenticated());
    QVERIFY(Api.ActiveRecipe().isEmpty());
    QVERIFY(Api.CurrentCookingSession().isEmpty());
    QCOMPARE(Api.Plans()->rowCount(), 0);
    QCOMPARE(Api.ServingCount(), 2);
    QTRY_VERIFY_WITH_TIMEOUT(
        QDir(ScopedCacheDirectory(CacheDirectory.path(), BaseUrl))
                .entryList(QStringList{QStringLiteral("plans.*.json")}, QDir::Files)
                .isEmpty(),
        5000);
    QThreadPool::globalInstance()->waitForDone(5000);
}

void ClientApiTest::ClearCacheBlocksStaleAsyncWrites() {
    QTemporaryDir CacheDirectory;
    QVERIFY(CacheDirectory.isValid());
    ThreadPoolDrainGuard DrainGuard;
    FixtureServer Fixture;
    QVERIFY(Fixture.Listen());

    ThreadPoolLimitGuard PoolGuard(QThreadPool::globalInstance());
    QThreadPool::globalInstance()->start([]() {
        QThread::msleep(500);
    });

    const QString BaseUrl = QStringLiteral("http://127.0.0.1:%1").arg(Fixture.Port());
    Menu::Client::ClientApi Api(BaseUrl, CacheDirectory.path());
    Api.LoadData();
    QTRY_COMPARE_WITH_TIMEOUT(Api.Recipes()->rowCount(), 1, 5000);
    Api.ClearCache();
    QThreadPool::globalInstance()->waitForDone(5000);

    QVERIFY(!QFile::exists(ScopedCachePath(
        CacheDirectory.path(), BaseUrl, QStringLiteral("recipes.json"))));
    QVERIFY(!QFile::exists(ScopedCachePath(
        CacheDirectory.path(), BaseUrl, QStringLiteral("ingredients.json"))));
    QVERIFY(!QFile::exists(ScopedCachePath(
        CacheDirectory.path(), BaseUrl, QStringLiteral("tonight.json"))));
    QCOMPARE(Api.Recipes()->rowCount(), 0);
    QCOMPARE(Api.Ingredients()->rowCount(), 0);
    QCOMPARE(Api.Recommendations()->rowCount(), 0);
}

void ClientApiTest::CancelsPendingLoadsWithoutLateModelWrites() {
    QTemporaryDir CacheDirectory;
    QVERIFY(CacheDirectory.isValid());
    ThreadPoolDrainGuard DrainGuard;
    Menu::Client::ClientApi Api(
        QStringLiteral("http://127.0.0.1:1"), CacheDirectory.path());
    QSignalSpy Cancelled(&Api, &Menu::Client::ClientApi::LoadCancelled);
    Api.LoadData();
    QVERIFY(Api.IsLoading());
    Api.CancelLoad();
    QVERIFY(!Api.IsLoading());
    QCOMPARE(Cancelled.count(), 1);
    QTest::qWait(100);
    QCOMPARE(Api.Recipes()->rowCount(), 0);
    QCOMPARE(Api.Ingredients()->rowCount(), 0);
    QCOMPARE(Api.Recommendations()->rowCount(), 0);
}

void ClientApiTest::DestroysWithInFlightRequests() {
    QTemporaryDir CacheDirectory;
    QVERIFY(CacheDirectory.isValid());
    {
        Menu::Client::ClientApi Api(
            QStringLiteral("http://127.0.0.1:1"), CacheDirectory.path());
        Api.LoadData();
        Api.Login(QStringLiteral("cook@example.com"),
            QStringLiteral("Menu-Cook-Password-2026"));
    }
    QTest::qWait(100);
    QThreadPool::globalInstance()->waitForDone(5000);
    QVERIFY(true);
}

QTEST_GUILESS_MAIN(ClientApiTest)
#include "ClientApiTest.moc"
