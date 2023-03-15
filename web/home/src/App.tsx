import { useState } from "react";
import Home from "./pages/Home/home";

function App() {
    const [theme, setTheme] = useState("light")
    return (
        <html className={"theme-"+theme}>
        <div className="min-h-screen bg-slate-900">
            <Home setTheme={setTheme}/>
        </div>
        </html>
    );
}

export default App;
