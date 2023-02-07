import React, { useEffect, useMemo, useRef, useState } from "react";
import { ArrowPathIcon, TrashIcon, PlayCircleIcon } from "@heroicons/react/24/outline";
import axiosServices from "../services/axios";
import { AxiosError, AxiosResponse } from "axios";
import alertHelper from "../helpers/alertHelpers";
import {Plugin} from "../interfaces/IPlugin";


function PluginBlock(props: Plugin) {
    //State to store error
    const [error, setError] = useState(<></>)
    //Callback to remove Error
    function removeError() :void {
        setError(<></>);
    }
    // Each Rerender we update fetch plugins data
    useEffect(() => {
        fetchPluginInfo();
        console.log("props UseEffect", props);
    }, [props]);
    //UseRef props
    let propsRef = useRef(props);
    const checkboxRef = useRef<HTMLInputElement>(null);
    const [toggleStatus, setStatus] = useState(props.isOnline)

    function setData(data: Plugin) {
        propsRef.current = data;
    }
    // fetch info from plugin
    // Not implemented atm
    async function fetchPluginInfo() : Promise<boolean> {
        try {
            const resultPluginInfo = await axiosServices.getpluginInfo(props.ID);           
        } catch (error) {           
            setError(alertHelper("error","Plugins infos failed to launch", removeError))
            return false                   
        }
        return true;
    }

    

    // Not implemented atm
    async function launchAndStopPlugins() {
        //Front end Update
        setStatus(!toggleStatus);
        // fetch info from plugin
        let resultPluginInfo : AxiosResponse<Plugin>;
        try {
            resultPluginInfo = await axiosServices.getpluginInfo(props.ID);           
        } catch (error) {           
            setError(alertHelper("error","Plugins infos failed to launch", removeError))
            return                   
        }
            setData(resultPluginInfo.data);
        let result : AxiosResponse<any>;
        if (!propsRef.current.isOnline) {
            console.log(props);
            // Launch plugin
            
            try {
                result = await axiosServices.manageLifePlugins(props.ID, "start");                
            } catch (error) {
                setError(alertHelper("error","Plugins infos failed to launch", removeError))
                return
            }
            // return result and change frontend if result is ok
            propsRef.current.isOnline = true
        } else {
            // Stop plugin
            try {
                result = await axiosServices.manageLifePlugins(props.ID, "stop");   
            } catch (error) {
                setError(alertHelper("error","Plugins infos failed to launch", removeError))
                return
            }
            return result;
        }
    }

    // // restart plugin
    // // Not implemented atm
    // function restartPlugins() {
    //     const result = axiosServices.manageLifePlugins(props.id, "restart");
    //         if (typeof result == typeof AxiosError){
    //             console.log("Error:", result);
    //             return                 
    //         }
    //         propsRef.current.online = true
    //         return result;
    // }
    // Update plugin
    // Not implemented atm
    function updatePlugins() {
        //Front end update
        // propsRef.current.updateAvailable = !propsRef.current.updateAvailable

        //TODO : Uncoment this when we have a update process
        //#####################################################
        // const result = axiosServices.uploadPlugins("filename");
        // if ( pluginsInfos.status && (pluginsInfos.status <= 200 || pluginsInfos.status >= 300)){
        //     setError(alertHelper("error","Plugins infos failed to launch", removeError))        
        // }
        // propsRef.current.online = true
        // return result;
        //#####################################################
        console.log("Update is Not implemented ATM")
    }
    // Open plugin homepage
    function openHomepagePlugins() {
        // TODO: Uncoment this when we have a url
        // window.open(propsRef.current.Url);
        console.log("OpenHomepage is Not implemented ATM")
    }
        // Uninstall plugin
        //Not implemented atm
        function removePlugins() {
            //TODO : Uncoment this when we have a remove process
            // const result = axiosServices.deletePlugins(props.id);
            // if ( pluginsInfos.status && (pluginsInfos.status <= 200 || pluginsInfos.status >= 300)){
            //     setError(alertHelper("error","Plugins infos failed to launch", removeError))        
            // }
            // }
            // return result;
            console.log("Remove is Not implemented ATM")
        }

    function minimizeString(str: string, length: number) {
        if (str.length > length) {
            return str.substring(0, length) + "...";
        } else {
            return str;
        }
    }

    function playStatus() {
        return props.isOnline ? "w-6 h-6 text-green-500" : "w-6 h-6 text-red-500";
    }

    function updateStatus() {
        return "w-6 h-6 text-yellow-500"
        //  return props.updateAvailable ? "w-6 h-6 text-yellow-500" : "w-6 h-6 text-green-500";
    }

    return (
        <section className="bg-slate-800 h-48 max-w-lg w-96 p-3 m-4 rounded-2xl">
            <div className=" flex-row h-full text-white ">
                <div className="flex">
                    <img className="w-10 h-10 pt-3 mx-2" src={props.logoPath} alt="Plugin Logo" />
                    <div className="w-full">
                        <h1 className="font-bold">{minimizeString(props.name, 90)}</h1>
                        <p className="font-light max-sm:text-sm">
                            {minimizeString(props.description, 100)}
                        </p>
                    </div>
                </div>
                    <div className="flex w-full pt-7 justify-around items-center">
                        <p className="font-light">V: {props.version}</p>
                        <input
                            type="checkbox"
                            className="toggle toggle-success"
                            checked={toggleStatus}
                            ref={checkboxRef}
                            onChange={launchAndStopPlugins}
                        />

                        <button>
                            <PlayCircleIcon className={playStatus()} onClick={openHomepagePlugins} />
                        </button>
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
