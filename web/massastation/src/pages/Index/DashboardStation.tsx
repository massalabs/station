import { Bridge } from './Dashboard/Bridge';
import { Massa } from './Dashboard/Massa';
import { Explorer } from './Dashboard/Explorer';
import { MassaEcosystem } from './Dashboard/MassaEcosystem';
import { NodeManager } from './Dashboard/NodeManager';
import { Deweb } from './Dashboard/Deweb';
import { Syntra } from './Dashboard/Syntra';
import { MassaGovernance } from './Dashboard/MassaGovernance';
import { BuyMas } from './Dashboard/BuyMas';
import { MASSA_WALLET, NODE_MANAGER } from '@/const';
import { MassaPluginModel, MassaStoreModel } from '@/models';
import { MassaWallet } from './Dashboard/MassaWallet';
import { useEffect, useState } from 'react';

import { usePluginState } from '@/custom/hooks/usePluginState';
import { usePost, useRefreshPlugins } from '@/custom/api';

export interface IDashboardStationProps {
  massaPlugins?: MassaPluginModel[] | undefined;
  availablePlugins?: MassaStoreModel[] | undefined;
}

export enum PluginStates {
  Active = 'Active',
  Inactive = 'Inactive',
  Updateable = 'Updateable',
}

export const PLUGIN_LIST = [MASSA_WALLET, NODE_MANAGER];

export function DashboardStation(props: IDashboardStationProps) {
  const { massaPlugins, availablePlugins } = props;
  const { refreshInstalledPlugins } = useRefreshPlugins();

  // Plugin installation logic
  const {
    mutate: installPlugin,
    isSuccess: installSuccess,
    isError: installError,
    isLoading: isInstalling,
  } = usePost<null>('plugin-manager');

  // Plugin installation states
  const [pluginsInstalled, setPluginsInstalled] = useState<Record<string, MassaPluginModel | undefined>>({});
  const [installingPlugin, setInstallingPlugin] = useState<string | null>(null);

  const installUrl = (pluginName: string) => 
    availablePlugins?.find((plugin: MassaStoreModel) => plugin.name === pluginName)?.file.url;
  
  // Initialize plugin states
  useEffect(() => {
    const installed: Record<string, MassaPluginModel | undefined> = {};

    // Check which plugins are installed
    PLUGIN_LIST.forEach(pluginName => {
      installed[pluginName] = massaPlugins?.find(
        (plugin: MassaPluginModel) => plugin.name === pluginName,
      );
    });

    setPluginsInstalled(installed);
  }, [massaPlugins, availablePlugins]);

  // Handle installation success/error
  useEffect(() => {
    if (installSuccess && installingPlugin) {
      // Invalidate the plugin-manager query to refresh massaPlugins
      refreshInstalledPlugins();
      
      // Mark only the specific plugin that was being installed as installed
      setPluginsInstalled(prev => ({
        ...prev,
        [installingPlugin]: massaPlugins?.find(
          (plugin: MassaPluginModel) => plugin.name === installingPlugin,
        ),
      }));
      setInstallingPlugin(null);
    }
    if (installError) {
      setInstallingPlugin(null);
    }
  }, [installSuccess, installError, installingPlugin, refreshInstalledPlugins, massaPlugins]);

  function handleInstallPlugin(url: string, pluginName?: string) {
    // Track which plugin is being installed
    if (pluginName) {
      setInstallingPlugin(pluginName);
    }
    const params = { source: url };
    installPlugin({ params });
  }

  // Plugin states
  const walletPlugin = usePluginState(pluginsInstalled[MASSA_WALLET]);
  const nodeManagerPlugin = usePluginState(pluginsInstalled[NODE_MANAGER]);

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
          isLoading={isInstalling}
          title={MASSA_WALLET}
          onClickActive={() =>
            window.open(pluginsInstalled[MASSA_WALLET]?.home, '_blank')
          }
          onClickInactive={() =>
            installUrl(MASSA_WALLET)
              ? handleInstallPlugin(installUrl(MASSA_WALLET)!, MASSA_WALLET)
              : null
          }
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
          isLoading={isInstalling}
          title={NODE_MANAGER}
          onClickActive={() =>
            window.open(pluginsInstalled[NODE_MANAGER]?.home, '_blank')
          }
          onClickInactive={() =>
            installUrl(NODE_MANAGER)
              ? handleInstallPlugin(installUrl(NODE_MANAGER)!, NODE_MANAGER)
              : console.log('No install URL for Node Manager')
          }
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
