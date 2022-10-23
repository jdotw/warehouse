import React, { useState } from "react";
import Form from "react-bootstrap/Form";
import Button from "react-bootstrap/Button";
import { useAuth0 } from "@auth0/auth0-react";
import "./AddCategory.css";

type Props = {
  onAdded: (error?: Error) => void;
  onCancelled: () => void;
};

const AddCategory = (props: Props) => {
  const { user, isAuthenticated, getAccessTokenSilently } = useAuth0();
  const [name, setName] = useState("");
  const { onAdded, onCancelled } = props;

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
      if (onAdded) onAdded();
    } catch (error: any) {
      console.log("ERROR: ", error.message);
      if (onAdded) onAdded(error);
    }

    console.log("CLICKED");
  };

  return (
    <div className="Container">
      <h2>Add New Category</h2>
      <Form onSubmit={formSubmitted}>
        <Form.Group controlId="formCategoryName" className="FieldRow">
          <Form.Label>Name</Form.Label>
          <Form.Control
            placeholder="Category Name"
            onChange={(e) => setName(e.currentTarget.value)}
          />
        </Form.Group>
        <div className="ButtonContainer">
          <Button
            variant="secondary"
            onClick={onCancelled}
            className="CancelButton"
          >
            Cancel
          </Button>
          <Button variant="primary" type="submit" className="SubmitButton">
            Submit
          </Button>
        </div>
      </Form>
    </div>
  );
};

export default AddCategory;
