import React, { MouseEventHandler, MouseEvent, useEffect, useMemo, useState } from 'react';
import thyraLogo from '../../assets/ThyraLogo-V0-Detailed.png';
import massaLogoLight from '../../assets/MASSA_LIGHT_Detailed.png';
import massaLogomark from '../../assets/massa_logomark_detailed.png';
import { useQuery } from 'react-query';
import gearingLogo from '../../assets/gearing.png';
import axios from 'axios';
import { PluginHomePage } from '../../../../shared/interfaces/IPlugin';
import { PluginCard } from '../../components/pluginCard';
import toggleTheme from '../../components/toggleTheme';
import Header from '../../components/Header';
import ArrowEntry from '../../assets/pictos/ArrowEntry.svg';
import { UIStore } from '../../store/UIStore';
import ManagePluginCard from '../../components/managePluginCard';
import grid1 from '../../assets/element/grid1.svg'
import wallet from '../../assets/logo/plugins/Wallet.svg'
import registry from '../../assets/logo/plugins/Registry.svg'
import webOnChain from '../../assets/logo/plugins/WebOnChain.svg'
import ArrowWhite6 from '../../assets/pictos/ArrowWhite6.svg'
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
        url = `/thyra/massa/${pluginName}`;
        break;
    }
    window.open(url, '_blank');
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
    //To delete when Api is merged.

    return res.data;
  };
  
  // Add the fake plugins
  useEffect(() => {
    getPlugins().then((res) => {
      res.forEach((element: PluginHomePage) => {
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

      <p className=" display flex-row flex justify-center text-font">
        <p className="text-brand">↳</p> Which plugin
      </p>

      {/* Display the plugins in a grid */}
      <div className="mt-24 gap-8 grid mx-auto w-fit rounded-lg grid-cols-4 place-items-center max-lg:grid-cols-3">
        {mapPluginList()}
        <>
          <ManagePluginCard />
        </>
      </div>
      {/* <img src={grid1} className="relative bg-scroll ">
      </img> */}
      </div>
  );
}

export default Home;
