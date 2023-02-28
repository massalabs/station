import { useEffect, useState } from 'react';
import axiosServices from '../services/axios';
import { Circles } from 'react-loader-spinner'

export interface InstallProps {
    errorHandler: (errorType: string, errorMessage: string) => void;
    getPluginsInfo: () => void;
}

function getOperatingSystem(): string | undefined {
    const userAgent = navigator.userAgent.toLowerCase();
    const platform = navigator.platform.toLowerCase();

    if (platform.includes("win")) {
        return "win";
    } else if (platform.includes("mac") || platform.includes("darwin")) {
        if (userAgent.includes("intel")) {
            return "macos amd64";
        } else if (userAgent.includes("apple")) {
            return "macos arm64";
        }
    } else if (platform.includes("linux")) {
        return "linux";
    }
}

function InstallNodeManager(p: InstallProps) {

    const [isInstalling, setIsInstalling] = useState<boolean>(false);

    const [os, setOs] = useState<string | undefined>(undefined);
    useEffect(() => setOs(getOperatingSystem()), []);

    async function handleInstallPlugin(_: any) {
        if (!os) {
            p.errorHandler("error", `Unable to detect operating system`);
            return
        }
        let url = "https://github.com/massalabs/thyra-node-manager-plugin/releases/latest/download/thyra-plugin-node-manager_";
        switch (os) {
            case "win": url =
                url = url + "windows-amd64.zip";
                break;
            case "macos amd64": url = url + "darwin-amd64.zip";
                break;
            case "macos arm64": url = url + "darwin-arm64.zip";
                break;
            case "linux": url = url + "linux-amd64.zip";
                break;
        }

        setIsInstalling(true)
        try {
            await axiosServices.installPlugin(url)
        } catch (err: any) {
            p.errorHandler("error", `Plugins installation failed: ${err.response?.data?.message}`);
        }
        p.errorHandler("success", "Plugin installation started");
        p.getPluginsInfo()
    }

    return (
        <section className="bg-slate-800 h-48 max-w-lg w-96 p-3 m-4 rounded-2xl">
            <div className="flex flex-col h-full text-white justify-center items-center">

                <div className="bg-slate-800 sm:px-6 sm:flex justify-center">
                    <h1 className="text-lg leading-6 font-bold" id="modal-headline">
                        Install Node Manager plugin
                    </h1>
                </div>
                <div className="flex w-full pt-2 justify-around items-center">
                    <button
                        type="submit"
                        hidden={isInstalling}
                        className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline"
                        onClick={handleInstallPlugin}
                    >
                        Install{os ? ` for ${os}` : ""}
                    </button>
                    {isInstalling ? "Installing..." : ""}
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

export default InstallNodeManager;
