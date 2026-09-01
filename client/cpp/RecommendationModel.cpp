#include "RecommendationModel.hpp"

#include <QJsonObject>

#include <utility>

namespace Menu::Client {
namespace {

QString MediaUrl(const QVariantMap& Recipe) {
    const QString ImagePath = Recipe.value(QStringLiteral("imagePath")).toString();
    if (ImagePath.startsWith(QStringLiteral("assets/media/")) &&
        !ImagePath.contains(QStringLiteral(".."))) {
        return QStringLiteral("qrc:/qt/qml/MenuClient/") + ImagePath;
    }
    return QStringLiteral("qrc:/qt/qml/MenuClient/assets/media/menu-placeholder.png");
}

}  // namespace

RecommendationModel::RecommendationModel(QObject* Parent)
    : QAbstractListModel(Parent) {}

int RecommendationModel::Count() const noexcept {
    return Recommendations.size();
}

int RecommendationModel::rowCount(const QModelIndex& Parent) const {
    return Parent.isValid() ? 0 : Recommendations.size();
}

QVariant RecommendationModel::data(const QModelIndex& Index, int RoleValue) const {
    if (!Index.isValid() || Index.row() < 0 || Index.row() >= Recommendations.size()) {
        return {};
    }
    const QVariantMap& Recommendation = Recommendations.at(Index.row());
    return Recommendation.value(QString::fromUtf8(RoleNamesMap().value(RoleValue)));
}

QHash<int, QByteArray> RecommendationModel::roleNames() const {
    return RoleNamesMap();
}

const QHash<int, QByteArray>& RecommendationModel::RoleNamesMap() {
    static const QHash<int, QByteArray> Roles = {
        {IdRole, "id"},
        {NameRole, "name"},
        {CuisineRole, "cuisine"},
        {DescriptionRole, "description"},
        {TotalMinutesRole, "totalMinutes"},
        {ServingsRole, "servings"},
        {DifficultyRole, "difficulty"},
        {ImagePathRole, "imagePath"},
        {MediaUrlRole, "mediaUrl"},
        {AvailableIngredientCountRole, "availableIngredientCount"},
        {MissingIngredientCountRole, "missingIngredientCount"},
        {ScoreRole, "score"},
        {IngredientsRole, "ingredients"},
        {StepsRole, "steps"},
    };
    return Roles;
}

void RecommendationModel::SetRecommendations(
    const QJsonArray& RecommendationsValue) {
    QVector<QVariantMap> NextRecommendations;
    NextRecommendations.reserve(RecommendationsValue.size());
    for (const QJsonValue& Value : RecommendationsValue) {
        if (!Value.isObject()) {
            continue;
        }
        const QVariantMap RecommendationMap = Value.toObject().toVariantMap();
        QVariantMap RecipeMap = RecommendationMap.value(QStringLiteral("recipe")).toMap();
        if (RecipeMap.isEmpty()) {
            continue;
        }
        RecipeMap.insert(QStringLiteral("mediaUrl"), MediaUrl(RecipeMap));
        RecipeMap.insert(
            QStringLiteral("availableIngredientCount"),
            RecommendationMap.value(QStringLiteral("availableIngredientCount")));
        RecipeMap.insert(
            QStringLiteral("missingIngredientCount"),
            RecommendationMap.value(QStringLiteral("missingIngredientCount")));
        RecipeMap.insert(
            QStringLiteral("score"), RecommendationMap.value(QStringLiteral("score")));
        NextRecommendations.push_back(std::move(RecipeMap));
    }

    beginResetModel();
    Recommendations = std::move(NextRecommendations);
    endResetModel();
    emit CountChanged();
}

QVariantMap RecommendationModel::RecommendationAt(int Index) const {
    if (Index < 0 || Index >= Recommendations.size()) {
        return {};
    }
    return Recommendations.at(Index);
}

}  // namespace Menu::Client
