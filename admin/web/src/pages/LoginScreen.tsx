import { CircleAlert, Menu as MenuIcon } from "lucide-react";
import { useState } from "react";
import { ApiError, AuthApi, type AuthResponse } from "../api";

export function LoginScreen({
  onAuthenticated,
}: {
  onAuthenticated: (value: AuthResponse) => void;
}) {
  const [mode, setMode] = useState<"login" | "register">("login");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [displayName, setDisplayName] = useState("");
  const [error, setError] = useState("");
  const [busy, setBusy] = useState(false);

  async function Submit(event: React.FormEvent) {
    event.preventDefault();
    setBusy(true);
    setError("");
    try {
      const Result = mode === "login"
        ? await AuthApi.login(email, password)
        : await AuthApi.register(email, password, displayName);
      if (!Result.user.isAdmin) {
        setError("该账号没有管理权限");
        return;
      }
      onAuthenticated(Result);
    } catch (Caught) {
      setError(Caught instanceof ApiError ? Caught.message : "网络暂时不可用");
    } finally {
      setBusy(false);
    }
  }

  return <main className="auth-shell">
    <section className="auth-panel" aria-label="Menu 管理台登录">
      <div className="brand-lockup"><span className="brand-mark"><MenuIcon size={18} /></span><span>Menu</span></div>
      <p className="eyebrow">家庭做饭执行器 · 内容工作台</p>
      <h1>{mode === "login" ? "进入管理台" : "创建首个管理员"}</h1>
      <p className="auth-caption">管理菜谱、食材和发布状态</p>
      <div className="mode-switch" role="tablist" aria-label="登录方式">
        <button className={mode === "login" ? "active" : ""} onClick={() => setMode("login")} type="button">登录</button>
        <button className={mode === "register" ? "active" : ""} onClick={() => setMode("register")} type="button">创建账户</button>
      </div>
      <form onSubmit={Submit} className="auth-form">
        {mode === "register" && <label>显示名称<input value={displayName} onChange={(Event) => setDisplayName(Event.target.value)} placeholder="例如：Menu 厨房" required /></label>}
        <label>邮箱<input type="email" value={email} onChange={(Event) => setEmail(Event.target.value)} placeholder="you@example.com" required /></label>
        <label>密码<input type="password" minLength={12} value={password} onChange={(Event) => setPassword(Event.target.value)} placeholder="至少 12 个字符" required /></label>
        {error && <div className="inline-error"><CircleAlert size={16} />{error}</div>}
        <button className="primary-button wide" disabled={busy} type="submit">
          {busy ? "处理中..." : mode === "login" ? "登录管理台" : "创建并进入"}
        </button>
      </form>
    </section>
  </main>;
}
