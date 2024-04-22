import { Theme } from '@massalabs/react-ui-kit';
import { ReactComponent as WalletActive } from '../../assets/wallet/walletActive.svg';
import { ReactComponent as WalletInactive } from '../../assets/wallet/walletInactive.svg';
import { Foundation } from './Dashboard/Foundation';
import { Bridge } from './Dashboard/Bridge';
import { MassaLabs } from './Dashboard/Massalabs';
import { Explorer } from './Dashboard/Explorer';
import { Purrfect } from './Dashboard/Purrfect';
import { Dusa } from './Dashboard/Dusa';
import { MASSA_WALLET, PLUGIN_UPDATE } from '@/const';
import { MassaPluginModel } from '@/models';
import { MassaWallet } from './Dashboard/MassaWallet';
import { useEffect, useState } from 'react';
import { useUpdatePlugin } from '@/custom/hooks/useUpdatePlugin';

export interface IDashboardStationProps {
  massaPlugins?: MassaPluginModel[] | undefined;
  pluginWalletIsInstalled: boolean;
  urlPlugin?: string | undefined;
  isLoading: boolean;
  handleInstallPlugin: (url: string) => void;
  theme?: Theme | undefined;
}

export enum WalletStates {
  Active = 'Active',
  Inactive = 'Inactive',
  Updateable = 'Updateable',
}

export function DashboardStation(props: IDashboardStationProps) {
  let {
    pluginWalletIsInstalled,
    urlPlugin,
    isLoading,
    handleInstallPlugin,
    massaPlugins,
  } = props;

  const id = massaPlugins?.find(
    (plugin: MassaPluginModel) => plugin.name === MASSA_WALLET,
  )?.id;

  const isUpdatable = massaPlugins?.find(
    (plugin: MassaPluginModel) => plugin.name === MASSA_WALLET,
  )?.updatable;

  const [walletState, setWalletState] = useState<WalletStates>();

  const { isExecuteSuccess, isExecuteLoading, updatePluginState } =
    useUpdatePlugin(id);

  useEffect(() => {
    if (pluginWalletIsInstalled && !isUpdatable) {
      setWalletState(WalletStates.Active);
    } else if (pluginWalletIsInstalled && isUpdatable) {
      setWalletState(WalletStates.Updateable);
    } else {
      setWalletState(WalletStates.Inactive);
    }
  }, [isUpdatable, pluginWalletIsInstalled]);

  useEffect(() => {
    if (isExecuteSuccess) {
      setWalletState(WalletStates.Active);
    }
  }, [isExecuteSuccess]);

  return (
    <div
      className="grid lg:grid-cols-10  grid-rows-3 gap-4 h-fit"
      data-testid="dashboard-station"
    >
      <div className="col-start-1 col-span-2 row-span-3">
        <MassaWallet
          key="wallet"
          state={walletState}
          status={
            massaPlugins?.find(
              (plugin: MassaPluginModel) => plugin.name === MASSA_WALLET,
            )?.status
          }
          isUpdating={isExecuteLoading}
          isLoading={isLoading}
          title="Massa Wallet"
          iconActive={<WalletActive />}
          iconInactive={<WalletInactive />}
          onClickActive={() =>
            window.open(
              '/plugin/massa-labs/massa-wallet/web-app/index',
              '_blank',
            )
          }
          onClickInactive={() =>
            urlPlugin ? handleInstallPlugin(urlPlugin) : null
          }
          onUpdateClick={() => updatePluginState(PLUGIN_UPDATE)}
        />
      </div>
      <div className="col-start-3 col-span-2 row-start-1 row-span-2">
        <Foundation />
      </div>
      <div className="col-start-5 col-span-2 row-start-1 row-span-2">
        <Bridge />
      </div>
      <div className="col-start-7 col-span-4 row-start-1 row-span-1">
        <MassaLabs />
      </div>
      <div className="col-start-3 col-span-4 row-start-3 row-span-1">
        <Explorer />
      </div>
      <div className="col-start-7 col-span-2 row-start-2 row-span-2">
        <Purrfect />
      </div>
      <div className="col-start-9 col-span-2 row-start-2 row-span-2">
        <Dusa />
      </div>
    </div>
  );
}
