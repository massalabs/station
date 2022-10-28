package executesc

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
)

const ExecuteSCOpID = 3

type OperationDetails struct {
	Data      []byte                  `json:"data"`
	MaxGas    uint64                  `json:"max_gas"`
	GasPrice  string                  `json:"gas_price"`
	DataStore map[[50]uint8][50]uint8 `json:"datastore"`
}

//nolint:tagliatelle
type Operation struct {
	ExecuteSC OperationDetails `json:"ExecuteSC"`
}

type ExecuteSC struct {
	data      []byte
	maxGas    uint64
	gasPrice  uint64
	dataStore map[[50]uint8][50]uint8
}

func New(data []byte, maxGas uint64, gasPrice uint64, dataStore map[[50]uint8][50]uint8) *ExecuteSC {
	return &ExecuteSC{
		data:      data,
		maxGas:    maxGas,
		gasPrice:  gasPrice,
		dataStore: dataStore,
	}
}

func (e *ExecuteSC) Content() interface{} {
	return &Operation{
		ExecuteSC: OperationDetails{
			Data:      e.data,
			MaxGas:    e.maxGas,
			GasPrice:  fmt.Sprint(e.gasPrice),
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

	// GasPrice
	nbBytes = binary.PutUvarint(buf, e.gasPrice)
	msg = append(msg, buf[:nbBytes]...)

	// data
	nbBytes = binary.PutUvarint(buf, uint64(len(e.data)))
	msg = append(msg, buf[:nbBytes]...)
	msg = append(msg, e.data...)

	// datastore

	bytesBuffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(bytesBuffer)
	err := encoder.Encode(e.dataStore)
	if err != nil {
		fmt.Print("Panicking...\n")
		panic(err)
	}
	binary.PutUvarint(buf, uint64(len(e.data)))
	bytesBufferLength := len(bytesBuffer.Bytes())
	nbBytes = binary.PutUvarint(buf, uint64(bytesBufferLength))
	msg = append(msg, buf[:nbBytes]...)

	fmt.Print("bytesBuffer len : \n")
	fmt.Print(bytesBufferLength)
	fmt.Print("\n")
	msg = append(msg, bytesBuffer.Bytes()...)
	//fmt.Print(msg)
	fmt.Print("Returning msg...\n")
	return msg
}
