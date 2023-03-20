import { UIStore } from '../store/UIStore'
import toggleTheme from './toggleTheme'
import logoLight from "../assets/logo/logoLight.svg";
import logoDark from "../assets/logo/logoDark.svg";

const Header = () => {
    const isThemeLight = UIStore.useState(s => (s.theme == "light" ? true : false));
  return (
      <div className="flex p-12 items-center content-center justify-between ">
        <button onClick={() => window.open("/thyra/home")}>
          <img className="max-h-6" src={isThemeLight ? logoLight : logoDark} alt="Thyra Logo" />
        </button>
        {toggleTheme()}
    </div>
  )
}

export default Header