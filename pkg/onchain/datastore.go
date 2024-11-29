package onchain

import (
	"bytes"
	"encoding/binary"

	"github.com/massalabs/station/pkg/convert"
)

type DatastoreContract struct {
	Data  []byte
	Args  []byte
	Coins uint64
}

type datastoreEntry struct {
	Key   []byte
	Value []byte
}

func argsKey(offset int) ([]byte, error) {
	// Create a buffer to serialize data
	var buf bytes.Buffer

	// Add U64 (64-bit unsigned integer)
	err := binary.Write(&buf, binary.LittleEndian, uint64(offset+1))
	if err != nil {
		return nil, err
	}

	// Add Uint8Array equivalent (assuming a single byte array with a 0 value)
	buf.Write([]byte{0})

	// Return the serialized data
	return buf.Bytes(), nil
}

func coinsKey(offset int) ([]byte, error) {
	// Create a buffer to serialize data
	var buf bytes.Buffer

	// Add U64 (64-bit unsigned integer)
	err := binary.Write(&buf, binary.LittleEndian, uint64(offset+1))
	if err != nil {
		return nil, err
	}

	// Add Uint8Array equivalent (assuming a single byte array with a 0 value)
	buf.Write([]byte{1})

	// Return the serialized data
	return buf.Bytes(), nil
}

func argsValue(args []byte) []byte {
	
	if len(args) == 0 {
		return convert.StrToBytes("")
	}

	return args
}

/** 
populateDatastore creates and serializes a datastore for the given contract.
*/
func populateDatastore(contract DatastoreContract) ([]byte, error) {
	var datastore []datastoreEntry

	// nomber of contracts to deploy
	numberOfContractsKey := []byte{0}
	numberOfContracts := convert.U64ToBytes(1) 
	datastore = append(datastore, datastoreEntry{Key: numberOfContractsKey, Value: numberOfContracts})

	// contract data
	contractKey := convert.U64ToBytes(1) 
	datastore = append(datastore, datastoreEntry{Key: contractKey, Value: contract.Data})

	// args data
	argsKey,err:= argsKey(1)  // TODO verify this key values
	if err != nil {
		return nil, err
	}
	datastore = append(datastore, datastoreEntry{Key: argsKey, Value: contract.Args})

	// coins data
	coinsKey,err := coinsKey(1)// TODO verify this key values
	if err != nil {
		return nil, err
	}
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
