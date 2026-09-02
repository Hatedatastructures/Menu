#pragma once

#include <Domain/Ingredient.hpp>
#include <Domain/Recipe.hpp>

#include <vector>

namespace Menu::Infrastructure {

class SeedCatalog final {
public:
    [[nodiscard]] static std::vector<Domain::Ingredient> Ingredients();
    [[nodiscard]] static std::vector<Domain::Recipe> Recipes();
};

}  // namespace Menu::Infrastructure
