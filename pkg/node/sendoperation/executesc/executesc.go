package executesc

import (
	"encoding/binary"
	"fmt"
)

const ExecuteSCOpID = 3

type OperationDetails struct {
	Data     []byte `json:"data"`
	MaxGas   uint64 `json:"max_gas"`
	GasPrice string `json:"gas_price"`
	Coins    string `json:"coins"`
}

type Operation struct {
	ExecuteSC OperationDetails `json:"ExecuteSC"`
}

type ExecuteSC struct {
	data     []byte
	maxGas   uint64
	gasPrice uint64
	Coins    uint64
}

func New(data []byte, maxGas uint64, gasPrice uint64, coins uint64) *ExecuteSC {
	return &ExecuteSC{
		data:     data,
		maxGas:   maxGas,
		gasPrice: gasPrice,
		Coins:    coins,
	}
}

func (e *ExecuteSC) Content() interface{} {
	return &Operation{
		ExecuteSC: OperationDetails{
			Data:     e.data,
			MaxGas:   e.maxGas,
			GasPrice: fmt.Sprint(e.gasPrice),
			Coins:    fmt.Sprint(e.Coins),
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

	// Coins
	nbBytes = binary.PutUvarint(buf, e.Coins)
	msg = append(msg, buf[:nbBytes]...)

	// GasPrice
	nbBytes = binary.PutUvarint(buf, e.gasPrice)
	msg = append(msg, buf[:nbBytes]...)

	// data
	nbBytes = binary.PutUvarint(buf, uint64(len(e.data)))
	msg = append(msg, buf[:nbBytes]...)
	msg = append(msg, e.data...)

	return msg
}
