import { ReactNode, useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { routeFor } from '../../utils';

import { useResource } from '../../custom/api';
import { AccountObject } from '../../models/AccountModel';
import { useAccountStore } from '../../store/store';
import { URL } from '../../const/url/url';

import { PAGES } from '../../const/pages/pages';

import {
  ThemeMode,
  StationLogo,
  Dropdown,
  Identicon,
  Button,
} from '@massalabs/react-ui-kit';
import { IMassaStore } from '../../../../shared/interfaces/IPlugin';

export interface LayoutStationProps {
  children?: ReactNode;
  navigator?: Navigator;
  onSetTheme?: () => void;
  storedTheme?: string;
  activePage: string;
}

export function LayoutStation({ ...props }) {
  const { children, navigator, onSetTheme, storedTheme, activePage } = props;

  const navigate = useNavigate();

  const [selectedTheme, setSelectedTheme] = useState(
    storedTheme || 'theme-dark',
  );

  const searchIsActive = activePage === PAGES.SEARCH;

  function handleSetTheme(theme: string) {
    setSelectedTheme(theme);

    onSetTheme?.(theme);
  }

  const { data: accounts = [] } = useResource<AccountObject[]>(
    `${URL.WALLET_BASE_API}/${URL.WALLET_ACCOUNTS}`,
  );

  const nickname = useAccountStore((state) => state.nickname);
  const setNickname = useAccountStore((state) => state.setNickname);

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

  const existingAccount: boolean = accounts.length > 0;

  const [pluginWalletIsInstalled, setPluginWalletIsInstalled] = useState(false);

  const { data: plugins, isSuccess } =
    useResource<IMassaStore[]>('plugin-manager');

  useEffect(() => {
    if (isSuccess) {
      plugins.forEach((plugin) => {
        if (plugin.name === 'Massa Wallet') {
          setPluginWalletIsInstalled(true);
        }
      });
    }
  }, [isSuccess]);

  return (
    <div
      data-testid="layout-station"
      className={`min-h-screen bg-primary px-20 pt-12 pb-8 $}`}
    >
      <div className="grid grid-cols-3">
        <div className="flex justify-start">
          <a href="/">
            <StationLogo theme={selectedTheme} />
          </a>
        </div>
        <div className="flex justify-center">
          {navigator && <div className="flex-row-reversed">{navigator}</div>}
        </div>
        <div className="flex justify-end items-start gap-20">
          {searchIsActive &&
            (pluginWalletIsInstalled ? (
              existingAccount ? (
                <div className="w-64">
                  <Dropdown
                    options={accountsItems}
                    select={selectedAccountKey}
                  />
                </div>
              ) : (
                <Button
                  customClass="w-64"
                  onClick={() =>
                    window.open(
                      '/plugin/massa-labs/massa-wallet/web-app/',
                      '_blank',
                    )
                  }
                >
                  Create Account
                </Button>
              )
            ) : (
              <Button
                customClass="w-64"
                onClick={() => navigate(routeFor('index'))}
              >
                Install Wallet
              </Button>
            ))}
          <ThemeMode onSetTheme={handleSetTheme} />
        </div>
      </div>
      {children}
    </div>
  );
}
