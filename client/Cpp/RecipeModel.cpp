#include "RecipeModel.hpp"

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

RecipeModel::RecipeModel(QObject* Parent)
    : QAbstractListModel(Parent) {}

int RecipeModel::Count() const noexcept {
    return Recipes.size();
}

int RecipeModel::rowCount(const QModelIndex& Parent) const {
    if (Parent.isValid()) {
        return 0;
    }
    return Recipes.size();
}

QVariant RecipeModel::data(const QModelIndex& Index, int RoleValue) const {
    if (!Index.isValid() || Index.row() < 0 || Index.row() >= Recipes.size()) {
        return {};
    }
    const QVariantMap& RecipeValue = Recipes.at(Index.row());
    const auto RoleName = RoleNamesMap().value(RoleValue);
    return RecipeValue.value(QString::fromUtf8(RoleName));
}

QHash<int, QByteArray> RecipeModel::roleNames() const {
    return RoleNamesMap();
}

const QHash<int, QByteArray>& RecipeModel::RoleNamesMap() {
    static const QHash<int, QByteArray> Roles = {
        {IdRole, "id"},
        {SlugRole, "slug"},
        {NameRole, "name"},
        {CuisineRole, "cuisine"},
        {DescriptionRole, "description"},
        {TotalMinutesRole, "totalMinutes"},
        {ServingsRole, "servings"},
        {DifficultyRole, "difficulty"},
        {ImagePathRole, "imagePath"},
        {MediaUrlRole, "mediaUrl"},
        {StatusRole, "status"},
        {IngredientsRole, "ingredients"},
        {StepsRole, "steps"},
    };
    return Roles;
}

void RecipeModel::SetRecipes(const QJsonArray& RecipesValue) {
    QVector<QVariantMap> NextRecipes;
    NextRecipes.reserve(RecipesValue.size());
    for (const QJsonValue& Value : RecipesValue) {
        if (!Value.isObject()) {
            continue;
        }
        QVariantMap RecipeMap = Value.toObject().toVariantMap();
        RecipeMap.insert(QStringLiteral("mediaUrl"), MediaUrl(RecipeMap));
        if (!RecipeMap.contains(QStringLiteral("totalMinutes"))) {
            RecipeMap.insert(
                QStringLiteral("totalMinutes"),
                RecipeMap.value(QStringLiteral("prepMinutes")).toInt() +
                    RecipeMap.value(QStringLiteral("cookMinutes")).toInt());
        }
        NextRecipes.push_back(std::move(RecipeMap));
    }

    beginResetModel();
    Recipes = std::move(NextRecipes);
    endResetModel();
    emit CountChanged();
}

QVariantMap RecipeModel::RecipeAt(int Index) const {
    if (Index < 0 || Index >= Recipes.size()) {
        return {};
    }
    return Recipes.at(Index);
}

}  // namespace Menu::Client
