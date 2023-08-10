package executesc

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
)

const ExecuteSCOpID = 3

type OperationDetails struct {
	Data     []byte `json:"data"`
	MaxGas   uint64 `json:"max_gas"`
	MaxCoins uint64 `json:"max_coins"`
	//nolint:tagliatelle
	DataStore []byte `json:"datastore"`
}

//nolint:tagliatelle
type Operation struct {
	ExecuteSC OperationDetails `json:"ExecuteSC"`
}

type ExecuteSC struct {
	data      []byte
	maxGas    uint64
	maxCoins  uint64
	dataStore []byte
}

/*
The dataStore parameter represents a storage that is accessible by the SC in the constructor
function when it gets deployed.
*/
func New(data []byte, maxGas uint64, maxCoins uint64, dataStore []byte) *ExecuteSC {
	return &ExecuteSC{
		data:      data,
		maxGas:    maxGas,
		maxCoins:  maxCoins,
		dataStore: dataStore,
	}
}

func (e *ExecuteSC) Content() (interface{}, error) {
	return &Operation{
		ExecuteSC: OperationDetails{
			Data:      e.data,
			MaxGas:    e.maxGas,
			MaxCoins:  e.maxCoins,
			DataStore: e.dataStore,
		},
	}, nil
}

// To date the datastore sent by the deploySC endpoint is always serialized.
// However the web on chain features make use of a non-serialized, nil datastore.
// Hence here we check that if datastore is not nil (and it means it comes from the deploySC endpoint)
// we do not encode it further but rather send it as is to the node.
func (e *ExecuteSC) Message() []byte {
	msg := make([]byte, 0)
	buf := make([]byte, binary.MaxVarintLen64)

	// operationId
	nbBytes := binary.PutUvarint(buf, ExecuteSCOpID)
	msg = append(msg, buf[:nbBytes]...)

	// maxGas
	nbBytes = binary.PutUvarint(buf, e.maxGas)
	msg = append(msg, buf[:nbBytes]...)

	nbBytes = binary.PutUvarint(buf, e.maxCoins)
	msg = append(msg, buf[:nbBytes]...)

	// data
	nbBytes = binary.PutUvarint(buf, uint64(len(e.data)))
	msg = append(msg, buf[:nbBytes]...)
	msg = append(msg, e.data...)

	// datastore
	// If the datastore is not nil, no need to serialize it.
	if e.dataStore != nil {
		msg = append(msg, e.dataStore...)

		return msg
	}

	// If the datastore is nil, we need to serialize it.
	// Number of entries in the datastore.
	nbBytes = binary.PutUvarint(buf, uint64(len(e.dataStore)))
	msg = append(msg, buf[:nbBytes]...)
	msg = append(msg, e.dataStore...)

	for key, value := range e.dataStore {
		compactAndAppendBytes(&msg, key)
		compactAndAppendBytes(&msg, value)
	}

	return msg
}

/*
This function serialize the content of the datastore in a byte array and should be used in the following way :

	for key, value := range dataStore {
		compactAndAppendBytes(&byteArray, key)
		compactAndAppendBytes(&byteArray, value)
	}
*/
func compactAndAppendBytes(msg *[]byte, value interface{}) {
	buf := make([]byte, binary.MaxVarintLen64)
	bytesBuffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(bytesBuffer)

	err := encoder.Encode(value)
	if err != nil {
		panic(err)
	}

	nbBytes := binary.PutUvarint(buf, uint64(bytesBuffer.Len()))
	// Value length
	*msg = append(*msg, buf[:nbBytes]...)
	// Value in bytes
	*msg = append(*msg, bytesBuffer.Bytes()...)
}
