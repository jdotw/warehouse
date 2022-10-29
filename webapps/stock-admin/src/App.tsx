import { useAuth0 } from "@auth0/auth0-react";
import { Button, Container, Nav, Navbar, NavDropdown } from "react-bootstrap";
import { Outlet } from "react-router-dom";
import "./App.css";

function AuthenticatedApp() {
  const { logout } = useAuth0();
  return (
    <div className="App">
      <div className="Nav">
        <Navbar bg="light" expand="lg">
          <Container>
            <Navbar.Brand href="/">Stock Admin</Navbar.Brand>
            <Navbar.Toggle aria-controls="basic-navbar-nav" />
            <Navbar.Collapse id="basic-navbar-nav">
              <Nav className="me-auto">
                <Nav.Link href="/">Home</Nav.Link>
                <Nav.Link href="/categories">Categories</Nav.Link>
                <Nav.Link
                  onClick={() => logout({ returnTo: window.location.origin })}
                >
                  Logout
                </Nav.Link>
              </Nav>
            </Navbar.Collapse>
          </Container>
        </Navbar>
      </div>
      <Outlet />
    </div>
  );
}

function UnAuthenticatedApp() {
  const { loginWithRedirect } = useAuth0();
  return (
    <div className="App">
      <div className="Nav">
        <Navbar bg="light" expand="lg">
          <Container>
            <Navbar.Brand href="/">Stock Admin</Navbar.Brand>
            <Navbar.Toggle aria-controls="basic-navbar-nav" />
            <Navbar.Collapse id="basic-navbar-nav">
              <Nav className="me-auto">
                <Nav.Link onClick={() => loginWithRedirect()}>Login</Nav.Link>
              </Nav>
            </Navbar.Collapse>
          </Container>
        </Navbar>
      </div>
      <div className="LogIn">
        <Button onClick={() => loginWithRedirect()}>Log In</Button>
      </div>
    </div>
  );
}

function App() {
  const { isAuthenticated, isLoading } = useAuth0();

  if (isLoading) {
    return <div>Loading ...</div>;
  }

  return isAuthenticated ? <AuthenticatedApp /> : <UnAuthenticatedApp />;
}

export default App;
