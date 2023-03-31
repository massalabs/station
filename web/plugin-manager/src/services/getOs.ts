export const getOs = () => {
    function getOperatingSystem(): string | undefined {
        const userAgent = navigator.userAgent.toLowerCase();
        const platform = navigator.platform.toLowerCase();

        if (platform.includes("win")) {
            return "windows";
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
    return getOperatingSystem();
};
