import { useEffect, useState } from "react";
import PluginBlock from "../components/pluginBlock";
import { Plugin } from "../interfaces/IPlugin";
import massaLogoLight from "../assets/MASSA_LIGHT_Detailed.png";
import axiosServices from "../services/axios";
import { AxiosError, AxiosResponse } from "axios";
import alertHelper from "../helpers/alertHelpers";
function Manager() {
    
    let pluginsInfos: Element | AxiosError | AxiosResponse<any, any>;
    
    //State to store error
    const [error, setError] = useState(<></>)
    //Callback to remove Error
    function removeError() :void {
        setError(<></>);
    }
    
    // Update plugin status each 2 seconds
    // Create a loop to fetch getPluginsInfo and update the status
    useEffect(() => {
        const interval = setInterval(async () => {
            pluginsInfos = await axiosServices.getPluginsInfo();
            console.log(pluginsInfos)
            if ( pluginsInfos.status && (pluginsInfos.status <= 200 || pluginsInfos.status >= 300)){
                setError(alertHelper("error","Plugins infos failed to launch", removeError))        
            }
        }, 5000);
        return () => clearInterval(interval);
    }, []);


    const mock: Plugin = {
        name: "Plugin 1",
        logoPath: "https://upload.wikimedia.org/wikipedia/fr/thumb/1/15/Audi_logo.svg/1280px-Audi_logo.svg.png",
        description: "T Becarefull",
        version: "1.0.0",
        isOnline: false,
        url: "/urlOfPlugin",
        // isUpdate: true,
        ID: 1,
    };
    const mock2: Plugin = {
        name: "Plugin 2",
        logoPath: "https://upload.wikimedia.org/wikipedia/fr/thumb/1/15/Audi_logo.svg/1280px-Audi_logo.svg.png",
        description:
            "This is a plugin Description BecarefullThis is a plugin Description BecarefullThis is a plugin Description BecarefullThis is a plugin Description Becarefull",
        version: "1.0.0",
        isOnline: false,
        url: "/urlOfPlugin",
        // isUpdate: false,
        ID: 2,
    };
    const mock3: Plugin = {
        name: "Plugin 3",
        logoPath: "https://upload.wikimedia.org/wikipedia/fr/thumb/1/15/Audi_logo.svg/1280px-Audi_logo.svg.png",
        description:
            "This is a plugin Descriaaaaaaaaaaaaatio  aaaaaaaaaaa n  aaaaaaaaaaa BecarefullThis is a plugin Description BecarefullThis is a plugin Description BecarefullThis is a plugin Description Becarefull",
        version: "1.0.0",
        isOnline: true,
        // isUpdate: true,
        url: "/urlOfPlugin",
        ID: 3,
    };
    const mock4: Plugin = {
        name: "Plugin 4",
        logoPath: "https://upload.wikimedia.org/wikipedia/fr/thumb/1/15/Audi_logo.svg/1280px-Audi_logo.svg.png",
        description:
            "This is a plugin Deaaaaaaaaaacription BecarefullThis is a plugin Description BecarefullThis is a plugin Description BecarefullThis is a plugin Description Becarefull",
        version: "1.0.0",
        isOnline: true,
        // isUpdate: false,
        url: "/urlOfPlugin",
        ID: 4,
    };
    const mock5: Plugin = {
        name: "Plugin 5",
        logoPath: "https://upload.wikimedia.org/wikipedia/fr/thumb/1/15/Audi_logo.svg/1280px-Audi_logo.svg.png",
        description:
            "This is a plugin Description BessssssssssscarefullThis is a plugin Description BecarefullThis is a plugin Description BecarefullThis is a plugin Description Becarefull",
        version: "1.0.0",
        isOnline: true,
        // isUpdate: true,
        url: "/urlOfPlugin",
        ID: 5,
    };
    const mock6: Plugin = {
        name: "Plugin 5",
        logoPath: "https://upload.wikimedia.org/wikipedia/fr/thumb/1/15/Audi_logo.svg/1280px-Audi_logo.svg.png",
        description:
            "This is a plugin Description BecarefullThis is a plugin Description BecarefullThis is a plugin Description BecarefullThis is a plugin Description Becarefull",
        version: "1.0.0",
        isOnline: true,
        // isUpdate: true,
        url: "/urlOfPlugin",
        ID: 6,
    };
    let mocks = [mock, mock2, mock3, mock4, mock5, mock6];
    // pluginsInfos ? mocks : mocks = pluginsInfos.data;
    let plugins = mocks.map((mock) => <PluginBlock {...mock} />);
    return (
        <>
            <div className="p-5 flex items-center">
                <img className="max-h-6" src={massaLogoLight} alt="Thyra Logo" />
                <h1 className="text-xl ml-6 font-bold text-white">Thyra</h1>
            </div>
            {/* FlexWrap is blocking align content in Plugin Block*/}
            <div className="flex flex-wrap mx-auto max-w-6xl justify-center content-center">
                {plugins}
                {error}
            </div>
        </>
    );
}

export default Manager;

