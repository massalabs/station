import React from "react";
import ReactDOM from "react-dom/client";
import "./index.css";
import App from "./App";
import GlobalStyles from "./styles/GlobalStyles";
ReactDOM.createRoot(document.getElementById("root") as HTMLElement).render(
    <>
        <GlobalStyles />
        <App />
    </>
);
