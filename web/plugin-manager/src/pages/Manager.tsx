import { useEffect, useState } from "react";
import PluginBlock from "../components/pluginBlock";
import { Plugin, PluginStatus } from "../../../shared/interfaces/IPlugin";
import axiosServices from "../services/axios";
import registry from "../assets/logo/plugins/Registry.svg";
import webOnChain from "../assets/logo/plugins/webOnChain.svg";
import notInstalled from "../assets/logo/plugins/notInstalledRed.png";
import grid1 from "../assets/element/grid1.svg";
import InstallPlugin from "../components/installPluginBlock";
import Header from "../components/Header";
import MainTitle from "../components/MainTitle";

function Manager() {
    const fakePluginsList: Plugin[] = [
        {
            name: "Web On Chain",
            description: "Buy your .massa domain and upload websites on the blockchain",
            id: "421",
            logo: webOnChain,
            status: PluginStatus.Up,
            version: "1.0.0",
            home: "/websiteUploader",
            isFake: true,
        },
        {
            name: "Registry",
            description: "Browse Massa blockchain and its .massa websites",
            id: "423",
            logo: registry,
            status: PluginStatus.Up,
            version: "1.0.0",
            home: "/search",
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
            getStorePlugins(combinedPlugins);
        } catch (error: any) {
            console.error("error", `Get plugins infos failed ,  error ${error.message} `);
        }
    }

    async function getStorePlugins(installedPlugins: Plugin[]) {
        try {
            // Get plugins that are not installed from the API
            const pluginsInfos = await axiosServices.getNotInstalledPlugins();
            // Transform the data to the format we need
            const combinedPlugins = pluginsInfos.data
                .map((element, index) => ({
                    name: element.name,
                    description: element.description,
                    url: element.file.url,
                    logo: notInstalled,
                    status: PluginStatus.NotInstalled,
                    id: 1000 + index.toString(),
                    home: "",
                    isNotInstalled: true,
                    isFake: false,
                }))
                // filter already installed based on the name
                .filter((p) => !installedPlugins.map((p) => p.name).includes(p.name));
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
        }, 3000);
        return () => clearInterval(interval);
    }, []);

    const filterSortList = (pluginsList: Plugin[]) => {
        return pluginsList.filter((p) => !!p.name).sort((a, b) => a.name.localeCompare(b.name));
    };

    const mapPluginList = (pluginsList: Plugin[]) => {
        return filterSortList(pluginsList).map((plugin) => {
            return <PluginBlock plugin={plugin} getPluginsInfo={getPluginsInfo} />;
        });
    };

    const GridStyle = `grid grid-flow-row mx-auto mt-3 gap-4 max-sm:grid-cols-1 sm:grid-cols-2 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4`;

    return (
        <div>
            <div className=" min-h-screen bg-img" style={{ backgroundImage: `url(${grid1})` }}>
                <Header />
                <div
                    className="mx-auto 
                max-sm:w-[300px] sm:w-[640px] md:w-[768px] lg:w-[980px] max-xl:w-[1024px] xl:w-[1280px]"
                >
                    <MainTitle title="Plugin Manager" />
                    <p className="Secondary mt-2 text-font ml-6">Installed</p>
                    <div className={GridStyle}>
                        {mapPluginList(fakePluginsList)}
                        {mapPluginList(plugins)}
                    </div>
                    <div className="divider mx-auto mt-8 w-3/4" />
                    <p className="Secondary mt-12 text-font ml-6">Plugin Store</p>
                    <div className={GridStyle}>
                        {mapPluginList(pluginsNotInstalled)}
                        <InstallPlugin plugins={plugins} getPluginsInfo={getPluginsInfo} />
                    </div>
                </div>
            </div>
        </div>
    );
}

export default Manager;
