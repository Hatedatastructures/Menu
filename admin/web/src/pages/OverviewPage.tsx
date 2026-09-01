import { BookOpen, Check, Clock3, FilePenLine, Send, ShoppingBasket, ChevronDown } from "lucide-react";
import type { Ingredient, Recipe } from "../api";
import { Metric, PageHeading, RecipeRow } from "../components/AdminPrimitives";

export function OverviewPage({
  recipes,
  ingredients,
  onOpenRecipes,
}: {
  recipes: Recipe[];
  ingredients: Ingredient[];
  onOpenRecipes: () => void;
}) {
  const Published = recipes.filter((RecipeValue) => RecipeValue.status === "published").length;
  const Drafts = recipes.filter((RecipeValue) => RecipeValue.status === "draft").length;
  const AverageMinutes = recipes.length
    ? Math.round(recipes.reduce((Sum, RecipeValue) => Sum + RecipeValue.totalMinutes, 0) / recipes.length)
    : 0;

  return <div className="page-wrap">
    <PageHeading eyebrow="工作台概览" title="今天，把内容准备好" action={
      <button className="primary-button" onClick={onOpenRecipes}><BookOpen size={17} />查看菜谱</button>
    } />
    <section className="metric-strip">
      <Metric icon={<Send size={17} />} label="已发布菜谱" value={Published} accent="green" />
      <Metric icon={<FilePenLine size={17} />} label="待编辑草稿" value={Drafts} accent="coral" />
      <Metric icon={<ShoppingBasket size={17} />} label="规范化食材" value={ingredients.length} accent="ink" />
      <Metric icon={<Clock3 size={17} />} label="平均制作时长" value={`${AverageMinutes} 分钟`} accent="gold" />
    </section>
    <div className="overview-grid">
      <section className="content-band">
        <div className="band-heading"><div><p className="section-kicker">最近更新</p><h2>菜谱状态</h2></div>
          <button className="text-button" onClick={onOpenRecipes}>全部菜谱<ChevronDown size={15} /></button>
        </div>
        {recipes.slice(0, 5).map((RecipeValue) => <RecipeRow key={RecipeValue.id} recipe={RecipeValue} onClick={onOpenRecipes} />)}
      </section>
      <section className="content-band prep-band">
        <p className="section-kicker">发布前检查</p><h2>让每道菜都能被执行</h2>
        <ul className="check-list">
          <li><span className="check-icon"><Check size={14} /></span>食材有中文显示字段</li>
          <li><span className="check-icon"><Check size={14} /></span>步骤顺序和耗时完整</li>
          <li><span className="check-icon"><Check size={14} /></span>发布状态可追踪</li>
        </ul>
      </section>
    </div>
  </div>;
}
