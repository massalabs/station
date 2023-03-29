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
function Manager() {
    const fakePluginsList: Plugin[] = [
        {
            name: "Massa's Wallet",
            description:
                "Create and manage your smart wallets to buy, sell, transfer and exchange tokens",
            id: "420",
            logo: wallet,
            status: PluginStatus.Up,
            version: "1.0.0",
            home: "/thyra/wallet",
            isFake: true,
        },
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

    const [plugins, setPlugins] = useState<Plugin[]>(fakePluginsList);
    const [pluginsNotInstalled, setPluginsNotInstalled] = useState<Plugin[]>([]);

    async function getPluginsInfo() {
        try {
            const pluginsInfos = await axiosServices.getPluginsInfo();
            let combinedPlugins = [...fakePluginsList, ...pluginsInfos.data];
            setPlugins(combinedPlugins);
            getPluginsInfoNotInstalled();
        } catch (error: any) {
            console.error("error", `Get plugins infos failed ,  error ${error.message} `);
        }
    }
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
            const pluginsInfos = await axiosServices.getNotInstalledPlugins();
            let combinedPlugins: Plugin[] = [];
            for (let index = 0; index < pluginsInfos.data.length; index++) {
                const element = pluginsInfos.data[index];

                combinedPlugins.push({
                    name: element.name,
                    description: element.description,
                    url: setDownloadLinkForPlatform(element.assets),
                    logo: notInstalled,
                    status: PluginStatus.NotInstalled,
                    id: 1000 + index.toString(),
                    home: "",
                    isNotInstalled: true,
                    isFake: false,
                });
            }
            setPluginsNotInstalled(combinedPlugins);
        } catch (error: any) {
            console.error("error", `Get plugins infos failed ,  error ${error.message} `);
        }
    }

    // Create a loop to fetch getPluginsInfo and update the status
    useEffect(() => {
        //Initialize Ui on first render
        getPluginsInfo();
        // Set interval to update plugin status periodically
        // let previousFetch: Plugin[] = [];
        const interval = setInterval(async () => {
            // Nice to have control on precedent fetch to avoid rerender if no change
            getPluginsInfo();

            // .then((res: Plugin[]) => {previousFetch = res});
        }, 5000);
        return () => clearInterval(interval);
    }, []);

    //To talk in code review is it better to doing stuff
    //like this or to repeat the filters mapping on bottom
    // const mapPluginList = () => {
    //     return plugins.map((plugin) => {
    //         return (
    //             <PluginBlock
    //                 {...{
    //                     plugin: {
    //                         id: plugin.id,
    //                         logo: plugin.logo ? plugin.logo : massaLogomark,
    //                         name: plugin ? plugin.name : "Plugin Problem",
    //                         description: plugin.description ? plugin.description : "Plugin Problem",
    //                         status: plugin.status ? plugin.status : PluginStatus.Down,
    //                         home: plugin.home ? plugin.home : "Plugin Problem",
    //                         version: plugin.version ? plugin.version : "Plugin Problem",
    //                     },
    //                     key: plugin.id,
    //                     errorHandler: errorHandler,
    //                     getPluginsInfo: getPluginsInfo,
    //                 }}
    //             />
    //         );
    //     });
    // };

    const setColsLength = (length: number) => {
        return length > 3 ? " grid-cols-4" : " grid-cols-3";
    };
    return (
        <div>
            <div className=" min-h-screen bg-img" style={{ backgroundImage: `url(${grid1})` }}>
                <Header />
                <MainTitle title="Plugin Manager" />
                <div className="w-[1307px] mx-auto">
                    <p className="Secondary mt-24 text-font ml-6">Installed</p>
                    <div
                        className={
                            "grid grid-flow-row-dense w-[1307px] mx-auto mt-3 gap-4 xs:max-xl:grid-cols-2 ml-6 xl: " +
                            setColsLength(plugins.length)
                        }
                    >
                        {plugins?.length ? (
                            plugins
                                .filter((p) => !!p.name)
                                .map((plugin) => (
                                    <PluginBlock plugin={plugin} getPluginsInfo={getPluginsInfo} />
                                ))
                        ) : (
                            <PuffLoader color="font" />
                        )}
                        <InstallPlugin plugins={plugins} getPluginsInfo={getPluginsInfo} />
                    </div>
                    <div className="divider mx-auto mt-8 w-2/3" />
                    <p className="Secondary mt-12 text-font ml-6">Not installed</p>
                    <div
                        className={
                            "grid grid-flow-row-dense w-[1307px] mx-auto my-3 gap-4 xs:max-xl:grid-cols-2 ml-6 xl: " +
                            setColsLength(pluginsNotInstalled.length)
                        }
                    >
                        {pluginsNotInstalled?.length ? (
                            pluginsNotInstalled
                                .filter((p) => !!p.name)
                                // sort plugins by names
                                .sort((a, b) => a.name.localeCompare(b.name))
                                .map((plugin) => (
                                    <PluginBlock plugin={plugin} getPluginsInfo={getPluginsInfo} />
                                ))
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
