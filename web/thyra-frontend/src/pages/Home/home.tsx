import React, { useEffect } from "react";
import thyraLogo from "../../assets/ThyraLogo-V0-Detailed.png";
import massaLogoLight from "../../assets/MASSA_LIGHT_Detailed.png";
import massaLogomark from "../../assets/massa_logomark_detailed.png";
type Props = {};

/**
 * Homepage of Thyra with a list of plugins installed
 * 
 */
function home({}: Props) {
    // Fetch plugins installed by calling get /plugin/manager
    useEffect(() => {
        //plugins = fetch('http://localhost:8080/plugin/manager')
    }, []);

    // Store the result in plugins
    // Mocked till we have the API
    let plugins = [
        {
            name: "WebOnChain",
            description: "A plugin for managing your Massa node",
        },
        {
            name: "NodeManager",
            description: "A plugin for managing your Thyra node",
        },
        {
            name: "NodeManager",
            description: "A plugin for managing your Thyra node",
        },
        {
            name: "NodeManager",
            description: "A plugin for managing your Thyra node",
        },
        {
            name: "NodeManager",
            description: "A plugin for managing your Thyra node",
        },
        {
            name: "Wallet",
            description: "A plugin for managing your Thyra node",
        },
        {
            name: "Wallet",
            description: "A plugin for managing your Thyra node",
        },
    ];

    // Map over the plugins and display them in a list
    const pluginList = plugins.map((plugin) => {
        return (
            <button className="flex flex-wrap rounded-lg p-5 m-5" key={plugin.name}>
                <div className="mx-auto">
                    <div className="tooltip" data-tip={plugin.description}>
                        <img
                            className="w-9 h-9 self-center bg-slate-800"
                            src={massaLogomark}
                        ></img>
                    </div>
                    <h1 className="text-xs text-center text-white">{plugin.name}</h1>
                </div>
            </button>
        );
    });

    return (
        <div className="p-3">
            <div className="flex items-center ">
                <img
                    className="max-h-6"
                    src={massaLogoLight}
                    alt="Thyra Logo"
                />
                <h1 className="text-xl ml-6 font-bold text-white">Thyra</h1>
            </div>
            <div className="flex">
                <div className="">
                    <p className="text-xl p-5 text-white">
                        θύρα (thýra) in ancient Greek means door, entrance. To pronounce "thoo-rah".
                    </p>
                </div>
            </div>
            <div className="">
                <img
                    className="max-w-32 max-h-32 mx-auto block mb-10"
                    src={thyraLogo}
                    alt="Thyra Logo"
                />
                <p className="text-center text-4xl font text-white">
                    A gateway to Massa blockchain
                </p>
            </div>
            {/* Display the plugins in a grid */}
            <div className="m-4 grid mx-auto w-fit grid-cols-2 rounded-lg sm:grid-cols-4">
                {pluginList}
            </div>
        </div>
    );
}

export default home;
