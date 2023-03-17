import { useEffect, useState } from "react";
import PluginBlock from "../components/pluginBlock";
import { Plugin, PluginHomePage } from "../../../shared/interfaces/IPlugin";
import axiosServices from "../services/axios";
import alertHelper from "../helpers/alertHelpers";
import { PuffLoader } from "react-spinners";
import InstallPlugin from "../components/installPluginBlock";
import InstallNodeManager from "../components/installNodeManager";
import Header from "../components/Header";
import wallet from '../assets/logo/plugins/Wallet.svg'
import registry from '../assets/logo/plugins/Registry.svg'
import webOnChain from '../assets/logo/plugins/WebOnChain.svg'
import massaLogomark from '../assets/massa_logomark_detailed.png';
import MainTitle from '../components/MainTitle';
function Manager() {
    
    const fakePluginsList:Plugin[] = [
        {
          name: "Massa's Wallet",
          description:
          'Create and manage your smart wallets to buy, sell, transfer and exchange tokens',
          id: '420',
          logo: wallet,
          status: '',
          version: '1.0.0',
          home:'/thyra/wallet'
        },
        {
          name: 'Web On Chain',
          description:
            'Buy your .massa domain and upload websites on the blockchain',
          id: '421',          logo: webOnChain,
          status: '',
          version: '1.0.0',
          home:'/thyra/websiteCreator'
        },
        {
          name: 'Registry',
          description: 'Browse Massa blockchain and its .massa websites',
          id: '423',
          logo: registry,
          status: '',
          version: '1.0.0',
          home:'/thyra/registry'
        },
      ];
    //State to store error
    const [error, setError] = useState(<></>);

    const [plugins, setPlugins] = useState<Plugin[]>(fakePluginsList);

    //Callback to remove Error
    function removeError(): void {
        setError(<></>);
    }

    function errorHandler(errorType: string, errorMessage: string): void {
        setError(alertHelper(errorType, errorMessage, removeError));
        setInterval(() => {
            removeError();
        }, 10000);
    }

    async function getPluginsInfo() {
        try {
            const pluginsInfos = await axiosServices.getPluginsInfo();
            setPlugins(pluginsInfos.data);
        } catch (error: any) {
            errorHandler("error", `Get plugins infos failed ,  error ${error.message} `);
        }
    };

    // Create a loop to fetch getPluginsInfo and update the status
    useEffect(() => {
        //Initialize Ui on first render
        getPluginsInfo();
        // Set interval to update plugin status periodically
        const interval = setInterval(async () => {
            getPluginsInfo();
        }, 1000);
        return () => clearInterval(interval);
    }, []);


      interface PluginHome {
        name: string;
        home: string
      }
      const findPluginHome = (pluginName:string) => {
        let home = '';
        plugins.forEach((element) => {
          if (element.name == pluginName) {
            home = element.home;
          }
        });
        return home;
      };
      const handleOpenPlugin = (pluginName: string) => {  
        window.open(findPluginHome(pluginName), '_blank');
      };

    const mapPluginList = () => {
        return plugins.map((plugin) => {
          return (
            <PluginBlock
              {...{
                plugin: {
                  id: plugin.id,
                  logo: plugin.logo ? plugin.logo : massaLogomark,
                  name: plugin ? plugin.name : 'Plugin Problem',
                  description: plugin.description
                    ? plugin.description
                    : 'Plugin Problem',
                  status: plugin.status ? plugin.status : 'Plugin Problem',
                  home: plugin.home ? plugin.home : 'Plugin Problem',
                  version : plugin.version ? plugin.version : 'Plugin Problem',
                },
                handleOpenPlugin: handleOpenPlugin,
                key: plugin.id,
                errorHandler: errorHandler,
                getPluginsInfo: getPluginsInfo,
              }}
            />
          );
        });
      };


    return (
        <>
            <Header />
            <MainTitle title="Plugin Manager" />
            <div className="">
                {plugins?.length ? plugins.filter(p => !!p.name)
                    // sort plugins by names
                    .sort((a, b) => a.name.localeCompare(b.name))
                    .map(plugin => (
                        <PluginBlock
                            plugin={plugin}
                            errorHandler={errorHandler}
                            getPluginsInfo={getPluginsInfo}
                            handleOpenPlugin={handleOpenPlugin}
                        />
                    ))
                    : <PuffLoader />
                }
                <InstallPlugin
                    errorHandler={errorHandler}
                    plugins={plugins}
                    getPluginsInfo={getPluginsInfo}
                />
                {plugins?.some(p => p.name === "Node Manager") ?
                    "" :
                    <InstallNodeManager
                        errorHandler={errorHandler}
                        getPluginsInfo={getPluginsInfo}
                    />
                }
                {error}
            </div>
        </>
    );
}

export default Manager;
