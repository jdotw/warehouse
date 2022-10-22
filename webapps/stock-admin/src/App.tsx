import { useAuth0 } from "@auth0/auth0-react";
import { Outlet } from "react-router-dom";
import "./App.css";

function App() {
  const { isAuthenticated, isLoading, loginWithRedirect } = useAuth0();

  if (isLoading) {
    return <div>Loading ...</div>;
  }

  return isAuthenticated ? (
    <div className="App">
      <h1>Stock Admin</h1>
      <Outlet />
    </div>
  ) : (
    <div className="App">
      <h1>Stock Admin</h1>
      <button onClick={() => loginWithRedirect()}>Log In</button>
    </div>
  );
}

export default App;
