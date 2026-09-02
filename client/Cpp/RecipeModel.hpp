#pragma once

#include <QAbstractListModel>
#include <QJsonArray>
#include <QVariantMap>

namespace Menu::Client {

class RecipeModel final : public QAbstractListModel {
    Q_OBJECT
    Q_PROPERTY(int Count READ Count NOTIFY CountChanged)

public:
    enum Role {
        IdRole = Qt::UserRole + 1,
        SlugRole,
        NameRole,
        CuisineRole,
        DescriptionRole,
        TotalMinutesRole,
        ServingsRole,
        DifficultyRole,
        ImagePathRole,
        MediaUrlRole,
        StatusRole,
        IngredientsRole,
        StepsRole,
    };
    Q_ENUM(Role)

    explicit RecipeModel(QObject* Parent = nullptr);

    [[nodiscard]] int Count() const noexcept;
    [[nodiscard]] int rowCount(
        const QModelIndex& Parent = QModelIndex()) const override;
    [[nodiscard]] QVariant data(
        const QModelIndex& Index,
        int RoleValue = Qt::DisplayRole) const override;
    [[nodiscard]] QHash<int, QByteArray> roleNames() const override;

    void SetRecipes(const QJsonArray& RecipesValue);
    Q_INVOKABLE QVariantMap RecipeAt(int Index) const;

signals:
    void CountChanged();

private:
    [[nodiscard]] static const QHash<int, QByteArray>& RoleNamesMap();

    QVector<QVariantMap> Recipes;
};

}  // namespace Menu::Client
