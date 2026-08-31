#pragma once

#include <Foundation/Result.hpp>
#include <Domain/Ingredient.hpp>

#include <vector>

namespace Menu::Application {

class IngredientRepository {
public:
    virtual ~IngredientRepository() = default;

    virtual Foundation::Result<std::vector<Domain::Ingredient>> ListAll() = 0;
};

}  // namespace Menu::Application
