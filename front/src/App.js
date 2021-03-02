import "./App.css";
import React from "react";
import AdminIf from "./AdminIf";
import SignIn from "./SignIn";
import { BrowserRouter as Router, Switch, Route, Link } from "react-router-dom";

function App() {
  const [email, setEmail] = React.useState("");
  return (
    <Router>
      <Switch>
        <Route
          exact
          path="/admin"
          render={(props) => (
            <AdminIf email={email} setEmail={setEmail} {...props} />
          )}
        />
        <Route
          path="/login"
          render={(props) => (
            <SignIn email={email} setEmail={setEmail} {...props} />
          )}
        />
        <Route
          path="/"
          render={(props) => (
            <SignIn email={email} setEmail={setEmail} {...props} />
          )}
        />
      </Switch>
    </Router>
  );
}

export default App;
