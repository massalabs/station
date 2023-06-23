import { Outlet, useLocation, useNavigate } from 'react-router-dom';
import { useLocalStorage } from '../../custom/useLocalStorage';
import { FiSun, FiMoon } from 'react-icons/fi';
import {
  Navigator,
  LayoutStation,
  Dropdown,
  Identicon,
  Toast,
} from '@massalabs/react-ui-kit';
import { FiCodepen, FiGlobe, FiHome } from 'react-icons/fi';
import { useEffect, useState } from 'react';
import { PAGES } from '../../const/pages/pages';
import { useResource } from '../../custom/api';
import { AccountObject } from '../../models/AccountModel';
import { useAccountStore } from '../../store/store';
import { URL } from '../../const/url/url';
import { useConfigStore } from '../../store/store';

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
  // Hooks
  const [theme, setThemeStorage] = useLocalStorage<string>(
    'massa-station-theme',
    'theme-dark',
  );

  const { pathname } = useLocation();
  const navigate = useNavigate();

  useEffect(() => {
    setActive(currentPage);
  }, [pathname]);

  const { data: accounts = [] } = useResource<AccountObject[]>(
    `${URL.WALLET_BASE_API}/${URL.WALLET_ACCOUNTS}`,
  );

  // State
  const currentPage = pathname.split('/').pop() || 'index';
  const [active, setActive] = useState(currentPage);

  // Store
  const nickname = useAccountStore((state) => state.nickname);
  const setNickname = useAccountStore((state) => state.setNickname);
  const setThemeStore = useConfigStore((s) => s.setTheme);

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

  const accountsItems = accounts.map((account) => ({
    icon: <Identicon username={account.nickname} size={32} />,
    item: account.nickname,
    onClick: () => setNickname(account.nickname),
  }));

  const selectedAccountKey: number = parseInt(
    Object.keys(accounts).find(
      (_, idx) => accounts[idx].nickname === nickname,
    ) || '0',
  );

  // Functions
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
    let toggledTheme = theme === 'theme-dark' ? 'theme-light' : 'theme-dark';

    setThemeStorage(toggledTheme);
    setThemeStore(toggledTheme);
  }

  // Template
  return (
    <div className={`${theme}`}>
      <LayoutStation navigator={navigator} onSetTheme={handleSetTheme}>
        <div className="absolute top-5 right-32 p-6">
          <div className="w-64">
            <Dropdown options={accountsItems} select={selectedAccountKey} />
          </div>
        </div>
        <Outlet />
        <Toast />
      </LayoutStation>
    </div>
  );
}
