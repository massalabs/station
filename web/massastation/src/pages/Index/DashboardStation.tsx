import { PluginWallet } from '@massalabs/react-ui-kit';
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
import { PluginExecuteRequest } from '../Store/StationSection/StationPlugin';
import { usePost } from '@/custom/api';
import { useEffect } from 'react';

export interface IDashboardStationProps {
  massaPlugins?: MassaPluginModel[] | undefined;
  pluginWalletIsInstalled: boolean;
  urlPlugin?: string | undefined;
  isLoading: boolean;
  handleInstallPlugin: (url: string) => void;
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

  const {
    mutate: mutateExecute,
    isSuccess: isExecuteSuccess,
    isLoading: isExecuteLoading,
  } = usePost<PluginExecuteRequest>(`plugin-manager/${id}/execute`);

  useEffect(() => {
    if (isUpdatable) {
      console.log('plugin is updatable');
    } else {
      console.log('plugin is not updatable');
    }
    if (isExecuteSuccess) {
      console.log('plugin updated');
    }
  }, [isUpdatable, isExecuteLoading, isExecuteSuccess]);

  function updatePluginState(command: string) {
    if (isExecuteLoading) return;
    const payload = { command } as PluginExecuteRequest;
    mutateExecute({ payload });
  }

  return (
    <div
      className="grid lg:grid-cols-10  grid-rows-3 gap-4 h-fit"
      data-testid="dashboard-station"
    >
      <div className="col-start-1 col-span-2  row-span-3">
        <MassaWallet
          key="wallet"
          isActive={pluginWalletIsInstalled}
          status={
            massaPlugins?.find(
              (plugin: MassaPluginModel) => plugin.name === MASSA_WALLET,
            )?.status
          }
          isUpdatable={isUpdatable}
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
