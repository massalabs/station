package executesc

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
)

const ExecuteSCOpID = 3

type OperationDetails struct {
	Data   []byte `json:"data"`
	MaxGas uint64 `json:"max_gas"`
	//nolint:tagliatelle
	DataStore map[[3]uint8][]uint8 `json:"datastore"`
}

//nolint:tagliatelle
type Operation struct {
	ExecuteSC OperationDetails `json:"ExecuteSC"`
}

type ExecuteSC struct {
	data      []byte
	maxGas    uint64
	dataStore map[[3]uint8][]uint8
}

/*
The dataStore parameter represents a storage that is accessible by the SC in the constructor
function when it gets deployed. For now it is not used by anyone.
*/
func New(data []byte, maxGas uint64, coins uint64, dataStore map[[3]uint8][]uint8) *ExecuteSC {
	gob.Register(map[[3]uint8]interface{}{})

	return &ExecuteSC{
		data:      data,
		maxGas:    maxGas,
		dataStore: dataStore,
	}
}

func (e *ExecuteSC) Content() interface{} {
	return &Operation{
		ExecuteSC: OperationDetails{
			Data:      e.data,
			MaxGas:    e.maxGas,
			DataStore: e.dataStore,
		},
	}
}

func (e *ExecuteSC) Message() []byte {
	msg := make([]byte, 0)
	buf := make([]byte, binary.MaxVarintLen64)

	// operationId
	nbBytes := binary.PutUvarint(buf, ExecuteSCOpID)
	msg = append(msg, buf[:nbBytes]...)

	// maxGas
	nbBytes = binary.PutUvarint(buf, e.maxGas)
	msg = append(msg, buf[:nbBytes]...)

	// data
	nbBytes = binary.PutUvarint(buf, uint64(len(e.data)))
	msg = append(msg, buf[:nbBytes]...)
	msg = append(msg, e.data...)

	// datastore
	// Number of entries in the datastore
	nbBytes = binary.PutUvarint(buf, uint64(len(e.dataStore)))
	msg = append(msg, buf[:nbBytes]...)

	for key, value := range e.dataStore {
		compactAndAppendBytes(&msg, key)
		compactAndAppendBytes(&msg, value)
	}

	return msg
}

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
