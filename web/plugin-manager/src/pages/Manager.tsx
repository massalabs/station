import { useEffect, useState } from "react";
import PluginBlock from "../components/pluginBlock";
import {
    Plugin,
    PluginNotInstalled,
    PluginStatus,
    PluginStoreAssets,
    PluginStoreItemRequest,
} from "../../../shared/interfaces/IPlugin";
import axiosServices from "../services/axios";
import alertHelper from "../helpers/alertHelpers";
import { PuffLoader } from "react-spinners";
import wallet from "../assets/logo/plugins/Wallet.svg";
import registry from "../assets/logo/plugins/Registry.svg";
import webOnChain from "../assets/logo/plugins/webOnChain.svg";
import notInstalled from "../assets/logo/plugins/notInstalledRed.png";
import massaLogomark from "../assets/massa_logomark_detailed.png";
import grid1 from "../assets/element/grid1.svg";
import InstallPlugin from "../components/installPluginBlock";
import Header from "../components/Header";
import MainTitle from "../components/MainTitle";
import { getOs } from "../services/getOs";
import {gridStyle , defineGridStyle} from "../../../shared/styles/grid";
function Manager() {
    const fakePluginsList: Plugin[] = [
        {
            name: "Web On Chain",
            description: "Buy your .massa domain and upload websites on the blockchain",
            id: "421",
            logo: webOnChain,
            status: PluginStatus.Up,
            version: "1.0.0",
            home: "/thyra/websiteCreator",
            isFake: true,
        },
        {
            name: "Registry",
            description: "Browse Massa blockchain and its .massa websites",
            id: "423",
            logo: registry,
            status: PluginStatus.Up,
            version: "1.0.0",
            home: "/thyra/registry",
            isFake: true,
        },
    ];
    const [plugins, setPlugins] = useState<Plugin[]>([]);
    const [pluginsNotInstalled, setPluginsNotInstalled] = useState<Plugin[]>([]);

    // This function get all the plugins info from the backend
    // The backend will return a list of plugins
    // The function will add the list of plugins to the plugins state
    // The function will call getPluginsInfoNotInstalled in order to get the plugins that are not installed
    async function getPluginsInfo() {
        try {
            const pluginsInfos = await axiosServices.getPluginsInfo();
            let combinedPlugins = [...pluginsInfos.data];
            setPlugins(combinedPlugins);
            getPluginsInfoNotInstalled();
        } catch (error: any) {
            console.error("error", `Get plugins infos failed ,  error ${error.message} `);
        }
    }
    // this function returns the url for the download link for the current platform

    const setDownloadLinkForPlatform = (pluginStoreItem: PluginStoreAssets) => {
        switch (getOs()) {
            case "windows":
                return pluginStoreItem.windows.url;
            case "macos arm64":
                return pluginStoreItem.macos_arm64.url;
            case "macos amd64":
                return pluginStoreItem.macos_amd64.url;
            case "linux":
                return pluginStoreItem.linux.url;
            default:
                break;
        }
    };
    async function getPluginsInfoNotInstalled() {
        try {
            // Get plugins that are not installed from the API
            const pluginsInfos = await axiosServices.getNotInstalledPlugins();
            // Create an empty array to store the plugins
            let combinedPlugins: Plugin[] = [];
            // Transform the data to the format we need
            combinedPlugins = pluginsInfos.data.map((element, index) => ({
                name: element.name,
                description: element.description,
                url: setDownloadLinkForPlatform(element.assets),
                logo: notInstalled,
                status: PluginStatus.NotInstalled,
                id: 1000 + index.toString(),
                home: "",
                isNotInstalled: true,
                isFake: false,
            }));
            // Store the plugins in the state
            setPluginsNotInstalled(combinedPlugins);
        } catch (error: any) {
            console.error("error", `Get plugins infos failed ,  error ${error.message} `);
        }
    }

    // Create a loop to fetch getPluginsInfo and update the status
    useEffect(() => {
        //Initialize Ui on first render
        const fetchData = async () => {
            await getPluginsInfo();
        };
        fetchData();
        // Set interval to update plugin status periodically
        const interval = setInterval(async () => {
            try {
                getPluginsInfo();
            } catch (error) {
                console.error("Error while updating plugins status", error);
            }
        }, 5000);
        return () => clearInterval(interval);
    }, []);

    const mapPluginList = (pluginsList: Plugin[], pluginListToCompare: Plugin[] = []) => {
        return pluginsList
            .filter((p) => !!p.name).filter((p) => !pluginListToCompare.find((p2) => p2.name === p.name))
            .sort((a, b) => a.name.localeCompare(b.name))
            .map((plugin) => {
                return <PluginBlock plugin={plugin} getPluginsInfo={getPluginsInfo} />;
            });
    };

    const defineGridStyle = (length: number) => {
        let styles = gridStyle;
        return styles += (length <= 3 ? " grid-cols-3 "  : setResponsiveGrid);
    };
    const gridStyle = " grid grid-flow-row mx-auto mt-3 gap-4 grid-cols-4"
    
    
    const setResponsiveGrid = " max-sm:grid-cols-1 sm:grid-cols-2 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 max-xl:grid-cols-4 "

    return (
        <div>
            <div className=" min-h-screen bg-img" style={{ backgroundImage: `url(${grid1})` }}>
                <Header />
                <div className="mx-auto 
                w-fit
                max-sm:w-[300px] sm:w-[640px] md:w-[768px] lg:w-[980px] max-xl:w-[1024px] xl:w-[1280px]"
                >
                <MainTitle title="Plugin Manager" />
                    <p className="Secondary mt-2 text-font ml-6">Installed</p>
                    <div
                        className={
                            
                            defineGridStyle(plugins.length+fakePluginsList.length)
                        }
                    >
                        
                            {mapPluginList(fakePluginsList)}
                            {mapPluginList(plugins)}
                        
                        <InstallPlugin plugins={plugins} getPluginsInfo={getPluginsInfo} />
                    </div>
                    <div className="divider mx-auto mt-8 w-3/4" />
                    <p className="Secondary mt-12 text-font ml-6">Not installed</p>
                    <div
                        className={
                            
                            defineGridStyle(pluginsNotInstalled.length) 
                        }
                    >
                        {pluginsNotInstalled?.length ? (
                            mapPluginList(pluginsNotInstalled)
                        ) : (
                            <PuffLoader />
                        )}
                    </div>
                </div>
            </div>
        </div>
    );
}

export default Manager;
