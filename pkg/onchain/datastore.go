package onchain

import (
	"fmt"

	"github.com/massalabs/station/pkg/convert"
)

type DatastoreContract struct {
	Data  []byte
	Args  []byte
	Coins uint64
}

// smartContractNumberKey := []bytes(0)

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

func populateDatastore(contracts []DatastoreContract) (map[string][]byte, error) {
	if len(contracts) == 0 {
		return nil, fmt.Errorf("contracts slice is empty with a length of: %v", len(contracts))
	}

	datastore := make(map[string][]byte)
	contractNumberKey := []byte{0}

	// Set the number of contracts in the first key of the datastore
	datastore[string(convert.ToString(contractNumberKey))] = convert.U64ToBytes(len(contracts))
	// TODO something is bugged here
	for i, contract := range contracts {
		datastore[convert.ToString(contractKey(i))] = contract.Data
		datastore[convert.ToString(argsKey(i))] = contract.Args
		datastore[convert.ToString(coinsKey(i))] = convert.U64ToBytes(int(contract.Coins))
	}

	return datastore, nil
}
