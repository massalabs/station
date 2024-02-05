import { ReactComponent as WalletActive } from '../../assets/wallet/walletActive.svg';
import { ReactComponent as WalletInactive } from '../../assets/wallet/walletInactive.svg';
import { ReactComponent as Image1Dark } from '../../assets/subduedImages/dark/1.svg';
import { ReactComponent as Image2Dark } from '../../assets/subduedImages/dark/2.svg';
import { ReactComponent as Image3Dark } from '../../assets/subduedImages/dark/3.svg';
import { ReactComponent as Image4Dark } from '../../assets/subduedImages/dark/4.svg';
import { ReactComponent as Image5Dark } from '../../assets/subduedImages/dark/5.svg';
import { ReactComponent as Image1Light } from '../../assets/subduedImages/light/1.svg';
import { ReactComponent as Image2Light } from '../../assets/subduedImages/light/2.svg';
import { ReactComponent as Image3Light } from '../../assets/subduedImages/light/3.svg';
import { ReactComponent as Image4Light } from '../../assets/subduedImages/light/4.svg';
import { ReactComponent as Image5Light } from '../../assets/subduedImages/light/5.svg';
import { Button, PluginWallet } from '@massalabs/react-ui-kit';
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
        <div className="overflow-auto">
          <div className="w-[70vw]">
            <div className="flex space-x-8 pb-10">
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
              theme={theme}
              imagesDark={[
                <Image1Dark />,
                <Image2Dark />,
                <Image3Dark />,
                <Image4Dark />,
                <Image5Dark />,
              ]}
              imagesLight={[
                <Image1Light />,
                <Image2Light />,
                <Image3Light />,
                <Image4Light />,
                <Image5Light />,
              ]}
              components={[
                <PluginWallet
                  key="wallet"
                  isActive={pluginWalletIsInstalled}
                  status={
                    massaPlugins?.find(
                      (plugin: MassaPluginModel) =>
                        plugin.name === MASSA_WALLET,
                    )?.status
                  }
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
                />,
              ]}
            />
          </div>
        </div>
      </div>
    </>
  );
}
