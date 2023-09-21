import { SyntheticEvent, useEffect } from 'react';
import { useDelete, usePost } from '@/custom/api';

import { Button, Certificate, Plugin, Tooltip } from '@massalabs/react-ui-kit';
import {
  FiAlertCircle,
  FiArrowUpRight,
  FiRefreshCw,
  FiTrash2,
} from 'react-icons/fi';
import {
  massalabsNomination,
  PLUGIN_START,
  PLUGIN_STOP,
  PLUGIN_UPDATE,
} from '@/const';

import { MassaPluginModel, PluginStatus } from '@/models';

interface PluginExecuteRequest {
  command: string;
}

const baseAPI = import.meta.env.VITE_BASE_API;

export function StationPlugin({
  plugin,
  refetch,
}: {
  plugin: MassaPluginModel;
  refetch: () => void;
}) {
  const { author, name, home, status, updatable, id, version } = plugin;

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
    if (deleteSuccess) {
      refetch();
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

  let args = {
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
    version: `v${version}`,
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

  switch (status) {
    case PluginStatus.Down:
      args = {
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
        subtitleIcon: massalabsNomination.includes(author) ? (
          <Certificate />
        ) : null,
        version: `v${version}`,
        content: [
          <Button variant="icon" disabled>
            <FiArrowUpRight />
          </Button>,
          <Button variant="icon" onClick={() => deletePlugin({})}>
            <FiTrash2 />
          </Button>,
        ],
      };
      break;
    case PluginStatus.Crashed:
      args = {
        preIcon: <img src={logoURL} alt="plugin-logo" />,
        topAction: (
          <Button
            disabled
            onClick={(e) => updatePluginState(e, PLUGIN_START)}
            customClass="bg-primary text-tertiary"
            variant="toggle"
          >
            off
          </Button>
        ),
        title: name,
        subtitle: author,
        subtitleIcon: massalabsNomination.includes(author) ? (
          <Certificate />
        ) : null,
        version: `v${version}`,
        content: [
          <Tooltip
            className="mas-tooltip"
            content="The plugin is not working, please uninstall it and install it again."
            icon={<FiAlertCircle className="text-s-error" />}
          />,
          <Button variant="icon" disabled>
            <FiArrowUpRight />
          </Button>,
          <Button variant="icon" onClick={() => deletePlugin({})}>
            <FiTrash2 />
          </Button>,
        ],
      };
      break;
  }

  return <Plugin {...args} />;
}

export default StationPlugin;
