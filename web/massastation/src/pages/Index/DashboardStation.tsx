import { PluginWallet, Theme } from '@massalabs/react-ui-kit';
import { ReactComponent as WalletActive } from '../../assets/wallet/walletActive.svg';
import { ReactComponent as WalletInactive } from '../../assets/wallet/walletInactive.svg';
import { Foundation } from './Dashboard/Foundation';
import { Bridge } from './Dashboard/Bridge';
import { MassaLabs } from './Dashboard/Massalabs';
import { Explorer } from './Dashboard/Explorer';
import { Purrfect } from './Dashboard/Purrfect';
import { Dusa } from './Dashboard/Dusa';
import { MASSA_WALLET } from '@/const';
import { MassaPluginModel } from '@/models';

export interface IDashboardStationProps {
  massaPlugins?: MassaPluginModel[] | undefined;
  pluginWalletIsInstalled: boolean;
  urlPlugin?: string | undefined;
  isLoading: boolean;
  handleInstallPlugin: (url: string) => void;
  theme?: Theme | undefined;
}

export function DashboardStation(props: IDashboardStationProps) {
  let {
    pluginWalletIsInstalled,
    urlPlugin,
    isLoading,
    handleInstallPlugin,
    massaPlugins,
  } = props;

  return (
    <div
      className="grid lg:grid-cols-10  grid-rows-3 gap-4 h-fit p-4"
      data-testid="dashboard-station"
    >
      <div className="col-start-1 col-span-2  row-span-3">
        <PluginWallet
          key="wallet"
          isActive={pluginWalletIsInstalled}
          status={
            massaPlugins?.find(
              (plugin: MassaPluginModel) => plugin.name === MASSA_WALLET,
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
        />
      </div>
      <div className="col-start-3 col-span-2 row-start-1 row-span-2">
        <Foundation />
      </div>
      <div className="col-start-5 col-span-2 row-start-1 row-span-2">
        <Bridge />
      </div>
      <div className="col-start-7 col-span-4 row-start-1 row-span-1">
        <MassaLabs />
      </div>
      <div className="col-start-3 col-span-4 row-start-3 row-span-1">
        <Explorer />
      </div>
      <div className="col-start-7 col-span-2 row-start-2 row-span-2">
        <Purrfect />
      </div>
      <div className="col-start-9 col-span-2 row-start-2 row-span-2">
        <Dusa />
      </div>
    </div>
  );
}
