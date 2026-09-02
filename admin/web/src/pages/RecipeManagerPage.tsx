import { BookOpen, FilePenLine, FlaskConical, Plus, Search } from "lucide-react";
import { useMemo, useState } from "react";
import type { Recipe, RecipeStatus } from "../api";
import { PageHeading, StatusBadge } from "../Components/AdminPrimitives";
import { FilterRecipes } from "../uiModel";

export function RecipeManagerPage({
  recipes,
  onCreate,
  onEdit,
}: {
  recipes: Recipe[];
  onCreate: () => void;
  onEdit: (recipe: Recipe) => void;
}) {
  const [Query, SetQuery] = useState("");
  const [Status, SetStatus] = useState<RecipeStatus | "all">("all");
  const Filtered = useMemo(
    () => FilterRecipes(recipes, Query, Status),
    [recipes, Query, Status],
  );

  return <div className="page-wrap">
    <PageHeading eyebrow="内容管理" title="菜谱" action={
      <button className="primary-button" onClick={onCreate}><Plus size={17} />新增菜谱</button>
    } />
    <section className="content-band table-band">
      <div className="toolbar">
        <label className="search-field"><Search size={17} /><input value={Query} onChange={(Event) => SetQuery(Event.target.value)} placeholder="搜索菜谱名称或菜系" /></label>
        <div className="filter-group">
          {(["all", "published", "draft", "archived"] as const).map((Item) =>
            <button key={Item} className={Status === Item ? "selected" : ""} onClick={() => SetStatus(Item)}>
              {Item === "all" ? "全部" : Item === "published" ? "已发布" : Item === "draft" ? "草稿" : "已归档"}
            </button>)}
        </div>
      </div>
      <div className="table-scroll"><table><thead><tr><th>菜谱</th><th>菜系</th><th>时间</th><th>难度</th><th>状态</th><th aria-label="操作" /></tr></thead>
        <tbody>{Filtered.map((RecipeValue) => <tr key={RecipeValue.id} onClick={() => onEdit(RecipeValue)}>
          <td><div className="recipe-cell"><div className="recipe-thumb"><BookOpen size={17} /></div><div><strong>{RecipeValue.name || "未命名菜谱"}</strong><span>{RecipeValue.ingredients.length} 种食材 · {RecipeValue.servings} 人份</span></div></div></td>
          <td>{RecipeValue.cuisine}</td><td>{RecipeValue.totalMinutes} 分钟</td>
          <td><span className="difficulty-dots">{[1, 2, 3, 4, 5].map((Value) => <i className={Value <= RecipeValue.difficulty ? "on" : ""} key={Value} />)}</span></td>
          <td><StatusBadge status={RecipeValue.status} /></td>
          <td><button className="icon-button" aria-label={`编辑 ${RecipeValue.name}`} onClick={(Event) => { Event.stopPropagation(); onEdit(RecipeValue); }}><FilePenLine size={17} /></button></td>
        </tr>)}</tbody>
      </table></div>
      {Filtered.length === 0 && <div className="empty-inline"><FlaskConical size={22} /><span>没有符合条件的菜谱</span></div>}
    </section>
  </div>;
}
