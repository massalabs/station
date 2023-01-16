import React, { useEffect } from "react";

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
                            src="https://s3.us-west-2.amazonaws.com/secure.notion-static.com/180e697b-a457-4ecf-9691-6192ba7eddec/MASSA_LIGHT.png?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Content-Sha256=UNSIGNED-PAYLOAD&X-Amz-Credential=AKIAT73L2G45EIPT3X45%2F20230113%2Fus-west-2%2Fs3%2Faws4_request&X-Amz-Date=20230113T135113Z&X-Amz-Expires=86400&X-Amz-Signature=3a067257d5f862d6c3e215fb22eb61c40b944b64b93abbc5a85b8f930ed35dc4&X-Amz-SignedHeaders=host&response-content-disposition=filename%3D%22MASSA_LIGHT.png%22&x-id=GetObject"
                        ></img>
                    </div>
                    <h1 className="text-xs text-center text-white">{plugin.name}</h1>
                </div>
            </button>
        );
    });

    return (
        <div className="p-3">
            <div className=" inline-flex">
                <img
                    className="max-w-12 max-h-12 self-center"
                    src="https://s3.us-west-2.amazonaws.com/secure.notion-static.com/180e697b-a457-4ecf-9691-6192ba7eddec/MASSA_LIGHT.png?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Content-Sha256=UNSIGNED-PAYLOAD&X-Amz-Credential=AKIAT73L2G45EIPT3X45%2F20230113%2Fus-west-2%2Fs3%2Faws4_request&X-Amz-Date=20230113T135113Z&X-Amz-Expires=86400&X-Amz-Signature=3a067257d5f862d6c3e215fb22eb61c40b944b64b93abbc5a85b8f930ed35dc4&X-Amz-SignedHeaders=host&response-content-disposition=filename%3D%22MASSA_LIGHT.png%22&x-id=GetObject"
                    alt="Thyra Logo"
                />
                <h1 className="text-2xl font-bold self-center text-white">Thyra</h1>
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
                    src="https://static.wixstatic.com/media/788ebc_5768179e35704336b41bce4f9cf246ea~mv2.jpg/v1/fill/w_560,h_522,al_c,q_80,usm_0.66_1.00_0.01,enc_auto/788ebc_5768179e35704336b41bce4f9cf246ea~mv2.jpg"
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
