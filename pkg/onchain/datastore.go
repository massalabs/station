package onchain

import (
	"bytes"
	"encoding/binary"

	"github.com/massalabs/station/pkg/convert"
)

type ContractDatastore struct {
	Data  []byte
	Args  []byte
	Coins uint64
}

type datastoreEntry struct {
	Key   []byte
	Value []byte
}

/*
creates and serializes a datastore for the given contract.
*/
func populateDatastore(contract ContractDatastore) ([]byte, error) {
	var datastore []datastoreEntry

	// nmber of contracts to deploy
	numberOfContractsKey := []byte{0}
	numberOfContracts := convert.U64ToBytes(1)
	datastore = append(datastore, datastoreEntry{Key: numberOfContractsKey, Value: numberOfContracts})

	// contract data
	contractKey := convert.U64ToBytes(1)
	datastore = append(datastore, datastoreEntry{Key: contractKey, Value: contract.Data})

	// args data

	// hardcoded for now, could be dynamix see: https://github.com/massalabs/massa-web3/blob/main/src/dataStore.ts
	argsKey := []byte{1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0}
	datastore = append(datastore, datastoreEntry{Key: argsKey, Value: contract.Args})

	// coins data
	//nolint: lll
	// 12/02/2024 single contract deployement is supported. Multiple not planned. see https://github.com/massalabs/station/issues/1364
	coinsKey := []byte{1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1}
	datastore = append(datastore, datastoreEntry{Key: coinsKey, Value: convert.U64ToBytes(contract.Coins)})

	// Serialize the datastore
	serializedDatastore, err := SerializeDatastore(datastore)
	if err != nil {
		return nil, err
	}

	return serializedDatastore, nil
}

// SerializeDatastore serializes the datastore into a []byte array.
func SerializeDatastore(datastore []datastoreEntry) ([]byte, error) {
	var buffer bytes.Buffer

	buf := make([]byte, binary.MaxVarintLen64)
	// Encode the number of key-value pairs
	datastoreSize := uint64(len(datastore))
	uDatastoreSize := binary.PutUvarint(buf, datastoreSize)

	buffer.Write(buf[:uDatastoreSize])

	// Encode each key-value pair
	for _, entry := range datastore {
		// Encode key
		keyLength := uint64(len(entry.Key))
		uKeyLength := binary.PutUvarint(buf, keyLength)
		keyLengthInBytes := buf[:uKeyLength]
		buffer.Write(keyLengthInBytes)
		buffer.Write(entry.Key)

		// Encode value
		valueLength := uint64(len(entry.Value))
		uValueLength := binary.PutUvarint(buf, valueLength)
		valueLengthInBytes := buf[:uValueLength]
		buffer.Write(valueLengthInBytes)
		buffer.Write(entry.Value)
	}

	return buffer.Bytes(), nil
}
