import { useState } from 'react';
import { Certificate, MassaWallet, Plugin } from '@massalabs/react-ui-kit';
import { IMassaPlugin } from './MyStore';
import { massalabsNomination } from '../../utils/massaConstants';
import { FiDownload } from 'react-icons/fi';

enum PluginStatus {
  // Up And Down are sent by the BE, we use the On/Off standard to stick to design
  On = 'Up',
  Off = 'Down',
}

export default function StorePlugin({ plugin }: { plugin: IMassaPlugin }) {
  const [myPlugin, setMyPlugin] = useState<IMassaPlugin>(plugin);

  let { author, name, logo, status, description } = myPlugin;

  const argsOn = {
    preIcon: massalabsNomination.includes(author) ? (
      <MassaWallet variant="rounded" />
    ) : (
      <img src={logo} />
    ),
    topAction: <FiDownload />,
    title: name,
    subtitle: author,
    subtitleIcon: massalabsNomination.includes(author) ? (
      <Certificate />
    ) : (
      <></>
    ),
    content: description,
  };
  const argsOff = {
    preIcon: massalabsNomination.includes(author) ? (
      <MassaWallet variant="rounded" />
    ) : (
      <img src={logo} />
    ),
    topAction: <FiDownload />,
    title: name,
    subtitle: author,
    subtitleIcon: massalabsNomination.includes(author) ? <Certificate /> : null,
    content: description,
  };
  return (
    <>
      <Plugin {...(status === PluginStatus.On ? argsOn : argsOff)} />
    </>
  );
}
