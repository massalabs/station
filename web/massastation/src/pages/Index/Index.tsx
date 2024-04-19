import { Button } from '@massalabs/react-ui-kit';
import { FiCodepen, FiGlobe } from 'react-icons/fi';
import { useNavigate } from 'react-router-dom';
import { useEffect, useState } from 'react';
import { MassaPluginModel, MassaStoreModel } from '@/models';
import Intl from '@/i18n/i18n';
import { routeFor } from '@/utils/utils';
import { useConfigStore } from '@/store/store';
import { usePost, useResource } from '../../custom/api';
import { UseQueryResult } from '@tanstack/react-query';
import { MASSA_WALLET } from '@/const';
import { DashboardStation } from './DashboardStation';

export function Index() {
  const plugins = useResource<MassaPluginModel[]>('plugin-manager');
  const store = useResource<MassaStoreModel[]>('plugin-store');
  const { isLoading: isPluginsLoading } = plugins;
  const { isLoading: isStoreLoading } = store;

  return isPluginsLoading || isStoreLoading ? (
    <>Loading</>
  ) : (
    <NestedIndex store={store} plugins={plugins}></NestedIndex>
  );
}

function NestedIndex({
  store,
  plugins,
}: {
  store: UseQueryResult<MassaStoreModel[]>;
  plugins: UseQueryResult<MassaPluginModel[]>;
}) {
  const navigate = useNavigate();
  const [pluginWalletIsInstalled, setPluginWalletIsInstalled] = useState(false);
  const [urlPlugin, setUrlPlugin] = useState<string | undefined>(undefined);
  const theme = useConfigStore((s) => s.theme);

  const { data: massaPlugins } = plugins;

  const availablePlugins = store.data;
  const {
    mutate: installPlugin,
    isSuccess: installSuccess,
    isError: installError,
    isLoading,
  } = usePost<null>('plugin-manager');

  useEffect(() => {
    const isWalletInstalled = massaPlugins?.some(
      (plugin: MassaPluginModel) => plugin.name === MASSA_WALLET,
    );
    if (isWalletInstalled) {
      setPluginWalletIsInstalled(true);
    } else {
      const storeWalletPlugin = availablePlugins?.find(
        (plugin: MassaStoreModel) => plugin.name === MASSA_WALLET,
      );
      setUrlPlugin(storeWalletPlugin?.file.url);
    }
  }, [plugins, availablePlugins, massaPlugins]);

  useEffect(() => {
    // we should check that installed plugin is actually the wallet
    if (installSuccess) {
      setPluginWalletIsInstalled(true);
    }
    if (installError) {
      setPluginWalletIsInstalled(false);
    }
  }, [installSuccess, installError]);

  function handleInstallPlugin(url: string) {
    const params = { source: url };
    installPlugin({ params });
  }

  return (
    <>
      <div className="bg-primary text-f-primary pt-24">
        <h1 className="mas-banner mb-10"> {Intl.t('index.title-banner')}</h1>
        <div className="w-fit max-w-[1760px] ">
          <div className="flex gap-8 pb-10">
            <Button
              preIcon={<FiGlobe />}
              customClass="w-96"
              onClick={() => navigate(routeFor('search'))}
            >
              <div className="flex items-center mas-buttons">
                {Intl.t('index.buttons.search')}
              </div>
            </Button>
            <Button
              variant="secondary"
              preIcon={<FiCodepen />}
              customClass="w-96"
              onClick={() => navigate(routeFor('store'))}
            >
              <div className="flex items-center mas-buttons">
                {Intl.t('index.buttons.explore')}
              </div>
            </Button>
          </div>
          <DashboardStation
            massaPlugins={massaPlugins}
            pluginWalletIsInstalled={pluginWalletIsInstalled}
            urlPlugin={urlPlugin}
            isLoading={isLoading}
            handleInstallPlugin={handleInstallPlugin}
            theme={theme}
          />
        </div>
      </div>
    </>
  );
}
