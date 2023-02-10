import { SetStateAction, useMemo, useState } from "react";
import { ArrowPathIcon, TrashIcon, PlayCircleIcon } from "@heroicons/react/24/outline";
import axiosServices from "../services/axios";
import { AxiosResponse } from "axios";
import { Plugin, PluginProps, PluginStatus } from "../../../shared/interfaces/IPlugin";
import { isStatusReady } from "../helpers/isStatusReady";

function PluginBlock(p: PluginProps) {
    // Callback to set error on parent
    function sendErrorData(errorType: string, errorMessage: string) {
        p.setErrorData(errorType, errorMessage);
    }
    // Data to display
    let dataMemoized = useMemo(() => {
        return p.props;
    }, [p.props]);

    // Toggle status state
    const [toggleStatus, setStatus] = useState(isStatusReady(p.props.status));

    const [playStatusClassName, setPlayStatusClassName] = useState(definePlayStatus());

    // fetch info from plugin to get fresh data on demand
    async function fetchPluginStatus(): Promise<AxiosResponse<string>> {
        let result: AxiosResponse<string> = {} as AxiosResponse<string>;
        try {
            return (result = await axiosServices.getpluginInfo(dataMemoized.id));
        } catch (error) {
            console.log(error);
            sendErrorData(
                "error",
                `Plugins infos failed to get infos from plugin name : ${dataMemoized.name} on id: ${dataMemoized.id}`
            );
        }
        return result;
    }

    // Launch or stop plugin
    async function launchAndStopPlugins() {
        let resultPluginInfo: AxiosResponse<string>;
        // fetch info from plugin to get fresh data
        resultPluginInfo = await fetchPluginStatus();
        // Update data
        setStatus(isStatusReady(resultPluginInfo.data));

        let result: AxiosResponse<number>;
        if (!toggleStatus) {
            // Launch plugin
            try {
                result = await axiosServices.manageLifePlugins(dataMemoized.id, "start");
                setStatus(!toggleStatus);
                forcePlayStatus(true)
                return result;
            } catch (error) {
                sendErrorData("error", "Start plugin failed");
            }
        } else {
            // Stop plugin
            try {
                result = await axiosServices.manageLifePlugins(dataMemoized.id, "stop");
                forcePlayStatus(false)
                setStatus(!toggleStatus);
                return result;
            } catch (error) {
                sendErrorData("error", "Stop plugin failed");
            }
        }
    }

    // Update plugin
    // Not implemented atm
    function updatePlugins() {
        //Front end update

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
        if (isStatusReady(dataMemoized.status)) window.open(dataMemoized.home);
        else {
            sendErrorData("error", "Plugin is not running can't be launched , launch it first");
        }
    }
    // Uninstall plugin
    function removePlugins() {
        try {
            axiosServices.deletePlugins(dataMemoized.id);
            p.triggerRefreshPluginList();
            sendErrorData("success", "Plugin removed");
        } catch (error) {
            sendErrorData("error", "Plugins failed to be removed");
        }
    }
    // Minimize string to fit in the block
    function minimize(str: string, length: number) {
        if (str.length > length) {
            return str.substring(0, length) + "...";
        } else {
            return str;
        }
    }

    // Change the play status icon color and update the status if we want to force it
    function definePlayStatus() {
            return isStatusReady(dataMemoized.status)
            ? "w-6 h-6 text-green-500"
            : "w-6 h-6 text-red-500";
    }

    function forcePlayStatus(status: boolean) {
        if (status) {
            setPlayStatusClassName("w-6 h-6 text-green-500");
        } else {
            setPlayStatusClassName("w-6 h-6 text-red-500");
        }
    }

    // Return the right icon for the update
    function defineUpdateStatus() {
        return "w-6 h-6 text-yellow-500";
        //Uncomment when update process is implemented
        //  return p.props.updateAvailable ? "w-6 h-6 text-yellow-500" : "w-6 h-6 text-green-500";
    }

    return (
        <section className="bg-slate-800 h-48 max-w-lg w-96 p-3 m-4 rounded-2xl">
            <div className=" flex-row h-full text-white ">
                {/* First block Display plugin name and description */}
                <div className="flex">
                    <img className="w-10 h-10 pt-3 mx-2" src={p.props.logo} alt="Plugin Logo" />
                    <div className="w-full">
                        <h1 className="font-bold">{minimize(p.props.name, 90)}</h1>
                        <p className="font-light max-sm:text-sm">
                            {minimize(p.props.description, 100)}
                        </p>
                    </div>
                </div>
                {/* Second Block with Icons  */}
                <div className="flex w-full pt-7 justify-around items-center">
                    <p className="hidden font-light">V: {p.props.version}</p>
                    <input
                        type="checkbox"
                        className="toggle toggle-success"
                        checked={toggleStatus}
                        onChange={launchAndStopPlugins}
                    />

                    <button>
                        <PlayCircleIcon
                            className={playStatusClassName}
                            onClick={openHomepagePlugins}
                        />
                    </button>
                    {/* Delete hidden when update process is set */}
                    <button className="hidden">
                        <ArrowPathIcon className={defineUpdateStatus()} onClick={updatePlugins} />
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
