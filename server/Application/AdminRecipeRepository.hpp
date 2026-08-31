#pragma once

#include <Foundation/Result.hpp>
#include <Domain/Recipe.hpp>

#include <string_view>
#include <vector>

namespace Menu::Application {

class AdminRecipeRepository {
public:
    virtual ~AdminRecipeRepository() = default;

    virtual Foundation::Result<std::vector<Domain::Recipe>> ListAll() = 0;

    virtual Foundation::Result<Domain::Recipe> Create(
        const Domain::Recipe& RecipeValue) = 0;

    virtual Foundation::Result<Domain::Recipe> Update(
        std::string_view Id,
        const Domain::Recipe& RecipeValue) = 0;

    virtual Foundation::Result<void> Delete(std::string_view Id) = 0;
};

}  // namespace Menu::Application
