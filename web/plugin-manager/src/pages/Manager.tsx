import { useEffect, useState } from "react";
import PluginBlock from "../components/pluginBlock";
import { Plugin, PluginProps } from "../../../shared/interfaces/IPlugin";
import massaLogoLight from "../assets/MASSA_LIGHT_Detailed.png";
import axiosServices from "../services/axios";
import { AxiosResponse } from "axios";
import alertHelper from "../helpers/alertHelpers";
import { PuffLoader } from "react-spinners";
function Manager() {
    let pluginsInfos: AxiosResponse<any, any> = {} as AxiosResponse<any, any>;

    //State to store error
    const [error, setError] = useState(<></>);
    //State to store plugins populated
    const [pluginsPopulated, setpluginsPopulated] = useState([<PuffLoader/>])
    //Callback to remove Error
    function removeError(): void {
        setError(<></>);
    }
    function setErrorHandler(errorType: string, errorMessage: string): void {
        setError(alertHelper(errorType, errorMessage, removeError));
        setInterval(() => {
            setError(<></>);
        }, 10000);
    }

    async function getPluginsInfo () {
        try {
            pluginsInfos = await axiosServices.getPluginsInfo();
            populatePlugins();
        } catch (error:any) {
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

    const mock: Plugin = {
        name: "Plugin 1",
        logo:massaLogoLight,
        description: "If you see this you probably have a problem with the plugin manager",
        version: "1.0.0",
        status: "Down",
        home: "/urlOfPlugin",
        // isUpdate: true,
        id: "1",
    };
    // Mocks in case we don't have the plugin manager
    let mocks = [mock];

    function populatePlugins () {
        if (pluginsInfos.status == 200) {
            setpluginsPopulated(pluginsInfos.data.map((mock: Plugin) => {
                let pluginProps: PluginProps = {
                    props: mock,
                    setErrorData: setErrorHandler,
                    triggerRefreshPluginList: function (): void {
                        getPluginsInfo();
                    }
                };
                return <PluginBlock {...pluginProps} />;
            }));
        } else {
            setpluginsPopulated (mocks.map((mock: Plugin) => {
                let pluginProps: PluginProps = {
                    props: mock,
                    setErrorData: setErrorHandler,
                    triggerRefreshPluginList: function (): void {
                        getPluginsInfo();
                    }
                };
                return <PluginBlock {...pluginProps} />;
            }));
        }
    }
    

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
                {pluginsPopulated}
                {error}
            </div>
        </>
    );
}

export default Manager;
