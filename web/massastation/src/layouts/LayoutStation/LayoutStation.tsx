import { ReactNode, useState, useEffect } from 'react';

import { useResource, usePost } from '../../custom/api';
import { URL } from '../../const/url/url';
import { NetworkModel } from '../../models';
import { useNetworkStore } from '../../store/store';

import {
  ThemeMode,
  StationLogo,
  Dropdown,
  Theme,
  toast,
} from '@massalabs/react-ui-kit';
import { useNavigate } from 'react-router-dom';
import Intl from '@/i18n/i18n';

interface LayoutStationProps {
  children?: ReactNode;
  navigator?: ReactNode;
  onSetTheme?: (theme: Theme) => void;
  storedTheme?: Theme | undefined;
}

interface NetworkRequest {
  network: string;
}

interface NetworkResponse {
  currentNetwork: string;
}

export function LayoutStation(props: LayoutStationProps) {
  const { children, navigator, onSetTheme, storedTheme } = props;

  const { data: version, isSuccess: getVersionSuccess } = useResource<string>(
    URL.VERSION,
  );

  const [selectedTheme, setSelectedTheme] = useState<Theme>(
    storedTheme || 'theme-dark',
  );

  function handleSetTheme(theme: Theme) {
    setSelectedTheme(theme);

    onSetTheme?.(theme);
  }

  const navigate = useNavigate();

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

  const {
    mutate: mutateUpdateNetwork,
    isSuccess: isSuccessUpdateNetwork,
    isError,
  } = usePost<NetworkRequest, NetworkResponse>(
    `${URL.PATH_NETWORKS}/${currentNetwork}`,
  );

  if (isError) {
    toast.error(Intl.t('unexpected-error.description'));
  }

  useEffect(() => {
    if (isSuccessUpdateNetwork) {
      navigate(0);
    }
  }, [isSuccessUpdateNetwork]);

  const availableNetworksItems = availableNetworks.map((network) => ({
    item: network,
    onClick: () => {
      setCurrentNetwork(network);
      mutateUpdateNetwork({});
    },
  }));

  return (
    <div
      data-testid="layout-station"
      className="min-h-screen bg-primary px-20 pt-12 pb-8"
    >
      <div className="grid grid-cols-3">
        <div className="flex justify-start items-center">
          <a href="/">
            <StationLogo theme={selectedTheme} />
          </a>
          {version && getVersionSuccess ? (
            <DisplayVersion version={version} />
          ) : null}
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

function DisplayVersion({ ...props }) {
  const { version } = props;

  return <p className="text-info ml-4 mas-body">{version}</p>;
}
