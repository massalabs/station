import { useEffect, useMemo, useState } from 'react';
import { MassaPluginModel, MassaLabsPlugins, MassaStoreModel } from '@/models';
import { PluginStates } from '@/pages/Index/DashboardStation';
import { PluginCommand, useUpdatePlugin } from './useUpdatePlugin';
import { usePost } from '../api';
import { useRefreshPlugins } from '../api/useRefreshPlugins';

export interface UsePluginStateReturn {
  state?: PluginStates;
  plugin?: MassaPluginModel;
  isUpdating: boolean;
  isInstalling: boolean;
  updatePlugin: () => void;
  installPlugin: () => void;
}

export function usePluginState(
  pluginName: MassaLabsPlugins,
  massaPlugins?: MassaPluginModel[],
  availablePlugins?: MassaStoreModel[],
): UsePluginStateReturn {

  const [state, setState] = useState<PluginStates>();
  const pluginInfos = useMemo(
    () =>
      massaPlugins?.find(
        (plugin: MassaPluginModel) => plugin.name === pluginName,
      ),
    [massaPlugins, pluginName],
  );

  const { isExecuteSuccess, isExecuteLoading, updatePluginState } =
    useUpdatePlugin(pluginInfos?.id);

  const { refreshInstalledPlugins } = useRefreshPlugins();
  const {
    mutate: installPluginMutate,
    isLoading: isInstallLoading,
  } = usePost<null>('plugin-manager');

  const isInstalled = !!pluginInfos;
  const isUpdatable = pluginInfos?.updatable;

  useEffect(() => {
    if (isInstalled && !isUpdatable) {
      setState(PluginStates.Active);
    } else if (isInstalled && isUpdatable) {
      setState(PluginStates.Updateable);
    } else {
      setState(PluginStates.Inactive);
    }
  }, [isUpdatable, isInstalled, pluginInfos]);

  useEffect(() => {
    if (isExecuteSuccess) {
      setState(PluginStates.Active);
    }
  }, [isExecuteSuccess]);

  function installPlugin() {
    const installUrl = availablePlugins?.find(
      (plugin: MassaStoreModel) => plugin.name === pluginName,
    )?.file.url;

    if (!installUrl || isInstallLoading) return;

    const params = { source: installUrl };
    installPluginMutate(
      { params },
      {
        onSuccess: () => {
          // Refresh installed plugins so that massaPlugins is updated upstream.
          refreshInstalledPlugins();
        },
      },
    );
  }

  return {
    state,
    plugin: pluginInfos,
    isUpdating: isExecuteLoading,
    isInstalling: isInstallLoading,
    updatePlugin: () => updatePluginState(PluginCommand.Update),
    installPlugin,
  };
}
