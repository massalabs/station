import React from "react";
import ReactDOM from "react-dom/client";
import App from "./App";
import "./index.css";
import {
    useQuery,
    useMutation,
    useQueryClient,
    QueryClient,
    QueryClientProvider,
} from "react-query";

ReactDOM.createRoot(document.getElementById("root") as HTMLElement).render(
    <React.StrictMode>
        <QueryClientProvider client={new QueryClient()}>
            <App />
        </QueryClientProvider>
    </React.StrictMode>
);
