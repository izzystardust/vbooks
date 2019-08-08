import React from "react";
import { render } from "react-dom";
// import { Router, Link } from "@reach/router";
// import SearchParams from "./SearchParams";
import { login } from "./lib/vbooks";

const App = () => {
  return (
    <React.StrictMode>
      <div>
        <header>
          <h1>Hello, world!</h1>
        </header>
      </div>
    </React.StrictMode>
  );
};
login("izzy", "test99**");
render(<App />, document.getElementById("root"));
