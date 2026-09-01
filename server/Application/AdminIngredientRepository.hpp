#pragma once

#include <Domain/Ingredient.hpp>
#include <Foundation/Result.hpp>

#include <string_view>
#include <vector>

namespace Menu::Application {

class AdminIngredientRepository {
public:
    virtual ~AdminIngredientRepository() = default;

    virtual Foundation::Result<std::vector<Domain::Ingredient>> ListAll() = 0;

    virtual Foundation::Result<Domain::Ingredient> Create(
        const Domain::Ingredient& IngredientValue) = 0;

    virtual Foundation::Result<Domain::Ingredient> Update(
        std::string_view Id,
        const Domain::Ingredient& IngredientValue) = 0;

    virtual Foundation::Result<void> Delete(std::string_view Id) = 0;
};

}  // namespace Menu::Application
