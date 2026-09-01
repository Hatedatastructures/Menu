import type { Recipe, RecipeStatus } from "./api";

export function FilterRecipes(
  recipes: Recipe[],
  query: string,
  status: RecipeStatus | "all",
): Recipe[] {
  const normalizedQuery = query.trim().toLocaleLowerCase();
  return recipes.filter((recipe) => {
    const matchesStatus = status === "all" || recipe.status === status;
    const searchable = `${recipe.name}${recipe.cuisine}`.toLocaleLowerCase();
    return matchesStatus && searchable.includes(normalizedQuery);
  });
}
