package executesc

import (
	"encoding/binary"
	"fmt"
)

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

func New(data []byte, maxGas uint64, gasPrice uint64, Coins uint64) *ExecuteSC {
	return &ExecuteSC{
		data:     data,
		maxGas:   maxGas,
		gasPrice: gasPrice,
		Coins:    Coins,
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

	ExecuteSCOperationID := uint64(3)

	// operationId
	n := binary.PutUvarint(buf, ExecuteSCOperationID)
	msg = append(msg, buf[:n]...)

	// maxGas
	n = binary.PutUvarint(buf, e.maxGas)
	msg = append(msg, buf[:n]...)

	// Coins
	n = binary.PutUvarint(buf, e.Coins)
	msg = append(msg, buf[:n]...)

	// GasPrice
	n = binary.PutUvarint(buf, e.gasPrice)
	msg = append(msg, buf[:n]...)

	// data
	msg = append(msg, e.data...)

	return msg
}
