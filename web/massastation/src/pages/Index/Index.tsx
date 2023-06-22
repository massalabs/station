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

export function Index() {
  const navigate = useNavigate();
  const [pluginWalletIsInstalled, setPluginWalletIsInstalled] = useState(false);
  const [urlPlugin, setUrlPlugin] = useState('');
  const [refreshPlugins, setRefreshPlugins] = useState(0);

  const { data: plugins } = useResource<PluginHomePage[]>('plugin-manager');

  const { data: availablePlugins } =
    useResource<PluginStoreItemRequest[]>('plugin-store');

  const { mutate, isSuccess, isError } = usePost<any, any>(
    `plugin-manager?source=${urlPlugin}`,
  );

  useEffect(() => {
    const isWalletInstalled = plugins?.some(
      (plugin: PluginHomePage) => plugin.name === 'wallet',
    );
    setPluginWalletIsInstalled(Boolean(isWalletInstalled));
    if (!isWalletInstalled && availablePlugins) {
      const walletPlugin = availablePlugins.find(
        (plugin: PluginStoreItemRequest) => plugin.name === 'MassaWallet',
      );
      if (walletPlugin) {
        setUrlPlugin(walletPlugin.file.url);
      }
    }
  }, [plugins, availablePlugins]);

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
  }, [pluginWalletIsInstalled]);

  function handleInstallPlugin() {
    try {
      mutate({});
    } catch (error) {
      console.error('Error installing plugin:', error);
    }
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
                  title="Massa Wallet"
                  iconActive={<WalletActive />}
                  iconInactive={<WalletInactive />}
                  onClickActive={() =>
                    navigate('/plugin/massa-labs/massa-wallet/web-app/index')
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
