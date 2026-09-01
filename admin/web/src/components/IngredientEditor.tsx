import { CircleAlert, Check, Trash2, X } from "lucide-react";
import { useState } from "react";
import { ApiError, IngredientApi, type Ingredient } from "../api";

export function IngredientEditor({
  isNew,
  token,
  ingredient: InitialIngredient,
  onClose,
  onSaved,
  onDeleted,
}: {
  isNew: boolean;
  token: string;
  ingredient: Ingredient;
  onClose: () => void;
  onSaved: (ingredient: Ingredient) => void;
  onDeleted: (id: string) => void;
}) {
  const [Ingredient, SetIngredient] = useState<Ingredient>(InitialIngredient);
  const [Busy, SetBusy] = useState(false);
  const [Error, SetError] = useState("");
  const [DeleteOpen, SetDeleteOpen] = useState(false);

  function Update<Key extends keyof Ingredient>(KeyValue: Key, Value: Ingredient[Key]) {
    SetIngredient((Current) => ({ ...Current, [KeyValue]: Value }));
  }

  async function Save() {
    SetBusy(true);
    SetError("");
    try {
      const Result = isNew
        ? await IngredientApi.create(token, Ingredient)
        : await IngredientApi.update(token, Ingredient);
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
      await IngredientApi.delete(token, Ingredient.id);
      onDeleted(Ingredient.id);
    } catch (Caught) {
      SetError(Caught instanceof ApiError ? Caught.message : "删除失败");
    } finally {
      SetBusy(false);
    }
  }

  return <div className="drawer-backdrop" onMouseDown={(Event) => { if (Event.target === Event.currentTarget) onClose(); }}>
    <aside className="editor-drawer">
      <header className="editor-header"><div><p className="section-kicker">{isNew ? "新增食材" : "编辑食材"}</p><h2>{Ingredient.name || "未命名食材"}</h2></div>
        <button className="icon-button" aria-label="关闭食材编辑器" onClick={onClose}><X size={19} /></button>
      </header>
      <div className="editor-scroll"><section className="form-section"><h3>食材信息</h3><div className="form-grid">
        <label className="span-2">稳定 ID<input value={Ingredient.id} disabled={!isNew} onChange={(Event) => Update("id", Event.target.value)} /></label>
        <label className="span-2">名称<input value={Ingredient.name} onChange={(Event) => Update("name", Event.target.value)} placeholder="例如：番茄" /></label>
        <label>类别<input value={Ingredient.category} onChange={(Event) => Update("category", Event.target.value)} placeholder="蔬菜" /></label>
        <label>默认单位<select value={Ingredient.defaultUnit} onChange={(Event) => Update("defaultUnit", Event.target.value)}><option>g</option><option>kg</option><option>ml</option><option>l</option><option>个</option><option>份</option><option>勺</option><option>茶匙</option></select></label>
        <label className="span-2">别名（用顿号分隔）<input value={Ingredient.aliases.join("、")} onChange={(Event) => Update("aliases", Event.target.value.split("、").map((Value) => Value.trim()).filter(Boolean))} placeholder="西红柿、番茄果" /></label>
        <label>替代组<input value={Ingredient.substituteGroup} onChange={(Event) => Update("substituteGroup", Event.target.value)} placeholder="可选" /></label>
        <label>未来 SKU 映射<input value={Ingredient.storeSkuMapping} onChange={(Event) => Update("storeSkuMapping", Event.target.value)} placeholder="可选" /></label>
        <label className="required-toggle"><input type="checkbox" checked={Ingredient.isPantryStaple} onChange={(Event) => Update("isPantryStaple", Event.target.checked)} />家中常备</label>
      </div></section>{Error && <div className="inline-error drawer-error"><CircleAlert size={16} />{Error}</div>}</div>
      <footer className="editor-footer">{!isNew && <button className="danger-button" onClick={() => SetDeleteOpen(true)} disabled={Busy}><Trash2 size={16} />删除</button>}
        <span className="footer-spacer" /><button className="secondary-button" onClick={onClose}>取消</button>
        <button className="primary-button" onClick={() => void Save()} disabled={Busy || !Ingredient.name}>{Busy ? "保存中..." : <><Check size={16} />保存食材</>}</button>
      </footer>
      {DeleteOpen && <div className="confirm-backdrop" role="presentation"><section className="confirm-dialog" role="dialog" aria-modal="true" aria-labelledby="ingredient-delete-dialog-title">
        <p className="section-kicker">删除确认</p><h3 id="ingredient-delete-dialog-title">删除这项食材？</h3><p>如果仍被菜谱使用，服务端会拒绝删除。</p>
        {Error && <div className="inline-error"><CircleAlert size={16} />{Error}</div>}
        <div className="confirm-actions"><button className="secondary-button" onClick={() => { SetDeleteOpen(false); SetError(""); }} disabled={Busy}>取消</button>
          <button className="danger-button" onClick={() => void Delete()} disabled={Busy}>{Busy ? "删除中..." : <><Trash2 size={16} />确认删除</>}</button>
        </div>
      </section></div>}
    </aside>
  </div>;
}
