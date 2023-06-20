import { Certificate, MassaWallet, Plugin } from '@massalabs/react-ui-kit';
import { IMassaPlugin } from './StoreSection';
import { massalabsNomination } from '../../../utils/massaConstants';
import { FiDownload } from 'react-icons/fi';

function StorePlugin({ plugin }: { plugin: IMassaPlugin }) {
  let { author, name, logo, description } = plugin;

  const argsStore = {
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
      <Plugin {...argsStore} />
    </>
  );
}

export default StorePlugin;
