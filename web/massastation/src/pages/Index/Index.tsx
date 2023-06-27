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
import {
  Button,
  DashboardStation,
  PluginWallet,
} from '@massalabs/react-ui-kit';
import { FiCodepen, FiGlobe } from 'react-icons/fi';
import { useNavigate } from 'react-router-dom';
import { useEffect, useState } from 'react';
import {
  IMassaPlugin,
  IMassaStore,
} from '../../../../shared/interfaces/IPlugin';
import Intl from '../../i18n/i18n';
import { routeFor } from '../../utils';
import { useConfigStore } from '../../store/store';
import { usePost, useResource } from '../../custom/api';
import { UseQueryResult } from '@tanstack/react-query';
import { MASSA_WALLET } from '../../const/const';

export function Index() {
  const plugins = useResource<IMassaPlugin[]>('plugin-manager');
  const store = useResource<IMassaStore[]>('plugin-store');
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
  store: UseQueryResult<IMassaStore[]>;
  plugins: UseQueryResult<IMassaPlugin[]>;
}) {
  const navigate = useNavigate();
  const [pluginWalletIsInstalled, setPluginWalletIsInstalled] = useState(false);
  const [urlPlugin, setUrlPlugin] = useState('');
  const [refreshPlugins, setRefreshPlugins] = useState(0);
  const theme = useConfigStore((s) => s.theme);

  const { data: massaPlugins } = plugins;

  const { data: availablePlugins } = useResource<IMassaStore[]>('plugin-store');

  const { mutate, isSuccess, isError, isLoading } =
    usePost<null>('plugin-manager');

  useEffect(() => {
    const isWalletInstalled = massaPlugins?.some(
      (plugin: IMassaPlugin) => plugin.name === MASSA_WALLET,
    );
    setPluginWalletIsInstalled(Boolean(isWalletInstalled));
    if (!isWalletInstalled && availablePlugins) {
      const walletPlugin = availablePlugins.find(
        (plugin: IMassaStore) => plugin.name === MASSA_WALLET,
      );
      if (walletPlugin) {
        setUrlPlugin(walletPlugin.file.url);
      }
    }
  }, [plugins, store]);

  useEffect(() => {
    if (isSuccess) {
      setPluginWalletIsInstalled(true);
    }
    if (isError) {
      setPluginWalletIsInstalled(false);
    }
  }, [isSuccess, isError]);

  useEffect(() => {
    setRefreshPlugins(refreshPlugins + 1);
  }, [pluginWalletIsInstalled, isLoading]);

  function handleInstallPlugin(url: string) {
    const params = { source: url };
    mutate({ params });
  }

  return (
    <>
      <div className="bg-primary text-f-primary pt-24">
        <h1 className="mas-banner mb-10"> {Intl.t('index.title-banner')}</h1>
        <div className="overflow-auto h-[65vh]">
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
              key={refreshPlugins}
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
                  onClickInactive={() => handleInstallPlugin(urlPlugin)}
                />,
              ]}
            />
          </div>
        </div>
      </div>
    </>
  );
}
