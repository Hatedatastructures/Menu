#include "ClientApi.hpp"

#include "ClientApiSupport.hpp"

#include <QFutureWatcher>
#include <QNetworkReply>
#include <QNetworkRequest>
#include <QThreadPool>
#include <QtConcurrent/QtConcurrentRun>

namespace Menu::Client {

void ClientApi::ReadCache(DataKind Kind, std::uint64_t RequestGenerationValue) {
    const QString Path = CacheFilePath(Kind);
    auto* Watcher = new QFutureWatcher<QByteArray>(this);
    connect(Watcher, &QFutureWatcher<QByteArray>::finished, this,
        [this, Kind, Watcher, RequestGenerationValue]() {
            const QByteArray Data = Watcher->result();
            Watcher->deleteLater();
            if (RequestGenerationValue != RequestGeneration) {
                return;
            }
            if (CancelRequested) {
                FinishDataOperation();
                return;
            }
            if (!Data.isEmpty() && ApplyData(Kind, Data, true)) {
                FinishDataOperation();
                return;
            }
            if (RecipeValues.rowCount() == 0 && IngredientValues.rowCount() == 0 &&
                RecommendationValues.rowCount() == 0 && PlanValues.rowCount() == 0) {
                SetError(QStringLiteral("网络不可用，且没有可用的离线数据"));
            }
            FinishDataOperation();
        });
    Watcher->setFuture(QtConcurrent::run([Path]() {
        return Support::ReadCacheFile(Path);
    }));
}

void ClientApi::HandleReply(
    DataKind Kind,
    QNetworkReply* Reply,
    std::uint64_t RequestGenerationValue) {
    const int StatusCode = Reply->attribute(
        QNetworkRequest::HttpStatusCodeAttribute).toInt();
    const bool IsSuccessful = Reply->error() == QNetworkReply::NoError &&
        StatusCode >= 200 && StatusCode < 300;
    const QByteArray Data = Reply->isOpen() ? Reply->readAll() : QByteArray();
    Reply->deleteLater();

    if (RequestGenerationValue != RequestGeneration) {
        return;
    }
    if (CancelRequested) {
        FinishDataOperation();
        return;
    }
    if (IsSuccessful && ApplyData(Kind, Data)) {
        PersistCache(Kind, Data);
        FinishDataOperation();
        return;
    }

    NetworkFailure = true;
    SetOffline(true);
    if (StatusCode == 401 && Session.IsAuthenticated()) {
        Logout();
        SetError(QStringLiteral("登录已过期，请重新登录"));
        return;
    }
    ReadCache(Kind, RequestGenerationValue);
}

}  // namespace Menu::Client
