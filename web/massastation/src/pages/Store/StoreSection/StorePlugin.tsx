import { useEffect } from 'react';
import { usePost } from '@/custom/api';

import { Certificate, Plugin, Tooltip } from '@massalabs/react-ui-kit';
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
  let {
    author,
    name,
    logo,
    description,
    file: { url },
    version,
    iscompatible: isCompatible,
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

  const incompatibleActions = [
    <div className="relative whitespace-nowrap">
      <Tooltip
        content={warningMessage}
        icon={<FiAlertTriangle className="text-s-warning" size={24} />}
      />
    </div>,
    <FiDownload className="w-6 h-10 text-tertiary" />,
  ];

  const downloadAction = [
    <FiDownload
      className="hover:cursor-pointer"
      size={24}
      onClick={() => mutate({ params })}
    />,
  ];

  const argsStore = {
    preIcon: <img src={logo} alt="plugin-logo" />,
    topActions: !isCompatible ? incompatibleActions : downloadAction,
    title: name,
    subtitle: author,
    subtitleIcon: massalabsNomination.includes(author) ? <Certificate /> : null,
    version: `v${version}`,
    content: description,
  };

  return (
    <>{isInstallLoading ? <PluginCardLoading /> : <Plugin {...argsStore} />}</>
  );
}

export default StorePlugin;
