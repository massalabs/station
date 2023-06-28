import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { UseQueryResult } from '@tanstack/react-query';
import { routeFor } from '../../utils';
import { useConfigStore } from '../../store/store';
import Intl from '../../i18n/i18n';
import { usePost, useResource } from '../../custom/api';

import {
  IMassaPlugin,
  IMassaStore,
} from '../../../../shared/interfaces/IPlugin';
import { MASSA_WALLET } from '../../const/const';

import {
  Button,
  DashboardStation,
  PluginWallet,
} from '@massalabs/react-ui-kit';
import { FiCodepen, FiGlobe } from 'react-icons/fi';
import {
  WalletActive,
  WalletInactive,
  Image1Dark,
  Image2Dark,
  Image3Dark,
  Image4Dark,
  Image5Dark,
  Image1Light,
  Image2Light,
  Image3Light,
  Image4Light,
  Image5Light,
} from './svg';

export function Index() {
  const plugins = useResource<IMassaPlugin[]>('plugin-manager');
  const store = useResource<IMassaStore[]>('plugin-store');
  const { isLoading: isPluginsLoading } = plugins;
  const { isLoading: isStoreLoading } = store;

  return (
    !isPluginsLoading &&
    !isStoreLoading && <NestedIndex store={store} plugins={plugins} />
  );
}

interface NestedIndexProps {
  store: UseQueryResult<IMassaStore[]>;
  plugins: UseQueryResult<IMassaPlugin[]>;
}

function NestedIndex(props: NestedIndexProps) {
  const { store, plugins } = props;

  const navigate = useNavigate();
  const theme = useConfigStore((s) => s.theme);
  const { data: availablePlugins } = useResource<IMassaStore[]>('plugin-store');

  const { mutate, isSuccess, isError, isLoading } =
    usePost<null>('plugin-manager');

  const [pluginWalletIsInstalled, setPluginWalletIsInstalled] =
    useState<boolean>(false);
  const [urlPlugin, setUrlPlugin] = useState<string>('');
  const [refreshPlugins, setRefreshPlugins] = useState<number>(0);

  const { data: massaPlugins } = plugins;

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
              <label className="mas-buttons">
                {Intl.t('index.buttons.search')}
              </label>
            </Button>
            <Button
              variant="secondary"
              preIcon={<FiCodepen />}
              customClass="w-96"
              onClick={() => navigate(routeFor('store'))}
            >
              <label className="mas-buttons">
                {Intl.t('index.buttons.explore')}
              </label>
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
  );
}
