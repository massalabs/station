import { useState } from "react";
import { Route, Routes } from "react-router";
import Home from "./pages/Home/home";
import Manager from "./pages/Plugin_Manager/manager";

function App() {
    return (
        <div className="min-h-screen bg-slate-900">
          <Routes>
            <Route path="/" element={<Home />} />
            <Route path="/manager" element={<Manager/>} />
          </Routes>
        </div>
    );
}

export default App;
