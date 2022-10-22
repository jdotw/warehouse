import { Auth0Provider } from "@auth0/auth0-react";
import ReactDOM from "react-dom/client";
import { BrowserRouter, Route, Routes } from "react-router-dom";
import App from "./App";
import AddCategory from "./Categories/AddCategory";
import Categories from "./Categories/Categories";
import Category from "./Categories/Category";
import CategoryList from "./Categories/CategoryList";
import Home from "./Home/Home";
import "./index.css";

ReactDOM.createRoot(document.getElementById("root") as HTMLElement).render(
  <Auth0Provider
    domain="warehouse.au.auth0.com"
    clientId="irEI1dVJObviJ67mZsYFEjP93EIlthys"
    redirectUri={window.location.origin}
    useRefreshTokens={true}
    cacheLocation="localstorage"
    scope="write:category"
    audience="http://localhost:8080/"
  >
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<App />}>
          <Route index element={<Home />} />
          <Route path="categories" element={<Categories />}>
            <Route path=":categoryID" element={<Category />} />
            <Route path="new" element={<AddCategory />} />
            <Route index element={<CategoryList />} />
          </Route>
        </Route>
      </Routes>
    </BrowserRouter>
  </Auth0Provider>
);
