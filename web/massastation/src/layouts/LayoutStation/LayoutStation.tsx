import { ReactNode, useState, useEffect } from 'react';

import { useResource, usePost } from '../../custom/api';
import { URL } from '../../const/url/url';
import { NetworkModel } from '../../models';
import { useNetworkStore } from '../../store/store';

import { ThemeMode, StationLogo, Dropdown } from '@massalabs/react-ui-kit';

export interface LayoutStationProps {
  children?: ReactNode;
  navigator?: Navigator;
  onSetTheme?: () => void;
  storedTheme?: string;
}

interface NetworkRequest {
  network: string;
}

export function LayoutStation({ ...props }) {
  const { children, navigator, onSetTheme, storedTheme } = props;

  const [selectedTheme, setSelectedTheme] = useState(
    storedTheme || 'theme-dark',
  );

  function handleSetTheme(theme: string) {
    setSelectedTheme(theme);

    onSetTheme?.(theme);
  }

  const { data: network, isSuccess: isSuccessNetwork } =
    useResource<NetworkModel>(URL.PATH_NETWORKS);

  const [
    currentNetwork,
    availableNetworks,
    setCurrentNetwork,
    setAvailableNetworks,
  ] = useNetworkStore((state) => [
    state.currentNetwork,
    state.availableNetworks,
    state.setCurrentNetwork,
    state.setAvailableNetworks,
  ]);

  useEffect(() => {
    if (isSuccessNetwork) {
      if (network.currentNetwork) setCurrentNetwork(network.currentNetwork);
      if (network.availableNetworks)
        setAvailableNetworks(network.availableNetworks);
    }
  }, [isSuccessNetwork]);

  const selectedNetworkKey: number = parseInt(
    Object.keys(availableNetworks).find(
      (_, idx) => availableNetworks[idx] === currentNetwork,
    ) || '0',
  );

  const { mutate: mutateUpdateNetwork } = usePost<NetworkRequest>(
    `${URL.PATH_NETWORKS}/${currentNetwork}`,
  );

  const availableNetworksItems = availableNetworks.map((network) => ({
    item: network,
    onClick: () => {
      setCurrentNetwork(network);
      mutateUpdateNetwork({});
      location.replace(location.href);
    },
  }));

  return (
    <div
      data-testid="layout-station"
      className={`min-h-screen bg-primary px-20 pt-12 pb-8 $}`}
    >
      <div className="grid grid-cols-3">
        <div className="flex justify-start items-center">
          <a href="/">
            <StationLogo theme={selectedTheme} />
          </a>
        </div>
        <div className="flex justify-center items-center">
          {navigator && <div className="flex-row-reversed">{navigator}</div>}
        </div>
        <div className="flex justify-end items-center gap-3">
          <div className="w-44">
            <Dropdown
              size="xs"
              options={availableNetworksItems}
              select={selectedNetworkKey}
            />
          </div>
          <ThemeMode onSetTheme={handleSetTheme} />
        </div>
      </div>
      {children}
    </div>
  );
}
