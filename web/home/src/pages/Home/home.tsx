import thyraLogo from "../../assets/ThyraLogo-V0-Detailed.png";
import massaLogoLight from "../../assets/MASSA_LIGHT_Detailed.png";
import massaLogomark from "../../assets/massa_logomark_detailed.png";
import { useQuery } from "react-query";
import gearingLogo from "../../assets/gearing.png";
import axios from "axios";
import { PluginHomePage } from "../../../../shared/interfaces/IPlugin";
import {PluginCard} from "../../components/PluginCard";
import toggleTheme from "../../components/toggleTheme";
import Header from "../../components/Header";
import ArrowEntry from "../../assets/pictos/ArrowEntry.svg";

import { UIStore } from "../../store/UIStore";
/**
 * Homepage of Thyra with a list of plugins installed
 *
 */
type Props = {
}

function Home(props : Props) {

    
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

        return res.data;
    };
    const { data, error, isError } = useQuery("plugins", getPlugins);
    if (isError) pluginList = [<> Error: {error} </>];
    // Store the result in plugins
    // Mocked till we have the API

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
                    <h1 className="text-xs text-center text-font">{plugin.name}</h1>
                </div>
            </button>
        );
    });

    return (
        <div className="">
            <Header/>
            <div className="mx-auto flex-row">
                <p className="display block text-center flex-row">
                <img
                    className="w-16 h-16"
                    src={ArrowEntry}
                    alt="Thyra Logo"
                />
                    Which plugin
                </p>
            </div>
            {/* Display the plugins in a grid */}
            <div className="m-4 grid mx-auto w-fit grid-cols-2 rounded-lg sm:grid-cols-4">
                {pluginList}
            </div>
            {/* Plugin Manager */}
            <div className="mx-auto cursor-pointer" onClick={() => {window.open("/thyra/plugin-manager")}}>
                    {/* Will change when manager page is done */}
                    <img
                        className="max-w-9 max-h-9 mx-auto block mb-2"
                        src={gearingLogo}
                        alt="Gearing Logo"
                        />
                    <p className="text-center text-m font text-font">Plugin Manager</p>
            </div>
                        <PluginCard {...{logo: massaLogoLight, name: "Hello", description:"holaholaholahol aholaholaholahol aholaholahola holaholaholaholaholaholaholaholaholah olaholaholaholaholaholaholaholaholaholaholahola" }}/>

        </div>
    );
}

export default Home;
