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
    const [errorPluginInstall, setErrorPluginInstall] = useState("")

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

    const verifyUrl = (url: string) => {
        const regex = new RegExp(
            '^(https?:\\/\\/)?' + // protocol
            '((([a-z\\d]([a-z\\d-]*[a-z\\d])*)\\.)+[a-z]{2,}|' + // domain name
            '((\\d{1,3}\\.){3}\\d{1,3}))' + // OR ip (v4) address
            '(\\?[;&a-z\\d%_.~+=-]*)?' + // query string
            '(\\#[-a-z\\d_]*)?$',
            'i'
        ); // fragment locator
        return regex.test(url);
    };

    const installPluginHandler = async (url: string) => {
        if (verifyUrl(url)) {
            setIsInstalling(true)
            await axiosServices.installPlugin(url);
            p.getPluginsInfo()
        } else {
            console.log("Invalid URL")
            setErrorPluginInstall('Invalid URL');
        }
    };
    
    

    return (
        <>
        <SmallCardExtended label2={'Install a plugin'} text3={'Manually install a plugin using .zip URL'} propsLabelButton={{
            callbackToParent: function (data: string): void {
                installPluginHandler(data);
            },
            label: 'Plugin link',
            placeholder: 'Set .zip Url Link',
            buttonValue: 'Install',
            error: errorPluginInstall
        }}/>
        </>)
    
    }
export default InstallPlugin;
