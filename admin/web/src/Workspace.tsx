import { BookOpen, LayoutDashboard, LogOut, Menu as MenuIcon, ShoppingBasket } from "lucide-react";
import { useEffect, useState } from "react";
import { AdminApi, ApiError, IngredientApi, type Ingredient, type Recipe } from "./api";
import { IngredientEditor } from "./Components/IngredientEditor";
import { RecipeEditor } from "./Components/RecipeEditor";
import { ErrorState, LoadingState, NavItem } from "./Components/AdminPrimitives";
import { EmptyRecipe } from "./contentModels";
import { IngredientManagerPage } from "./Pages/IngredientManagerPage";
import { OverviewPage } from "./Pages/OverviewPage";
import { RecipeManagerPage } from "./Pages/RecipeManagerPage";

type View = "overview" | "recipes" | "ingredients";

export function Workspace({
  token,
  userName,
  onSignOut,
}: {
  token: string;
  userName: string;
  onSignOut: () => void;
}) {
  const [ViewValue, SetView] = useState<View>("overview");
  const [Recipes, SetRecipes] = useState<Recipe[]>([]);
  const [Ingredients, SetIngredients] = useState<Ingredient[]>([]);
  const [Loading, SetLoading] = useState(true);
  const [Error, SetError] = useState("");
  const [Selected, SetSelected] = useState<Recipe | null>(null);
  const [EditorOpen, SetEditorOpen] = useState(false);
  const [Creating, SetCreating] = useState(false);

  async function LoadWorkspace() {
    SetLoading(true);
    SetError("");
    try {
      const [RecipeResult, IngredientResult] = await Promise.all([
        AdminApi.listRecipes(token),
        IngredientApi.list(),
      ]);
      SetRecipes(RecipeResult);
      SetIngredients(IngredientResult);
    } catch (Caught) {
      SetError(Caught instanceof ApiError ? Caught.message : "无法加载工作台");
    } finally {
      SetLoading(false);
    }
  }

  useEffect(() => { void LoadWorkspace(); }, [token]);

  function OpenCreate() {
    const Timestamp = Date.now();
    SetSelected({ ...EmptyRecipe, id: `recipe.new-${Timestamp}`, slug: `new-recipe-${Timestamp}` });
    SetCreating(true);
    SetEditorOpen(true);
  }

  function CloseEditor() {
    SetCreating(false);
    SetEditorOpen(false);
    SetSelected(null);
  }

  return <div className="app-shell">
    <aside className="sidebar"><div className="brand-lockup"><span className="brand-mark"><MenuIcon size={18} /></span><span>Menu</span></div>
      <div className="workspace-label">内容工作台</div><nav className="side-nav" aria-label="主导航">
        <NavItem icon={<LayoutDashboard size={18} />} label="概览" active={ViewValue === "overview"} onClick={() => SetView("overview")} />
        <NavItem icon={<BookOpen size={18} />} label="菜谱" active={ViewValue === "recipes"} onClick={() => SetView("recipes")} />
        <NavItem icon={<ShoppingBasket size={18} />} label="食材" active={ViewValue === "ingredients"} onClick={() => SetView("ingredients")} />
      </nav>
      <div className="sidebar-bottom"><div className="account-row"><span className="avatar">{userName.slice(0, 1) || "管"}</span><span className="account-name">{userName || "管理员"}</span></div>
        <button className="side-action" onClick={onSignOut}><LogOut size={17} />退出登录</button>
      </div>
    </aside>
    <main className="main-content"><header className="topbar"><div className="mobile-brand"><MenuIcon size={19} />Menu</div><div className="topbar-meta"><span className="connection-dot" />服务正常</div></header>
      {Loading ? <LoadingState /> : Error ? <ErrorState message={Error} onRetry={() => void LoadWorkspace()} /> :
        ViewValue === "overview" ? <OverviewPage recipes={Recipes} ingredients={Ingredients} onOpenRecipes={() => SetView("recipes")} /> :
          ViewValue === "recipes" ? <RecipeManagerPage recipes={Recipes} onCreate={OpenCreate} onEdit={(RecipeValue) => { SetCreating(false); SetSelected(RecipeValue); SetEditorOpen(true); }} /> :
            <IngredientManagerPage token={token} ingredients={Ingredients} onChanged={() => void LoadWorkspace()} />}
    </main>
    {EditorOpen && Selected && <RecipeEditor isNew={Creating} token={token} recipe={Selected} ingredients={Ingredients} onClose={CloseEditor} onSaved={(RecipeValue) => { SetRecipes((Current) => { const Index = Current.findIndex((Item) => Item.id === RecipeValue.id); if (Index < 0) return [RecipeValue, ...Current]; const Next = [...Current]; Next[Index] = RecipeValue; return Next; }); CloseEditor(); }} onDeleted={(RecipeId) => { SetRecipes((Current) => Current.filter((Item) => Item.id !== RecipeId)); CloseEditor(); }} />}
  </div>;
}
