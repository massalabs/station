package executesc

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
)

const OpType = 3

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
	OperationType uint64
	MaxGas        uint64
	MaxCoins      uint64
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

func (e *ExecuteSC) Message() []byte {
	msg := make([]byte, 0)
	buf := make([]byte, binary.MaxVarintLen64)

	// operationId
	nbBytes := binary.PutUvarint(buf, OpType)
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

	msg = append(msg, e.dataStore...)

	return msg
}

/*
This function serialize the content of the datastore in a byte array and should be used in the following way :

	for key, value := range dataStore {
		compactAndAppendBytes(&byteArray, key)
		compactAndAppendBytes(&byteArray, value)
	}
*/
//nolint:unused
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

	// Read operation type
	operationType, err := binary.ReadUvarint(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to read operation type: %w", err)
	}

	operationContent.OperationType = operationType

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
