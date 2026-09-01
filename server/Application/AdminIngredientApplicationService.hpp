#pragma once

#include <Application/AdminIngredientRepository.hpp>

#include <memory>
#include <string_view>

namespace Menu::Application {

class AdminIngredientApplicationService final {
public:
    explicit AdminIngredientApplicationService(
        std::unique_ptr<AdminIngredientRepository> RepositoryValue);

    Foundation::Result<std::vector<Domain::Ingredient>> ListAll();

    Foundation::Result<Domain::Ingredient> Create(
        const Domain::Ingredient& IngredientValue);

    Foundation::Result<Domain::Ingredient> Update(
        std::string_view Id,
        const Domain::Ingredient& IngredientValue);

    Foundation::Result<void> Delete(std::string_view Id);

private:
    std::unique_ptr<AdminIngredientRepository> Repository;
};

}  // namespace Menu::Application
