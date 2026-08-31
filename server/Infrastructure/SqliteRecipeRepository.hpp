#pragma once

#include <Application/RecipeRepository.hpp>
#include <Infrastructure/SqliteDatabase.hpp>

namespace Menu::Infrastructure {

class SqliteRecipeRepository final : public Application::RecipeRepository {
public:
    explicit SqliteRecipeRepository(SqliteDatabase& DatabaseValue)
        : Database(DatabaseValue) {}

    Foundation::Result<std::vector<Domain::Recipe>> ListPublished() override;

    Foundation::Result<std::optional<Domain::Recipe>> FindPublishedById(
        std::string_view Id) override;

private:
    SqliteDatabase& Database;
};

}  // namespace Menu::Infrastructure
