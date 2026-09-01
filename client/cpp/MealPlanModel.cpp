#include "MealPlanModel.hpp"

#include <QJsonObject>

#include <utility>

namespace Menu::Client {

MealPlanModel::MealPlanModel(QObject* Parent)
    : QAbstractListModel(Parent) {}

int MealPlanModel::Count() const noexcept {
    return Plans.size();
}

int MealPlanModel::rowCount(const QModelIndex& Parent) const {
    return Parent.isValid() ? 0 : Plans.size();
}

QVariant MealPlanModel::data(const QModelIndex& Index, int RoleValue) const {
    if (!Index.isValid() || Index.row() < 0 || Index.row() >= Plans.size()) {
        return {};
    }
    const QVariantMap& Plan = Plans.at(Index.row());
    return Plan.value(QString::fromUtf8(RoleNamesMap().value(RoleValue)));
}

QHash<int, QByteArray> MealPlanModel::roleNames() const {
    return RoleNamesMap();
}

const QHash<int, QByteArray>& MealPlanModel::RoleNamesMap() {
    static const QHash<int, QByteArray> Roles = {
        {IdRole, "id"},
        {PlanDateRole, "planDate"},
        {ItemsRole, "items"},
        {CombinedIngredientsRole, "combinedIngredients"},
    };
    return Roles;
}

void MealPlanModel::SetPlans(const QJsonArray& PlansValue) {
    QVector<QVariantMap> NextPlans;
    NextPlans.reserve(PlansValue.size());
    for (const QJsonValue& Value : PlansValue) {
        if (Value.isObject()) {
            NextPlans.push_back(Value.toObject().toVariantMap());
        }
    }
    beginResetModel();
    Plans = std::move(NextPlans);
    endResetModel();
    emit CountChanged();
}

QVariantMap MealPlanModel::PlanAt(int Index) const {
    if (Index < 0 || Index >= Plans.size()) {
        return {};
    }
    return Plans.at(Index);
}

}  // namespace Menu::Client
