import { DashboardStation, PluginWallet } from '@massalabs/react-ui-kit';

export function Index() {
  return (
    <DashboardStation
      imagesDark={[]}
      imagesLight={[]}
      components={[
        <PluginWallet
          isActive={true}
          title={'Massawallet'}
          iconActive={null}
          iconInactive={null}
          onClickActive={() => console.log('install')}
          onClickInactive={() => console.log('launch')}
        />,
      ]}
    />
  );
}
