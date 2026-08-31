#pragma once

#include <Domain/Recipe.hpp>

#include <initializer_list>

namespace Menu::Tests::RecipeFixtures {

inline Domain::Recipe SingleRecipeWithIngredients(
    int Servings,
    std::initializer_list<Domain::RecipeIngredient> Ingredients) {
    Domain::Recipe Recipe;
    Recipe.Id = "recipe.fixture";
    Recipe.Slug = "recipe-fixture";
    Recipe.Name = "测试菜谱";
    Recipe.Cuisine = "中餐";
    Recipe.Description = "测试用菜谱";
    Recipe.PrepMinutes = 10;
    Recipe.CookMinutes = 20;
    Recipe.Servings = Servings;
    Recipe.Difficulty = 2;
    Recipe.Status = "published";
    Recipe.Ingredients.assign(Ingredients.begin(), Ingredients.end());
    Recipe.Steps = {
        Domain::RecipeStep{1, "准备", "准备食材", 60, false},
        Domain::RecipeStep{2, "烹饪", "完成烹饪", 1200, true},
    };
    Recipe.Cookware = {"炒锅"};
    return Recipe;
}

}  // namespace Menu::Tests::RecipeFixtures
