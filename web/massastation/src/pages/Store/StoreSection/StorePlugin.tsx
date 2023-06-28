import { Button, Certificate, Plugin } from '@massalabs/react-ui-kit';
import { massalabsNomination } from '../../../utils/massaConstants';
import { FiDownload } from 'react-icons/fi';
import { usePost } from '../../../custom/api';
import { useEffect } from 'react';
import { IMassaStore } from '../../../../../shared/interfaces/IPlugin';

interface StorePluginProps {
  plugin: IMassaStore;
  refetch: () => void;
}

function LoadingDownload() {
  return (
    <div className={`rounded-full animate-pulse blur-sm`}>
      <FiDownload />
    </div>
  );
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
      <Button onClick={() => mutate({ params })} disabled={isInstallLoading}>
        <FiDownload />
      </Button>
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
