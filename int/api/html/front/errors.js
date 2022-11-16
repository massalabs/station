const errorCodes = new Map([
    ["Wallet-0001", "That nickname is taken. Try Another"],
    ["Wallet-0002", "Wrong password. Try again"],
    ["Wallet-0003", "Error while retrieving that wallet. Try again"],
    ["Wallet-0004", "Wallet not canceled"],
    ["Wallet-0005", "You just canceled your action. Try again"],    
    ["Wallet-0006", "Password is Empty. Try again"],

    ["Wallet-1001", "Enter a wallet nickname"],
    ["Wallet-1002", "Enter a wallet password"],
    ["Wallet-1003", "Error while creating your wallet. Try again"],
    ["Wallet-2001", "Select a wallet to delete"],
    [
        "Wallet-2002",
        "Error while deleting that wallet. Close all programs using this wallet and try again",
    ],
    ["Wallet-3001", "Error while importing this wallet. Try again"],
    [
        "Wallet-4001",
        "Error while connecting all your wallets. Reconnect all your wallets and try again",
    ],
    [
        "Wallet-4002",
        "Error while fetching your balance. Reconnect all your wallets and try again",
    ],

    ["Wallet-5001", "Please select a wallet to be able to see your domains"],
    ["Wallet-5002", "Please select a wallet to be able to perform that action"],

    ["Domains-0001", "Error while looking for your domain names"],
    ["Domains-0002", "Error while connecting your domain and smart contract's address"],

    ["WebCreator-0001", "Error while creating your website container"],
    ["WebCreator-1001", "Impossible to read you ZIP file. Try again"],
    ["WebCreator-1002", "The upload of your ZIP file failed. Try again"],
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
