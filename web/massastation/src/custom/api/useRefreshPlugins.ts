import { useQueryClient } from '@tanstack/react-query';

export function useRefreshPlugins() {
  const queryClient = useQueryClient();

  const refreshInstalledPlugins = () => {
    queryClient.invalidateQueries({
      queryKey: ['', `${import.meta.env.VITE_BASE_API}/plugin-manager`]
    });
  };

  const refreshStorePlugins = () => {
    queryClient.invalidateQueries({
      queryKey: ['', `${import.meta.env.VITE_BASE_API}/plugin-store`]
    });
  };

  const refreshAllPlugins = () => {
    refreshInstalledPlugins();
    refreshStorePlugins();
  };

  return {
    refreshInstalledPlugins,
    refreshStorePlugins,
    refreshAllPlugins,
  };
}
