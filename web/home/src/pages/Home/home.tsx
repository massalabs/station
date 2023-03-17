import {  useEffect, useState } from 'react';

import massaLogomark from '../../assets/massa_logomark_detailed.png';

import axios from 'axios';
import { PluginHomePage } from '../../../../shared/interfaces/IPlugin';
import { PluginCard } from '../../components/pluginCard';

import Header from '../../components/Header';

import ManagePluginCard from '../../components/managePluginCard';
import grid1 from '../../assets/element/grid1.svg'
import wallet from '../../assets/logo/plugins/Wallet.svg'
import registry from '../../assets/logo/plugins/Registry.svg'
import webOnChain from '../../assets/logo/plugins/WebOnChain.svg'
import MainTitle from '../../components/MainTitle';

/**
 * Homepage of Thyra with a list of plugins installed
 *
 */
type Props = {};

function Home(props: Props) {
  // Fetch plugins installed by calling get /plugin/manager
  
  const fakePluginsList:PluginHomePage[] = [
    {
      name: "Massa's Wallet",
      description:
      'Create and manage your smart wallets to buy, sell, transfer and exchange tokens',
      id: '420',
      logo: wallet,
      status: '',
    },
    {
      name: 'Web On Chain',
      description:
        'Buy your .massa domain and upload websites on the blockchain',
      id: '421',
      logo: webOnChain,
      status: '',
    },
    {
      name: 'Registry',
      description: 'Browse Massa blockchain and its .massa websites',
      id: '423',
      logo: registry,
      status: '',
    },
  ];
  const [plugins, setPlugins] = useState<PluginHomePage[]>(fakePluginsList);
  interface PluginHome {
    name: string;
    home: string
  }
  const [pluginsHomeName, setPluginsHomeName] = useState<PluginHome[]>([{name:'',home:''}]);
  const handleOpenPlugin = (pluginName: string) => {
    let url;
    // Handle Fake plugins for now and only for massa plugins
    // TODO: Remove this when we have the API with authors of plugins
    
    switch (pluginName) {
      case 'Registry':
        url = '/thyra/registry';
        break;
      case 'Web On Chain':
        url = '/thyra/websiteCreator';
        break;
      case "Massa's Wallet":
        url = '/thyra/wallet';
        break;
      default:
        // If it's not a special case we just redirect to the plugin's home caller
        url = findPluginHome(pluginName);
        break;
    }
    window.open(url, '_blank');
  };
  const findPluginHome = (pluginName:string) => {
    let home = '';
    pluginsHomeName.forEach((element) => {
      if (element.name == pluginName) {
        home = element.home;
      }
    });
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
    getPlugins().then((res) => {
      res.forEach((element: PluginHomePage) => {
        setPluginsHomeName((prev) => [...prev, {name:element.name,home:element.home || ''}]);
        setPlugins((prev) => [...prev, element]);
      });
    });
    
  }, []);

  const mapPluginList = () => {
    return plugins.map((plugin) => {
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
    });
  };

  return (
    <div className=" min-h-screen bg-img" style={{backgroundImage: `url(${grid1})`}} >

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
