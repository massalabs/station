// Used in web\plugin-manager\src\components\pluginBlock.tsx
export interface Plugin {
    ID: number;
    name: string;
    description: string;
    version: string;
    url: string;
    // isUpdate: boolean;
    isOnline: boolean;
    logoPath: string;
}

export interface PluginProps {
    props: Plugin;
    setErrorData: (errorType: string, errorMessage: string) => void;
  }