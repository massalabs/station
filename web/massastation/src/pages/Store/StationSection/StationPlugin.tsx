import { useEffect, useState } from 'react';

import { Button, Certificate, Plugin } from '@massalabs/react-ui-kit';
import { FiArrowUpRight, FiRefreshCcw, FiTrash2 } from 'react-icons/fi';
import { IMassaPlugin } from './StationSection';
import {
  massalabsNomination,
  PLUGIN_START,
  PLUGIN_STOP,
  PLUGIN_UPDATE,
} from '../../../utils/massaConstants';

import { useDelete, usePost, useResource } from '../../../custom/api';

enum PluginStatus {
  Up = 'Up',
  Down = 'Down',
}

interface PluginPostMethod {
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
    isRefetching,
  } = useResource<IMassaPlugin>(`plugin-manager/${id}`);

  const { mutate, isSuccess } = usePost<PluginPostMethod>(
    `plugin-manager/${id}/execute`,
  );

  const { mutate: deletePlugin, isSuccess: deleteSuccess } = useDelete(
    `plugin-manager/${id}`,
  );

  const logoURL = `${baseAPI}/plugin-manager/${id}/logo`;

  useEffect(() => {
    if (isSuccess) {
      refetch();
    }
  }, [isSuccess]);

  useEffect(() => {
    if (!isRefetching && newPlugin) {
      setMyPlugin(newPlugin);
    }
  }, [isRefetching]);

  useEffect(() => {
    if (deleteSuccess) {
      fetchPlugins();
    }
  }, [deleteSuccess]);

  function updatePluginState(command: string) {
    const payload = { command };
    mutate({ payload });
  }
  const argsOn = {
    preIcon: <img src={logoURL} alt="plugin-logo" />,
    topAction: (
      <Button onClick={() => updatePluginState(PLUGIN_STOP)} variant="toggle">
        on
      </Button>
    ),
    title: name,
    subtitle: author,
    subtitleIcon: massalabsNomination.includes(author) ? <Certificate /> : null,
    content: [
      updatable && (
        <Button variant="icon">
          <FiRefreshCcw
            className="text-s-warning"
            onClick={() => updatePluginState(PLUGIN_UPDATE)}
          />
        </Button>
      ),
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
        onClick={() => updatePluginState(PLUGIN_START)}
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
