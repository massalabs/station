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
function StorePlugin(props: StorePluginProps) {
  const { plugin, refetch } = props;
  let {
    author,
    name,
    logo,
    description,
    file: { url },
  } = plugin;
  const { mutate, isSuccess, isLoading } = usePost(`plugin-manager`);

  const params = { source: url };

  useEffect(() => {
    if (isSuccess) {
      refetch();
    }
  }, [isSuccess]);

  const argsStore = {
    preIcon: <img src={logo} alt="plugin-logo" />,
    topAction: (
      <Button
        customClass={isLoading ? 'animate-pulse blur-sm' : ''}
        onClick={() => mutate({ params })}
        disabled={isLoading}
      >
        <div
          className={`rounded-full ${isLoading ? 'animate-pulse blur-sm' : ''}`}
        >
          <FiDownload />
        </div>
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
