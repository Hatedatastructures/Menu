#pragma once

#include <Application/AdminRecipeRepository.hpp>

#include <memory>
#include <string_view>

namespace Menu::Application {

class AdminRecipeApplicationService final {
public:
    explicit AdminRecipeApplicationService(
        std::unique_ptr<AdminRecipeRepository> RepositoryValue);

    AdminRecipeApplicationService(const AdminRecipeApplicationService&) = delete;
    AdminRecipeApplicationService& operator=(const AdminRecipeApplicationService&) = delete;

    Foundation::Result<std::vector<Domain::Recipe>> ListAllRecipes();

    Foundation::Result<Domain::Recipe> CreateRecipe(const Domain::Recipe& RecipeValue);

    Foundation::Result<Domain::Recipe> UpdateRecipe(
        std::string_view Id,
        const Domain::Recipe& RecipeValue);

    Foundation::Result<void> DeleteRecipe(std::string_view Id);

private:
    std::unique_ptr<AdminRecipeRepository> Repository;
};

}  // namespace Menu::Application
