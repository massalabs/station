// Used in web\plugin-manager\src\components\pluginBlock.tsx
export interface Plugin {
    id: string;
    name: string;
    description: string;
    version: string;
    home: string;
    // isUpdate: boolean;
    status: string;
    logo: string;
    isNotInstalled?: boolean;
    url?: string;
    isFake?: boolean;
}

export interface PluginProps {
    plugin: Plugin;
    errorHandler: (errorType: string, errorMessage: string) => void;
    getPluginsInfo: () => void;
}

export enum PluginStatus {
    Down = "Down",
    Up = "Up",
    Starting = "Starting",
    Stopping = "Stopping",
    Error = "Error",
}
export interface PluginHomePage {
    id:string;
    name: string;
    description: string;
    logo: string;
    status: string;
    home?: string;
}
export interface PluginNotInstalled {

    name: string;
    description: string;
    url: string;
    logo: string;
}

export interface PluginStoreItemRequest{
    name:string
    description:string
    version: string
    url:string
    assets:{
        windows: pluginStoreItemFile
        linux: pluginStoreItemFile
        macos_arm64: pluginStoreItemFile
        macos_amd64: pluginStoreItemFile
    }
}
export interface pluginStoreItemFile{
    url:string
    checksum:string
}