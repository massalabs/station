import { ReactNode, useState, useEffect } from 'react';
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
  const [availableNetworks, setAvailableNetworks] = useState([]);
  const [currentNetwork, setCurrentNetwork] = useState('');

  function handleSetTheme(theme: string) {
    setSelectedTheme(theme);
    onSetTheme?.(theme);
  }

  const { data: networkConfig = {} } = useResource<NetworksObject>(
    `${URL.PATH_NETWORKS}`,
  );

  useEffect(() => {
    if (networkConfig.availableNetworks) {
      setAvailableNetworks(
        networkConfig.availableNetworks.map((n) => ({ item: n })),
      );
    }
    if (networkConfig.currentNetwork) {
      setCurrentNetwork(networkConfig.currentNetwork);
    }
  }, [networkConfig]);

  const handleNetworkChange = (selectedNetwork: string) => {
    // Make the network switch request here
    // Assuming you are using fetch() for API requests
    fetch(`/network/${selectedNetwork}`)
      .then((response) => response.json())
      .then((data) => {
        // Handle the response if necessary
        console.log('Current Network:', data.currentNetwork);
      })
      .catch((error) => {
        // Handle the error if necessary
        console.error(error);
      });
  };

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
              options={availableNetworks}
              select={currentNetwork}
              onChange={handleNetworkChange}
            />
          </div>
          <ThemeMode onSetTheme={handleSetTheme} />
        </div>
      </div>
      {children}
    </div>
  );
}
