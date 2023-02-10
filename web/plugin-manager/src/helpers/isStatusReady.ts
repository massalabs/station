export const isStatusReady = (status: string) => {
    switch (status) {
        case "Up":
            return true;
        case "Down":
            return false;
        case "Starting":
            return true;
        case "Stopping":
            return false;
        case "Error":
            return false;
        default:
            return false;
    }
};