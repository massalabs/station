import { useEffect, useState } from 'react';

import massaLogomark from '../../assets/massa_logomark_detailed.png';

import axios from 'axios';
import {
  PluginHomePage,
  PluginStatus,
} from '../../../../shared/interfaces/IPlugin';
import { PluginCard } from '../../components/pluginCard';

import Header from '../../components/Header';

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
      home: '/websiteUploader',
    },
    {
      name: 'Registry',
      description: 'Browse Massa blockchain and its .massa websites',
      id: '423',
      logo: registry,
      status: 'Up',
      home: '/search',
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

  // Contains the list of plugins populated by fakeplugins
  const [plugins, setPlugins] = useState<PluginHomePage[]>(fakePluginsList);
  // Keep track of the previous fetch
  const [previousFetch, setPreviousFetch] = useState<PluginHomePage[]>([]);
  const displayPlugins = async () => {
    await getPlugins()
      .then((res: PluginHomePage[]) => {
        let combinedPlugins: PluginHomePage[] = [...fakePluginsList, ...res];
        // If the list of plugins has changed, update the state
        if (JSON.stringify(previousFetch) !== JSON.stringify(res)) {
          setPlugins(combinedPlugins);
          setPreviousFetch(res);
        }
      })
      .catch((err) => {
        console.log(err);
        // If there is an error, use the previous list
        setPlugins((previousFetch) =>
          previousFetch.length > 0 ? previousFetch : fakePluginsList,
        );
      });
  };
  // fetch plugins every 10 seconds
  useEffect(() => {
    // fetch plugins
    const fetchData = async () => {
      await displayPlugins();
    };
    fetchData();
    const interval = setInterval(() => {
      displayPlugins();
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
                  name: plugin ? plugin.name : '',
                  description: plugin.description ? plugin.description : '',
                  status: plugin.status ? plugin.status : '',
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
  const GridStyle = `grid grid-flow-row mx-auto mt-3 gap-4 max-sm:grid-cols-1 sm:grid-cols-2 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4`;

  return (
    <div>
      <div
        className=" min-h-screen bg-img"
        style={{ backgroundImage: `url(${grid1})` }}
      >
        <Header />

        {/* Display the plugins in a grid */}
        <div className="mx-auto max-sm:w-[300px] sm:w-[640px] md:w-[768px] lg:w-[980px] max-xl:w-[1024px] xl:w-[1280px]">
          <MainTitle title="Which Plugins" />
          <div className={GridStyle}>{mapPluginList()}</div>
        </div>
      </div>
    </div>
  );
}

export default Home;
