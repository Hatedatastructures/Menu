#include "ClientApi.hpp"

#include "ConnectionEndpoint.hpp"
#include "ClientApiSupport.hpp"

#include <QStandardPaths>

#include <utility>

namespace Menu::Client {

ClientApi::ClientApi(
    QString BaseUrlValue,
    QString CacheDirectoryValue,
    QObject* Parent)
    : QObject(Parent),
      ServiceUrl(Support::NormalizeBaseUrl(
          BaseUrlValue.isEmpty() ? DefaultBaseUrl()
                                 : BaseUrlValue)),
      CacheDirectory(
          CacheDirectoryValue.isEmpty()
              ? qEnvironmentVariable("MENU_CACHE_DIR")
              : std::move(CacheDirectoryValue)),
      Transport(this) {
    if (CacheDirectory.isEmpty()) {
        CacheDirectory = QStandardPaths::writableLocation(
            QStandardPaths::AppLocalDataLocation);
    }
    RecipeValues.setParent(this);
    IngredientValues.setParent(this);
    RecommendationValues.setParent(this);
    PlanValues.setParent(this);
    ConnectionState = QStringLiteral("未测试");
    if (BaseUrlValue.isEmpty()) {
        LoadConnectionSettings();
    }
    Transport.SetBaseUrl(ServiceUrl);
    LoadPreferences();
}

ClientApi::~ClientApi() {
    ++RequestGeneration;
    ++AuthenticationGeneration;
    ++ConnectionTestGeneration;
    CacheEpoch->fetch_add(1, std::memory_order_acq_rel);
    PreferenceGeneration->fetch_add(1, std::memory_order_acq_rel);
    ConnectionGeneration->fetch_add(1, std::memory_order_acq_rel);
    Transport.AbortAll();
}

RecipeModel* ClientApi::Recipes() const noexcept {
    return const_cast<RecipeModel*>(&RecipeValues);
}

IngredientModel* ClientApi::Ingredients() const noexcept {
    return const_cast<IngredientModel*>(&IngredientValues);
}

RecommendationModel* ClientApi::Recommendations() const noexcept {
    return const_cast<RecommendationModel*>(&RecommendationValues);
}

MealPlanModel* ClientApi::Plans() const noexcept {
    return const_cast<MealPlanModel*>(&PlanValues);
}

QVariantMap ClientApi::ActiveRecipe() const {
    return ActiveRecipeValue;
}

QVariantMap ClientApi::CurrentCookingSession() const {
    return CurrentCookingSessionValue;
}

bool ClientApi::IsAuthenticated() const noexcept {
    return Session.IsAuthenticated();
}

QString ClientApi::DisplayName() const {
    return Session.DisplayName();
}

int ClientApi::ServingCount() const noexcept {
    return Preferences.ServingCount();
}

QStringList ClientApi::PreferredCuisines() const {
    return Preferences.PreferredCuisines();
}

QStringList ClientApi::Allergies() const {
    return Preferences.Allergies();
}

int ClientApi::AvailableMinutes() const noexcept {
    return Preferences.AvailableMinutes();
}

QStringList ClientApi::PantryIngredientIds() const {
    return Preferences.PantryIngredientIds();
}

QStringList ClientApi::Cookware() const {
    return Preferences.Cookware();
}

bool ClientApi::RemindersEnabled() const noexcept {
    return Preferences.RemindersEnabled();
}

bool ClientApi::IsLoading() const noexcept {
    return Loading;
}

bool ClientApi::IsOffline() const noexcept {
    return Offline;
}

QString ClientApi::ErrorMessage() const {
    return Error;
}

QString ClientApi::BaseUrl() const {
    return ServiceUrl;
}

void ClientApi::SetBaseUrl(const QString& BaseUrlValue) {
    (void)ApplyBaseUrl(BaseUrlValue);
}

void ClientApi::SetServingCount(int ServingCountValueInput) {
    if (!Preferences.SetServingCount(ServingCountValueInput)) {
        return;
    }
    emit PreferencesChanged();
    PersistPreferences();
}

void ClientApi::SetPreferredCuisines(const QStringList& CuisinesValue) {
    if (!Preferences.SetPreferredCuisines(CuisinesValue)) {
        return;
    }
    emit PreferencesChanged();
    PersistPreferences();
}

void ClientApi::SetAllergies(const QStringList& AllergiesValueInput) {
    if (!Preferences.SetAllergies(AllergiesValueInput)) {
        return;
    }
    emit PreferencesChanged();
    PersistPreferences();
}

void ClientApi::SetAvailableMinutes(int AvailableMinutesValueInput) {
    if (!Preferences.SetAvailableMinutes(AvailableMinutesValueInput)) {
        return;
    }
    emit PreferencesChanged();
    PersistPreferences();
}

void ClientApi::SetPantryIngredientIds(
    const QStringList& PantryIngredientIdsValueInput) {
    if (!Preferences.SetPantryIngredientIds(PantryIngredientIdsValueInput)) {
        return;
    }
    emit PreferencesChanged();
    PersistPreferences();
}

void ClientApi::SetCookware(const QStringList& CookwareValueInput) {
    if (!Preferences.SetCookware(CookwareValueInput)) {
        return;
    }
    emit PreferencesChanged();
    PersistPreferences();
}

void ClientApi::SetRemindersEnabled(bool EnabledValue) {
    if (!Preferences.SetRemindersEnabled(EnabledValue)) {
        return;
    }
    emit PreferencesChanged();
    PersistPreferences();
}

}  // namespace Menu::Client
