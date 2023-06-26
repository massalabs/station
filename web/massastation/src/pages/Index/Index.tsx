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
  PluginHomePage,
  PluginStoreItemRequest,
} from '../../../../shared/interfaces/IPlugin';
import { usePost, useResource } from '../../custom/api';
import Intl from '../../i18n/i18n';
import { routeFor } from '../../utils';
import { useConfigStore } from '../../store/store';

export function Index() {
  const navigate = useNavigate();
  const [pluginWalletIsInstalled, setPluginWalletIsInstalled] = useState(false);
  const [urlPlugin, setUrlPlugin] = useState('intit');
  const [refreshPlugins, setRefreshPlugins] = useState(0);
  const walletName = 'Massa Wallet';

  const theme = useConfigStore((s) => s.theme);

  const { data: plugins } = useResource<PluginHomePage[]>('plugin-manager');

  const { data: availablePlugins } =
    useResource<PluginStoreItemRequest[]>('plugin-store');

  const { mutate, isSuccess, isError, isLoading } =
    usePost<null>('plugin-manager');

  function checkIfWalletIsInstalled() {
    console.log('checking');
    if (plugins) {
      setPluginWalletIsInstalled(
        plugins.some((plugin: PluginHomePage) => plugin.name === walletName) ||
          false,
      );
    }
  }

  function availablePluginsList(name: string) {
    const walletPlugin =
      availablePlugins?.find((plugin) => plugin.name === name) || null;
    console.log('walletPlugin', walletPlugin);
    setUrlPlugin(walletPlugin?.file.url || 'reinit');
  }

  useEffect(() => {
    checkIfWalletIsInstalled();
    console.log(pluginWalletIsInstalled, availablePlugins);
    if (!pluginWalletIsInstalled && availablePlugins) {
      availablePluginsList(walletName);
    }
    console.log('urlPlugin', urlPlugin);
  }, [pluginWalletIsInstalled, availablePlugins, urlPlugin]);

  useEffect(() => {
    if (isSuccess) {
      setPluginWalletIsInstalled(true);
    }
    if (isError) {
      setPluginWalletIsInstalled(false);
    }
  }, [isLoading]);

  useEffect(() => {
    setRefreshPlugins(refreshPlugins + 1);
  }, [pluginWalletIsInstalled, isLoading]);

  function handleInstallPlugin() {
    console.log(urlPlugin);
    const params = { source: urlPlugin };
    console.log('params', params);
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
            {/* <div>
              <Button onClick={handleInstallPlugin}>install</Button>
            </div> */}
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
                  onClickInactive={handleInstallPlugin}
                />,
              ]}
            />
          </div>
        </div>
      </div>
    </>
  );
}
