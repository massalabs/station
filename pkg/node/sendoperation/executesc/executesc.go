package executesc

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
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

// MessageContent stores essential fields extracted from the message during the sign operation.
type MessageContent struct {
	MaxGas      uint64
	MaxCoins    uint64
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

// DecodeMessage decodes a byte slice,
// It extracts the necessary fields: operationID, maxGas, and MaxCoin for display in the Wails pop-up.
func DecodeMessage(data []byte) (*MessageContent, error) {
	operationContent := &MessageContent{}
	buf := bytes.NewReader(data)

	// Skip the  operationId
	_, err := binary.ReadUvarint(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to read ExecuteSCOpID: %w", err)
	}
	// Read maxGas
	maxGas, err := binary.ReadUvarint(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to read maxGas: %w", err)
	}

	operationContent.MaxGas = maxGas

	// Read maxCoins
	maxCoins, err := binary.ReadUvarint(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to read maxCoins: %w", err)
	}

	operationContent.MaxCoins = maxCoins

	return operationContent, nil
}
