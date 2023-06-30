import { SyntheticEvent, useEffect, useState } from 'react';
import { useDelete, usePost, useResource } from '@/custom/api';

import { Button, Certificate, Plugin } from '@massalabs/react-ui-kit';
import { FiArrowUpRight, FiRefreshCw, FiTrash2 } from 'react-icons/fi';
import {
  massalabsNomination,
  PLUGIN_START,
  PLUGIN_STOP,
  PLUGIN_UPDATE,
} from '@/const';

import { IMassaPlugin, PluginStatus } from '../../../../../shared/';

interface PluginExecuteRequest {
  command: string;
}

const baseAPI = import.meta.env.VITE_BASE_API;

export function StationPlugin({
  plugin,
  fetchPlugins,
}: {
  plugin: IMassaPlugin;
  fetchPlugins: () => void;
}) {
  const [myPlugin, setMyPlugin] = useState<IMassaPlugin>(plugin);
  const { author, name, home, status, updatable, id } = myPlugin;
  const {
    data: newPlugin,
    refetch,
    isLoading,
    isRefetching,
  } = useResource<IMassaPlugin>(`plugin-manager/${id}`);

  const {
    mutate: mutateExecute,
    isSuccess: isExecuteSuccess,
    isLoading: isExecuteLoading,
  } = usePost<PluginExecuteRequest>(`plugin-manager/${id}/execute`);

  const { mutate: deletePlugin, isSuccess: deleteSuccess } = useDelete(
    `plugin-manager/${id}`,
  );

  const logoURL = `${baseAPI}/plugin-manager/${id}/logo`;

  useEffect(() => {
    if (isExecuteSuccess) {
      refetch();
    }
  }, [isExecuteSuccess]);

  useEffect(() => {
    if (newPlugin && !isRefetching && !isLoading) {
      setMyPlugin(newPlugin);
    }
  }, [isRefetching]);

  useEffect(() => {
    if (deleteSuccess) {
      fetchPlugins();
    }
  }, [deleteSuccess]);

  function updatePluginState(e: SyntheticEvent, command: string) {
    e.preventDefault();
    if (isExecuteLoading) return;
    const payload = { command } as PluginExecuteRequest;
    mutateExecute({ payload });
  }

  function UpdateLoading() {
    return <FiRefreshCw className={`text-s-warning animate-spin`} />;
  }
  const argsOn = {
    preIcon: <img src={logoURL} alt="plugin-logo" />,
    topAction: (
      <Button
        onClick={(e) => updatePluginState(e, PLUGIN_STOP)}
        variant="toggle"
      >
        on
      </Button>
    ),
    title: name,
    subtitle: author,
    subtitleIcon: massalabsNomination.includes(author) ? <Certificate /> : null,
    content: [
      updatable &&
        (isExecuteLoading ? (
          <UpdateLoading />
        ) : (
          <Button
            variant="icon"
            onClick={(e) => updatePluginState(e, PLUGIN_UPDATE)}
            disabled={isExecuteLoading}
          >
            <FiRefreshCw className="text-s-warning" />
          </Button>
        )),
      <Button variant="icon">
        <FiArrowUpRight onClick={() => window.open(home, '_blank')} />
      </Button>,
      <Button variant="icon" onClick={() => deletePlugin({})}>
        <FiTrash2 />
      </Button>,
    ],
  };

  const argsOff = {
    preIcon: <img src={logoURL} alt="plugin-logo" />,
    topAction: (
      // we use customClass because "disabled" doesn't let us click on the button to turn it back on
      <Button
        onClick={(e) => updatePluginState(e, PLUGIN_START)}
        customClass="bg-primary text-tertiary"
        variant="toggle"
      >
        off
      </Button>
    ),
    title: name,
    subtitle: author,
    subtitleIcon: massalabsNomination.includes(author) ? <Certificate /> : null,
    content: [
      <Button variant="icon" disabled>
        <FiArrowUpRight />
      </Button>,
      <Button variant="icon" onClick={() => deletePlugin({})}>
        <FiTrash2 />
      </Button>,
    ],
  };
  return <Plugin {...(status === PluginStatus.Up ? argsOn : argsOff)} />;
}

export default StationPlugin;
