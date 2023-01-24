import React from "react";
import ReactDOM from "react-dom/client";
import App from "./App";
import "./index.css";
import {
    QueryClient,
    QueryClientProvider,
} from "react-query";
import { BrowserRouter } from "react-router-dom";

ReactDOM.createRoot(document.getElementById("root") as HTMLElement).render(
    <React.StrictMode>
        <BrowserRouter>
            <QueryClientProvider client={new QueryClient()}>
                <App />
            </QueryClientProvider>
        </BrowserRouter>
    </React.StrictMode>
);
