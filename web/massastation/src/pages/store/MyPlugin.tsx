import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';

import {
  Button,
  Certificate,
  MassaWallet,
  Plugin,
} from '@massalabs/react-ui-kit';
import { FiArrowUpRight, FiRefreshCcw, FiTrash2 } from 'react-icons/fi';
import { IMassaPlugin } from './MyStation';
import {
  massalabsNomination,
  PLUGIN_START,
  PLUGIN_STOP,
  PLUGIN_UPDATE,
} from '../../utils/massaConstants';
import { usePost, useResource, useDelete } from '../../custom/api';

enum PluginStatus {
  Up = 'Up',
  Down = 'Down',
}

interface PluginPostMethod {
  command: string;
}

export function MyPlugin({
  plugin,
  fetchPlugins,
}: {
  plugin: IMassaPlugin;
  fetchPlugins: () => void;
}) {
  const navigate = useNavigate();
  const [myPlugin, setMyPlugin] = useState<IMassaPlugin>(plugin);
  const { author, name, home, logo, status, updatable, id } = myPlugin;
  const {
    data: newPlugin,
    refetch,
    isRefetching,
  } = useResource<IMassaPlugin>(`plugin-manager/${id}`);

  const { mutate, isSuccess, isLoading } = usePost<PluginPostMethod>(
    `plugin-manager/${id}/execute`,
  );

  const { mutate: deletePlugin, isSuccess: deleteSuccess } = useDelete(
    `plugin-manager/${id}`,
  );

  useEffect(() => {
    if (isSuccess && !isLoading) {
      refetch();
    }
  }, [isSuccess, isLoading, refetch]);

  useEffect(() => {
    if (!isRefetching && newPlugin) {
      console.log('newPlugin', newPlugin);
      setMyPlugin(newPlugin);
    }
  }, [isRefetching, newPlugin]);

  useEffect(() => {
    if (deleteSuccess) {
      fetchPlugins();
    }
  }, [deleteSuccess]);

  function changePluginState(command: string) {
    mutate({ command } as PluginPostMethod);
  }

  function handleDelete() {
    deletePlugin({});
  }

  const argsOn = {
    preIcon: massalabsNomination.includes(author) ? (
      <MassaWallet variant="rounded" />
    ) : (
      <img src={logo} alt="Plugin Logo" />
    ),
    topAction: (
      <Button onClick={() => changePluginState(PLUGIN_STOP)} variant="toggle">
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
            onClick={() => changePluginState(PLUGIN_UPDATE)}
          />
        </Button>
      ),
      <Button variant="icon">
        <FiArrowUpRight onClick={() => navigate(home)} />
      </Button>,
      <Button variant="icon" onClick={handleDelete}>
        <FiTrash2 />
      </Button>,
    ],
  };

  const argsOff = {
    preIcon: massalabsNomination.includes(author) ? (
      <MassaWallet variant="rounded" />
    ) : (
      <img src={logo} />
    ),
    topAction: (
      // we use customClass because "disabled" doesn't let us click on the button to turn it back on
      <Button
        onClick={() => changePluginState(PLUGIN_START)}
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
      <Button variant="icon" onClick={handleDelete}>
        <FiTrash2 />
      </Button>,
    ],
  };
  return <Plugin {...(status === PluginStatus.Up ? argsOn : argsOff)} />;
}
