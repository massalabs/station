import { PluginExecuteRequest } from '@/pages/Store/StationSection/StationPlugin';
import { usePost } from '../api';


export const enum PluginCommand {
  Start = 'start',
  Stop = 'stop',
  Update = 'update',
  Restart = 'restart',
}

export function useUpdatePlugin(id: string | undefined) {
  const {
    mutate: mutateExecute,
    isSuccess: isExecuteSuccess,
    isLoading: isExecuteLoading,
  } = usePost<PluginExecuteRequest>(`plugin-manager/${id}/execute`);

  function updatePluginState(command: PluginCommand) {
    if (isExecuteLoading) return;
    const payload = { command } as PluginExecuteRequest;
    mutateExecute({ payload });
  }

  return { isExecuteSuccess, isExecuteLoading, updatePluginState };
}
