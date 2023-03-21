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
    handleOpenPlugin: (pluginName:string) => void;
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
    id: string;
    name: string;
    description: string;
    url: string;
    logo: string;
}