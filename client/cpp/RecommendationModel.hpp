#pragma once

#include <QAbstractListModel>
#include <QJsonArray>
#include <QVariantMap>

namespace Menu::Client {

class RecommendationModel final : public QAbstractListModel {
    Q_OBJECT
    Q_PROPERTY(int Count READ Count NOTIFY CountChanged)

public:
    enum Role {
        IdRole = Qt::UserRole + 1,
        NameRole,
        CuisineRole,
        DescriptionRole,
        TotalMinutesRole,
        ServingsRole,
        DifficultyRole,
        ImagePathRole,
        MediaUrlRole,
        AvailableIngredientCountRole,
        MissingIngredientCountRole,
        ScoreRole,
        IngredientsRole,
        StepsRole,
    };
    Q_ENUM(Role)

    explicit RecommendationModel(QObject* Parent = nullptr);

    [[nodiscard]] int Count() const noexcept;
    [[nodiscard]] int rowCount(
        const QModelIndex& Parent = QModelIndex()) const override;
    [[nodiscard]] QVariant data(
        const QModelIndex& Index,
        int RoleValue = Qt::DisplayRole) const override;
    [[nodiscard]] QHash<int, QByteArray> roleNames() const override;

    void SetRecommendations(const QJsonArray& RecommendationsValue);
    Q_INVOKABLE QVariantMap RecommendationAt(int Index) const;

signals:
    void CountChanged();

private:
    [[nodiscard]] static const QHash<int, QByteArray>& RoleNamesMap();

    QVector<QVariantMap> Recommendations;
};

}  // namespace Menu::Client
