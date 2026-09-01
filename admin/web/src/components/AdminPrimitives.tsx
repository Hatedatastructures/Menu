import {
  BookOpen,
  Check,
  CircleAlert,
  FilePenLine,
  LayoutDashboard,
  ShoppingBasket,
} from "lucide-react";
import type { ReactNode } from "react";
import type { Recipe, RecipeStatus } from "../api";

export function NavItem({
  icon,
  label,
  active,
  onClick,
}: {
  icon: ReactNode;
  label: string;
  active: boolean;
  onClick: () => void;
}) {
  return <button className={`nav-item ${active ? "active" : ""}`} onClick={onClick}>
    {icon}<span>{label}</span>{active && <span className="nav-indicator" />}
  </button>;
}

export function LoadingState() {
  return <section className="state-panel"><div className="spinner" /><span>正在加载工作台</span></section>;
}

export function ErrorState({ message, onRetry }: { message: string; onRetry: () => void }) {
  return <section className="state-panel error-state">
    <CircleAlert size={24} /><strong>{message}</strong>
    <button className="secondary-button" onClick={onRetry}>重新加载</button>
  </section>;
}

export function Metric({
  icon,
  label,
  value,
  accent,
}: {
  icon: ReactNode;
  label: string;
  value: number | string;
  accent: string;
}) {
  return <div className={`metric metric-${accent}`}>
    <span className="metric-icon">{icon}</span>
    <div><span className="metric-label">{label}</span><strong>{value}</strong></div>
  </div>;
}

export function PageHeading({
  eyebrow,
  title,
  action,
}: {
  eyebrow: string;
  title: string;
  action?: ReactNode;
}) {
  return <div className="page-heading"><div><p className="section-kicker">{eyebrow}</p><h1>{title}</h1></div>{action}</div>;
}

export function StatusBadge({ status }: { status: RecipeStatus }) {
  return <span className={`status-badge status-${status}`}><i />
    {status === "published" ? "已发布" : status === "draft" ? "草稿" : "已归档"}
  </span>;
}

export function RecipeRow({ recipe, onClick }: { recipe: Recipe; onClick: () => void }) {
  return <button className="recipe-row" onClick={onClick}>
    <div className="recipe-row-main"><span className="recipe-row-icon"><BookOpen size={16} /></span>
      <span><strong>{recipe.name}</strong><small>{recipe.cuisine} · {recipe.totalMinutes} 分钟</small></span>
    </div><StatusBadge status={recipe.status} />
  </button>;
}

export function EmptyIngredientState() {
  return <div className="empty-inline"><ShoppingBasket size={22} /><span>没有符合条件的食材</span></div>;
}

export const NavigationIcons = {
  Overview: <LayoutDashboard size={18} />,
  Recipes: <BookOpen size={18} />,
  Ingredients: <ShoppingBasket size={18} />,
  Published: <Check size={14} />,
  RecipeEdit: <FilePenLine size={17} />,
};
