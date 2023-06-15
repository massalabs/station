import { useState } from 'react';

import {
  Button,
  Certificate,
  MassaWallet,
  Plugin,
} from '@massalabs/react-ui-kit';
import { FiArrowUpRight, FiRefreshCcw, FiTrash2 } from 'react-icons/fi';
import { IMassaPlugin } from './StationSection';
import {
  massalabsNomination,
  PLUGIN_START,
  PLUGIN_STOP,
} from '../../../utils/massaConstants';

enum PluginStatus {
  // Up And Down are sent by the BE, we use the On/Off standard to stick to design
  On = 'Up',
  Off = 'Down',
}

function StationPlugin({ plugin }: { plugin: IMassaPlugin }) {
  const [myPlugin, setMyPlugin] = useState<IMassaPlugin>(plugin);
  let { author, name, logo, status, updatable } = myPlugin;

  function handlePluginState(method: string) {
    // TODO :
    // const {mutate } = usePost<any>("plugin-manager");
    // const ok = mutate({method:"${method}"})
    const newStatus =
      method === PLUGIN_START ? PluginStatus.On : PluginStatus.Off;
    const updatedPlugin = { ...myPlugin, status: newStatus };
    setMyPlugin(updatedPlugin);
  }
  const argsOn = {
    preIcon: massalabsNomination.includes(author) ? (
      <MassaWallet variant="rounded" />
    ) : (
      <img src={logo} />
    ),
    topAction: (
      <Button onClick={() => handlePluginState(PLUGIN_STOP)} variant="toggle">
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
        onClick={() => handlePluginState(PLUGIN_START)}
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

export default StationPlugin;
