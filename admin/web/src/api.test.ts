import { describe, expect, it } from "vitest";
import { SerializeRecipe } from "./api";
import type { Recipe } from "./api";

describe("SerializeRecipe", () => {
  it("removes display-only recipe and ingredient fields before writes", () => {
    const recipe = {
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
      imagePath: "assets/media/one.png",
      status: "draft",
      allergens: [],
      cookware: [],
      ingredients: [{
        ingredientId: "ingredient.tomato",
        ingredientName: "番茄",
        ingredientCategory: "蔬菜",
        ingredientDefaultUnit: "g",
        ingredientIsPantryStaple: false,
        quantity: 200,
        unit: "g",
        required: true,
        servingFactor: 1,
        preparation: "切块",
      }],
      steps: [],
    } as Recipe;

    const payload = SerializeRecipe(recipe) as Record<string, unknown>;
    const ingredient = (payload.ingredients as Array<Record<string, unknown>>)[0];

    expect(payload).not.toHaveProperty("totalMinutes");
    expect(ingredient).not.toHaveProperty("ingredientName");
    expect(ingredient).not.toHaveProperty("ingredientCategory");
    expect(ingredient).toMatchObject({ ingredientId: "ingredient.tomato", quantity: 200 });
  });
});
