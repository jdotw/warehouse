import { useAuth0 } from "@auth0/auth0-react";
import React, { SyntheticEvent, useEffect, useState } from "react";
import { Alert, Button, Modal, Table } from "react-bootstrap";
import { Link } from "react-router-dom";
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

  const [categoryToDelete, setCategoryToDelete] = useState<Category>();
  const [showDeleteConfirmation, setShowDeleteConfirmation] = useState(false);

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

  const deleteClicked = (category: Category) => {
    setCategoryToDelete(category);
    setShowDeleteConfirmation(true);
  };

  const deleteCategoryCancelled = () => {
    setShowDeleteConfirmation(false);
  };

  const deleteCategoryConfirmed = async () => {
    setShowDeleteConfirmation(false);
    if (!categoryToDelete) {
      console.log("ERROR: categoryToDelete is nil");
      return;
    }
    const domain = "localhost:8080";
    try {
      const accessToken = await getAccessTokenSilently({
        audience: `http://${domain}/`,
        scope: "write:category",
      });
      const deleteCategoryURL = `/api/categories/${categoryToDelete?.id}`;
      const response = await fetch(deleteCategoryURL, {
        headers: {
          Authorization: `Bearer ${accessToken}`,
        },
        method: "DELETE",
      });
      const response_json = await response.json();
      console.log("DELETE JSON: ", response_json);
      setCategories(categories?.filter((x) => x.id != categoryToDelete.id));
    } catch (error: any) {
      console.log("ERROR: ", error.message);
      alert(error.message);
    } finally {
      setCategoryToDelete(undefined);
    }
  };

  return (
    <>
      <div className="CategoryList">
        <div className="HeaderContainer">
          <h1>Categories</h1>
          <div className="TopButtonBar">
            <Button onClick={() => setShowAdd(true)}>Create Category</Button>
          </div>
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
                  <td>
                    <div className="CategoryRowContainer">
                      <div className="CategoryName">
                        <Link to={c.id}>{c.name}</Link>
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
      <Modal show={showDeleteConfirmation} onHide={deleteCategoryCancelled}>
        <Modal.Header closeButton>
          <Modal.Title>Delete Category '{categoryToDelete?.name}'</Modal.Title>
        </Modal.Header>
        <Modal.Body>
          Are you sure you want to delete the category named '
          {categoryToDelete?.name}'
        </Modal.Body>
        <Modal.Footer>
          <Button variant="secondary" onClick={deleteCategoryCancelled}>
            Cancel
          </Button>
          <Button variant="primary" onClick={deleteCategoryConfirmed}>
            Delete
          </Button>
        </Modal.Footer>
      </Modal>
    </>
  );
};

export default CategoryList;
