import { useState } from 'react';
import axiosServices from '../services/axios';

export interface InstallProps {
    setErrorData: (errorType: string, errorMessage: string) => void;
    getPluginsInfo: () => void;
}

function InstallPlugin(p: InstallProps) {
    const [pluginUrl, setPluginUrl] = useState('');

    function handlePluginUrlChange(event: any) {
        setPluginUrl(event.target.value);
    }

    async function handleInstallPlugin(event: any) {
        event.preventDefault();

        try {
            await axiosServices.installPlugin(pluginUrl);
            p.setErrorData("success", "Plugin installed");

        } catch (error: any) {
            p.setErrorData("error", `Plugins installation failed: ${error.message}`);
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
                                Plugin URL:
                            </label>
                            <input
                                id="plugin-url"
                                name="plugin-url"
                                type="text"
                                className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
                                value={pluginUrl}
                                onChange={handlePluginUrlChange}
                            />
                        </div>
                        <div className="flex justify-end">
                            <button
                                type="submit"
                                className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline"
                            >
                                Install
                            </button>
                        </div>
                    </form>
                </div>
                </div>
        </section>
    );
}

export default InstallPlugin;
