import { chromium } from "../admin/web/node_modules/playwright/index.mjs";

const BaseUrl = "http://127.0.0.1:5173";
const Email = process.env.MENU_E2E_EMAIL ?? "codex-admin-20260901@example.com";
const Password = process.env.MENU_E2E_PASSWORD ?? "Menu-Cook-Local-2026!";
const RecipeName = `Playwright验收菜谱-${Date.now()}`;
const IngredientName = `Playwright验收食材-${Date.now()}`;
const IngredientAlias = `验收香草-${Date.now()}`;

async function ExpectVisible(Page, Text) {
  await Page.getByText(Text, { exact: true }).first().waitFor({ state: "visible" });
}

async function SaveScreenshot(Page, Width, Height, Path) {
  await Page.setViewportSize({ width: Width, height: Height });
  await Page.screenshot({ path: Path });
  const Overflow = await Page.evaluate(() => ({
    Horizontal: document.documentElement.scrollWidth > document.documentElement.clientWidth,
    Vertical: document.documentElement.scrollHeight > document.documentElement.clientHeight,
  }));
  if (Overflow.Horizontal) {
    throw new Error(`页面存在横向溢出：${Width}x${Height}`);
  }
}

async function RemovePreviousPlaywrightRecipes(Page) {
  await Page.getByRole("button", { name: "菜谱", exact: true }).click();
  const Rows = Page.getByRole("row").filter({ hasText: "Playwright验收菜谱" });
  while (await Rows.count() > 0) {
    await Rows.first().getByRole("button").first().click();
    await Page.getByRole("button", { name: "删除", exact: true }).click();
    await Page.getByRole("dialog", { name: "删除这道菜谱？" }).getByRole("button", { name: "确认删除" }).click();
    await Rows.first().waitFor({ state: "detached" });
  }
  await Page.getByRole("button", { name: "概览", exact: true }).click();
}

async function RemovePreviousPlaywrightIngredients(Page) {
  await Page.getByRole("button", { name: "食材", exact: true }).click();
  const Rows = Page.getByRole("row").filter({ hasText: "Playwright验收食材" });
  while (await Rows.count() > 0) {
    await Rows.first().getByRole("button").first().click();
    await Page.getByRole("button", { name: "删除", exact: true }).click();
    await Page.getByRole("dialog", { name: "删除这项食材？" }).getByRole("button", { name: "确认删除" }).click();
    await Rows.first().waitFor({ state: "detached" });
  }
  await Page.getByRole("button", { name: "概览", exact: true }).click();
}

