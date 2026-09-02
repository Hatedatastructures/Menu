import { Check, FilePenLine, Plus, Search, ShoppingBasket } from "lucide-react";
import { useState } from "react";
import type { Ingredient } from "../api";
import { PageHeading } from "../Components/AdminPrimitives";
import { EmptyIngredient, } from "../contentModels";
import { IngredientEditor } from "../Components/IngredientEditor";

export function IngredientManagerPage({
  token,
  ingredients,
  onChanged,
}: {
  token: string;
  ingredients: Ingredient[];
  onChanged: () => void;
}) {
  const [Query, SetQuery] = useState("");
  const [Selected, SetSelected] = useState<Ingredient | null>(null);
  const [EditorOpen, SetEditorOpen] = useState(false);
  const Filtered = ingredients.filter((Item) =>
    (Item.name + Item.category + Item.aliases.join("")).toLocaleLowerCase().includes(Query.toLocaleLowerCase()),
  );

  function OpenCreate() {
    SetSelected({ ...EmptyIngredient, id: `ingredient.new-${Date.now()}` });
    SetEditorOpen(true);
  }

  return <div className="page-wrap">
    <PageHeading eyebrow="基础数据" title="食材" action={<button className="primary-button" onClick={OpenCreate}><Plus size={17} />新增食材</button>} />
    <section className="content-band table-band">
      <div className="toolbar"><label className="search-field"><Search size={17} /><input value={Query} onChange={(Event) => SetQuery(Event.target.value)} placeholder="搜索食材、别名或类别" /></label><span className="toolbar-count">{Filtered.length} 项</span></div>
      <div className="table-scroll"><table><thead><tr><th>名称</th><th>别名</th><th>类别</th><th>默认单位</th><th>常备</th><th>替代组</th><th aria-label="操作" /></tr></thead>
        <tbody>{Filtered.map((Item) => <tr key={Item.id}><td><strong>{Item.name}</strong></td><td className="muted-cell">{Item.aliases.join("、") || "暂无"}</td><td>{Item.category}</td><td>{Item.defaultUnit}</td><td>{Item.isPantryStaple ? <span className="pantry-yes"><Check size={14} />常备</span> : <span className="muted-cell">按需</span>}</td><td className="muted-cell">{Item.substituteGroup || "暂无"}</td><td><button className="icon-button" aria-label={`编辑 ${Item.name}`} onClick={() => { SetSelected(Item); SetEditorOpen(true); }}><FilePenLine size={17} /></button></td></tr>)}</tbody>
      </table></div>
      {Filtered.length === 0 && <div className="empty-inline"><ShoppingBasket size={22} /><span>没有符合条件的食材</span></div>}
    </section>
    {EditorOpen && Selected && <IngredientEditor isNew={Selected.id.startsWith("ingredient.new-")} token={token} ingredient={Selected} onClose={() => SetEditorOpen(false)} onSaved={() => { SetEditorOpen(false); onChanged(); }} onDeleted={() => { SetEditorOpen(false); onChanged(); }} />}
  </div>;
}
