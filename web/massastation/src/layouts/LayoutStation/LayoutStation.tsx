import { ReactNode, useState } from 'react';

import { useResource } from '../../custom/api';
import { AccountObject } from '../../models/AccountModel';
import { useAccountStore } from '../../store/store';
import { URL } from '../../const/url/url';

import {
  ThemeMode,
  StationLogo,
  Dropdown,
  Identicon,
} from '@massalabs/react-ui-kit';

export interface LayoutStationProps {
  children?: ReactNode;
  navigator?: Navigator;
  onSetTheme?: () => void;
  storedTheme?: string;
  activePage: string;
}

export function LayoutStation({ ...props }) {
  const { children, navigator, onSetTheme, storedTheme, activePage } = props;

  const [selectedTheme, setSelectedTheme] = useState(
    storedTheme || 'theme-dark',
  );

  const searchIsActive = activePage === 'search';

  function handleSetTheme(theme: string) {
    setSelectedTheme(theme);

    onSetTheme?.(theme);
  }

  const { data: accounts = [] } = useResource<AccountObject[]>(
    `${URL.WALLET_BASE_API}/${URL.WALLET_ACCOUNTS}`,
  );

  // Store
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
          {searchIsActive && (
            <div className="w-64">
              <Dropdown options={accountsItems} select={selectedAccountKey} />
            </div>
          )}
          <ThemeMode onSetTheme={handleSetTheme} />
        </div>
      </div>
      {children}
    </div>
  );
}
