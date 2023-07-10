import { useEffect } from 'react';
import { usePost } from '@/custom/api';

import { Certificate, Plugin } from '@massalabs/react-ui-kit';
import { massalabsNomination } from '@/const';
import { FiDownload } from 'react-icons/fi';
import { IMassaStore } from '@/shared/interfaces/IPlugin';

interface StorePluginProps {
  plugin: IMassaStore;
  refetch: () => void;
}

function LoadingDownload() {
  return <FiDownload className="animate-ping" />;
}

function StorePlugin(props: StorePluginProps) {
  const { plugin, refetch } = props;
  let {
    author,
    name,
    logo,
    description,
    file: { url },
  } = plugin;

  const {
    mutate,
    isSuccess: isInstallSuccess,
    isLoading: isInstallLoading,
  } = usePost(`plugin-manager`);

  const params = { source: url };

  useEffect(() => {
    if (isInstallSuccess) {
      refetch();
    }
  }, [isInstallSuccess]);

  const argsStore = {
    preIcon: <img src={logo} alt="plugin-logo" />,
    topAction: isInstallLoading ? (
      <LoadingDownload />
    ) : (
      <FiDownload onClick={() => mutate({ params })} />
    ),
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
