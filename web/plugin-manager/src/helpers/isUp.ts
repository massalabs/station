import { PluginStatus } from "../../../shared/interfaces/IPlugin";
export const isUp = (status: string) => {
    switch (status) {
        case PluginStatus.Up:
        case PluginStatus.Starting:
            return true;
        case PluginStatus.Down:
        case PluginStatus.Stopping:
        case PluginStatus.Error:
        case PluginStatus.NotInstalled:
        default:
            return false;
    }
};
