#include "ClientApi.hpp"

#include <QNetworkReply>

#include <utility>

namespace Menu::Client {

void ClientApi::CancelLoad() {
    if (!Loading) {
        return;
    }
    ++RequestGeneration;
    const QList<QNetworkReply*> Replies = ActiveReplies;
    for (QNetworkReply* Reply : Replies) {
        if (Reply != nullptr) {
            Reply->abort();
        }
    }
    ActiveReplies.clear();
    PendingOperations = 0;
    DataLoadPending = false;
    PlanLoadPending = false;
    CancelRequested = false;
    SetLoading(false);
    SetError(QStringLiteral("已取消加载"));
    emit LoadCancelled();
}

void ClientApi::ClearError() {
    SetError({});
}

void ClientApi::FinishDataOperation() {
    --PendingOperations;
    if (PendingOperations > 0) {
        return;
    }
    SetOffline(NetworkFailure);
    SetLoading(false);
    if (DataLoadPending) {
        DataLoadPending = false;
        LoadData();
        return;
    }
    if (PlanLoadPending && IsAuthenticated()) {
        PlanLoadPending = false;
        LoadPlans();
        return;
    }
    if (CancelRequested) {
        SetError(QStringLiteral("已取消加载"));
        emit LoadCancelled();
    } else {
        emit DataReady();
    }
}

void ClientApi::SetLoading(bool LoadingValue) {
    if (Loading == LoadingValue) {
        return;
    }
    Loading = LoadingValue;
    emit LoadingChanged();
}

void ClientApi::SetOffline(bool OfflineValue) {
    if (Offline == OfflineValue) {
        return;
    }
    Offline = OfflineValue;
    emit OfflineChanged();
}

void ClientApi::SetError(QString ErrorValue) {
    if (Error == ErrorValue) {
        return;
    }
    Error = std::move(ErrorValue);
    emit ErrorMessageChanged();
}

}  // namespace Menu::Client
