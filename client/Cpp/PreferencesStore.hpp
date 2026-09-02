#pragma once

#include <QByteArray>
#include <QStringList>

#include <optional>

namespace Menu::Client {

struct PreferencesSnapshot final {
    int ServingCount = 2;
    QStringList PreferredCuisines;
    QStringList Allergies;
    int AvailableMinutes = 45;
    QStringList PantryIngredientIds;
    QStringList Cookware;
    bool RemindersEnabled = true;
};

class PreferencesStore final {
public:
    [[nodiscard]] static std::optional<PreferencesSnapshot> Decode(
        const QByteArray& Data);
    [[nodiscard]] static QByteArray Encode(
        const PreferencesSnapshot& Snapshot);
};

}  // namespace Menu::Client
