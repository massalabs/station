const errorCodes = new Map([
	['Wallet-0001', 'A wallet with the same nickname already exists on your computer'],
	['Wallet-0002', 'Wrong wallet password'],
	['Wallet-0003', 'Error while retrieving the specified wallet'],
	['Wallet-1001', 'A wallet cannot be created without nickname'],
	['Wallet-1002', 'A wallet cannot be created without password'],
	['Wallet-1003', 'A technical issue occured while creating the wallet'],
	['Wallet-2001', 'Specify the wallet wallet address to delete'],
	['Wallet-2002', 'Error while trying to delete the wallet, may be opened by an other program'],
	['Wallet-3001', 'Error while trying to import the wallet file'],
	['Wallet-4001', 'Error while trying to fetch all your wallets'],

	['Domains-0001', 'Error while fetching DNS owned'],
	['Domains-0002', 'Error while fetching addresses of DNS'],

	['WebCreator-0001', 'Creation of Webcreator Smart Contract failed'],
	['WebCreator-1001', 'Reading of ZIP file failed'],
	['WebCreator-1002', 'Sending of ZIP file on blockchain failed'],
]);

function getErrorMessage(errorCode) {
	return errorCodes.get(errorCode);
}
