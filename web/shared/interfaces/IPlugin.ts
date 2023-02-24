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
}

export interface PluginProps {
    plugin: Plugin;
    setErrorData: (errorType: string, errorMessage: string) => void;
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
    home: string;
    description: string;
    logo: string;
    status: string;
}