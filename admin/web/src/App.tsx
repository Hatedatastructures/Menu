import { useEffect, useState } from "react";
import { LoginScreen } from "./Pages/LoginScreen";
import { Workspace } from "./Workspace";

function App() {
  const [Token, SetToken] = useState(() => sessionStorage.getItem("menu-access-token"));
  const [UserName, SetUserName] = useState(() => sessionStorage.getItem("menu-user-name") ?? "");

  useEffect(() => {
    const Expire = () => {
      sessionStorage.removeItem("menu-access-token");
      sessionStorage.removeItem("menu-user-name");
      SetToken(null);
      SetUserName("");
    };
    window.addEventListener("menu-auth-expired", Expire);
    return () => window.removeEventListener("menu-auth-expired", Expire);
  }, []);

  if (!Token) {
    return <LoginScreen onAuthenticated={(Value) => {
      sessionStorage.setItem("menu-access-token", Value.accessToken);
      sessionStorage.setItem("menu-user-name", Value.user.displayName);
      SetToken(Value.accessToken);
      SetUserName(Value.user.displayName);
    }} />;
  }

  return <Workspace token={Token} userName={UserName} onSignOut={() => {
    sessionStorage.clear();
    SetToken(null);
    SetUserName("");
  }} />;
}

export default App;
