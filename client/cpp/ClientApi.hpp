#pragma once

#include "IngredientModel.hpp"
#include "MealPlanModel.hpp"
#include "RecommendationModel.hpp"
#include "RecipeModel.hpp"

#include <QObject>
#include <QList>
#include <QString>
#include <QStringList>
#include <QVariantList>
#include <QVariantMap>

#include <atomic>
#include <cstdint>
#include <memory>

class QNetworkAccessManager;
class QNetworkReply;
class QNetworkRequest;

namespace Menu::Client {

class ClientApi final : public QObject {
    Q_OBJECT
    Q_PROPERTY(RecipeModel* Recipes READ Recipes CONSTANT)
    Q_PROPERTY(IngredientModel* Ingredients READ Ingredients CONSTANT)
    Q_PROPERTY(RecommendationModel* Recommendations READ Recommendations CONSTANT)
    Q_PROPERTY(MealPlanModel* Plans READ Plans CONSTANT)
    Q_PROPERTY(QVariantMap ActiveRecipe READ ActiveRecipe NOTIFY ActiveRecipeChanged)
    Q_PROPERTY(QVariantMap CurrentCookingSession READ CurrentCookingSession NOTIFY CookingSessionChanged)
    Q_PROPERTY(bool IsAuthenticated READ IsAuthenticated NOTIFY AuthenticationChanged)
    Q_PROPERTY(QString DisplayName READ DisplayName NOTIFY AuthenticationChanged)
    Q_PROPERTY(int ServingCount READ ServingCount WRITE SetServingCount NOTIFY PreferencesChanged)
    Q_PROPERTY(QStringList PreferredCuisines READ PreferredCuisines WRITE SetPreferredCuisines NOTIFY PreferencesChanged)
    Q_PROPERTY(QStringList Allergies READ Allergies WRITE SetAllergies NOTIFY PreferencesChanged)
    Q_PROPERTY(int AvailableMinutes READ AvailableMinutes WRITE SetAvailableMinutes NOTIFY PreferencesChanged)
    Q_PROPERTY(QStringList PantryIngredientIds READ PantryIngredientIds WRITE SetPantryIngredientIds NOTIFY PreferencesChanged)
    Q_PROPERTY(QStringList Cookware READ Cookware WRITE SetCookware NOTIFY PreferencesChanged)
    Q_PROPERTY(bool RemindersEnabled READ RemindersEnabled WRITE SetRemindersEnabled NOTIFY PreferencesChanged)
    Q_PROPERTY(bool IsLoading READ IsLoading NOTIFY LoadingChanged)
    Q_PROPERTY(bool IsOffline READ IsOffline NOTIFY OfflineChanged)
    Q_PROPERTY(QString ErrorMessage READ ErrorMessage NOTIFY ErrorMessageChanged)
    Q_PROPERTY(QString BaseUrl READ BaseUrl WRITE SetBaseUrl NOTIFY BaseUrlChanged)

public:
    explicit ClientApi(
        QString BaseUrlValue = {},
        QString CacheDirectoryValue = {},
        QObject* Parent = nullptr);

    [[nodiscard]] RecipeModel* Recipes() const noexcept;
    [[nodiscard]] IngredientModel* Ingredients() const noexcept;
    [[nodiscard]] RecommendationModel* Recommendations() const noexcept;
    [[nodiscard]] MealPlanModel* Plans() const noexcept;
    [[nodiscard]] QVariantMap ActiveRecipe() const;
    [[nodiscard]] QVariantMap CurrentCookingSession() const;
    [[nodiscard]] bool IsAuthenticated() const noexcept;
    [[nodiscard]] QString DisplayName() const;
    [[nodiscard]] int ServingCount() const noexcept;
    [[nodiscard]] QStringList PreferredCuisines() const;
    [[nodiscard]] QStringList Allergies() const;
    [[nodiscard]] int AvailableMinutes() const noexcept;
    [[nodiscard]] QStringList PantryIngredientIds() const;
    [[nodiscard]] QStringList Cookware() const;
    [[nodiscard]] bool RemindersEnabled() const noexcept;
    [[nodiscard]] bool IsLoading() const noexcept;
    [[nodiscard]] bool IsOffline() const noexcept;
    [[nodiscard]] QString ErrorMessage() const;
    [[nodiscard]] QString BaseUrl() const;

    void SetBaseUrl(const QString& BaseUrlValue);
    void SetServingCount(int ServingCountValue);
    void SetPreferredCuisines(const QStringList& CuisinesValue);
    void SetAllergies(const QStringList& AllergiesValue);
    void SetAvailableMinutes(int AvailableMinutesValue);
    void SetPantryIngredientIds(const QStringList& PantryIngredientIdsValue);
    void SetCookware(const QStringList& CookwareValue);
    void SetRemindersEnabled(bool EnabledValue);

