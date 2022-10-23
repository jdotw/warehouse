import { useAuth0 } from "@auth0/auth0-react";
import React, { useEffect, useState } from "react";
import { Button, Modal, Table } from "react-bootstrap";
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
  const [itemToDelete, setItemToDelete] = useState<Category>();
  const [showDeleteConfirmation, setShowDeleteConfirmation] = useState(false);

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
      setItemLoadError(undefined);
      setLoadingCategory(false);
    } catch (error: any) {
      console.log("ERROR: ", error.message);
      setCategory(undefined);
      setItemLoadError(error);
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

  const deleteClicked = (category: Category) => {
    setItemToDelete(category);
    setShowDeleteConfirmation(true);
  };

  const deleteItemCancelled = () => {
    setItemToDelete(undefined);
    setShowDeleteConfirmation(false);
  };

  const deleteItemConfirmed = async () => {
    setShowDeleteConfirmation(false);
    if (!itemToDelete) {
      console.log("ERROR: itemToDelete is nil");
      return;
    }
    const domain = "localhost:8080";
    try {
      const accessToken = await getAccessTokenSilently({
        audience: `http://${domain}/`,
        scope: "write:category",
      });
      const response = await fetch(`/api/items/${itemToDelete?.id}`, {
        headers: {
          Authorization: `Bearer ${accessToken}`,
        },
        method: "DELETE",
      });
      const response_json = await response.json();
      console.log("DELETE JSON: ", response_json);
      setItems(items?.filter((x) => x.id != itemToDelete.id));
    } catch (error: any) {
      console.log("ERROR: ", error.message);
      alert(error.message);
    } finally {
      setItemToDelete(undefined);
    }
  };

  if (loadingCategory || loadingItems) return <Loading />;

  if (!category) return <div>Category not found</div>;

  return (
    <>
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
                      <Button variant="danger" onClick={() => deleteClicked(i)}>
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
      <Modal show={showDeleteConfirmation} onHide={deleteItemCancelled}>
        <Modal.Header closeButton>
          <Modal.Title>Delete Item '{itemToDelete?.name}'</Modal.Title>
        </Modal.Header>
        <Modal.Body>
          Are you sure you want to delete the item named '{itemToDelete?.name}'
        </Modal.Body>
        <Modal.Footer>
          <Button variant="secondary" onClick={deleteItemCancelled}>
            Cancel
          </Button>
          <Button variant="primary" onClick={deleteItemConfirmed}>
            Delete
          </Button>
        </Modal.Footer>
      </Modal>
    </>
  );
};

export default Category;
