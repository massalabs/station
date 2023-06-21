import { Outlet, useLocation, useNavigate } from 'react-router-dom';
import { useLocalStorage } from '../../custom/useLocalStorage';
import { FiSun, FiMoon } from 'react-icons/fi';
import { Navigator, LayoutStation } from '@massalabs/react-ui-kit';
import { FiCodepen, FiGlobe, FiHome } from 'react-icons/fi';
import { useEffect, useState } from 'react';
import { PAGES } from '../../const/pages/pages';

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

interface INavigatorSteps {
  [key: string]: object;
}

const navigatorSteps: INavigatorSteps = {
  index: {
    previous: null,
    next: PAGES.STORE,
  },
  store: {
    previous: PAGES.INDEX,
    next: PAGES.SEARCH,
  },
  search: {
    previous: PAGES.STORE,
    next: null,
  },
};

export function Base() {
  const [theme, setTheme] = useLocalStorage<string>(
    'massa-station-theme',
    'theme-dark',
  );
  const { pathname } = useLocation();
  const currentPage = pathname.split('/').pop() || 'index';
  const [active, setActive] = useState(currentPage);
  const context = { handleSetTheme };
  const navigate = useNavigate();
  const navigator = (
    <Navigator
      items={[
        {
          icon: <FiHome />,
          isActive: PAGES.INDEX === active,
        },
        {
          icon: <FiCodepen />,
          isActive: PAGES.STORE === active,
        },
        {
          icon: <FiGlobe />,
          isActive: PAGES.SEARCH === active,
        },
      ]}
      onClickNext={handleNext}
      onClickBack={handlePrevious}
    />
  );
  const STEP = navigatorSteps[currentPage] as INavigatorSteps;

  useEffect(() => {
    setActive(currentPage);
  }, [pathname]);

  function handleNext() {
    let { next } = STEP;

    setActive(next.toString());
    navigate(next);
  }

  function handlePrevious() {
    let { previous } = STEP;

    setActive(previous.toString());
    navigate(previous);
  }

  function handleSetTheme() {
    setTheme(theme === 'theme-dark' ? 'theme-light' : 'theme-dark');
  }

  return (
    <div className={`${theme}`}>
      <LayoutStation navigator={navigator} onSetTheme={handleSetTheme}>
        <Outlet context={context} />
      </LayoutStation>
    </div>
  );
}
