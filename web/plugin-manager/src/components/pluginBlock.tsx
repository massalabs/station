import { useEffect, useState } from "react";
import { ArrowPathIcon, TrashIcon, PlayCircleIcon } from "@heroicons/react/24/outline";
import axiosServices from "../services/axios";
import { PluginProps } from "../../../shared/interfaces/IPlugin";
import { isUp } from "../helpers/isUp";
import TogglePlugin from "./TogglePlugin";
import plugin from "@tailwindcss/typography";

function PluginBlock(props: PluginProps) {

    // Donc on va faire un truc si c'est un fake plugin 
    // on le fait quand même passé mais on met un bool? pour capter que c'est un fake plugin coté front
    // par exemple button delete disabled si fake plugins pareil pour le toggle button 

    // Not installed on a NodeManager et Hello World

    // Installed on à 
    // 1. Wallet
    // 2. Registry
    // 3. Web on chain


    const [isPluginUp, setStatus] = useState(isUp(props.plugin.status));
    useEffect(() => setStatus(isUp(props.plugin.status)), [props.plugin.status]);

    const handleCardClick = () => {
        props.handleOpenPlugin(props.plugin.name);
      };
    // fetch info from plugin to get fresh data on demand
    async function getpluginInfo(): Promise<string | undefined> {
        try {
            const res = await axiosServices.getpluginInfo(props.plugin.id);
            const status = res.data.status
            setStatus(isUp(status))

            return status
        } catch (error: any) {
            props.errorHandler(
                "error",
                `Plugins infos failed to get infos from 
                plugin name : ${props.plugin.name} on id: ${props.plugin.id}, error: ${error.message}}`
            );
        }
    }

    // Launch or stop plugin
    async function launchOrStop() {
        setStatus(!isPluginUp);

        // Launch plugin
        try {
            await axiosServices.manageLifePlugins(props.plugin.id, isPluginUp ? "stop" : "start");
        } catch (error: any) {
            props.errorHandler("error", `Start plugin failed , error :${error.message}`);
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
            window.open(props.plugin.home);
        else {
            props.errorHandler("error", "Plugin is not running, launch it first");
        }
    }
    // Uninstall plugin
    async function removePlugins() {
        try {
            await axiosServices.deletePlugins(props.plugin.id);
            props.getPluginsInfo();
            props.errorHandler("success", "Plugin removed");
        } catch (error: any) {
            props.errorHandler("error", `Plugins failed to be removed , error ${error.message}`);
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
        //  return props.props.updateAvailable ? "w-6 h-6 text-yellow-500" : "w-6 h-6 text-green-500";
    }

    return (
            <div className="flex flex-col justify-center items-start p-5 gap-4 w-64 h-72 
                    border-[1px] border-solid border-border rounded-2xl bg-bgCard cursor-pointer hover:bg-hoverbgCard">
                {/* First block Display plugin name and description */}
                <div className="flex flex-row items-center justify-between w-full">
                <img
        src={props.plugin.logo}
        alt="Album"
        className="rounded-3xl w-10 h-10"
      />
                          <input
                        type="checkbox"
                        className="toggle toggle-success"
                        checked={isPluginUp}
                        onChange={launchOrStop}
                    />
                    <TogglePlugin handleChange={launchOrStop} checked={isPluginUp} />

                </div>
                    <div className="w-full">
                        <h1 className="font-bold">{minimize(props.plugin.name, 90)}</h1>
                        <p className="font-light max-sm:text-sm">
                            {minimize(props.plugin.description, 100)}
                        </p>
                    </div>
                {/* Second Block with Icons  */}
                <div className="flex w-full pt-7 justify-around items-center">
                    {/* Delete hidden when version will be send through the API */}
                    <p className="hidden font-light">V: {props.plugin.version ?? "0.0.0"}</p>

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
    );
}

export default PluginBlock;
