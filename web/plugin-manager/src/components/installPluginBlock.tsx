import { useState } from "react";
import axiosServices from "../services/axios";

import { Plugin } from "../../../shared/interfaces/IPlugin";

import SmallCardExtended from "./SmallCardExtended";
export interface InstallZipProps extends InstallProps {
    plugins: Plugin[];
}

export interface InstallProps {
    getPluginsInfo: () => void;
}

function InstallPlugin(p: InstallZipProps) {
    const [errorPluginInstall, setErrorPluginInstall] = useState("");
    const [installIsPending, setInstallIsPending] = useState(false);
    const isValidUrl = (url: string) => {
        const regex = /^(http)[^\s]*\.zip$/i; // fragment locator
        return regex.test(url);
    };

    const installPluginHandler = async (url: string) => {
        if (!isValidUrl(url)) {
            setErrorPluginInstall("Invalid URL");
            return;
        }
        try {
            setInstallIsPending(true);
            await axiosServices.installPlugin(url);
            p.getPluginsInfo();
            setInstallIsPending(false);
        } catch (err: any) {
            console.error(err.response?.data?.message);
            setErrorPluginInstall("Error while installing plugin");
            setInstallIsPending(false);
        }
    };

    const debouncedInstallPluginHandler = debounce(installPluginHandler, 300);

    function debounce<T extends (...args: any[]) => any>(
        func: T,
        delay: number
    ): (...args: Parameters<T>) => void {
        let timeoutId: ReturnType<typeof setTimeout>;
        return function (this: any, ...args: Parameters<T>): void {
            clearTimeout(timeoutId);
            timeoutId = setTimeout(() => func.apply(this, args), delay);
        };
    }

    return (
        <>
            <SmallCardExtended
                label2={"Install a plugin"}
                text3={"Install a plugin using .zip URL"}
                propsLabelButton={{
                    callbackToParent: function (data: string): void {
                        debouncedInstallPluginHandler(data);
                    },
                    label: "Plugin link",
                    placeholder: "Set .zip Url Link",
                    buttonValue: "Install",
                    error: errorPluginInstall,
                    processIsPending: installIsPending,
                }}
            />
        </>
    );
}
export default InstallPlugin;
