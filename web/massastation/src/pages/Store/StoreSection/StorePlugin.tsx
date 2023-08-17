import { useEffect, useState } from 'react';
import { usePost } from '@/custom/api';

import { Certificate, Plugin } from '@massalabs/react-ui-kit';
import { massalabsNomination } from '@/const';
import { FiAlertTriangle, FiDownload } from 'react-icons/fi';
import { PluginCardLoading } from '../PluginCardLoading';
import { MassaStoreModel } from '@/models';
import Intl from '@/i18n/i18n';
interface StorePluginProps {
  plugin: MassaStoreModel;
  refetch: () => void;
}

function StorePlugin(props: StorePluginProps) {
  const { plugin, refetch } = props;
  const [message, setMessage] = useState<boolean>(false);
  let {
    author,
    name,
    logo,
    description,
    file: { url },
    iscompatible,
    massastationMinVersion,
  } = plugin;

  const warningMessage = Intl.t('store.massa-station-incompatible', {
    version: massastationMinVersion.slice(2),
  });
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

  const argsStoreError = [
    <div>
      <div
        onMouseEnter={() => setMessage(true)}
        onMouseLeave={() => setMessage(false)}
      >
        <FiAlertTriangle className="text-s-warning w-6 h-10" />
      </div>
      {message && (
        <div className="w-fit absolute z-10 l-10 bg-tertiary p-3 rounded-lg text-neutral">
          {warningMessage}
        </div>
      )}
    </div>,
    <FiDownload className="w-6 h-10 text-tertiary" />,
  ];

  const argsStore = {
    preIcon: <img src={logo} alt="plugin-logo" />,
    topAction: <FiDownload onClick={() => mutate({ params })} />,
    title: name,
    subtitle: author,
    subtitleIcon: massalabsNomination.includes(author) ? <Certificate /> : null,
    content: description,
    topActions: !iscompatible ? argsStoreError : undefined,
  };

  return (
    <>{isInstallLoading ? <PluginCardLoading /> : <Plugin {...argsStore} />}</>
  );
}

export default StorePlugin;
