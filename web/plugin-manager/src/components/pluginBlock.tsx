import { useEffect, useState } from "react";
import { ArrowPathIcon, TrashIcon, PlayCircleIcon } from "@heroicons/react/24/outline";
import axiosServices from "../services/axios";
import { PluginProps } from "../../../shared/interfaces/IPlugin";
import { isUp } from "../helpers/isUp";

function PluginBlock(p: PluginProps) {
    if (!p) {
        return <div></div>
    }

    const [isPluginUp, setStatus] = useState(isUp(p.plugin.status));
    useEffect(() => setStatus(isUp(p.plugin.status)), [p.plugin.status]);

    // fetch info from plugin to get fresh data on demand
    async function getpluginInfo(): Promise<string | undefined> {
        try {
            const res = await axiosServices.getpluginInfo(p.plugin.id);
            setStatus(isUp(res.data.status))

            return res.data?.status
        } catch (error: any) {
            p.errorHandler(
                "error",
                `Plugins infos failed to get infos from 
                plugin name : ${p.plugin.name} on id: ${p.plugin.id}, error: ${error.message}}`
            );
        }
    }

    // Launch or stop plugin
    async function launchOrStop() {
        setStatus(!isPluginUp);

        // Launch plugin
        try {
            await axiosServices.manageLifePlugins(p.plugin.id, isPluginUp ? "stop" : "start");
        } catch (error: any) {
            p.errorHandler("error", `Start plugin failed , error :${error.message}`);
        }

        getpluginInfo()
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
        if (isPluginUp)
            window.open(p.plugin.home);
        else {
            p.errorHandler("error", "Plugin is not running, launch it first");
        }
    }
    // Uninstall plugin
    async function removePlugins() {
        try {
            await axiosServices.deletePlugins(p.plugin.id);
            p.getPluginsInfo();
            p.errorHandler("success", "Plugin removed");
        } catch (error: any) {
            p.errorHandler("error", `Plugins failed to be removed , error ${error.message}`);
        }
    }
    //Truncate the string so that it fits in the given lenght if needed.
    function minimize(str: string, length: number) {
        if (!str) {
            return ""
        }
        if (str.length > length) {
            return str.substring(0, length) + "...";
        } else {
            return str;
        }
    }

    // Change the play status icon color and update the status if we want to force it
    function setRunStatusColor() {
        if (isPluginUp) {
            return "w-6 h-6 text-green-500";
        } else {
            return "w-6 h-6 text-red-500";
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
                    <img className="w-10 h-10 pt-3 mx-2" src={p.plugin.logo} alt="Plugin Logo" />
                    <div className="w-full">
                        <h1 className="font-bold">{minimize(p.plugin.name, 90)}</h1>
                        <p className="font-light max-sm:text-sm">
                            {minimize(p.plugin.description, 100)}
                        </p>
                    </div>
                </div>
                {/* Second Block with Icons  */}
                <div className="flex w-full pt-7 justify-around items-center">
                    {/* Delete hidden when version will be send through the API */}
                    <p className="hidden font-light">V: {p.plugin.version ?? "0.0.0"}</p>
                    <input
                        type="checkbox"
                        className="toggle toggle-success"
                        checked={isPluginUp}
                        onChange={launchOrStop}
                    />
                    <button>
                        <PlayCircleIcon
                            className={setRunStatusColor()}
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
