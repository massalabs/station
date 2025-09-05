import { useEffect, useState } from 'react';
import { MassaPluginModel } from '@/models';
import { PluginStates } from '@/pages/Index/DashboardStation';
import { PluginCommand, useUpdatePlugin } from './useUpdatePlugin';

export interface UsePluginStateReturn {
  state?: PluginStates;
  plugin?: MassaPluginModel;
  isUpdating: boolean;
  updatePlugin: () => void;
}

export function usePluginState(
  pluginInfos: MassaPluginModel | undefined,
): UsePluginStateReturn {

  const [state, setState] = useState<PluginStates>();
  const { isExecuteSuccess, isExecuteLoading, updatePluginState } =
    useUpdatePlugin(pluginInfos?.id);

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

  return {
    state,
    plugin: pluginInfos,
    isUpdating: isExecuteLoading,
    updatePlugin: () => updatePluginState(PluginCommand.Update),
  };
}
