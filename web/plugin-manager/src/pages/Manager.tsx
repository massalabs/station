import { useEffect, useState } from "react";
import PluginBlock from "../components/pluginBlock";
import { Plugin } from "../../../shared/interfaces/IPlugin";
import massaLogoLight from "../assets/MASSA_LIGHT_Detailed.png";
import axiosServices from "../services/axios";
import alertHelper from "../helpers/alertHelpers";
import { PuffLoader } from "react-spinners";
import InstallPlugin from "../components/installPluginBlock";
function Manager() {
    //State to store error
    const [error, setError] = useState(<></>);

    const [plugins, setPlugins] = useState<Plugin[]>([]);

    //Callback to remove Error
    function removeError(): void {
        setError(<></>);
    }
    function setErrorHandler(errorType: string, errorMessage: string): void {
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
            setErrorHandler("error", `Get plugins infos failed ,  error ${error.message} `);
        }
    };


    // Update plugin status each 10 seconds
    // Create a loop to fetch getPluginsInfo and update the status
    useEffect(() => {
        //Initialize Ui on first render
        getPluginsInfo();
        // Set interval to update plugin status periodically
        const interval = setInterval(async () => {
            getPluginsInfo();
        }, 10000);
        return () => clearInterval(interval);
    }, []);

    return (
        <>
            <div className="p-5 flex items-center">
                <img className="max-h-6" src={massaLogoLight} alt="Thyra Logo" />
                <h1 className="text-xl ml-6 font-bold text-white">Thyra</h1>
            </div>
            {/* FlexWrap is blocking align content in Plugin Block*/}
            {/* Good First Issue For Community : Rework Css Classname to align bottom line of icon on bottom of container
            Need to delete FlexWrap and rework the container */}
            <div className="flex flex-wrap mx-auto max-w-6xl justify-center content-center">
                {plugins?.length ? plugins
                    // sort plugins by Id
                    .sort((a, b) => parseInt(a.id) > parseInt(b.id) ? 1 : 0)
                    .map(plugin => (
                        <PluginBlock
                            plugin={plugin}
                            setErrorData={setErrorHandler}
                            getPluginsInfo={getPluginsInfo}
                        />
                    ))
                    : <PuffLoader />
                }
                <InstallPlugin
                    setErrorData={setErrorHandler}
                    getPluginsInfo={getPluginsInfo}
                />
                {error}
            </div>
        </>
    );
}

export default Manager;
