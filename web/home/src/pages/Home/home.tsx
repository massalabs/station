import { useEffect, useState } from 'react';

import massaLogomark from '../../assets/massa_logomark_detailed.png';

import axios from 'axios';
import {
  PluginHomePage,
  PluginStatus,
} from '../../../../shared/interfaces/IPlugin';
import { PluginCard } from '../../components/pluginCard';

import Header from '../../components/Header';

import ManagePluginCard from '../../components/managePluginCard';
import grid1 from '../../assets/element/grid1.svg';
import registry from '../../assets/logo/plugins/Registry.svg';
import webOnChain from '../../assets/logo/plugins/webOnChain.svg';
import MainTitle from '../../components/MainTitle';

/**
 * Homepage of Thyra with a list of plugins installed
 *
 */

function Home() {
  // Fetch plugins installed by calling get /plugin/manager

  const fakePluginsList: PluginHomePage[] = [
    {
      name: 'Web On Chain',
      description:
        'Buy your .massa domain and upload websites on the blockchain',
      id: '421',
      logo: webOnChain,
      status: 'Up',
      home: '/thyra/websiteCreator',
    },
    {
      name: 'Registry',
      description: 'Browse Massa blockchain and its .massa websites',
      id: '423',
      logo: registry,
      status: 'Up',
      home: '/thyra/registry',
    },
    {
      name: 'Registrvy',
      description: 'Browse Massa blockchain and its .massa websites',
      id: '423',
      logo: registry,
      status: 'Up',
      home: '/thyra/registry',
    },
    {
      name: 'Registrya',
      description: 'Browse Massa blockchain and its .massa websites',
      id: '423',
      logo: registry,
      status: 'Up',
      home: '/thyra/registry',
    },
  ];

  const handleOpenPlugin = (pluginName: string) => {
    window.open(findPluginHome(pluginName));
  };
  const findPluginHome = (pluginName: string) => {
    const plugin = plugins.find((element) => element.name === pluginName);
    if (plugin) {
      return plugin.home;
    } else {
      console.log('Link to the plugin not found');
      return 'error';
    }
  };

  // List of plugins
  const getPlugins = async () => {
    const init = {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    };
    const res = await axios.get(`/plugin-manager`, init);
    return res.data;
  };

  // Add the fake plugins

  const [plugins, setPlugins] = useState<PluginHomePage[]>(fakePluginsList);

  // fetch plugins every 10 seconds
  useEffect(() => {
    // used to check if plugin list has changed
    let previousFetch: PluginHomePage[] = [];
    // fetch plugins
    const interval = setInterval(() => {
      getPlugins()
        .then((res: PluginHomePage[]) => {
          let combinedPlugins: PluginHomePage[] = [...fakePluginsList, ...res];
          // If the list of plugins has changed, update the state
          if (JSON.stringify(previousFetch) !== JSON.stringify(res)) {
            setPlugins(combinedPlugins);
            previousFetch = res;
          }
        })
        .catch((err) => {
          console.log(err);
          // If there is an error, use the previous list
          setPlugins((previousFetch) =>
            previousFetch.length > 0 ? previousFetch : fakePluginsList,
          );
        });
    }, 10000);

    return () => clearInterval(interval);
  }, []);

  const mapPluginList = () => {
    return plugins
      .filter((p) => !!p.name)
      .sort((a, b) => a.name.localeCompare(b.name))
      .map((plugin) => {
        if (
          plugin.status == PluginStatus.Starting ||
          plugin.status == PluginStatus.Up
        ) {
          return (
            <PluginCard
              {...{
                plugin: {
                  id: plugin.id,
                  logo: plugin.logo ? plugin.logo : massaLogomark,
                  name: plugin ? plugin.name : 'Plugin Problem',
                  description: plugin.description
                    ? plugin.description
                    : 'Plugin Problem',
                  status: plugin.status ? plugin.status : 'Plugin Problem',
                },
                handleOpenPlugin: handleOpenPlugin,
                key: plugin.id,
              }}
            />
          );
        } else {
          return <></>;
        }
      });
  };

  const setColsLength = (length: number) => {
    const cols = length > 3 ? 'grid-cols-4 ' : ' grid-cols-3 ';
    return (
      'max-sm:grid-cols-1 sm:grid-cols-2 md:grid-cols-2 lg:grid-cols-3 xl:' +
      cols
    );
  };

  return (
    <div
      className=" min-h-screen bg-img"
      style={{ backgroundImage: `url(${grid1})` }}
    >
      <Header />

      <MainTitle title="Which Plugins" />

      {/* Display the plugins in a grid */}
      <div className="mx-auto max-sm:w-[300px] sm:w-[640px] md:w-[768px] lg:w-[980px] xl:w-[1280px]">
        <div
          className={
            'grid grid-flow-row-dense w-[1650px] grid-cols-4 mx-auto mt-6 gap-4 ' +
            setColsLength(plugins.length + fakePluginsList.length) +
            ' max-sm:w-[300px] sm:w-[640px] md:w-[720px] lg:w-[980px] max-xl:grid-cols-3 w-[1280px] xl:w-[1280px]'
          }
        >
          {mapPluginList()}
          <ManagePluginCard />
        </div>
      </div>
    </div>
  );
}

export default Home;
