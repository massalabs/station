import { useState } from "react";
import { Route, Routes } from "react-router";
import Home from "./pages/Home/home";

function App() {
    return (
        <div className="min-h-screen bg-slate-900">
            <Home/>
        </div>
    );
}

export default App;
