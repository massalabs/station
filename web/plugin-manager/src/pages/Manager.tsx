import { useEffect, useState } from "react";
import PluginBlock from "../components/pluginBlock";
import { Plugin } from "../../../shared/interfaces/IPlugin";
import massaLogoLight from "../assets/MASSA_LIGHT_Detailed.png";
import axiosServices from "../services/axios";
import alertHelper from "../helpers/alertHelpers";
import { PuffLoader } from "react-spinners";
import InstallPlugin from "../components/installPluginBlock";
import InstallNodeManager from "../components/installNodeManager";
function Manager() {
    //State to store error
    const [error, setError] = useState(<></>);

    const [plugins, setPlugins] = useState<Plugin[]>([]);

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

    return (
        <div>
            <div className="p-5 flex items-center">
                <img className="max-h-6" src={massaLogoLight} alt="Thyra Logo" />
                <h1 className="text-xl ml-6 font-bold text-font">Thyra</h1>
            </div>
            {/* FlexWrap is blocking align content in Plugin Block*/}
            {/* Good First Issue For Community : Rework Css Classname to align bottom line of icon on bottom of container
            Need to delete FlexWrap and rework the container */}
            <div className="flex flex-wrap mx-auto max-w-6xl justify-center content-center">
                {plugins?.length ? plugins.filter(p => !!p.name)
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
        </div>
    );
}

export default Manager;
