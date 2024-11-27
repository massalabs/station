package onchain

import (
	"fmt"
	"strings"

	"github.com/massalabs/station/pkg/convert"
)

/**
------------------------------------------------------------------------------------------------------------------------
TYPES
------------------------------------------------------------------------------------------------------------------------
*/
type DatastoreContract struct {
	Data  []byte
	Args  []byte
	Coins uint64
}

type JSONableSliceMap []byte


type DatastoreSCEntry struct {
	Entry DatastoreData `json:"entry"`
	Bytes DatastoreData `json:"bytes"`
}

type DatastoreData struct {
	Entry JSONableSliceMap `json:"entry"`
	Bytes JSONableSliceMap `json:"bytes"`
}

func (u JSONableSliceMap) MarshalJSON() ([]byte, error) {
	var result string
	if u == nil {
		result = "null"
	} else {
		result = strings.Join(strings.Fields(fmt.Sprintf("%d", u)), ",")
	}

	return []byte(result), nil
}

/**
------------------------------------------------------------------------------------------------------------------------
UTILITY FUNCTIONS 
------------------------------------------------------------------------------------------------------------------------
*/

/**
* Used to create a new datastore entry object.
*/
func NewDeployScDatastoreData(entry []byte, bytes []byte) DatastoreData {
	return DatastoreData{
		Entry: entry,
		Bytes: bytes,
	}
}

/**
use to create a map of each datastore entry
*/
func NewDeployScDatastoreEntry(entry DatastoreData, bytes DatastoreData) DatastoreSCEntry {
	return DatastoreSCEntry{
		Entry: entry,
		Bytes: bytes,
	}
}

/**
------------------------------------------------------------------------------------------------------------------------
KEY GENERATION FUNCTIONS
------------------------------------------------------------------------------------------------------------------------
*/

/**
 * Generates a key for coin data in the datastore.
 *
 * @param offset - The offset to use when generating the key.
 * @returns A Uint8Array representing the key.
 */
func coinsKey(offset int) []byte {
	return convert.U64ToBytes(offset+1)
}

/**
 * Generates a key for args data in the datastore.
 *
 * @param offset - The offset to use when generating the key.
 * @returns A Uint8Array representing the key.
 */
func argsKey(offset int) []byte {
	return convert.U64ToBytes(offset+1)
}

/**
 * Generates a key for contract data in the datastore.
 *
 * @param offset - The offset to use when generating the key.
 * @returns A Uint8Array representing the key.
 */
func contractKey(offset int) []byte {
	return convert.U64ToBytes(offset+1)
}


/**
------------------------------------------------------------------------------------------------------------------------
POPULATE DATASTORE FUNCTION
------------------------------------------------------------------------------------------------------------------------

--- Datastore Format ---
	[[key],[value]] 

	=== 

	[
		[
			[KEY_BYTE_ARRAY_LENGTH], [KEY_BYTE_ARRAY_DATA]]
		], 
		[
			[VALUE_BYTE_ARRAY_LENGTH], [VALUE_BYTE_ARRAY_DATA]
		]
	]

*/

/**
 * Populates the datastore with the contracts.
 *
 * @remarks
 * This function is to be used in conjunction with the deployer smart contract.
 * The deployer smart contract expects to have an execution datastore in a specific state.
 * This function populates the datastore according to that expectation.
 *
 * @param contracts - The contracts to populate the datastore with.
 *
 * @returns The populated datastore.
 */
func populateDatastore(contract DatastoreContract) ([]DatastoreSCEntry, error) {
	//TODO bug -> four empty datastore entries

	// IMPORTANT we assume ATM that there is only one contract to deploy

	// number of entries in the datastore: number of contracts, and the contract data, args, and coins
	datastore := []DatastoreSCEntry{} 
	
	// contractsNumberKey := []byte{0}
	contractlength := convert.U64ToBytes(1) // assuming there is one contract to deploy

	//number of contracts to deploy
	numberOfContracts  := NewDeployScDatastoreEntry(
		NewDeployScDatastoreData(
			convert.U64ToBytes(len(convert.U64ToBytes(1))),  // length of the key
			convert.U64ToBytes(1)), // value in bytes
		NewDeployScDatastoreData(
			convert.U64ToBytes(len(contractlength)), 
			convert.U64ToBytes(1),
		))

	_dataStore := append(datastore, numberOfContracts)	

	//byteCode of the smartContract to be appended to the deployer 
	contractData := NewDeployScDatastoreEntry(
			NewDeployScDatastoreData(
				convert.U64ToBytes(len(contractKey(0))),
				contractKey(0),
			),
			NewDeployScDatastoreData(
				convert.U64ToBytes(len(contract.Data)),
				contract.Data,
			),
		)

	_dataStore = append(_dataStore, contractData)


	contractArgs := NewDeployScDatastoreEntry(
			NewDeployScDatastoreData(
				convert.U64ToBytes(len(argsKey(0))),
				argsKey(0),
			),
			NewDeployScDatastoreData(
				convert.U64ToBytes(len(contract.Args)),
				contract.Args,
			),
		)

	_dataStore = append(_dataStore, contractArgs)

	contractCoins := NewDeployScDatastoreEntry(
			NewDeployScDatastoreData(
				convert.U64ToBytes(len(coinsKey(0))),
				coinsKey(0),
			),
			NewDeployScDatastoreData(
				convert.U64ToBytes(len(convert.U64ToBytes(int(contract.Coins)))),
				convert.U64ToBytes(int(contract.Coins)),
			),
		)

	_dataStore = append(_dataStore, contractCoins)

	return _dataStore, nil
}

