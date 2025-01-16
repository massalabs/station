package onchain

import (
	"bytes"
	"encoding/binary"

	"github.com/massalabs/station/pkg/convert"
)

const (
	byteArrayLengthPrefix = 1
	contractByteCodeKey   = 1
)

type ContractDatastore struct {
	ByteCode []byte
	Args     []byte
	Coins    uint64
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

	// number of contracts to deploy
	numberOfContractsKey := []byte{0}
	numberOfContracts := convert.U64ToBytes(1)
	datastore = append(datastore, datastoreEntry{Key: numberOfContractsKey, Value: numberOfContracts})

	// contract bytecode
	datastore = append(datastore, datastoreEntry{Key: getContractByteCodeKey(), Value: contract.ByteCode})

	// args data
	// hardcoded for now, could be dynamix see: https://github.com/massalabs/massa-web3/blob/main/src/dataStore.ts
	datastore = append(datastore, datastoreEntry{Key: getArgsKey(), Value: contract.Args})

	// coins data
	//nolint: lll
	// 12/02/2024 single contract deployement is supported. Multiple not planned. see https://github.com/massalabs/station/issues/1364
	datastore = append(datastore, datastoreEntry{Key: getCoinsKey(), Value: convert.U64ToBytes(contract.Coins)})

	// Serialize the datastore
	serializedDatastore, err := SerializeDatastore(datastore)
	if err != nil {
		return nil, err
	}

	return serializedDatastore, nil
}

// getContractByteCodeKey returns the key for the deployed contract bytecode in the datastore
func getContractByteCodeKey() []byte {
	return convert.U64ToBytes(contractByteCodeKey)
}

// getArgsKey returns the key for the deployed contract constructor parameters in the datastore
func getArgsKey() []byte {
	lengthPrefix := convert.U32ToBytes(byteArrayLengthPrefix)
	argsKeySuffix := []byte{0}
	tempKey := append(getContractByteCodeKey(), lengthPrefix...)
	return append(tempKey, argsKeySuffix...)
}

// getCoinsKey returns the key in the datastore for the amount of MAS to be sent to the deployed contract
func getCoinsKey() []byte {
	lengthPrefix := convert.U32ToBytes(byteArrayLengthPrefix)
	coinsKeySuffix := []byte{1}
	tempKey := append(getContractByteCodeKey(), lengthPrefix...)
	return append(tempKey, coinsKeySuffix...)
}

// DatastoreToDeployedContract If the datastore is a valid datastore for a deployed contract, it will return the contract's bytecode, args and coins
// If the datastore is not a valid datastore for a deployed contract, it will return an empty ContractDatastore and success = false
func DatastoreToDeployedContract(datastore []datastoreEntry) (contractDatastore ContractDatastore, success bool) {
	if !isContractDeployDatastore(datastore) {
		return ContractDatastore{}, false
	}

	var contract ContractDatastore
	contract.ByteCode = datastore[1].Value
	contract.Args = datastore[2].Value

	var err error
	contract.Coins, err = convert.BytesToU64(datastore[3].Value)
	if err != nil {
		return ContractDatastore{}, false
	}

	return contract, true
}

// isContractDeployDatastore checks if the provided datastore entries correspond to a contract deployment.
func isContractDeployDatastore(datastore []datastoreEntry) bool {
	if len(datastore) == 4 &&
		bytes.Equal(datastore[0].Key, []byte{0}) &&
		bytes.Equal(datastore[1].Key, getContractByteCodeKey()) &&
		bytes.Equal(datastore[2].Key, getArgsKey()) &&
		bytes.Equal(datastore[3].Key, getCoinsKey()) {
		return true
	}
	return false
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

func DeSerializeDatastore(datastore []byte) ([]datastoreEntry, error) {
	if len(datastore) == 0 {
		return nil, nil
	}

	var entries []datastoreEntry

	reader := bytes.NewReader(datastore)

	// Decode the number of key-value pairs
	datastoreSize, err := binary.ReadUvarint(reader)
	if err != nil {
		return nil, err
	}

	// Decode each key-value pair
	for i := uint64(0); i < datastoreSize; i++ {
		/* Decode key*/
		// get the key length
		keyLength, err := binary.ReadUvarint(reader)
		if err != nil {
			return nil, err
		}

		// get the key
		key := make([]byte, keyLength)
		_, err = reader.Read(key)
		if err != nil {
			return nil, err
		}

		/* Decode value*/
		// get value length
		valueLength, err := binary.ReadUvarint(reader)
		if err != nil {
			return nil, err
		}

		// get the value
		value := make([]byte, valueLength)
		_, err = reader.Read(value)
		if err != nil {
			return nil, err
		}

		entries = append(entries, datastoreEntry{Key: key, Value: value})
	}

	return entries, nil
}
