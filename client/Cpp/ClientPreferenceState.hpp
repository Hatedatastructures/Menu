#pragma once

#include "PreferencesStore.hpp"

#include <QStringList>

namespace Menu::Client {

class ClientPreferenceState final {
public:
    ClientPreferenceState();

    [[nodiscard]] int ServingCount() const noexcept;
    [[nodiscard]] const QStringList& PreferredCuisines() const noexcept;
    [[nodiscard]] const QStringList& Allergies() const noexcept;
    [[nodiscard]] int AvailableMinutes() const noexcept;
    [[nodiscard]] const QStringList& PantryIngredientIds() const noexcept;
    [[nodiscard]] const QStringList& Cookware() const noexcept;
    [[nodiscard]] bool RemindersEnabled() const noexcept;

    [[nodiscard]] bool SetServingCount(int Value) noexcept;
    [[nodiscard]] bool SetPreferredCuisines(const QStringList& Values);
    [[nodiscard]] bool SetAllergies(const QStringList& Values);
    [[nodiscard]] bool SetAvailableMinutes(int Value) noexcept;
    [[nodiscard]] bool SetPantryIngredientIds(const QStringList& Values);
    [[nodiscard]] bool SetCookware(const QStringList& Values);
    [[nodiscard]] bool SetRemindersEnabled(bool Value) noexcept;

    [[nodiscard]] bool Apply(const PreferencesSnapshot& Snapshot) noexcept;
    [[nodiscard]] PreferencesSnapshot Snapshot() const;
    void Reset() noexcept;

private:
    static QStringList Normalize(const QStringList& Values);

    int ServingCountValue = 2;
    QStringList PreferredCuisinesValue;
    QStringList AllergiesValue;
    int AvailableMinutesValue = 45;
    QStringList PantryIngredientIdsValue;
    QStringList CookwareValue;
    bool RemindersEnabledValue = true;
};

}  // namespace Menu::Client
