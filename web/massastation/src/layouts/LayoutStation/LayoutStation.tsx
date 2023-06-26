import { ReactNode, useState } from 'react';
import { Dropdown, StationLogo, ThemeMode } from '@massalabs/react-ui-kit';
import { URL } from '../../const/url/url';
import { useResource } from '../../custom/api';
import { NetworksObject } from '../../models/AccountModel';

export interface LayoutStationProps {
  children?: ReactNode;
  onSetTheme?: () => void;
  activePage: string;
}

export function LayoutStation({
  children,
  onSetTheme,
  activePage,
}: LayoutStationProps) {
  const [selectedTheme, setSelectedTheme] = useState('theme-dark');

  function handleSetTheme(theme: string) {
    setSelectedTheme(theme);
    onSetTheme?.(theme);
  }
  const { data: networkConfig = {} } = useResource<NetworksObject>(
    `${URL.PATH_NETWORKS}`,
  );

  const availableNetworks = networkConfig.availableNetworks?.map((n) => ({
    item: n,
  }));
  const currentNetwork = networkConfig.currentNetwork;

  return (
    <div
      data-testid="layout-station"
      className="min-h-screen bg-primary px-20 pt-12 pb-8"
    >
      <div className="grid grid-cols-3">
        <div className="flex justify-start">
          <a href="/">
            <StationLogo theme={selectedTheme} />
          </a>
        </div>
        <div className="flex justify-center"></div>
        <div className="flex justify-end items-start gap-20">
          <div className="w-64">
            <Dropdown
              options={availableNetworks || []}
              select={currentNetwork}
            />
          </div>
          <ThemeMode onSetTheme={handleSetTheme} />
        </div>
      </div>
      {children}
    </div>
  );
}
