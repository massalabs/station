import { UIStore } from '../store/UIStore';
import toggleTheme from './toggleTheme';
import logoLight from '../assets/logo/logoLight.svg';
import logoDark from '../assets/logo/logoDark.svg';
import {THEMES} from "../../../shared/consts/theme";
const Header = () => {
  const isThemeLight = UIStore.useState((s) =>
    s.theme == THEMES.LIGHT ? true : false,
  );
  return (
    <div className="flex p-12 items-center content-center justify-between ">
      <img
        className="max-h-6"
        src={isThemeLight ? logoLight : logoDark}
        alt="Thyra Logo"
      />
      {toggleTheme()}
    </div>
  );
};

export default Header;
