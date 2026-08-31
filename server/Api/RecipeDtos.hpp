#pragma once

#include <Domain/Ingredient.hpp>
#include <Domain/Recipe.hpp>
#include <Domain/Recommendation.hpp>

#include <boost/json/array.hpp>
#include <boost/json/object.hpp>

namespace Menu::Api {

class RecipeDtos final {
public:
    static boost::json::object ToObject(const Domain::Recipe& RecipeValue);
    static boost::json::object ToObject(const Domain::Ingredient& IngredientValue);
    static boost::json::object ToObject(const Domain::Recommendation& RecommendationValue);
    static boost::json::array ToArray(const std::vector<Domain::Recipe>& Recipes);
    static boost::json::array ToArray(const std::vector<Domain::Ingredient>& Ingredients);
    static boost::json::array ToArray(
        const std::vector<Domain::Recommendation>& Recommendations);
};

}  // namespace Menu::Api
