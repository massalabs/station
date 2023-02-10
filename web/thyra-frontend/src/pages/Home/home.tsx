import thyraLogo from "../../assets/ThyraLogo-V0-Detailed.png";
import massaLogoLight from "../../assets/MASSA_LIGHT_Detailed.png";
import massaLogomark from "../../assets/massa_logomark_detailed.png";
import { useQuery } from "react-query";
import gearingLogo from "../../assets/gearing.png";
import axios from "axios";
import { Plugin, PluginHomePage } from "../../../../shared/interfaces/IPlugin";
import { MouseEventHandler } from "react";

/**
 * Homepage of Thyra with a list of plugins installed
 *
 */
function Home() {
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
        const res = await axios.get(`/plugin-manager`, init);
        //To delete when Api is merged.
        console.log(res.data);
        return res.data;
    };
    const { isLoading, data, error, isError } = useQuery("plugins", getPlugins);
    if (isError) pluginList = [<> Error: {error} </>];
    // Store the result in plugins
    // Mocked till we have the API

    const openWindows = (url: string): undefined | MouseEventHandler<HTMLButtonElement> => {
        window.open(url);
        return undefined;
    }

    let plugins: PluginHomePage[] = [];

    if (data) plugins = data;

    plugins.push(
        {
            name: "Wallet",
            description: "Wallet plugin for managing your funds",
            id: "420",
            home: "/thyra/wallet",
            logo: "",
            status: "",
        },
        {
            name: "Web On Chain",
            description: "Web On Chain is a plugin for managing websites on the blockchain",
            id: "421",
            home: "/thyra/websiteCreator",
            logo: "",
            status: "",
        },
        // {
        //     name: "Node Manager",
        //     description: "A plugin for managing your local node",
        //     id: "422",
        //     home: "/:4200",
        //     logo: "",
        //     status: "",
        // },
        {
            name: "Registry",
            description: "Registry page for accessing websites on the blockchain",
            id: "423",
            home: "/thyra/registry",
            logo: "",
            status: "",
        }
    );

    // Map over the plugins and display them in a list
    pluginList = plugins.map((plugin) => {
        return (
            <button
                className="flex flex-wrap rounded-lg p-5 m-5"
                key={plugin.name}
                onClick={() => {window.open(plugin.home)}}
            >
                <div className="mx-auto">
                    <div className="tooltip" data-tip={plugin.description}>
                        <img className="w-9 h-9 self-center bg-slate-800" src={massaLogomark}></img>
                    </div>
                    <h1 className="text-xs text-center text-white">{plugin.name}</h1>
                </div>
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
            <div className="mx-auto cursor-pointer" onClick={() => {window.open("/thyra/plugin-manager")}}>
                    {/* Will change when manager page is done */}
                    <img
                        className="max-w-9 max-h-9 mx-auto block mb-2"
                        src={gearingLogo}
                        alt="Gearing Logo"
                    />
                    <p className="text-center text-m font text-white">Plugin Manager</p>
            </div>
        </div>
    );
}

export default Home;
