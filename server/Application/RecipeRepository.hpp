#pragma once

#include <Foundation/Result.hpp>
#include <Domain/Recipe.hpp>

#include <optional>
#include <string_view>
#include <vector>

namespace Menu::Application {

class RecipeRepository {
public:
    virtual ~RecipeRepository() = default;

    virtual Foundation::Result<std::vector<Domain::Recipe>> ListPublished() = 0;

    virtual Foundation::Result<std::optional<Domain::Recipe>> FindPublishedById(
        std::string_view Id) = 0;
};

}  // namespace Menu::Application
