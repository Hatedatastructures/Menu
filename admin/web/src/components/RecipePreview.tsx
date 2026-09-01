import { Clock3, Users } from "lucide-react";
import type { Ingredient, Recipe } from "../api";
import { StatusBadge } from "./AdminPrimitives";

export function RecipePreview({ recipe, ingredients }: { recipe: Recipe; ingredients: Ingredient[] }) {
  return <div className="preview-scroll">
    <div className="preview-hero"><span className="preview-cuisine">{recipe.cuisine}</span><h3>{recipe.name || "未命名菜谱"}</h3>
      <p>{recipe.description || "还没有简介"}</p><div className="preview-meta"><span><Clock3 size={15} />{recipe.totalMinutes} 分钟</span><span><Users size={15} />{recipe.servings} 人份</span><StatusBadge status={recipe.status} /></div>
    </div>
    <div className="preview-section"><p className="section-kicker">食材清单</p>
      {recipe.ingredients.length ? recipe.ingredients.map((Item, Index) => {
        const Catalog = ingredients.find((IngredientValue) => IngredientValue.id === Item.ingredientId);
        const DisplayName = Item.ingredientName || Catalog?.name || "未选择食材";
        const DisplayCategory = Item.ingredientCategory || Catalog?.category || "";
        return <div className="preview-ingredient" key={`${Item.ingredientId}-${Index}`}><span><strong>{DisplayName}</strong>{DisplayCategory && <small>{DisplayCategory}</small>}</span><strong>{Item.quantity} {Item.unit}</strong></div>;
      }) : <p className="muted-cell">尚未添加食材</p>}
    </div>
    <div className="preview-section"><p className="section-kicker">执行步骤</p>{recipe.steps.map((Step, Index) =>
      <div className="preview-step" key={Step.stepOrder}><span>{Index + 1}</span><div><strong>{Step.title}</strong><p>{Step.instruction || "尚未填写步骤说明"}</p></div></div>)}
    </div>
  </div>;
}
