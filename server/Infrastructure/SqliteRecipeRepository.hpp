#pragma once

#include <Application/AdminRecipeRepository.hpp>
#include <Application/RecipeRepository.hpp>
#include <Infrastructure/SqliteDatabase.hpp>

namespace Menu::Infrastructure {

class SqliteRecipeRepository final
    : public Application::RecipeRepository,
      public Application::AdminRecipeRepository {
public:
    explicit SqliteRecipeRepository(SqliteDatabase& DatabaseValue)
        : Database(DatabaseValue) {}

    Foundation::Result<std::vector<Domain::Recipe>> ListPublished() override;

    Foundation::Result<std::optional<Domain::Recipe>> FindPublishedById(
        std::string_view Id) override;

    Foundation::Result<std::vector<Domain::Recipe>> ListAll() override;

    Foundation::Result<Domain::Recipe> Create(
        const Domain::Recipe& RecipeValue) override;

    Foundation::Result<Domain::Recipe> Update(
        std::string_view Id,
        const Domain::Recipe& RecipeValue) override;

    Foundation::Result<void> Delete(std::string_view Id) override;

private:
    SqliteDatabase& Database;
};

}  // namespace Menu::Infrastructure
