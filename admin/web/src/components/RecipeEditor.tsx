import { CircleAlert, Check, Plus, Trash2, X } from "lucide-react";
import { useState } from "react";
import {
  AdminApi,
  ApiError,
  type Ingredient,
  type Recipe,
  type RecipeIngredient,
  type RecipeStep,
} from "../api";
import { RecipePreview } from "./RecipePreview";

export function RecipeEditor({
  isNew,
  token,
  recipe: InitialRecipe,
  ingredients,
  onClose,
  onSaved,
  onDeleted,
}: {
  isNew: boolean;
  token: string;
  recipe: Recipe;
  ingredients: Ingredient[];
  onClose: () => void;
  onSaved: (recipe: Recipe) => void;
  onDeleted: (id: string) => void;
}) {
  const [Recipe, SetRecipe] = useState<Recipe>(InitialRecipe);
  const [Tab, SetTab] = useState<"edit" | "preview">("edit");
  const [Busy, SetBusy] = useState(false);
  const [Error, SetError] = useState("");
  const [DeleteOpen, SetDeleteOpen] = useState(false);

  function Update<Key extends keyof Recipe>(KeyValue: Key, Value: Recipe[Key]) {
    SetRecipe((Current) => ({
      ...Current,
      [KeyValue]: Value,
      ...(KeyValue === "prepMinutes" || KeyValue === "cookMinutes"
        ? {
            totalMinutes:
              (KeyValue === "prepMinutes" ? Number(Value) : Current.prepMinutes) +
              (KeyValue === "cookMinutes" ? Number(Value) : Current.cookMinutes),
          }
        : {}),
    }));
  }

  function UpdateIngredient(Index: number, Value: Partial<RecipeIngredient>) {
    SetRecipe((Current) => ({
      ...Current,
      ingredients: Current.ingredients.map((Item, ItemIndex) =>
        ItemIndex === Index ? { ...Item, ...Value } : Item,
      ),
    }));
  }

  function UpdateStep(Index: number, Value: Partial<RecipeStep>) {
    SetRecipe((Current) => ({
      ...Current,
      steps: Current.steps.map((Item, ItemIndex) =>
        ItemIndex === Index ? { ...Item, ...Value } : Item,
      ),
    }));
  }

  async function Save() {
    SetBusy(true);
    SetError("");
    try {
      const Result = isNew
        ? await AdminApi.createRecipe(token, Recipe)
        : await AdminApi.updateRecipe(token, Recipe);
      onSaved(Result);
    } catch (Caught) {
      SetError(Caught instanceof ApiError ? Caught.message : "保存失败，请稍后重试");
    } finally {
      SetBusy(false);
    }
  }

  async function Delete() {
    SetBusy(true);
    SetError("");
    try {
      await AdminApi.deleteRecipe(token, Recipe.id);
      onDeleted(Recipe.id);
    } catch (Caught) {
      SetError(Caught instanceof ApiError ? Caught.message : "删除失败");
    } finally {
      SetBusy(false);
    }
  }

  return <div className="drawer-backdrop" onMouseDown={(Event) => { if (Event.target === Event.currentTarget) onClose(); }}>
    <aside className="editor-drawer">
      <header className="editor-header"><div><p className="section-kicker">{isNew ? "新建菜谱" : "编辑菜谱"}</p><h2>{Recipe.name || "未命名菜谱"}</h2></div>
        <button className="icon-button" aria-label="关闭编辑器" onClick={onClose}><X size={19} /></button>
      </header>
      <div className="editor-tabs"><button className={Tab === "edit" ? "active" : ""} onClick={() => SetTab("edit")}>编辑</button><button className={Tab === "preview" ? "active" : ""} onClick={() => SetTab("preview")}>预览</button></div>
      {Tab === "preview" ? <RecipePreview recipe={Recipe} ingredients={ingredients} /> : <div className="editor-scroll">
        <section className="form-section"><h3>基本信息</h3><div className="form-grid">
          <label className="span-2">名称<input value={Recipe.name} onChange={(Event) => Update("name", Event.target.value)} placeholder="例如：番茄炒蛋" /></label>
          <label>菜系<select value={Recipe.cuisine} onChange={(Event) => Update("cuisine", Event.target.value)}><option>中餐</option><option>西餐</option><option>日系</option></select></label>
          <label>状态<select value={Recipe.status} onChange={(Event) => Update("status", Event.target.value as Recipe["status"])}><option value="draft">草稿</option><option value="published">已发布</option><option value="archived">已归档</option></select></label>
          <label>准备（分钟）<input type="number" min="0" value={Recipe.prepMinutes} onChange={(Event) => Update("prepMinutes", Number(Event.target.value))} /></label>
          <label>烹饪（分钟）<input type="number" min="0" value={Recipe.cookMinutes} onChange={(Event) => Update("cookMinutes", Number(Event.target.value))} /></label>
          <label>份量<input type="number" min="1" value={Recipe.servings} onChange={(Event) => Update("servings", Number(Event.target.value))} /></label>
          <label>难度<select value={Recipe.difficulty} onChange={(Event) => Update("difficulty", Number(Event.target.value))}><option value="1">简单</option><option value="2">适中</option><option value="3">需要专注</option><option value="4">较难</option><option value="5">挑战</option></select></label>
          <label className="span-2">简介<textarea rows={3} value={Recipe.description} onChange={(Event) => Update("description", Event.target.value)} placeholder="一句话说清这道菜为什么适合今晚" /></label>
        </div></section>
        <section className="form-section"><div className="section-line"><h3>食材 <span>{Recipe.ingredients.length}</span></h3><button className="small-action" onClick={() => SetRecipe((Current) => ({ ...Current, ingredients: [...Current.ingredients, { ingredientId: ingredients[0]?.id ?? "", quantity: 1, unit: ingredients[0]?.defaultUnit ?? "g", required: true, servingFactor: 1, preparation: "" }] }))}><Plus size={15} />添加</button></div>
          <div className="ingredient-editor">{Recipe.ingredients.map((Item, Index) => <div className="ingredient-line" key={`${Item.ingredientId}-${Index}`}>
            <select value={Item.ingredientId} onChange={(Event) => { const SelectedIngredient = ingredients.find((IngredientValue) => IngredientValue.id === Event.target.value); UpdateIngredient(Index, { ingredientId: Event.target.value, unit: SelectedIngredient?.defaultUnit ?? Item.unit }); }}><option value="">选择食材</option>{ingredients.map((IngredientValue) => <option value={IngredientValue.id} key={IngredientValue.id}>{IngredientValue.name}</option>)}</select>
            <input className="quantity-input" type="number" min="0.01" step="0.01" value={Item.quantity} onChange={(Event) => UpdateIngredient(Index, { quantity: Number(Event.target.value) })} />
            <input className="unit-input" value={Item.unit} onChange={(Event) => UpdateIngredient(Index, { unit: Event.target.value })} />
            <label className="required-toggle"><input type="checkbox" checked={Item.required} onChange={(Event) => UpdateIngredient(Index, { required: Event.target.checked })} />必需</label>
            <button className="icon-button danger" aria-label="删除食材" onClick={() => SetRecipe((Current) => ({ ...Current, ingredients: Current.ingredients.filter((_, ItemIndex) => ItemIndex !== Index) }))}><Trash2 size={15} /></button>
          </div>)}</div>
        </section>
        <section className="form-section"><div className="section-line"><h3>步骤 <span>{Recipe.steps.length}</span></h3><button className="small-action" onClick={() => SetRecipe((Current) => ({ ...Current, steps: [...Current.steps, { stepOrder: Current.steps.length + 1, title: "新步骤", instruction: "", durationSeconds: 60, hasTimer: false }] }))}><Plus size={15} />添加</button></div>
          <div className="step-editor">{Recipe.steps.map((Step, Index) => <div className="step-line" key={Step.stepOrder}><span className="step-number">{Index + 1}</span><div className="step-fields">
            <input value={Step.title} onChange={(Event) => UpdateStep(Index, { title: Event.target.value })} placeholder="步骤标题" />
            <textarea rows={2} value={Step.instruction} onChange={(Event) => UpdateStep(Index, { instruction: Event.target.value })} placeholder="写清楚动作、状态和判断标准" />
            <div className="step-meta"><label>耗时（秒）<input type="number" min="0" value={Step.durationSeconds} onChange={(Event) => UpdateStep(Index, { durationSeconds: Number(Event.target.value) })} /></label><label className="required-toggle"><input type="checkbox" checked={Step.hasTimer} onChange={(Event) => UpdateStep(Index, { hasTimer: Event.target.checked })} />需要计时</label></div>
          </div><button className="icon-button danger" aria-label="删除步骤" onClick={() => SetRecipe((Current) => ({ ...Current, steps: Current.steps.filter((_, ItemIndex) => ItemIndex !== Index).map((Item, ItemIndex) => ({ ...Item, stepOrder: ItemIndex + 1 })) }))}><Trash2 size={15} /></button></div>)}</div>
        </section>
        {Error && <div className="inline-error drawer-error"><CircleAlert size={16} />{Error}</div>}
      </div>}
      <footer className="editor-footer">{!isNew && <button className="danger-button" onClick={() => SetDeleteOpen(true)} disabled={Busy}><Trash2 size={16} />删除</button>}
        <span className="footer-spacer" /><button className="secondary-button" onClick={onClose}>取消</button><button className="primary-button" onClick={() => void Save()} disabled={Busy || !Recipe.name}>{Busy ? "保存中..." : <><Check size={16} />保存菜谱</>}</button>
      </footer>
      {DeleteOpen && <div className="confirm-backdrop" role="presentation"><section className="confirm-dialog" role="dialog" aria-modal="true" aria-labelledby="delete-dialog-title"><p className="section-kicker">删除确认</p><h3 id="delete-dialog-title">删除这道菜谱？</h3><p>删除后管理台和用户端都不会再显示它。</p>{Error && <div className="inline-error"><CircleAlert size={16} />{Error}</div>}<div className="confirm-actions"><button className="secondary-button" onClick={() => { SetDeleteOpen(false); SetError(""); }} disabled={Busy}>取消</button><button className="danger-button" onClick={() => void Delete()} disabled={Busy}>{Busy ? "删除中..." : <><Trash2 size={16} />确认删除</>}</button></div></section></div>}
    </aside>
  </div>;
}
