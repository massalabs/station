import { ReactNode, useState, useEffect } from 'react';

import { useResource, usePost } from '../../custom/api';
import { URL } from '../../const/url/url';
import { NetworkModel } from '../../models';
import { useNetworkStore } from '../../store/store';

import { ThemeMode, StationLogo, Dropdown, Theme, toast } from '@massalabs/react-ui-kit';
import { useNavigate } from 'react-router-dom';
import Intl from '@/i18n/i18n';
import { THEME_STORAGE_KEY } from '@/const';

interface LayoutStationProps {
  children?: ReactNode;
  navigator?: ReactNode;
  onSetTheme?: (theme: Theme) => void;
  storedTheme?: Theme | undefined;
}

// Typings for switch network
type SwitchNetworkResponse = NetworkModel; // server returns NetworkManagerItem shape
type SwitchNetworkBody = Record<string, never>; // no body expected

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

  const [currentNetwork, setCurrentNetwork] = useNetworkStore((state) => [
    state.currentNetwork,
    state.setCurrentNetwork,
  ]);

  useEffect(() => {
    if (isSuccessNetwork) {
      if (network.currentNetwork) setCurrentNetwork(network.currentNetwork);
    }
  }, [isSuccessNetwork, network, setCurrentNetwork]);

  const infos = network?.availableNetworkInfos || [];

  const [targetNetwork, setTargetNetwork] = useState<string | null>(null);
  const {
    mutate: mutateSwitch,
    isSuccess: isSuccessSwitch,
    isError: isErrorSwitch,
  } = usePost<SwitchNetworkBody, SwitchNetworkResponse>(
    `${URL.PATH_NETWORKS}/${encodeURIComponent(targetNetwork ?? currentNetwork ?? '')}`,
  );

  useEffect(() => {
    if (isErrorSwitch) {
      toast.error(Intl.t('unexpected-error.description'));
    }
  }, [isErrorSwitch]);

  useEffect(() => {
    if (targetNetwork) {
      mutateSwitch({});
    }
  }, [targetNetwork, mutateSwitch]);

  useEffect(() => {
    if (isSuccessSwitch && targetNetwork) {
      setCurrentNetwork(targetNetwork);
      setTargetNetwork(null);
      navigate(0);
    }
  }, [isSuccessSwitch, targetNetwork, setCurrentNetwork, navigate]);

  const availableNetworksItems = infos.map((nfo) => ({
    item: nfo.status === 'down'
      ? `${nfo.name} (Offline)`
      : nfo.name,
    title: nfo.status === 'down'
      ? `${nfo.name} (Offline) - ${nfo.url}`
      : `${nfo.name} - ${nfo.url}`,
    onClick: () => {
      if (nfo.status !== 'down') {
        setTargetNetwork(nfo.name);
      }
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
          <div className="min-w-44 max-w-80 w-auto">
            <Dropdown
              size="xs"
              options={availableNetworksItems}
              select={currentNetwork ? infos.findIndex((n) => n.name === currentNetwork) : 0}
            />
          </div>
          <ThemeMode
            onSetTheme={handleSetTheme}
            storageKey={THEME_STORAGE_KEY}
          />
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
