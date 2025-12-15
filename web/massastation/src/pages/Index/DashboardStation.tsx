import { Bridge } from './Dashboard/Bridge';
import { Massa } from './Dashboard/Massa';
import { Explorer } from './Dashboard/Explorer';
import { MassaEcosystem } from './Dashboard/MassaEcosystem';
import { NodeManager } from './Dashboard/NodeManager';
import { Deweb } from './Dashboard/Deweb';
import { Syntra } from './Dashboard/Syntra';
import { MassaGovernance } from './Dashboard/MassaGovernance';
import { BuyMas } from './Dashboard/BuyMas';
import { MassaLabsPlugins, MassaPluginModel, MassaStoreModel } from '@/models';
import { MassaWallet } from './Dashboard/MassaWallet';
import { usePluginState } from '@/custom/hooks/usePluginState';

export interface IDashboardStationProps {
  massaPlugins?: MassaPluginModel[] | undefined;
  availablePlugins?: MassaStoreModel[] | undefined;
}

export enum PluginStates {
  Active = 'Active',
  Inactive = 'Inactive',
  Updateable = 'Updateable',
}

export function DashboardStation(props: IDashboardStationProps) {
  const { massaPlugins, availablePlugins } = props;
  // Plugin states
  const walletPlugin = usePluginState(
    MassaLabsPlugins.MassaWallet,
    massaPlugins,
    availablePlugins,
  );
  const nodeManagerPlugin = usePluginState(
    MassaLabsPlugins.NodeManager,
    massaPlugins,
    availablePlugins,
  );

  return (
    <div
      className="grid lg:grid-cols-15  grid-rows-2 gap-4 h-fit"
      data-testid="dashboard-station"
    >
      <div className="col-start-1 col-span-2 row-start-1 row-span-1">
        <MassaWallet
          key="wallet"
          state={walletPlugin.state}
          status={walletPlugin.plugin?.status}
          isUpdating={walletPlugin.isUpdating}
          isLoading={walletPlugin.isInstalling}
          title={MassaLabsPlugins.MassaWallet}
          onClickActive={() =>
            walletPlugin.plugin &&
            window.open(walletPlugin.plugin.home, '_blank')
          }
          onClickInactive={walletPlugin.installPlugin}
          onUpdateClick={walletPlugin.updatePlugin}
        />
      </div>
      <div className="col-start-1 col-span-2 row-start-2 row-span-1">
        <BuyMas />
      </div>
      <div className="col-start-3 col-span-2 row-start-2 row-span-1">
        <Bridge />
      </div>
      <div className="col-start-9 col-span-2 row-start-1 row-span-2">
        <Massa />
      </div>
      <div className="col-start-7 col-span-2 row-start-2 row-span-1">
        <Explorer />
      </div>
      <div className="col-start-3 col-span-2 row-start-1 row-span-1">
        <MassaEcosystem />
      </div>
      <div className="col-start-5 col-span-2 row-start-2 row-span-1">
        <NodeManager
          key="node-manager"
          state={nodeManagerPlugin.state}
          status={nodeManagerPlugin.plugin?.status}
          isUpdating={nodeManagerPlugin.isUpdating}
          isLoading={nodeManagerPlugin.isInstalling}
          title={MassaLabsPlugins.NodeManager}
          onClickActive={() =>
            nodeManagerPlugin.plugin &&
            window.open(nodeManagerPlugin.plugin.home, '_blank')
          }
          onClickInactive={nodeManagerPlugin.installPlugin}
          onUpdateClick={nodeManagerPlugin.updatePlugin}
        />
      </div>
      <div className="col-start-5 col-span-4 row-start-1 row-span-1">
        <Deweb />
      </div>
      <div className="col-start-11 col-span-2 row-start-2 row-span-1">
        <Syntra />
      </div>
      <div className="col-start-11 col-span-2 row-start-1 row-span-1">
        <MassaGovernance />
      </div>
    </div>
  );
}
