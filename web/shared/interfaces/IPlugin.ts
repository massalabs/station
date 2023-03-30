// Used in web\plugin-manager\src\components\pluginBlock.tsx
export interface Plugin {
    id: string;
    name: string;
    description: string;
    version?: string;
    home: string;
    // isUpdate: boolean;
    status: PluginStatus;
    logo: string;
    isNotInstalled?: boolean;
    url?: string;
    isFake?: boolean;
}

export interface PluginProps {
    plugin: Plugin;
    getPluginsInfo: () => void;
}

export enum PluginStatus {
    Down = "Down",
    Up = "Up",
    Starting = "Starting",
    Stopping = "Stopping",
    Error = "Error",
    NotInstalled = "NotInstalled",
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
    id: string;
    name: string;
    description: string;
    url: string;
    logo: string;
    status: PluginStatus;
}

export interface PluginStoreItemRequest{
    name:string
    description:string
    version: string
    url:string
    assets:PluginStoreAssets
}

export interface PluginStoreAssets{
    windows: pluginStoreItemFile
    linux: pluginStoreItemFile
    macos_arm64: pluginStoreItemFile
    macos_amd64: pluginStoreItemFile
}

export interface pluginStoreItemFile{
    url:string
    checksum:string
}