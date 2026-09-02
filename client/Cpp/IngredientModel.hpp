#pragma once

#include <QAbstractListModel>
#include <QJsonArray>
#include <QVariantMap>

namespace Menu::Client {

class IngredientModel final : public QAbstractListModel {
    Q_OBJECT
    Q_PROPERTY(int Count READ Count NOTIFY CountChanged)
    Q_PROPERTY(int count READ Count NOTIFY countChanged)

public:
    enum Role {
        IdRole = Qt::UserRole + 1,
        NameRole,
        CategoryRole,
        DefaultUnitRole,
        IsPantryStapleRole,
        AliasesRole,
        SubstituteGroupRole,
        StoreSkuMappingRole,
    };
    Q_ENUM(Role)

    explicit IngredientModel(QObject* Parent = nullptr);

    [[nodiscard]] int Count() const noexcept;
    [[nodiscard]] int rowCount(
        const QModelIndex& Parent = QModelIndex()) const override;
    [[nodiscard]] QVariant data(
        const QModelIndex& Index,
        int RoleValue = Qt::DisplayRole) const override;
    [[nodiscard]] QHash<int, QByteArray> roleNames() const override;

    Q_INVOKABLE QVariantMap IngredientAt(int Index) const;

    void SetIngredients(const QJsonArray& IngredientsValue);

signals:
    void CountChanged();
    void countChanged();

private:
    [[nodiscard]] static const QHash<int, QByteArray>& RoleNamesMap();

    QVector<QVariantMap> Ingredients;
};

}  // namespace Menu::Client
