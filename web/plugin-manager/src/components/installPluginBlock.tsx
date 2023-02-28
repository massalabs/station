import { useEffect, useState } from 'react';
import axiosServices from '../services/axios';
import { Circles } from 'react-loader-spinner'
import { Plugin } from '../../../shared/interfaces/IPlugin';
import { InstallProps } from './installNodeManager';

export interface InstallZipProps extends InstallProps {
    plugins: Plugin[];
}

function InstallPlugin(p: InstallZipProps) {
    const [pluginUrl, setPluginUrl] = useState('');
    const [isInstalling, setIsInstalling] = useState<boolean>(false);

    let nbPluginsTmp = 0;
    useEffect(() => {
        // Wait for the plugin to be installed to clean things
        if(isInstalling && p.plugins.length > nbPluginsTmp) {
            setIsInstalling(false)
            setPluginUrl("");
        }
    }, [p.plugins]);

    function handlePluginUrlChange(event: any) {
        setPluginUrl(event.target.value);
    }

    async function handleInstallPlugin(event: any) {
        event.preventDefault();

        setIsInstalling(true)
        nbPluginsTmp = p.plugins.length;
        try {
            await axiosServices.installPlugin(pluginUrl);
            p.errorHandler("success", "Plugin installed");

        } catch (error: any) {
            p.errorHandler("error", `Plugins installation failed: ${error.response?.data?.message}`);
        }
        p.getPluginsInfo()
    }

    return (
        <section className="bg-slate-800 h-48 max-w-lg w-96 p-3 m-4 rounded-2xl">
            <div className=" flex-row h-full text-white ">

                <div className="bg-slate-800 sm:px-6 sm:flex justify-center">
                    <h1 className="text-lg leading-6 font-bold" id="modal-headline">
                        Install Plugin from zip
                    </h1>
                </div>
                <div className="flex w-full pt-2 justify-around items-center">
                    <form onSubmit={handleInstallPlugin}>
                        <div className="mb-4">
                            <label htmlFor="plugin-url" className="block font-bold mb-2">
                                {isInstalling ? "Installing..." : "Plugin URL:"}
                            </label>
                            <input
                                id="plugin-url"
                                name="plugin-url"
                                type="text"
                                className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
                                value={pluginUrl}
                                onChange={handlePluginUrlChange}
                                hidden={isInstalling}
                            />
                        </div>
                        <div className="flex justify-end">
                            <button
                                type="submit"
                                hidden={isInstalling}
                                className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline"
                            >
                                Install
                            </button>
                        </div>
                    </form>
                    <Circles
                        height="60"
                        width="60"
                        color="#dc091e"
                        ariaLabel="circles-loading"
                        wrapperStyle={{}}
                        wrapperClass=""
                        visible={isInstalling}
                    />
                </div>
                </div>
        </section>
    );
}

export default InstallPlugin;
