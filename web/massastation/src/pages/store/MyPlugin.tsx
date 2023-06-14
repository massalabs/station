import { useState } from 'react';

import {
  Button,
  Certificate,
  MassaWallet,
  Plugin,
} from '@massalabs/react-ui-kit';
import { FiArrowUpRight, FiRefreshCcw, FiTrash2 } from 'react-icons/fi';
import { IMassaPlugin } from './MyStation';

enum PluginStatus {
  Up = 'Up',
  Down = 'Down',
}

function MyPlugin({ plugin }: { plugin: IMassaPlugin }) {
  const [myPlugin, setMyPlugin] = useState<IMassaPlugin>(plugin);
  function handlePluginState(method: string) {
    // TODO :
    // const {mutate } = usePost<any>("plugin-manager");
    // const ok = mutate({method:"${method}"})
    const newStatus = method === 'start' ? PluginStatus.Up : PluginStatus.Down;
    const updatedPlugin = { ...myPlugin, status: newStatus };
    setMyPlugin(updatedPlugin);
  }
  const argsOn = {
    preIcon:
      myPlugin.author === 'MassaLabs' ? (
        <MassaWallet variant="rounded" />
      ) : (
        <img src={myPlugin.logo} />
      ),
    topAction: (
      <Button onClick={() => handlePluginState('stop')} variant="toggle">
        on
      </Button>
    ),
    title: `${myPlugin.name} `,
    subtitle: `${myPlugin.author}`,
    subtitleIcon: myPlugin.author === 'MassaLabs' ? <Certificate /> : <></>,
    content: [
      // conditionally render the update icon based on the updatable property
      myPlugin.updatable && (
        <Button variant="icon" onClick={() => console.log('reload')}>
          <FiRefreshCcw className="text-s-warning" />
        </Button>
      ),
      <Button variant="icon" onClick={() => console.log('arrow')}>
        <FiArrowUpRight />
      </Button>,
      <Button variant="icon" onClick={() => console.log('trash')}>
        <FiTrash2 />
      </Button>,
    ],
  };

  const argsOff = {
    preIcon:
      myPlugin.author === 'MassaLabs' ? (
        <MassaWallet variant="rounded" />
      ) : (
        <img src={myPlugin.logo} />
      ),
    topAction: (
      // we use customClass because "disabled" doesn't let us click on the button to turn it back on
      <Button
        onClick={() => handlePluginState('start')}
        customClass="bg-primary text-tertiary"
        variant="toggle"
      >
        off
      </Button>
    ),
    title: `${myPlugin.name}`,
    subtitle: `${myPlugin.author}`,
    subtitleIcon: myPlugin.author === 'MassaLabs' ? <Certificate /> : null,
    content: [
      <Button variant="icon" disabled>
        <FiArrowUpRight />
      </Button>,
      <Button variant="icon">
        <FiTrash2 />
      </Button>,
    ],
  };
  return <Plugin {...(myPlugin.status === 'Up' ? argsOn : argsOff)} />;
}

export default MyPlugin;