    Q_INVOKABLE void LoadData();
    Q_INVOKABLE void LoadPlans();
    Q_INVOKABLE void SavePlan(const QString& PlanDate, const QVariantList& RecipeIds);
    Q_INVOKABLE void ClearCache();
    Q_INVOKABLE void Login(const QString& Email, const QString& Password);
    Q_INVOKABLE void RefreshSession();
    Q_INVOKABLE void Register(
        const QString& Email,
        const QString& Password,
        const QString& DisplayNameValue);
    Q_INVOKABLE void Logout();
    Q_INVOKABLE void CancelLoad();
    Q_INVOKABLE void OpenRecipe(int Index);
    Q_INVOKABLE void CloseRecipe();
    Q_INVOKABLE void StartCookingSession();
    Q_INVOKABLE void UpdateCookingSession(int CurrentStepOrder, const QString& State);
    Q_INVOKABLE void SubmitFeedback(
        const QString& Outcome,
        const QVariantList& Tags,
        const QString& Comment);
    Q_INVOKABLE void ClearError();

signals:
    void LoadingChanged();
    void OfflineChanged();
    void ErrorMessageChanged();
    void BaseUrlChanged();
    void ActiveRecipeChanged();
    void CookingSessionChanged();
    void AuthenticationChanged();
    void PreferencesChanged();
    void DataReady();
    void LoadCancelled();

private:
    enum class DataKind {
        Recipes,
        Ingredients,
        Recommendations,
        Plans,
    };

    void RequestData(DataKind Kind);
    void HandleReply(
        DataKind Kind,
        QNetworkReply* Reply,
        std::uint64_t RequestGenerationValue);
    void ReadCache(DataKind Kind, std::uint64_t RequestGenerationValue);
    bool ApplyData(DataKind Kind, const QByteArray& Data, bool IsCached = false);
    void PersistCache(DataKind Kind, const QByteArray& Data);
    void FinishDataOperation();
    void SetLoading(bool LoadingValue);
    void SetOffline(bool OfflineValue);
    void SetError(QString ErrorValue);
    void ConfigureReply(QNetworkReply* Reply);
    [[nodiscard]] QNetworkRequest CreateRequest(const QString& Resource) const;
    void LoadPreferences();
    bool ApplyPreferences(const QByteArray& Data);
    void PersistPreferences() const;
    void HandleAuthenticationReply(
        QNetworkReply* Reply,
        bool IsRegistration,
        bool IsRefresh,
        std::uint64_t RequestGenerationValue);
    void HandleCookingSessionReply(
        QNetworkReply* Reply,
        std::uint64_t RequestGenerationValue);
    void HandleFeedbackReply(
        QNetworkReply* Reply,
        std::uint64_t RequestGenerationValue);
    void SetAuthentication(
        QString AccessTokenValue,
        QString RefreshTokenValue,
        QString DisplayNameValue,
        QString UserIdValue);
    void ResetPrivateState(const QString& UserIdToRemove, bool RemovePlanCache);
    void ResetPreferencesToDefaults();
    [[nodiscard]] QByteArray RecommendationRequestBody() const;
    [[nodiscard]] QString PreferencesFilePath() const;
    [[nodiscard]] QString CacheFilePath(DataKind Kind) const;

    QString ServiceUrl;
    QString CacheDirectory;
    RecipeModel RecipeValues;
    IngredientModel IngredientValues;
    RecommendationModel RecommendationValues;
    MealPlanModel PlanValues;
    QVariantMap ActiveRecipeValue;
    QVariantMap CurrentCookingSessionValue;
    QString AccessTokenValue;
    QString RefreshTokenValue;
    QString DisplayNameValue;
    QString UserIdValue;
    int ServingCountValue = 2;
    QStringList PreferredCuisinesValue = {QStringLiteral("中餐"), QStringLiteral("西餐"), QStringLiteral("日系")};
    QStringList AllergiesValue;
    int AvailableMinutesValue = 45;
    QStringList PantryIngredientIdsValue;
    QStringList CookwareValue = {QStringLiteral("炒锅"), QStringLiteral("平底锅"), QStringLiteral("烤箱")};
    bool RemindersEnabledValue = true;
    QNetworkAccessManager* NetworkManager = nullptr;
    bool Loading = false;
    bool Offline = false;
    bool NetworkFailure = false;
    bool CancelRequested = false;
    bool PreferencesLoading = false;
    QString Error;
    int PendingOperations = 0;
    bool DataLoadPending = false;
    bool PlanLoadPending = false;
    QList<QNetworkReply*> ActiveReplies;
    QNetworkReply* AuthenticationReply = nullptr;
    QNetworkReply* CookingSessionReply = nullptr;
    QNetworkReply* FeedbackReply = nullptr;
    std::uint64_t RequestGeneration = 0;
    std::shared_ptr<std::atomic<std::uint64_t>> CacheEpoch =
        std::make_shared<std::atomic<std::uint64_t>>(0);
    std::shared_ptr<std::atomic<std::uint64_t>> PreferenceGeneration =
        std::make_shared<std::atomic<std::uint64_t>>(0);
};

}  // namespace Menu::Client