const Browser = await chromium.launch({ headless: true });
const Page = await Browser.newPage({ viewport: { width: 1440, height: 900 } });
try {
  await Page.goto(BaseUrl, { waitUntil: "networkidle" });
  if (await Page.getByRole("button", { name: "登录管理台" }).isVisible()) {
    await Page.getByRole("textbox", { name: "邮箱" }).fill(Email);
    await Page.getByRole("textbox", { name: "密码" }).fill(Password);
    await Page.getByRole("button", { name: "登录管理台" }).click();
    await Page.getByRole("heading", { name: "今天，把内容准备好" }).waitFor();
  }
  await ExpectVisible(Page, "今天，把内容准备好");
  await RemovePreviousPlaywrightRecipes(Page);
  await RemovePreviousPlaywrightIngredients(Page);
  await ExpectVisible(Page, "今天，把内容准备好");

  await SaveScreenshot(Page, 1440, 900, "I:/code/Menu/docs/evidence-admin-1440x900.png");
  await SaveScreenshot(Page, 1280, 800, "I:/code/Menu/docs/evidence-admin-1280x800.png");
  await SaveScreenshot(Page, 390, 844, "I:/code/Menu/docs/evidence-admin-390x844.png");

  await Page.getByRole("button", { name: "菜谱", exact: true }).click();
  await Page.getByRole("button", { name: "新增菜谱" }).click();
  await Page.getByRole("textbox", { name: "名称" }).fill(RecipeName);
  await Page.getByRole("textbox", { name: "简介" }).fill("Playwright 本地真实 CRUD 验收");
  await Page.getByRole("button", { name: "添加", exact: true }).nth(0).click();
  await Page.getByRole("combobox").nth(3).selectOption({ label: "番茄" });
  await Page.getByRole("textbox", { name: "写清楚动作、状态和判断标准" }).fill("清洗番茄并切块");
  await Page.getByRole("button", { name: "保存菜谱" }).click();
  await Page.getByRole("row", { name: new RegExp(`${RecipeName}.*草稿`) }).waitFor();

  await Page.getByRole("button", { name: `编辑 ${RecipeName}` }).click();
  await Page.getByRole("combobox", { name: "状态" }).selectOption({ label: "已发布" });
  await Page.getByRole("button", { name: "保存菜谱" }).click();
  await Page.getByRole("row", { name: new RegExp(`${RecipeName}.*已发布`) }).waitFor();

  await Page.getByRole("button", { name: `编辑 ${RecipeName}` }).click();
  await Page.getByRole("button", { name: "预览", exact: true }).click();
  await ExpectVisible(Page, "执行步骤");
  const PreviewText = await Page.locator(".preview-scroll").innerText();
  if (!PreviewText.includes("番茄") || PreviewText.includes("ingredient.")) {
    throw new Error("菜谱预览未使用食材中文显示字段");
  }
  await Page.screenshot({ path: "I:/code/Menu/docs/evidence-admin-playwright-preview.png" });
  await Page.getByRole("button", { name: "删除", exact: true }).click();
  const DeleteDialog = Page.getByRole("dialog", { name: "删除这道菜谱？" });
  await DeleteDialog.getByRole("button", { name: "取消", exact: true }).click();
  if (await DeleteDialog.isVisible()) {
    throw new Error("取消删除后确认框仍可见");
  }
  await ExpectVisible(Page, RecipeName);

  await Page.getByRole("button", { name: "删除", exact: true }).click();
  await Page.getByRole("dialog", { name: "删除这道菜谱？" }).getByRole("button", { name: "确认删除" }).click();
  await Page.getByRole("row", { name: new RegExp(RecipeName) }).waitFor({ state: "detached" });

  await Page.route("**/api/v1/admin/recipes", (Route) => Route.abort());
  await Page.reload({ waitUntil: "domcontentloaded" });
  await ExpectVisible(Page, "无法加载工作台");
  await ExpectVisible(Page, "重新加载");
  await Page.unroute("**/api/v1/admin/recipes");

  await Page.reload({ waitUntil: "networkidle" });
  await Page.getByRole("button", { name: "食材", exact: true }).click();
  await Page.getByRole("button", { name: "新增食材" }).click();
  await Page.getByRole("textbox", { name: "稳定 ID" }).fill(`ingredient.playwright-${Date.now()}`);
  await Page.getByRole("textbox", { name: "名称" }).fill(IngredientName);
  await Page.getByRole("textbox", { name: "类别" }).fill("验收香料");
  await Page.getByRole("textbox", { name: "别名（用顿号分隔）" }).fill(IngredientAlias);
  await Page.getByRole("button", { name: "保存食材" }).click();
  const IngredientRow = Page.getByRole("row", { name: new RegExp(IngredientName) });
  await IngredientRow.waitFor();
  await Page.getByRole("button", { name: "关闭食材编辑器" }).waitFor({ state: "detached" });
  await Page.getByRole("button", { name: `编辑 ${IngredientName}` }).click();
  await Page.getByRole("checkbox", { name: "家中常备" }).check();
  await Page.getByRole("button", { name: "保存食材" }).click();
  await IngredientRow.waitFor();
  await Page.getByRole("button", { name: "关闭食材编辑器" }).waitFor({ state: "detached" });
  await Page.getByRole("button", { name: `编辑 ${IngredientName}` }).click();
  await Page.getByRole("button", { name: "删除", exact: true }).click();
  const IngredientDialog = Page.getByRole("dialog", { name: "删除这项食材？" });
  await IngredientDialog.getByRole("button", { name: "取消", exact: true }).click();
  if (await IngredientDialog.isVisible()) {
    throw new Error("取消食材删除后确认框仍可见");
  }
  await Page.getByRole("button", { name: "删除", exact: true }).click();
  await Page.getByRole("dialog", { name: "删除这项食材？" }).getByRole("button", { name: "确认删除" }).click();
  await IngredientRow.waitFor({ state: "detached" });

  process.stdout.write(JSON.stringify({
    Login: true,
    CreateDraft: true,
    Publish: true,
    Preview: true,
    CancelDelete: true,
    ConfirmDelete: true,
    ErrorState: true,
    IngredientCrud: true,
  }));
} finally {
  await Browser.close();
}
