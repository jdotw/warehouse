import { useAuth0 } from "@auth0/auth0-react";
import React, { useEffect, useState } from "react";
import { Button, Table } from "react-bootstrap";
import { Link, useParams } from "react-router-dom";
import Loading from "../components/Loading";
import AddItem from "./AddItem";
import "./Category.css";

type Item = {
  id: string;
  name: string;
  categoryID: string;
};

type Category = {
  id: string;
  name: string;
};

const Category = () => {
  const { categoryID } = useParams();
  const { getAccessTokenSilently } = useAuth0();
  const [category, setCategory] = useState<Category>();
  const [loadingCategory, setLoadingCategory] = useState(true);
  const [items, setItems] = useState<Item[]>();
  const [loadingItems, setLoadingItems] = useState(true);
  const [showAdd, setShowAdd] = useState(false);
  const [itemLoadError, setItemLoadError] = useState<Error>();

  const loadCategory = async (categoryID: string) => {
    const domain = "localhost:8080";
    try {
      const accessToken = await getAccessTokenSilently({
        audience: `http://${domain}/`,
        scope: "read:category",
      });
      const categoryURL = `/api/categories/${categoryID}`;
      const response = await fetch(categoryURL, {
        headers: {
          Authorization: `Bearer ${accessToken}`,
        },
      });
      const response_json = await response.json();
      console.log("JSON: ", response_json);
      setCategory(response_json);
      setLoadingCategory(false);
    } catch (error: any) {
      console.log("ERROR: ", error.message);
      setCategory(undefined);
      setLoadingCategory(false);
    }
  };

  const loadItems = async (categoryID: string) => {
    const domain = "localhost:8080";
    try {
      const accessToken = await getAccessTokenSilently({
        audience: `http://${domain}/`,
        scope: "read:item",
      });
      const response = await fetch(`/api/categories/${categoryID}/items`, {
        headers: {
          Authorization: `Bearer ${accessToken}`,
        },
      });
      const response_json = await response.json();
      console.log("JSON: ", response_json);
      setItems(response_json);
      setLoadingItems(false);
    } catch (error: any) {
      console.log("ERROR: ", error.message);
      setItems(undefined);
      setLoadingItems(false);
    }
  };

  useEffect(() => {
    if (categoryID) {
      loadCategory(categoryID);
      loadItems(categoryID);
    }
    return () => {};
  }, [categoryID]);

  const onAdded = (err?: Error) => {
    setShowAdd(false);
    loadItems(categoryID!);
  };

  if (loadingCategory || loadingItems) return <Loading />;

  if (!category) return <div>Category not found</div>;

  return (
    <div>
      <div className="HeaderContainer">
        <h1>{category.name}</h1>
        <div className="TopButtonBar">
          <Button onClick={() => setShowAdd(true)}>Add Item</Button>
        </div>
      </div>
      {showAdd && (
        <AddItem
          categoryID={category.id}
          onAdded={(err?: Error) => onAdded(err)}
          onCancelled={() => setShowAdd(false)}
        />
      )}
      {itemLoadError ? (
        <div>Failed to load Items: {itemLoadError.message}</div>
      ) : (
        <Table striped bordered hover>
          <thead>
            <tr>
              <th>Name</th>
            </tr>
          </thead>
          <tbody>
            {items!.map((i) => (
              <tr key={i.id}>
                <td>
                  <div className="ItemRowContainer">
                    <div className="ItemName">
                      <Link to={"items/" + i.id}>{i.name}</Link>
                    </div>
                    <Button variant="danger" onClick={() => deleteClicked(c)}>
                      Delete
                    </Button>
                  </div>
                </td>
              </tr>
            ))}
          </tbody>
        </Table>
      )}
    </div>
  );
};

export default Category;
