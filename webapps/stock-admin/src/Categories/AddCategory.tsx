import React, { useState } from "react";
import Form from "react-bootstrap/Form";
import Button from "react-bootstrap/Button";
import { useAuth0 } from "@auth0/auth0-react";

const AddCategory = () => {
  const { user, isAuthenticated, getAccessTokenSilently } = useAuth0();
  const [name, setName] = useState("");

  const formSubmitted = async (e: React.SyntheticEvent) => {
    e.preventDefault();
    console.log("CREATING: ", name);
    const domain = "localhost:8080";
    try {
      const accessToken = await getAccessTokenSilently({
        audience: `http://${domain}/`,
        scope: "write:category",
      });
      console.log("TOKEN: ", accessToken);

      const addCategoryURL = `/api/categories`;

      const response = await fetch(addCategoryURL, {
        headers: {
          Authorization: `Bearer ${accessToken}`,
        },
        body: JSON.stringify({
          name: name,
        }),
        method: "POST",
      });

      const response_json = await response.json();
      console.log("JSON: ", response_json);
    } catch (error) {
      console.log("ERROR: ", error.message);
    }

    console.log("CLICKED");
  };

  return (
    <div>
      <h2>Add New Category</h2>
      <Form onSubmit={formSubmitted}>
        <Form.Group controlId="formCategoryName">
          <Form.Label>Name</Form.Label>
          <Form.Control
            placeholder="Category Name"
            onChange={(e) => setName(e.currentTarget.value)}
          />
        </Form.Group>
        <Button variant="primary" type="submit">
          Submit
        </Button>
      </Form>
    </div>
  );
};

export default AddCategory;
