package onchain

import (
	"fmt"
	"strings"

	"github.com/massalabs/station/pkg/convert"
)

type DatastoreContract struct {
	Data  []byte
	Args  []byte
	Coins uint64
}

type JSONableSlice []byte

type DatastoreSCEntry struct {
	Entry JSONableSlice `json:"entry"`
	Bytes JSONableSlice `json:"bytes"`
}

func (u JSONableSlice) MarshalJSON() ([]byte, error) {
	var result string
	if u == nil {
		result = "null"
	} else {
		result = strings.Join(strings.Fields(fmt.Sprintf("%d", u)), ",")
	}

	return []byte(result), nil
}


func NewDeployScDatastoreEntry(entry []byte, bytes []byte) DatastoreSCEntry {
	return DatastoreSCEntry{
		Entry: entry,
		Bytes: bytes,
	}
}

/**
 * Generates a key for coin data in the datastore.
 *
 * @param offset - The offset to use when generating the key.
 * @returns A Uint8Array representing the key.
 */
func coinsKey(offset int) []byte {
	byteArray := []byte{}
	byteArray = append(byteArray, convert.U64ToBytes(offset+1)...)
	byteArray = append(byteArray, []byte{1}...)
	return byteArray
}

/**
 * Generates a key for args data in the datastore.
 *
 * @param offset - The offset to use when generating the key.
 * @returns A Uint8Array representing the key.
 */
func argsKey(offset int) []byte {
	byteArray := []byte{}
	byteArray = append(byteArray, convert.U64ToBytes(offset+1)...)
	byteArray = append(byteArray, []byte{0}...)
	return byteArray
}

/**
 * Generates a key for contract data in the datastore.
 *
 * @param offset - The offset to use when generating the key.
 * @returns A Uint8Array representing the key.
 */
func contractKey(offset int) []byte {
	byteArray := []byte{}
	return append(byteArray, convert.U64ToBytes(offset+1)...)
}

// TODO implement correct datastore structure

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
func populateDatastore(contracts []DatastoreContract) ([]DatastoreSCEntry, error) {
	if len(contracts) == 0 {
		return nil, fmt.Errorf("contracts slice is empty with a length of: %v", len(contracts))
	}

	// number of entries in the datastore: number of contracts, and the contract data, args, and coins
	datastore := make([]DatastoreSCEntry, 4)
	
	contractsNumberKey := []byte{0}
	/**
	// data store entry
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

	// length of the byte array | byte array data 
	// length of the byte array | byte array data 
	datastore[0] = NewDeployScDatastoreEntry(convert.U64ToBytes(len(contractsNumberKey)), contractsNumberKey)

	for i, contract := range contracts {
		datastore[1] = NewDeployScDatastoreEntry(contractKey(i), contract.Data)
		datastore[2] = NewDeployScDatastoreEntry(argsKey(i), contract.Args)
		datastore[3] = NewDeployScDatastoreEntry(coinsKey(i), convert.U64ToBytes(int(contract.Coins)))
	}

	return datastore, nil
}

