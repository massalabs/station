import { Certificate, MassaWallet, Plugin } from '@massalabs/react-ui-kit';
import { massalabsNomination } from '../../../utils/massaConstants';
import { FiDownload } from 'react-icons/fi';
import { usePost } from '../../../custom/api';
import { useEffect } from 'react';
import { IMassaStore } from './StoreSection';

function StorePlugin({
  plugin,
  refetch,
}: {
  plugin: IMassaStore;
  refetch: () => void;
}) {
  let {
    author,
    name,
    logo,
    description,
    file: { url },
  } = plugin;
  const { mutate, isSuccess } = usePost(`plugin-manager?source=${url}`);
  const argsStore = {
    preIcon: massalabsNomination.includes(author) ? (
      <MassaWallet variant="rounded" />
    ) : (
      <img src={logo} />
    ),
    topAction: <FiDownload onClick={() => mutate({})} />,
    title: name,
    subtitle: author,
    subtitleIcon: massalabsNomination.includes(author) ? <Certificate /> : null,
    content: description,
  };
  useEffect(() => {
    if (isSuccess) {
      refetch();
    }
  }, [isSuccess]);

  return (
    <>
      <Plugin {...argsStore} />
    </>
  );
}

export default StorePlugin;
