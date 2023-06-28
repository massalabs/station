import { ReactNode, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';

import { useNetworkStore } from '../../store/store';

import { ThemeMode, StationLogo, Dropdown } from '@massalabs/react-ui-kit';
import { routeFor } from '../../utils';

export interface LayoutStationProps {
  children?: ReactNode;
  navigator?: Navigator;
  onSetTheme?: () => void;
  storedTheme?: string;
}

export function LayoutStation({ ...props }) {
  const { children, navigator, onSetTheme, storedTheme } = props;

  const { network } = useParams();
  const navigate = useNavigate();

  const [selectedTheme, setSelectedTheme] = useState(
    storedTheme || 'theme-dark',
  );

  function handleSetTheme(theme: string) {
    setSelectedTheme(theme);

    onSetTheme?.(theme);
  }

  const [availableNetworks] = useNetworkStore((state) => [
    state.availableNetworks,
  ]);

  const selectedNetworkKey: number = parseInt(
    Object.keys(availableNetworks).find(
      (_, idx) => availableNetworks[idx] === network,
    ) || '0',
  );

  const availableNetworksItems = availableNetworks.map((network) => ({
    item: network,
    onClick: () => {
      navigate(routeFor(network));
    },
  }));

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
          <div className="w-64">
            <Dropdown
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
