import { useState } from "react";
import { Route, Routes } from "react-router";
import reactLogo from "./assets/react.svg";
import Home from "./pages/Home/home";
import Manager from "./pages/Plugin_Manager/manager";

function App() {
    const [count, setCount] = useState(0);

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
