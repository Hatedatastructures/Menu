export type RecipeStatus = "draft" | "published" | "archived";

export type Ingredient = {
  id: string;
  name: string;
  aliases: string[];
  category: string;
  defaultUnit: string;
  isPantryStaple: boolean;
  substituteGroup: string;
  storeSkuMapping: string;
};

export type RecipeIngredient = {
  ingredientId: string;
  quantity: number;
  unit: string;
  required: boolean;
  servingFactor: number;
  preparation: string;
  ingredientName?: string;
  ingredientCategory?: string;
  ingredientDefaultUnit?: string;
  ingredientIsPantryStaple?: boolean;
};

export type RecipeStep = {
  stepOrder: number;
  title: string;
  instruction: string;
  durationSeconds: number;
  hasTimer: boolean;
};

export type Recipe = {
  id: string;
  slug: string;
  name: string;
  cuisine: string;
  description: string;
  prepMinutes: number;
  cookMinutes: number;
  totalMinutes: number;
  servings: number;
  difficulty: number;
  imagePath: string;
  status: RecipeStatus;
  allergens: string[];
  cookware: string[];
  ingredients: RecipeIngredient[];
  steps: RecipeStep[];
};

type ApiErrorBody = {
  error?: { code?: string; message?: string; requestId?: string };
};

export class ApiError extends Error {
  readonly status: number;
  readonly code: string;

  constructor(status: number, code: string, message: string) {
    super(message);
    this.name = "ApiError";
    this.status = status;
    this.code = code;
  }
}

export function SerializeRecipe(recipe: Recipe) {
  return {
    id: recipe.id,
    slug: recipe.slug,
    name: recipe.name,
    cuisine: recipe.cuisine,
    description: recipe.description,
    prepMinutes: recipe.prepMinutes,
    cookMinutes: recipe.cookMinutes,
    servings: recipe.servings,
    difficulty: recipe.difficulty,
    imagePath: recipe.imagePath,
    status: recipe.status,
    allergens: recipe.allergens,
    cookware: recipe.cookware,
    ingredients: recipe.ingredients.map((ingredient) => ({
      ingredientId: ingredient.ingredientId,
      quantity: ingredient.quantity,
      unit: ingredient.unit,
      required: ingredient.required,
      servingFactor: ingredient.servingFactor,
      preparation: ingredient.preparation,
    })),
    steps: recipe.steps.map((step) => ({
      stepOrder: step.stepOrder,
      title: step.title,
      instruction: step.instruction,
      durationSeconds: step.durationSeconds,
      hasTimer: step.hasTimer,
    })),
  };
}

export function SerializeIngredient(ingredient: Ingredient) {
  return {
    id: ingredient.id,
    name: ingredient.name,
    aliases: ingredient.aliases,
    category: ingredient.category,
    defaultUnit: ingredient.defaultUnit,
    isPantryStaple: ingredient.isPantryStaple,
    substituteGroup: ingredient.substituteGroup,
    storeSkuMapping: ingredient.storeSkuMapping,
  };
}

const ApiBaseUrl = import.meta.env.VITE_API_BASE_URL ?? "";

async function Request<T>(path: string, init?: RequestInit): Promise<T> {
  const response = await fetch(`${ApiBaseUrl}${path}`, {
    ...init,
    headers: {
      "Content-Type": "application/json",
      ...(init?.headers ?? {}),
    },
  });
  const text = await response.text();
  let body: unknown = undefined;
  if (text) {
    try {
      body = JSON.parse(text);
    } catch {
      body = undefined;
    }
  }
  if (!response.ok) {
    const errorBody = (body ?? {}) as ApiErrorBody;
    if (response.status === 401) {
      window.dispatchEvent(new Event("menu-auth-expired"));
    }
    throw new ApiError(
      response.status,
      errorBody.error?.code ?? "request_failed",
      errorBody.error?.message ?? "请求失败",
    );
  }
  return body as T;
}

export type AuthResponse = {
  user: { id: string; email: string; displayName: string; isAdmin: boolean };
  accessToken: string;
  refreshToken: string;
  accessTokenExpiresInSeconds: number;
};

export const AuthApi = {
  login(email: string, password: string) {
    return Request<AuthResponse>("/api/v1/auth/login", {
      method: "POST",
      body: JSON.stringify({ email, password }),
    });
  },
  register(email: string, password: string, displayName: string) {
    return Request<AuthResponse>("/api/v1/auth/register", {
      method: "POST",
      body: JSON.stringify({ email, password, displayName }),
    });
  },
};

function AuthorizedInit(token: string, init?: RequestInit): RequestInit {
  return {
    ...init,
    headers: {
      ...(init?.headers ?? {}),
      Authorization: `Bearer ${token}`,
    },
  };
}

export const AdminApi = {
  listRecipes(token: string) {
    return Request<Recipe[]>("/api/v1/admin/recipes", AuthorizedInit(token));
  },
  createRecipe(token: string, recipe: Recipe) {
    return Request<Recipe>("/api/v1/admin/recipes", AuthorizedInit(token, {
      method: "POST",
      body: JSON.stringify(SerializeRecipe(recipe)),
    }));
  },
  updateRecipe(token: string, recipe: Recipe) {
    return Request<Recipe>(`/api/v1/admin/recipes/${encodeURIComponent(recipe.id)}`, AuthorizedInit(token, {
      method: "PATCH",
      body: JSON.stringify(SerializeRecipe(recipe)),
    }));
  },
  deleteRecipe(token: string, id: string) {
    return Request<void>(`/api/v1/admin/recipes/${encodeURIComponent(id)}`, AuthorizedInit(token, {
      method: "DELETE",
    }));
  },
};

export const IngredientApi = {
  list() {
    return Request<Ingredient[]>("/api/v1/ingredients");
  },
  listAdmin(token: string) {
    return Request<Ingredient[]>("/api/v1/admin/ingredients", AuthorizedInit(token));
  },
  create(token: string, ingredient: Ingredient) {
    return Request<Ingredient>("/api/v1/admin/ingredients", AuthorizedInit(token, {
      method: "POST",
      body: JSON.stringify(SerializeIngredient(ingredient)),
    }));
  },
  update(token: string, ingredient: Ingredient) {
    return Request<Ingredient>(`/api/v1/admin/ingredients/${encodeURIComponent(ingredient.id)}`, AuthorizedInit(token, {
      method: "PATCH",
      body: JSON.stringify(SerializeIngredient(ingredient)),
    }));
  },
  delete(token: string, id: string) {
    return Request<void>(`/api/v1/admin/ingredients/${encodeURIComponent(id)}`, AuthorizedInit(token, {
      method: "DELETE",
    }));
  },
};
