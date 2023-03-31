import { UIStore } from "../store/UIStore";
import toggleTheme from "./toggleTheme";
import logoLight from "../assets/logo/logoLight.svg";
import logoDark from "../assets/logo/logoDark.svg";

const Header = () => {
    const isThemeLight = UIStore.useState((s) => s.theme === "light");
    return (
        <div className="flex p-12 items-center content-center justify-between cursor-pointer ">
            <img onClick={() => window.open("/thyra/home")} className="max-h-6" src={isThemeLight ? logoLight : logoDark} alt="Thyra Logo" />
            {toggleTheme()}
        </div>
    );
};

export default Header;
