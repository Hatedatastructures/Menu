#include "ClientApi.hpp"

#include "ClientApiSupport.hpp"

#include <QNetworkAccessManager>
#include <QStandardPaths>

#include <utility>

namespace Menu::Client {
namespace {

QStringList NormalizePreferenceValues(const QStringList& Values) {
    QStringList NormalizedValues;
    for (const QString& RawValue : Values) {
        const QString Value = RawValue.trimmed();
        if (Value.isEmpty() || Value.size() > 32 || NormalizedValues.contains(Value)) {
            continue;
        }
        bool HasControlCharacter = false;
        for (const QChar Character : Value) {
            if (Character.unicode() < 0x20U) {
                HasControlCharacter = true;
                break;
            }
        }
        if (HasControlCharacter) {
            continue;
        }
        NormalizedValues.push_back(Value);
        if (NormalizedValues.size() == 32) {
            break;
        }
    }
    return NormalizedValues;
}

}  // namespace

ClientApi::ClientApi(
    QString BaseUrlValue,
    QString CacheDirectoryValue,
    QObject* Parent)
    : QObject(Parent),
      ServiceUrl(Support::NormalizeBaseUrl(
          BaseUrlValue.isEmpty() ? QStringLiteral("http://127.0.0.1:8080")
                                 : std::move(BaseUrlValue))),
      CacheDirectory(
          CacheDirectoryValue.isEmpty()
              ? qEnvironmentVariable("MENU_CACHE_DIR")
              : std::move(CacheDirectoryValue)) {
    if (CacheDirectory.isEmpty()) {
        CacheDirectory = QStandardPaths::writableLocation(
            QStandardPaths::AppLocalDataLocation);
    }
    NetworkManager = new QNetworkAccessManager(this);
    RecipeValues.setParent(this);
    IngredientValues.setParent(this);
    RecommendationValues.setParent(this);
    PlanValues.setParent(this);
    LoadPreferences();
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
    return !AccessTokenValue.isEmpty();
}

QString ClientApi::DisplayName() const {
    return DisplayNameValue;
}

int ClientApi::ServingCount() const noexcept {
    return ServingCountValue;
}

QStringList ClientApi::PreferredCuisines() const {
    return PreferredCuisinesValue;
}

QStringList ClientApi::Allergies() const {
    return AllergiesValue;
}

int ClientApi::AvailableMinutes() const noexcept {
    return AvailableMinutesValue;
}

QStringList ClientApi::PantryIngredientIds() const {
    return PantryIngredientIdsValue;
}

QStringList ClientApi::Cookware() const {
    return CookwareValue;
}

bool ClientApi::RemindersEnabled() const noexcept {
    return RemindersEnabledValue;
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
    const QString NormalizedValue = Support::NormalizeBaseUrl(BaseUrlValue);
    if (NormalizedValue.isEmpty() || NormalizedValue == ServiceUrl) {
        return;
    }
    ServiceUrl = NormalizedValue;
    emit BaseUrlChanged();
}

void ClientApi::SetServingCount(int ServingCountValueInput) {
    const int NormalizedValue = qBound(1, ServingCountValueInput, 24);
    if (ServingCountValue == NormalizedValue) {
        return;
    }
    ServingCountValue = NormalizedValue;
    emit PreferencesChanged();
    PersistPreferences();
}

void ClientApi::SetPreferredCuisines(const QStringList& CuisinesValue) {
    if (PreferredCuisinesValue == CuisinesValue) {
        return;
    }
    PreferredCuisinesValue = CuisinesValue;
    emit PreferencesChanged();
    PersistPreferences();
}

void ClientApi::SetAllergies(const QStringList& AllergiesValueInput) {
    const QStringList NormalizedValue = NormalizePreferenceValues(AllergiesValueInput);
    if (AllergiesValue == NormalizedValue) {
        return;
    }
    AllergiesValue = NormalizedValue;
    emit PreferencesChanged();
    PersistPreferences();
}

void ClientApi::SetAvailableMinutes(int AvailableMinutesValueInput) {
    const int NormalizedValue = qBound(10, AvailableMinutesValueInput, 240);
    if (AvailableMinutesValue == NormalizedValue) {
        return;
    }
    AvailableMinutesValue = NormalizedValue;
    emit PreferencesChanged();
    PersistPreferences();
}

void ClientApi::SetPantryIngredientIds(
    const QStringList& PantryIngredientIdsValueInput) {
    const QStringList NormalizedValue =
        NormalizePreferenceValues(PantryIngredientIdsValueInput);
    if (PantryIngredientIdsValue == NormalizedValue) {
        return;
    }
    PantryIngredientIdsValue = NormalizedValue;
    emit PreferencesChanged();
    PersistPreferences();
}

void ClientApi::SetCookware(const QStringList& CookwareValueInput) {
    if (CookwareValue == CookwareValueInput) {
        return;
    }
    CookwareValue = CookwareValueInput;
    emit PreferencesChanged();
    PersistPreferences();
}

void ClientApi::SetRemindersEnabled(bool EnabledValue) {
    if (RemindersEnabledValue == EnabledValue) {
        return;
    }
    RemindersEnabledValue = EnabledValue;
    emit PreferencesChanged();
    PersistPreferences();
}

}  // namespace Menu::Client
