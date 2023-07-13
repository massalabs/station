const errorCodes = new Map([
    ["Wallet-5001", "Please select a wallet to be able to see your domains"],
    ["Wallet-5002", "Please select a wallet to be able to perform that action"],

    ["Domains-0001", "Error while looking for your domain names"],
    ["Domains-0002", "Error while connecting your domain and smart contract's address"],

    ["WebCreator-0001", "Error while creating your website container"],
    ["WebCreator-1001", "Impossible to read you ZIP file. Try again"],
    ["WebCreator-1002", "The upload of your ZIP file failed. Try again"],
    ["WebCreator-1005", ".html doesn not exist in source of Website ZIP. Try again"],
    ["Unknown-0001", "An unknown error occured. Please try again"],
]);

function getErrorMessage(errorCode) {
    return errorCodes.get(errorCode) || errorCodes.get("Unknown-0001");
}

// If the error is from Thyra, we display the error to the user and log the details in the console.
// Otherwise, we simply display the details in the console.
function handleAPIError(error) {
    if (error.response && error.response.data) {
        if (error.response.data.code) {
            errorAlert(getErrorMessage(error.response.data.code));
        }
        console.error("Thyra error:", error.response.data);
    } else {
        errorAlert(getErrorMessage("Unknown-0001"));
        console.error(error);
    }
}
