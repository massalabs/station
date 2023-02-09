import { useMemo, useState } from "react";
import { ArrowPathIcon, TrashIcon, PlayCircleIcon } from "@heroicons/react/24/outline";
import axiosServices from "../services/axios";
import { AxiosResponse } from "axios";

import { Plugin, PluginProps, PluginStatus } from "../interfaces/IPlugin";
import { statusHelper } from "../helpers/statusHelpers";

function PluginBlock(p: PluginProps) {
    // Callback to set error on parent
    function sendErrorData(errorType: string, errorMessage: string) {
        p.setErrorData(errorType, errorMessage);
    }
    function setData(data: Plugin) {
        dataMemoized = data;
        setStatus(statusHelper(data.status));
    }

    let dataMemoized = useMemo(() => {
        console.log("p.props useMemo", p.props);
        return p.props;
    }, [p.props]);

    //UseRef p.props
    const [toggleStatus, setStatus] = useState(p.props.status == "Up");

    // fetch info from plugin
    // Not implemented atm
    async function fetchPluginInfo(): Promise<AxiosResponse<Plugin>> {
        let result: AxiosResponse<Plugin> = {} as AxiosResponse<Plugin>;
        try {
            return (result = await axiosServices.getpluginInfo(dataMemoized.id));
        } catch (error) {
            sendErrorData(
                "error",
                `Plugins infos failed to get infos from plugin name : ${dataMemoized.name} on id: ${dataMemoized.id}`
            );
        }
        return result;
    }

    // Not implemented atm
    async function launchAndStopPlugins() {
        //Front end Update with fresh data
        let resultPluginInfo: AxiosResponse<Plugin>;
        // fetch info from plugin
        try {
            resultPluginInfo = await fetchPluginInfo();
        } catch (error) {
            sendErrorData("error", "Plugins infos failed to launch");
            return;
        }
        setData(resultPluginInfo.data);
        let result: AxiosResponse<number>;
        if (!toggleStatus) {
            // Launch plugin
            try {
                result = await axiosServices.manageLifePlugins(dataMemoized.id, "start");
                dataMemoized.status = PluginStatus.Up;
                setStatus(!toggleStatus);
                return result;
            } catch (error) {
                sendErrorData("error", "Start plugin failed");
                dataMemoized.status = PluginStatus.Error;
                return;
            }
        } else {
            // Stop plugin
            try {
                result = await axiosServices.manageLifePlugins(dataMemoized.id, "stop");
                dataMemoized.status = PluginStatus.Down;
                setStatus(!toggleStatus);
                return result;
            } catch (error) {
                sendErrorData("error", "Stop plugin failed");
                dataMemoized.status = PluginStatus.Error;
                return;
            }
        }
    }

    // Update plugin
    // Not implemented atm
    function updatePlugins() {
        //Front end update
        // propsRef.current.updateAvailable = !propsRef.current.updateAvailable

        //TODO : Uncoment this when we have a update process
        //#####################################################
        // const result = axiosServices.uploadPlugins("filename");
        // if ( pluginsInfos.status && (pluginsInfos.status <= 200 || pluginsInfos.status >= 300)){
        //     sendErrorData("error","Plugins infos failed to launch")
        // }
        // propsRef.current.online = true
        // return result;
        //#####################################################
        console.log("Update is Not implemented ATM");
    }
    // Open plugin homepage
    function openHomepagePlugins() {
        if (statusHelper(dataMemoized.status)) window.open(dataMemoized.home);
        else {
            sendErrorData("error", "Plugin is not running can't be launched , Launch it first");
        }
    }
    // Uninstall plugin
    function removePlugins() {
        try {
            axiosServices.deletePlugins(dataMemoized.id);
            sendErrorData("success", "Plugin removed");
        } catch (error) {
            sendErrorData("error", "Plugins failed to be removed");
        }
    }

    function minimizeString(str: string, length: number) {
        if (str.length > length) {
            return str.substring(0, length) + "...";
        } else {
            return str;
        }
    }

    function playStatus() {
        return statusHelper(dataMemoized.status)
            ? "w-6 h-6 text-green-500"
            : "w-6 h-6 text-red-500";
    }

    function updateStatus() {
        return "w-6 h-6 text-yellow-500";
        //Uncomment when update process is implemented
        //  return p.props.updateAvailable ? "w-6 h-6 text-yellow-500" : "w-6 h-6 text-green-500";
    }

    return (
        <section className="bg-slate-800 h-48 max-w-lg w-96 p-3 m-4 rounded-2xl">
            <div className=" flex-row h-full text-white ">
                <div className="flex">
                    <img className="w-10 h-10 pt-3 mx-2" src={p.props.logo} alt="Plugin Logo" />
                    <div className="w-full">
                        <h1 className="font-bold">{minimizeString(p.props.name, 90)}</h1>
                        <p className="font-light max-sm:text-sm">
                            {minimizeString(p.props.description, 100)}
                        </p>
                    </div>
                </div>
                <div className="flex w-full pt-7 justify-around items-center">
                    <p className="font-light">V: {p.props.version}</p>
                    <input
                        type="checkbox"
                        className="toggle toggle-success"
                        checked={toggleStatus}
                        onChange={launchAndStopPlugins}
                    />

                    <button>
                        <PlayCircleIcon className={playStatus()} onClick={openHomepagePlugins} />
                    </button>
                    {/* Delete hidden when update process is set */}
                    <button className="hidden">
                        <ArrowPathIcon className={updateStatus()} onClick={updatePlugins} />
                    </button>
                    <button>
                        <TrashIcon className="w-6 h-6 " onClick={removePlugins} />
                    </button>
                </div>
            </div>
        </section>
    );
}

export default PluginBlock;
