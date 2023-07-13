import { useEffect, useState } from "react";
import axiosServices from "../services/axios";
import { PluginProps } from "../../../shared/interfaces/IPlugin";
import { isUp } from "../helpers/isUp";
import TogglePlugin from "./TogglePlugin";
import Arrow6 from "../assets/pictos/arrow6.svg";
import ArrowWhite6 from "../assets/pictos/ArrowWhite6.svg";
import PrimaryButton from "./buttons/PrimaryButton";
import SecondaryButton from "./buttons/SecondaryButton";
import { BarLoader } from "react-spinners";

function PluginBlock(props: PluginProps) {
    const [isPluginUp, setStatus] = useState(isUp(props.plugin.status));
    const [isInstalling, setisInstalling] = useState(false);
    useEffect(() => setStatus(isUp(props.plugin.status)), [props.plugin.status]);

    // fetch info from plugin to get fresh data on demand
    async function getpluginInfo(): Promise<string | undefined> {
        try {
            const res = await axiosServices.getpluginInfo(props.plugin.id);
            const status = res.data.status;
            setStatus(isUp(status));

            return status;
        } catch (error: any) {
            console.error(
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
            console.error("error", `Start plugin failed , error :${error.message}`);
        }

        getpluginInfo();
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
        if (isPluginUp) window.open(props.plugin.home);
        else {
            console.error("error", "Plugin is not running, launch it first");
        }
    }
    // Uninstall plugin
    async function removePlugins() {
        try {
            await axiosServices.deletePlugins(props.plugin.id);
            props.getPluginsInfo();
        } catch (error: any) {
            console.log(`Plugins failed to be removed , error ${error.message}`);
        }
    }
    // Download plugin
    async function downloadPlugins() {
        try {
            if (props.plugin.url === undefined) return console.error("Plugin url is undefined");
            setisInstalling(true);
            (await axiosServices.installPlugin(props.plugin.url)).status === (200 || 500) &&
                setisInstalling(false);
            props.getPluginsInfo();
        } catch (error: any) {
            setisInstalling(false);
            console.log(`Plugins failed to be downloaded , error ${error.message}`);
        }
    }

    return (
        <div
            className="flex flex-col justify-center items-start p-6 gap-4 w-72 h-56 
                    border-[1px] border-solid border-border rounded-2xl bg-bgCard "
        >
            {/* First block Display plugin name and description */}
            <div className="flex flex-row items-center justify-between w-full">
                <img src={props.plugin.logo} alt="Album" className="rounded-3xl w-10 h-10" />
                {!props.plugin.isFake && !props.plugin.isNotInstalled && (
                    <TogglePlugin handleChange={launchOrStop} checked={isPluginUp} />
                )}
            </div>
            <div className="w-full h-16">
                <h1 className={`label2 text-font h-8 minimize`}>{props.plugin.name}</h1>
                <p className="text3 text-font h-8 minimize max-sm:text-sm">
                    {props.plugin.description}
                </p>
            </div>
            <p className="hidden text3 text-font">V: {props.plugin.version ?? "0.0.0"}</p>
            {/* Second Block with Icons  */}
            <div className="flex w-full">
                {/* Delete hidden when version will be send through the API */}
                <div className="flex w-full content-between justify-between mx-auto gap-4">
                    {props.plugin.isNotInstalled ? (
                        !isInstalling ? (
                            <SecondaryButton
                                label={"Download"}
                                onClick={downloadPlugins}
                                style={" w-full"}
                            />
                        ) : (
                            <BarLoader width={"100%"} color="hsl(var(--twc-brand))" />
                        )
                    ) : (
                        <>
                            <PrimaryButton
                                label={"Open"}
                                onClick={openHomepagePlugins}
                                iconPathDark={Arrow6}
                                iconPathLight={ArrowWhite6}
                                isDisabled={isPluginUp ? false : true}
                            />

                            <SecondaryButton
                                label={"Delete"}
                                onClick={removePlugins}
                                isDisabled={props.plugin.isFake}
                            />
                        </>
                    )}
                </div>
            </div>
        </div>
    );
}

export default PluginBlock;
