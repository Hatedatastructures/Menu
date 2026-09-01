import { describe, expect, it } from "vitest";
import { FilterRecipes } from "./uiModel";
import type { Recipe } from "./api";

const Recipes: Recipe[] = [
  {
    id: "recipe.one",
    slug: "one",
    name: "番茄炒蛋",
    cuisine: "中餐",
    description: "",
    prepMinutes: 5,
    cookMinutes: 10,
    totalMinutes: 15,
    servings: 2,
    difficulty: 1,
    imagePath: "",
    status: "published",
    allergens: [],
    cookware: [],
    ingredients: [],
    steps: [],
  },
  {
    id: "recipe.two",
    slug: "two",
    name: "柠檬三文鱼",
    cuisine: "西餐",
    description: "",
    prepMinutes: 5,
    cookMinutes: 20,
    totalMinutes: 25,
    servings: 2,
    difficulty: 2,
    imagePath: "",
    status: "draft",
    allergens: [],
    cookware: [],
    ingredients: [],
    steps: [],
  },
];

describe("FilterRecipes", () => {
  it("filters by text and publication status", () => {
    expect(FilterRecipes(Recipes, "番茄", "published").map((recipe) => recipe.id)).toEqual(["recipe.one"]);
    expect(FilterRecipes(Recipes, "", "draft").map((recipe) => recipe.id)).toEqual(["recipe.two"]);
  });
});
