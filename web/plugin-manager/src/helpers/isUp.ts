export const isUp = (status: string) => {
    switch (status) {
        case "Up":
        case "Starting":
            return true;
        case "Down":
        case "Stopping":
        case "Error":
        default:
            return false;
    }
};
