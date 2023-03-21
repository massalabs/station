import { useEffect, useState } from 'react';
import axiosServices from '../services/axios';
import { Circles } from 'react-loader-spinner'
import { Plugin } from '../../../shared/interfaces/IPlugin';
import { InstallProps } from './installNodeManager';
import Arrow2 from '../assets/pictos/Arrow2.svg';
import SmallCardExtended from './SmallCardExtended';
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
        <>
        <SmallCardExtended label2={'Install a plugin'} text3={'Manually install a plugin using .zip URL'} propsLabelButton={{
            callbackToParent: function (data: string): void {
                setIsInstalling(true)
            },
            label: 'Plugin link',
            placeholder: 'Set .zip Url Link',
            buttonValue: 'Install',
            axiosCall: function (data: string): void {
                axiosServices.installPlugin(data);
            }
        }}/>
        </>)
        // <section className="h-28 w-80 max-w-lg p-6 gap-3 border-[1px] border-solid border-border rounded-2xl bg-bgCard">
            /* <div className='flex flex-row'>
            <div className='flex flex-col gap-3'>
                <p className='label2 text-font'>
                    Install a plugin
                </p>
                <p className='text3 text-font'>
                    Install a plugin using .zip URL
                </p>
            </div>
            <div className='flex self-center mx-auto'>
                <img className='w-8 h-4 hover:animate-rotate180Smooth' src={Arrow2} alt="" />
                {/* contains the icon to grow the container */}
            {/* </div> */}
            {/* </div> */}
            {/* <div className=" flex-row h-full text-font ">
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
                                className="bg-blue-500 hover:bg-blue-700 text-font font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline"
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
                </div> */}
        // </section>


export default InstallPlugin;
