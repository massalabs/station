import { Outlet } from 'react-router-dom';
import { useLocalStorage } from '../../custom/useLocalStorage';
import { FiSun, FiMoon } from 'react-icons/fi';

type ThemeSettings = {
  [key: string]: {
    icon: JSX.Element;
    label: string;
  };
};

export interface IOutletContextType {
  themeIcon: JSX.Element;
  themeLabel: string;
  theme: string;
  handleSetTheme: () => void;
}

export const themeSettings: ThemeSettings = {
  'theme-dark': {
    icon: <FiSun />,
    label: 'light theme',
  },
  'theme-light': {
    icon: <FiMoon />,
    label: 'dark theme',
  },
};

function Base() {
  const [theme, setTheme] = useLocalStorage<string>(
    'massa-station-theme',
    'theme-dark',
  );
  const context = { handleSetTheme };

  function handleSetTheme() {
    setTheme(theme === 'theme-dark' ? 'theme-light' : 'theme-dark');
  }

  return (
    // TODO
    // remove theme-dark
    // this needs to be removed as soon we fix the steps to create an account
    <div className={`${theme}`}>
      <Outlet context={context} />
    </div>
  );
}

export default Base;
