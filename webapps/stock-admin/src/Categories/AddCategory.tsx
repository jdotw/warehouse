import React from "react";
import Form from "react-bootstrap/Form";
import Button from "react-bootstrap/Button";

const AddCategory = () => {
  const formSubmitted = (e: React.SyntheticEvent): void => {
    console.log("CLICKED");
    e.preventDefault();
  };
  return (
    <div>
      <h2>Add New Category</h2>
      <Form onSubmit={formSubmitted}>
        <Form.Group controlId="formCategoryName">
          <Form.Label>Name</Form.Label>
          <Form.Control placeholder="Category Name" />
        </Form.Group>
        <Button variant="primary" type="submit">
          Submit
        </Button>
      </Form>
    </div>
  );
};

export default AddCategory;
