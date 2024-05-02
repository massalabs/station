import { useEffect, useState } from 'react';
import { Outlet, useLocation, useNavigate } from 'react-router-dom';
import { useLocalStorage } from '@/custom/useLocalStorage';
import { useConfigStore } from '@/store/store';

import { Navigator, Theme, Toast } from '@massalabs/react-ui-kit';
import { FiCodepen, FiGlobe, FiHome, FiSun, FiMoon } from 'react-icons/fi';
import { LayoutStation } from '@/layouts/LayoutStation/LayoutStation';

import { PAGES } from '@/const/pages/pages';
import { THEME_STORAGE_KEY } from '@/const';

type ThemeSettings = {
  [key: string]: {
    icon: JSX.Element;
    label: string;
  };
};

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
  const [theme, setThemeStorage] = useLocalStorage<Theme>(
    THEME_STORAGE_KEY,
    'theme-dark',
  );

  const { pathname } = useLocation();
  const navigate = useNavigate();

  const currentPage = pathname.split('/').pop() || 'index';
  const [activePage, setActivePage] = useState(currentPage);

  const setThemeStore = useConfigStore((s) => s.setTheme);

  useEffect(() => {
    setActivePage(currentPage);
  }, [setActivePage, pathname, currentPage]);

  const navigator = (
    <Navigator
      items={[
        {
          icon: <FiHome />,
          isActive: PAGES.INDEX === activePage,
        },
        {
          icon: <FiCodepen />,
          isActive: PAGES.STORE === activePage,
        },
        {
          icon: <FiGlobe />,
          isActive: PAGES.SEARCH === activePage,
        },
      ]}
      onClickNext={handleNext}
      onClickBack={handlePrevious}
    />
  );
  const STEP = navigatorSteps[currentPage] as INavigatorSteps;

  // Functions
  function handleNext() {
    let { next } = STEP;

    setActivePage(next.toString());
    navigate(next);
  }

  function handlePrevious() {
    let { previous } = STEP;

    setActivePage(previous.toString());
    navigate(previous);
  }

  function handleSetTheme() {
    let toggledTheme: Theme =
      theme === 'theme-dark' ? 'theme-light' : 'theme-dark';

    setThemeStorage(toggledTheme);
    setThemeStore(toggledTheme);
  }

  // Template
  return (
    <div className={theme}>
      <LayoutStation
        navigator={navigator}
        onSetTheme={handleSetTheme}
        storedTheme={theme}
      >
        <Outlet />
        <Toast />
      </LayoutStation>
    </div>
  );
}
