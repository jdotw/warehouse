import { useAuth0 } from "@auth0/auth0-react";
import React, { useEffect, useState } from "react";
import { Alert, Button, Table } from "react-bootstrap";
import Loading from "../components/Loading";
import AddCategory from "./AddCategory";
import "./CategoryList.css";

type Category = {
  id: string;
  name: string;
};

const CategoryList = () => {
  const { user, isAuthenticated, getAccessTokenSilently } = useAuth0();
  const [loadingCategories, setLoadingCategories] = useState(true);
  const [categories, setCategories] = useState<Category[]>();
  const [categoryLoadError, setCategoryLoadError] = useState<Error>();
  const [showAdd, setShowAdd] = useState(false);

  const loadCategories = async () => {
    const domain = "localhost:8080";
    try {
      const accessToken = await getAccessTokenSilently({
        audience: `http://${domain}/`,
        scope: "read:category",
      });
      console.log("TOKEN: ", accessToken);

      const addCategoryURL = `/api/categories`;

      const response = await fetch(addCategoryURL, {
        headers: {
          Authorization: `Bearer ${accessToken}`,
        },
      });

      const response_json = await response.json();
      console.log("JSON: ", response_json);
      setCategories(response_json);
      setCategoryLoadError(undefined);
      setLoadingCategories(false);
    } catch (error: any) {
      console.log("ERROR: ", error.message);
      setCategories([]);
      setCategoryLoadError(error);
      setLoadingCategories(false);
    }
  };

  useEffect(() => {
    loadCategories();
  }, []);

  if (loadingCategories) {
    return <Loading />;
  }

  const onAdded = (err?: Error) => {
    setShowAdd(false);
    loadCategories();
  };

  return (
    <div className="CategoryList">
      <div className="TopButtonBar">
        <Button onClick={() => setShowAdd(true)}>Create Category</Button>
      </div>
      {showAdd && <AddCategory onAdded={(err?: Error) => onAdded(err)} />}
      {categoryLoadError ? (
        <div>Failed to load Categories: {categoryLoadError.message}</div>
      ) : (
        <Table striped bordered hover>
          <thead>
            <tr>
              <th>Name</th>
            </tr>
          </thead>
          <tbody>
            {categories!.map((c) => (
              <tr key={c.id}>
                <td>{c.name}</td>
              </tr>
            ))}
          </tbody>
        </Table>
      )}
    </div>
  );
};

export default CategoryList;
