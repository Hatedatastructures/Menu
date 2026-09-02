#include "PreferencesStore.hpp"

#include "ClientApiSupport.hpp"

#include <QJsonArray>
#include <QJsonDocument>
#include <QJsonObject>
#include <QJsonParseError>

namespace Menu::Client {
namespace {

std::optional<QStringList> DecodeStringList(const QJsonValue& Value) {
    if (!Value.isArray() || Value.toArray().size() > 32) {
        return std::nullopt;
    }
    QStringList Values;
    for (const QJsonValue& Item : Value.toArray()) {
        const QString Text = Item.toString();
        if (!Item.isString() || Text.isEmpty() || Text.size() > 32) {
            return std::nullopt;
        }
        Values.push_back(Text);
    }
    return Values;
}

QJsonArray EncodeStringList(const QStringList& Values) {
    QJsonArray Array;
    for (const QString& Value : Values) {
        Array.append(Value);
    }
    return Array;
}

}  // namespace

std::optional<PreferencesSnapshot> PreferencesStore::Decode(
    const QByteArray& Data) {
    QJsonParseError ParseError;
    const QJsonDocument Document = QJsonDocument::fromJson(Data, &ParseError);
    if (ParseError.error != QJsonParseError::NoError || !Document.isObject()) {
        return std::nullopt;
    }
    const QJsonObject Object = Document.object();
    if (Object.value(QStringLiteral("schemaVersion")).toInt(-1) !=
            Support::CacheSchemaVersion ||
        Object.value(QStringLiteral("kind")).toString() !=
            QStringLiteral("preferences")) {
        return std::nullopt;
    }

    const QJsonValue ServingCountValue = Object.value(QStringLiteral("servings"));
    const QJsonValue AvailableMinutesValue = Object.value(
        QStringLiteral("availableMinutes"));
    const QJsonValue RemindersValue = Object.value(QStringLiteral("remindersEnabled"));
    if (!ServingCountValue.isDouble() || !AvailableMinutesValue.isDouble() ||
        !RemindersValue.isBool()) {
        return std::nullopt;
    }
    const int ServingCount = ServingCountValue.toInt(-1);
    const int AvailableMinutes = AvailableMinutesValue.toInt(-1);
    if (ServingCount < 1 || ServingCount > 24 ||
        AvailableMinutes < 10 || AvailableMinutes > 240) {
        return std::nullopt;
    }

    const auto Cuisines = DecodeStringList(Object.value(QStringLiteral("cuisines")));
    const auto Allergies = DecodeStringList(Object.value(QStringLiteral("allergies")));
    const auto PantryIngredientIds = DecodeStringList(
        Object.value(QStringLiteral("pantryIngredientIds")));
    const auto Cookware = DecodeStringList(Object.value(QStringLiteral("cookware")));
    if (!Cuisines.has_value() || !Allergies.has_value() ||
        !PantryIngredientIds.has_value() || !Cookware.has_value()) {
        return std::nullopt;
    }

    PreferencesSnapshot Snapshot;
    Snapshot.ServingCount = ServingCount;
    Snapshot.PreferredCuisines = Cuisines.value();
    Snapshot.Allergies = Allergies.value();
    Snapshot.AvailableMinutes = AvailableMinutes;
    Snapshot.PantryIngredientIds = PantryIngredientIds.value();
    Snapshot.Cookware = Cookware.value();
    Snapshot.RemindersEnabled = RemindersValue.toBool();
    return Snapshot;
}

QByteArray PreferencesStore::Encode(const PreferencesSnapshot& Snapshot) {
    QJsonObject Object;
    Object.insert(QStringLiteral("schemaVersion"), Support::CacheSchemaVersion);
    Object.insert(QStringLiteral("kind"), QStringLiteral("preferences"));
    Object.insert(QStringLiteral("servings"), Snapshot.ServingCount);
    Object.insert(QStringLiteral("cuisines"), EncodeStringList(Snapshot.PreferredCuisines));
    Object.insert(QStringLiteral("allergies"), EncodeStringList(Snapshot.Allergies));
    Object.insert(QStringLiteral("availableMinutes"), Snapshot.AvailableMinutes);
    Object.insert(QStringLiteral("pantryIngredientIds"),
                  EncodeStringList(Snapshot.PantryIngredientIds));
    Object.insert(QStringLiteral("cookware"), EncodeStringList(Snapshot.Cookware));
    Object.insert(QStringLiteral("remindersEnabled"), Snapshot.RemindersEnabled);
    return QJsonDocument(Object).toJson(QJsonDocument::Compact);
}

}  // namespace Menu::Client
