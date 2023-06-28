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
export interface IMassaPlugin {
    id:string;
    name: string;
    author:string;
    description: string;
    logo: string;
    status: string;
    home: string;
    updatable?: boolean;
    version?: string;
}
export interface IMassaStore {
    name: string;
    author: string;
    description: string;
    version: string;
    url: string;
    logo: string;
    file: pluginStoreItemFile
    os: string;
}

export interface PluginStoreItemRequest{
    name: string
    description: string
    version: string
    url: string
    file: pluginStoreItemFile
    os: string
}

export interface pluginStoreItemFile{
    url: string
    checksum: string
}