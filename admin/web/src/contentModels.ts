import type { Ingredient, Recipe } from "./api";

export const EmptyRecipe: Recipe = {
  id: "recipe.new",
  slug: "new-recipe",
  name: "",
  cuisine: "中餐",
  description: "",
  prepMinutes: 10,
  cookMinutes: 20,
  totalMinutes: 30,
  servings: 2,
  difficulty: 1,
  imagePath: "assets/media/menu-placeholder.png",
  status: "draft",
  allergens: [],
  cookware: ["炒锅"],
  ingredients: [],
  steps: [{
    stepOrder: 1,
    title: "准备",
    instruction: "",
    durationSeconds: 60,
    hasTimer: false,
  }],
};

export const EmptyIngredient: Ingredient = {
  id: "ingredient.new",
  name: "",
  aliases: [],
  category: "蔬菜",
  defaultUnit: "g",
  isPantryStaple: false,
  substituteGroup: "",
  storeSkuMapping: "",
};
