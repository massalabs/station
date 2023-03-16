import React, { MouseEventHandler, MouseEvent} from 'react';
import thyraLogo from '../../assets/ThyraLogo-V0-Detailed.png';
import massaLogoLight from '../../assets/MASSA_LIGHT_Detailed.png';
import massaLogomark from '../../assets/massa_logomark_detailed.png';
import { useQuery } from 'react-query';
import gearingLogo from '../../assets/gearing.png';
import axios from 'axios';
import { PluginHomePage } from '../../../../shared/interfaces/IPlugin';
import { PluginCard } from '../../components/PluginCard';
import toggleTheme from '../../components/toggleTheme';
import Header from '../../components/Header';
import ArrowEntry from '../../assets/pictos/ArrowEntry.svg';
import registry from '../../assets/logo/plugins/Registry.svg';
import { UIStore } from '../../store/UIStore';
import ManagePluginCard from '../../components/managePluginCard';
/**
 * Homepage of Thyra with a list of plugins installed
 *
 */
type Props = {};

function Home(props: Props) {
  // Fetch plugins installed by calling get /plugin/manager

  const handleOpenPlugin = (event: MouseEvent<HTMLDivElement>)  => {
    // let url;
    // // Handle Fake plugins for now and only for massa plugins 
    // // TODO: Remove this when we have the API with authors of plugins 
    // switch (pluginName) {
    //     case 'Registry':
    //         url = '/thyra/registry';
    //         break;
    //     case 'Web On Chain':
    //         url = '/thyra/websiteCreator';
    //         break;
    //     case 'Wallet':
    //         url = '/thyra/wallet';
    //         break;
    //     default:
    //         // If it's not a special case we just redirect to the plugin's home caller
    //         url = `/thyra/massa/${pluginName}`;
    //         break;
    //     }
    //     window.open(url, '_blank');
    };

  // List of plugins
  let pluginList: JSX.Element[] = [<> Loading... </>];
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
  const { data, error, isError } = useQuery('plugins', getPlugins);
  if (isError) pluginList = [<> Error: {error} </>];
  // Store the result in plugins
  // Mocked till we have the API

  let plugins: PluginHomePage[] = [];

  if (data) plugins = data;

  plugins.push(
    {
      name: "Massa's Wallet",
      description: "Create and manage your smart wallets to buy, sell, transfer and exchange tokens",
      id: '420',
      logo: '',
      status: '',
    },
    {
      name: 'Web On Chain',
      description:
        'Buy your .massa domain and upload websites on the blockchain',
      id: '421',
      logo: '',
      status: '',
    },
    // {
    //     name: "Node Manager",
    //     description: "A plugin for managing your local node",
    //     id: "422",
    //     home: "/:4200",
    //     logo: "",
    //     status: "",
    // },
    {
      name: 'Registry',
      description:  "Browse Massa blockchain and its .massa websites",
      id: '423',
      logo: '',
      status: '',
    },
  );

  // Map over the plugins and display them in a list
  pluginList = plugins.map((plugin) => {
    return (
        <PluginCard
        
        {...{plugin:{
            id: plugin.id,
            logo: plugin.logo ? plugin.logo : massaLogomark,
            name: plugin ? plugin.name : "Plugin Problem",
            description:plugin.description ? plugin.description : "Plugin Problem",
            status: plugin.status ? plugin.status : "Plugin Problem",
        },
          handleOpenPlugin: handleOpenPlugin,
          key:plugin.id
        }} 
      />
    );
  });

  return (
    <div className="">
      <Header />

      <p className=" display flex-row flex justify-center">
        <p className="text-brand">↳</p> Which plugin
      </p>

      {/* Display the plugins in a grid */}
      <div className="m-4 grid mx-auto w-fit grid-cols-2 rounded-lg sm:grid-cols-4">
        {pluginList}
      <>
        <ManagePluginCard />
      </>
      </div>
    </div>
  );
}

export default Home;
