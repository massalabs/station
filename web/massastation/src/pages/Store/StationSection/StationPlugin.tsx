import { SyntheticEvent, useEffect } from 'react';
import { useDelete, usePost } from '@/custom/api';

import {
  Button,
  Certificate,
  Plugin,
  PluginProps,
  Tag,
  Tooltip,
} from '@massalabs/react-ui-kit';
import {
  FiAlertCircle,
  FiArrowUpRight,
  FiRefreshCw,
  FiTrash2,
} from 'react-icons/fi';
import {
  massalabsNomination,
  NODE_MANAGER,
  PLUGIN_START,
  PLUGIN_STOP,
  PLUGIN_UPDATE,
  pluginList,
} from '@/const';
import Intl from '@/i18n/i18n';

import { MassaPluginModel, PluginStatus } from '@/models';

export interface PluginExecuteRequest {
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
  }, [isExecuteSuccess, refetch]);

  useEffect(() => {
    if (deleteSuccess) {
      refetch();
    }
  }, [deleteSuccess, refetch]);

  function updatePluginState(e: SyntheticEvent, command: string) {
    e.preventDefault();
    if (isExecuteLoading) return;
    const payload = { command } as PluginExecuteRequest;
    mutateExecute({ payload });
  }

  function UpdateLoading() {
    return <FiRefreshCw className={`text-s-warning animate-spin`} />;
  }

  let pluginArgs = {
    preIcon: <img src={logoURL} alt="plugin-logo" />,
    title: name,
    subtitle: author,
    tag:
      pluginList.includes(name) && name === NODE_MANAGER ? (
        <Tag type="warning">{Intl.t('store.beta')}</Tag>
      ) : null,
    subtitleIcon: massalabsNomination.includes(author) ? <Certificate /> : null,
    version: version,
  } as PluginProps;

  switch (status) {
    case PluginStatus.Down:
      pluginArgs = {
        ...pluginArgs,
        topAction: (
          <Button
            onClick={(e) => updatePluginState(e, PLUGIN_START)}
            customClass="bg-primary text-tertiary"
            variant="toggle"
          >
            off
          </Button>
        ),
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
      pluginArgs = {
        ...pluginArgs,
        topAction: (
          <Button
            onClick={(e) => updatePluginState(e, PLUGIN_START)}
            customClass="bg-primary text-tertiary"
            variant="toggle"
          >
            off
          </Button>
        ),
        content: [
          <Tooltip
            className="mas-tooltip"
            body={Intl.t('store.crashed-module')}
          >
            <FiAlertCircle className="text-s-error" />
          </Tooltip>,
          <Button variant="icon" disabled>
            <FiArrowUpRight />
          </Button>,
          <Button variant="icon" onClick={() => deletePlugin({})}>
            <FiTrash2 />
          </Button>,
        ],
      };
      break;
    default:
      pluginArgs = {
        ...pluginArgs,
        topAction: (
          <Button
            onClick={(e) => updatePluginState(e, PLUGIN_STOP)}
            variant="toggle"
          >
            on
          </Button>
        ),
        content: [
          updatable &&
            (isExecuteLoading ? (
              <UpdateLoading />
            ) : (
              <Button
                variant="icon"
                hoverText={Intl.t('store.update')}
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
      break;
  }

  return <Plugin {...pluginArgs} />;
}

export default StationPlugin;
