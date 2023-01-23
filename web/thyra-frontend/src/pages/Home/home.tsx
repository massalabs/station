import React, { useEffect, useMemo } from "react";
import thyraLogo from "../../assets/ThyraLogo-V0-Detailed.png";
import massaLogoLight from "../../assets/MASSA_LIGHT_Detailed.png";
import massaLogomark from "../../assets/massa_logomark_detailed.png";
import { useQuery, useMutation, useQueryClient } from "react-query";
import gearingLogo from "../../assets/gearing.png";
import { Link, Route, Routes } from "react-router-dom";
import Manager from "../Plugin_Manager/manager";
type Props = {};

/**
 * Homepage of Thyra with a list of plugins installed
 *
 */
function Home({}: Props) {
    // Fetch plugins installed by calling get /plugin/manager

    // List of plugins
    let pluginList: JSX.Element[] = [<> Loading... </>];
    const getPlugins = async () => {
        const init = {
            method: "GET",
            headers: {
                "Content-Type": "application/json",
            },
        };
        const res = await fetch("https://my.massa/thyra/plugin-manager", init);
        //To delete when Api is merged.
        console.log(res.json())
        return res.json();
    };
    const { isLoading, data, error, isError } = useQuery("plugins", getPlugins);
    if (isError) pluginList = [<> Error: {error} </>];
    // Store the result in plugins
    // Mocked till we have the API
    let plugins = [
        {
            name: "WebOnChain1",
            description: "A plugin for managing your Massa node",
        },
        {
            name: "NodeManager2",
            description: "A plugin for managing your Thyra node",
        },
        {
            name: "NodeManager3",
            description: "A plugin for managing your Thyra node",
        },
        {
            name: "NodeManager4",
            description: "A plugin for managing your Thyra node",
        },
        {
            name: "NodeManager5",
            description: "A plugin for managing your Thyra node",
        },
        {
            name: "Wallet6",
            description: "A plugin for managing your Thyra node",
        },
        {
            name: "Wallet7",
            description: "A plugin for managing your Thyra node",
        },
    ];

    if (data) plugins = data;

    // Map over the plugins and display them in a list
    pluginList = plugins.map((plugin) => {
        return (
            <button className="flex flex-wrap rounded-lg p-5 m-5" key={plugin.name}>
                {/* Uncomment when url is ready */}
                {/* <a href={`https://localhost/${plugin.authorname}/${plugin.name}`} target="_blank"> */}
                <div className="mx-auto">
                    <div className="tooltip" data-tip={plugin.description}>
                        <img className="w-9 h-9 self-center bg-slate-800" src={massaLogomark}></img>
                    </div>
                    <h1 className="text-xs text-center text-white">{plugin.name}</h1>
                </div>
                {/* </a> */}
            </button>
        );
    });

    return (
        <div className="p-3">
            <div className="flex items-center ">
                <img className="max-h-6" src={massaLogoLight} alt="Thyra Logo" />
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
            <div className="mx-auto">
                <Link to="/thyra/manager">
                    <link className="" href="../Plugin_Manager/manager.tsx"/>
                        <img
                            className="max-w-9 max-h-9 mx-auto block mb-2"
                            src={gearingLogo}
                            alt="Gearing Logo"
                        />
                        <p className="text-center text-m font text-white">Plugin Manager</p>                
                </Link>
            </div>
        </div>
    );
}

export default Home;
