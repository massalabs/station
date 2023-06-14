import { Button, Certificate, Plugin } from '@massalabs/react-ui-kit';
import { IMassaPlugin } from './MyStation';
import { FiArrowUpRight, FiRefreshCcw, FiTrash2 } from 'react-icons/fi';
import { useState } from 'react';
enum PluginStatus {
  Up = 'Up',
  Down = 'Down',
}
function MyPlugin({ plugin }: { plugin: IMassaPlugin }) {
  const [myPlugin, setMyPlugin] = useState<IMassaPlugin>(plugin);
  function handlePluginStateChange(method: string) {
    // TODO :
    // const {mutate } = usePost<any>("plugin-manager");
    // const ok = mutate({method:"${method}"})
    const newStatus = method === 'start' ? PluginStatus.Up : PluginStatus.Down;
    const updatedPlugin = { ...myPlugin, status: newStatus };
    setMyPlugin(updatedPlugin);
  }
  const argsOn = {
    preIcon: <img src={myPlugin.logo} />,
    topAction: (
      <Button onClick={() => handlePluginStateChange('stop')} variant="toggle">
        on
      </Button>
    ),
    title: `${myPlugin.name} `,
    subtitle: `${myPlugin.author}`,
    subtitleIcon: myPlugin.author === 'MassaLabs' ? <Certificate /> : <></>,
    content: [
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
    preIcon: <img src={myPlugin.logo} />,
    topAction: (
      <Button
        onClick={() => handlePluginStateChange('start')}
        customClass="bg-primary text-tertiary"
        variant="toggle"
      >
        off
      </Button>
    ),
    title: `${myPlugin.name}`,
    subtitle: `${myPlugin.author}`,
    subtitleIcon: myPlugin.author === 'MassaLabs' ? <Certificate /> : <></>,
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
