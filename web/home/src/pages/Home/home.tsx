import { useEffect, useState } from 'react';

import massaLogomark from '../../assets/massa_logomark_detailed.png';

import axios from 'axios';
import { PluginHomePage } from '../../../../shared/interfaces/IPlugin';
import { PluginCard } from '../../components/pluginCard';

import Header from '../../components/Header';

import ManagePluginCard from '../../components/managePluginCard';
import grid1 from '../../assets/element/grid1.svg';
import wallet from '../../assets/logo/plugins/Wallet.svg';
import registry from '../../assets/logo/plugins/Registry.svg';
import webOnChain from '../../assets/logo/plugins/WebOnChain.svg';
import MainTitle from '../../components/MainTitle';

/**
 * Homepage of Thyra with a list of plugins installed
 *
 */

function Home() {
  // Fetch plugins installed by calling get /plugin/manager

  const fakePluginsList: PluginHomePage[] = [
    {
      name: "Massa's Wallet",
      description:
        'Create and manage your smart wallets to buy, sell, transfer and exchange tokens',
      id: '420',
      logo: wallet,
      status: 'Up',
      home: '/thyra/wallet'
    },
    {
      name: 'Web On Chain',
      description:
        'Buy your .massa domain and upload websites on the blockchain',
      id: '421',
      logo: webOnChain,
      status: 'Up',
      home: '/thyra/websiteCreator'
    },
    {
      name: 'Registry',
      description: 'Browse Massa blockchain and its .massa websites',
      id: '423',
      logo: registry,
      status: 'Up',
      home: '/thyra/registry'
    },
  ];
  const [plugins, setPlugins] = useState<PluginHomePage[]>(fakePluginsList);

  const handleOpenPlugin = (pluginName: string) => {
    window.open(findPluginHome(pluginName), '_blank');
  };
  const findPluginHome = (pluginName: string) => {
    let home = '';
    plugins.forEach((element) => {
      if (element.name == pluginName && element.home) {
        return element.home;
      }
    });
    console.log('Link to the plugin not found')
    return home;
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
  useEffect(() => {
    // Fetch the plugins on first render
    let previousFetch: PluginHomePage[] = [];
    getPlugins().then((res: PluginHomePage[]) => {
      previousFetch = res;
      res.forEach((element: PluginHomePage) => {
        setPlugins((prev) => [...prev, element]);
      });
    });
    
    setInterval(() => {
    getPlugins().then((res : PluginHomePage[]) => {
      let combinedPlugins = [ ...fakePluginsList, ...res];
      // If the list of plugins has changed, update the state
      if (JSON.stringify(previousFetch) !== JSON.stringify(res)) {
        setPlugins(combinedPlugins);
        previousFetch = res;
      }
    });
      
    }, 10000);
  }, []);

  const mapPluginList = () => {
    return plugins.map((plugin) => {
      if (plugin.status == 'Starting' || plugin.status == 'Up'){
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
      }
      else{
        return (<></>)
      }
    })
      
  }

  return (
    <div
      className=" min-h-screen bg-img"
      style={{ backgroundImage: `url(${grid1})` }}
    >
      <Header />

      <MainTitle title="Which Plugins" />

      {/* Display the plugins in a grid */}
      <div className="mt-24 gap-8 grid mx-auto w-fit rounded-lg grid-cols-4 place-items-center max-lg:grid-cols-3">
        {mapPluginList()}
        <>
          <ManagePluginCard />
        </>
      </div>
    </div>
  );
}

export default Home;
