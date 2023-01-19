import { useState } from "react";
import reactLogo from "./assets/react.svg";
import Home from "./pages/Home/home";

function App() {
    const [count, setCount] = useState(0);

    return (
        <div className="h-screen  bg-slate-900">
          <Home />
        </div>
    );
}

export default App;
