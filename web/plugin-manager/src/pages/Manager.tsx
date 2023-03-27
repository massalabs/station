import { useEffect, useState } from "react";
import PluginBlock from "../components/pluginBlock";
import { Plugin, PluginNotInstalled } from "../../../shared/interfaces/IPlugin";
import axiosServices from "../services/axios";
import alertHelper from "../helpers/alertHelpers";
import { PuffLoader } from "react-spinners";
import wallet from '../assets/logo/plugins/Wallet.svg'
import registry from '../assets/logo/plugins/Registry.svg'
import webOnChain from '../assets/logo/plugins/WebOnChain.svg'
import notInstalled from '../assets/logo/plugins/notInstalledRed.png';
import massaLogomark from '../assets/massa_logomark_detailed.png';
import grid1 from '../assets/element/grid1.svg';
import InstallPlugin from "../components/installPluginBlock";
import Header from "../components/Header";
import MainTitle from '../components/MainTitle';
function Manager() {
    
    const fakePluginsList:Plugin[] = [
        {
          name: "Massa's Wallet",
          description:
          'Create and manage your smart wallets to buy, sell, transfer and exchange tokens',
          id: '420',
          logo: wallet,
          status: 'Up',
          version: '1.0.0',
          home:'/thyra/wallet',
          isFake: true
        },
        {
          name: 'Web On Chain',
          description:
            'Buy your .massa domain and upload websites on the blockchain',
          id: '421',          
          logo: webOnChain,
          status: 'Up',
          version: '1.0.0',
          home:'/thyra/websiteCreator',
          isFake: true
        },
        {
          name: 'Registry',
          description: 'Browse Massa blockchain and its .massa websites',
          id: '423',
          logo: registry,
          status: 'Up',
          version: '1.0.0',
          home:'/thyra/registry',
          isFake: true
        },
      ]
    //State to store error
    const [error, setError] = useState(<></>);

    const [plugins, setPlugins] = useState<Plugin[]>(fakePluginsList);
    const [pluginsNotInstalled, setPluginsNotInstalled] = useState<PluginNotInstalled[]>([]);

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
            let combinedPlugins = [ ...fakePluginsList, ...pluginsInfos.data];
            setPlugins(combinedPlugins);
            getPluginsInfoNotInstalled();
        } catch (error: any) {
            errorHandler("error", `Get plugins infos failed ,  error ${error.message} `);
        }
    };

    async function getPluginsInfoNotInstalled() {
      try {
          const pluginsInfos = await axiosServices.getNotInstalledPlugins();
          let combinedPlugins : PluginNotInstalled[] = [];
          pluginsInfos.data.forEach(element => {
            combinedPlugins.push({
              name: element.name,
              description: element.description,
              url: element.url,
              logo: notInstalled
            })
          });
          setPluginsNotInstalled(combinedPlugins);

      } catch (error: any) {
          errorHandler("error", `Get plugins infos failed ,  error ${error.message} `);
      }
  };

    // Create a loop to fetch getPluginsInfo and update the status
    useEffect(() => {
        //Initialize Ui on first render
        getPluginsInfo();
        // Set interval to update plugin status periodically
        // let previousFetch: Plugin[] = [];
        const interval = setInterval(async () => {
          // Nice to have control on precedent fetch to avoid rerender if no change
            getPluginsInfo()
            
            // .then((res: Plugin[]) => {previousFetch = res});
        }, 5000);
        return () => clearInterval(interval);
    }, []);

    //To talk in code review is it better to doing stuff 
    //like this or to repeat the filters mapping on bottom
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
                key: plugin.id,
                errorHandler: errorHandler,
                getPluginsInfo: getPluginsInfo,
              }}
            />
          );
        });
      };

    const setColsLength = (length:number) => {
        return length > 3 ? " grid-cols-4" : " grid-cols-3";
      };
    return (
        <div>
              <div
      className=" min-h-screen bg-img"
      style={{ backgroundImage: `url(${grid1})` }}
    >
            <Header />
            <MainTitle title="Plugin Manager" />
            <div className="w-[1307px] mx-auto">

            <p className="Secondary mt-24 text-font">Installed</p>
            <div className= {"grid grid-flow-row-dense w-[1307px] mx-auto mt-3 gap-4 " + setColsLength(plugins.length)}>
                {plugins?.length ? plugins.filter(p => !!p.name)
                    .map(plugin => (
                        <PluginBlock
                            plugin={plugin}
                            errorHandler={errorHandler}
                            getPluginsInfo={getPluginsInfo}
                            />
                    ))
                    : <PuffLoader color="font" />
                }
                                <InstallPlugin
                          errorHandler={errorHandler}
                          plugins={plugins}
                          getPluginsInfo={getPluginsInfo}
                      />
            </div>
            <div className="divider mx-auto mt-8 w-2/3"/>
            <p className="Secondary mt-12 text-font">Not installed</p> 
            <div className= {"grid grid-flow-row-dense w-[1307px] mx-auto my-3" + setColsLength(pluginsNotInstalled.length)}>
                {pluginsNotInstalled?.length ? plugins.filter(p => !!p.name)
                    // sort plugins by names
                    .sort((a, b) => a.name.localeCompare(b.name))
                    .map(plugin => (
                        <PluginBlock
                            plugin={plugin}
                            errorHandler={errorHandler}
                            getPluginsInfo={getPluginsInfo}
                            />
                    ))
                    : <PuffLoader />
                  }
            </div>
                  </div>
        </div>
        </div>
    );
}

export default Manager;
