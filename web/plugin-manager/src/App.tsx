import Manager from "./pages/Manager"
import { UIStore } from "./store/UIStore";

function App() {

  return (
    <html className={"theme-"+UIStore.useState(s => (s.theme))}>
    <div className="min-h-screen bg-primaryBG">
      <Manager/>
    </div>
    </html>
  )
}

export default App
