import { useState } from "react";
import Home from "./pages/Home/home";
import { UIStore } from "./store/UIStore";

function App() {
    return (
        <html className={"theme-"+UIStore.useState(s => (s.theme))}>
        <div className="min-h-screen bg-primaryBG">
            <Home/>
        </div>
        </html>
    );
}

export default App;
