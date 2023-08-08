import { useEffect } from 'react';
import { usePost } from '@/custom/api';

import { Certificate, Plugin } from '@massalabs/react-ui-kit';
import { massalabsNomination } from '@/const';
import { FiDownload } from 'react-icons/fi';
import { PluginCardLoading } from '../PluginCardLoading';
import { MassaStoreModel } from '@/models';

interface StorePluginProps {
  plugin: MassaStoreModel;
  refetch: () => void;
}

function StorePlugin(props: StorePluginProps) {
  const { plugin, refetch } = props;
  let {
    author,
    name,
    logo,
    description,
    file: { url },
    iscompatible,
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
    topAction: (
      <FiDownload onClick={() => !iscompatible && mutate({ params })} />
    ),
    title: name,
    subtitle: author,
    subtitleIcon: massalabsNomination.includes(author) ? <Certificate /> : null,
    content: description,
  };

  return (
    <>{isInstallLoading ? <PluginCardLoading /> : <Plugin {...argsStore} />}</>
  );
}

export default StorePlugin;
