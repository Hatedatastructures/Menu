#include "IngredientModel.hpp"

#include <QJsonObject>

#include <utility>

namespace Menu::Client {

IngredientModel::IngredientModel(QObject* Parent)
    : QAbstractListModel(Parent) {}

int IngredientModel::Count() const noexcept {
    return Ingredients.size();
}

int IngredientModel::rowCount(const QModelIndex& Parent) const {
    return Parent.isValid() ? 0 : Ingredients.size();
}

QVariant IngredientModel::data(const QModelIndex& Index, int RoleValue) const {
    if (!Index.isValid() || Index.row() < 0 || Index.row() >= Ingredients.size()) {
        return {};
    }
    const QVariantMap& Ingredient = Ingredients.at(Index.row());
    return Ingredient.value(QString::fromUtf8(RoleNamesMap().value(RoleValue)));
}

QHash<int, QByteArray> IngredientModel::roleNames() const {
    return RoleNamesMap();
}

const QHash<int, QByteArray>& IngredientModel::RoleNamesMap() {
    static const QHash<int, QByteArray> Roles = {
        {IdRole, "id"},
        {NameRole, "name"},
        {CategoryRole, "category"},
        {DefaultUnitRole, "defaultUnit"},
        {IsPantryStapleRole, "isPantryStaple"},
        {AliasesRole, "aliases"},
        {SubstituteGroupRole, "substituteGroup"},
        {StoreSkuMappingRole, "storeSkuMapping"},
    };
    return Roles;
}

void IngredientModel::SetIngredients(const QJsonArray& IngredientsValue) {
    QVector<QVariantMap> NextIngredients;
    NextIngredients.reserve(IngredientsValue.size());
    for (const QJsonValue& Value : IngredientsValue) {
        if (Value.isObject()) {
            NextIngredients.push_back(Value.toObject().toVariantMap());
        }
    }

    beginResetModel();
    Ingredients = std::move(NextIngredients);
    endResetModel();
    emit CountChanged();
}

}  // namespace Menu::Client
