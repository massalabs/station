import { useState } from "react";
import { Route, Routes } from "react-router";
import Home from "./pages/Home/home";
import Manager from "./pages/Plugin_Manager/manager";

function App() {
    return (
        <div className="min-h-screen bg-slate-900">
            <Home/>
        </div>
    );
}

export default App;
