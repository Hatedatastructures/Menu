#include "ClientPreferenceState.hpp"

#include <QChar>

namespace Menu::Client {

ClientPreferenceState::ClientPreferenceState() {
    Reset();
}

int ClientPreferenceState::ServingCount() const noexcept {
    return ServingCountValue;
}

const QStringList& ClientPreferenceState::PreferredCuisines() const noexcept {
    return PreferredCuisinesValue;
}

const QStringList& ClientPreferenceState::Allergies() const noexcept {
    return AllergiesValue;
}

int ClientPreferenceState::AvailableMinutes() const noexcept {
    return AvailableMinutesValue;
}

const QStringList& ClientPreferenceState::PantryIngredientIds() const noexcept {
    return PantryIngredientIdsValue;
}

const QStringList& ClientPreferenceState::Cookware() const noexcept {
    return CookwareValue;
}

bool ClientPreferenceState::RemindersEnabled() const noexcept {
    return RemindersEnabledValue;
}

bool ClientPreferenceState::SetServingCount(int Value) noexcept {
    const int Normalized = qBound(1, Value, 24);
    if (ServingCountValue == Normalized) {
        return false;
    }
    ServingCountValue = Normalized;
    return true;
}

bool ClientPreferenceState::SetPreferredCuisines(const QStringList& Values) {
    const QStringList Normalized = Normalize(Values);
    if (PreferredCuisinesValue == Normalized) {
        return false;
    }
    PreferredCuisinesValue = Normalized;
    return true;
}

bool ClientPreferenceState::SetAllergies(const QStringList& Values) {
    const QStringList Normalized = Normalize(Values);
    if (AllergiesValue == Normalized) {
        return false;
    }
    AllergiesValue = Normalized;
    return true;
}

bool ClientPreferenceState::SetAvailableMinutes(int Value) noexcept {
    const int Normalized = qBound(10, Value, 240);
    if (AvailableMinutesValue == Normalized) {
        return false;
    }
    AvailableMinutesValue = Normalized;
    return true;
}

bool ClientPreferenceState::SetPantryIngredientIds(const QStringList& Values) {
    const QStringList Normalized = Normalize(Values);
    if (PantryIngredientIdsValue == Normalized) {
        return false;
    }
    PantryIngredientIdsValue = Normalized;
    return true;
}

bool ClientPreferenceState::SetCookware(const QStringList& Values) {
    const QStringList Normalized = Normalize(Values);
    if (CookwareValue == Normalized) {
        return false;
    }
    CookwareValue = Normalized;
    return true;
}

bool ClientPreferenceState::SetRemindersEnabled(bool Value) noexcept {
    if (RemindersEnabledValue == Value) {
        return false;
    }
    RemindersEnabledValue = Value;
    return true;
}

bool ClientPreferenceState::Apply(const PreferencesSnapshot& SnapshotValue) noexcept {
    const bool Changed = ServingCountValue != SnapshotValue.ServingCount ||
        PreferredCuisinesValue != SnapshotValue.PreferredCuisines ||
        AllergiesValue != SnapshotValue.Allergies ||
        AvailableMinutesValue != SnapshotValue.AvailableMinutes ||
        PantryIngredientIdsValue != SnapshotValue.PantryIngredientIds ||
        CookwareValue != SnapshotValue.Cookware ||
        RemindersEnabledValue != SnapshotValue.RemindersEnabled;
    ServingCountValue = SnapshotValue.ServingCount;
    PreferredCuisinesValue = SnapshotValue.PreferredCuisines;
    AllergiesValue = SnapshotValue.Allergies;
    AvailableMinutesValue = SnapshotValue.AvailableMinutes;
    PantryIngredientIdsValue = SnapshotValue.PantryIngredientIds;
    CookwareValue = SnapshotValue.Cookware;
    RemindersEnabledValue = SnapshotValue.RemindersEnabled;
    return Changed;
}

PreferencesSnapshot ClientPreferenceState::Snapshot() const {
    PreferencesSnapshot Result;
    Result.ServingCount = ServingCountValue;
    Result.PreferredCuisines = PreferredCuisinesValue;
    Result.Allergies = AllergiesValue;
    Result.AvailableMinutes = AvailableMinutesValue;
    Result.PantryIngredientIds = PantryIngredientIdsValue;
    Result.Cookware = CookwareValue;
    Result.RemindersEnabled = RemindersEnabledValue;
    return Result;
}

void ClientPreferenceState::Reset() noexcept {
    ServingCountValue = 2;
    PreferredCuisinesValue = {
        QStringLiteral("中餐"), QStringLiteral("西餐"), QStringLiteral("日系")};
    AllergiesValue.clear();
    AvailableMinutesValue = 45;
    PantryIngredientIdsValue.clear();
    CookwareValue = {
        QStringLiteral("炒锅"), QStringLiteral("平底锅"), QStringLiteral("烤箱")};
    RemindersEnabledValue = true;
}

QStringList ClientPreferenceState::Normalize(const QStringList& Values) {
    QStringList Result;
    for (const QString& RawValue : Values) {
        const QString Value = RawValue.trimmed();
        if (Value.isEmpty() || Value.size() > 32 || Result.contains(Value)) {
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
        Result.push_back(Value);
        if (Result.size() == 32) {
            break;
        }
    }
    return Result;
}

}  // namespace Menu::Client
