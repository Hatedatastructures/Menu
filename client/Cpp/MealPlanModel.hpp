#pragma once

#include <QAbstractListModel>
#include <QJsonArray>
#include <QVariantMap>

namespace Menu::Client {

class MealPlanModel final : public QAbstractListModel {
    Q_OBJECT
    Q_PROPERTY(int Count READ Count NOTIFY CountChanged)

public:
    enum Role {
        IdRole = Qt::UserRole + 1,
        PlanDateRole,
        ItemsRole,
        CombinedIngredientsRole,
    };
    Q_ENUM(Role)

    explicit MealPlanModel(QObject* Parent = nullptr);

    [[nodiscard]] int Count() const noexcept;
    [[nodiscard]] int rowCount(
        const QModelIndex& Parent = QModelIndex()) const override;
    [[nodiscard]] QVariant data(
        const QModelIndex& Index,
        int RoleValue = Qt::DisplayRole) const override;
    [[nodiscard]] QHash<int, QByteArray> roleNames() const override;

    void SetPlans(const QJsonArray& PlansValue);
    Q_INVOKABLE QVariantMap PlanAt(int Index) const;

signals:
    void CountChanged();

private:
    [[nodiscard]] static const QHash<int, QByteArray>& RoleNamesMap();

    QVector<QVariantMap> Plans;
};

}  // namespace Menu::Client
