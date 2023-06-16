import { useEffect, useState } from 'react';

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
} from '../../utils/massaConstants';
import { usePost } from '../../custom/api';

enum PluginStatus {
  // Up And Down are sent by the BE, we use the On/Off standard to stick to design
  On = 'Up',
  Off = 'Down',
}
interface PluginPostMethod {
  method: string;
}
interface IDataReturn {
  plugin: IMassaPlugin;
}

export function MyPlugin({ plugin }: { plugin: IMassaPlugin }) {
  const [myPlugin, setMyPlugin] = useState<IMassaPlugin>(plugin);
  let { author, name, logo, status, updatable, id } = myPlugin;

  const { mutate, data, isSuccess, isLoading } = usePost<
    PluginPostMethod,
    IDataReturn
  >(`plugin-manager/${id}`);

  useEffect(() => {
    if (!isLoading) {
      if (isSuccess && data) {
        console.log(data);
        setMyPlugin(data.plugin);
      }
    }
  }),
  [isSuccess, data];

  function changePluginStatus(method: string) {
    mutate({ method: method });
  }

  const argsOn = {
    preIcon: massalabsNomination.includes(author) ? (
      <MassaWallet variant="rounded" />
    ) : (
      <img src={logo} />
    ),
    topAction: (
      <Button onClick={() => changePluginStatus(PLUGIN_STOP)} variant="toggle">
        on
      </Button>
    ),
    title: name,
    subtitle: author,
    subtitleIcon: massalabsNomination.includes(author) ? (
      <Certificate />
    ) : (
      <></>
    ),
    content: [
      updatable && (
        <Button variant="icon">
          <FiRefreshCcw className="text-s-warning" />
        </Button>
      ),
      <Button variant="icon">
        <FiArrowUpRight />
      </Button>,
      <Button variant="icon">
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
        onClick={() => changePluginStatus(PLUGIN_START)}
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
      <Button variant="icon">
        <FiTrash2 />
      </Button>,
    ],
  };
  return <Plugin {...(status === PluginStatus.On ? argsOn : argsOff)} />;
}
