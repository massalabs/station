package callsc

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"

	"github.com/massalabs/thyra/pkg/node/base58"
)

type OperationDetails struct {
	MaxGaz          int64       `json:"max_gas"`
	GazPrice        string      `json:"gas_price"`
	ParallelCoins   string      `json:"parallel_coins"`
	SequentialCoins string      `json:"sequential_coins"`
	TargetAddr      string      `json:"target_addr"`
	TargetFunc      string      `json:"target_func"`
	Param           interface{} `json:"param"`
}

type Operation struct {
	CallSC OperationDetails `json:"CallSC"`
}

type CallSC struct {
	address           []byte
	function          string
	parameters        []byte
	gazPrice          uint64
	gazLimit          uint64
	nbSequentialCoins uint64
	nbParallelCoins   uint64
}

func New(address []byte, function string, parameters []byte, gazPrice uint64, gazLimit uint64, nbSequentialCoins uint64,
	nbParallelCoins uint64) *CallSC {
	return &CallSC{
		address: address, function: function, parameters: parameters,
		gazPrice: gazPrice, gazLimit: gazLimit, nbSequentialCoins: nbSequentialCoins,
		nbParallelCoins: nbParallelCoins,
	}
}

func (c *CallSC) Content() interface{} {
	return &Operation{
		CallSC: OperationDetails{
			MaxGaz:          int64(c.gazLimit),
			GazPrice:        fmt.Sprint(c.gazPrice),
			ParallelCoins:   fmt.Sprint(c.nbParallelCoins),
			SequentialCoins: fmt.Sprint(c.nbSequentialCoins),
			TargetAddr:      "A" + base58.CheckEncode(append(make([]byte, 1), c.address...)),
			TargetFunc:      c.function,
			Param:           hex.EncodeToString(c.parameters),
		},
	}
}

func (c *CallSC) Message() []byte {
	msg := make([]byte, 0)
	buf := make([]byte, binary.MaxVarintLen64)

	//operationId
	n := binary.PutUvarint(buf, 4)
	msg = append(msg, buf[:n]...)

	//maxGaz
	n = binary.PutUvarint(buf, c.gazLimit)
	msg = append(msg, buf[:n]...)

	//ParallelCoins
	n = binary.PutUvarint(buf, c.nbParallelCoins)
	msg = append(msg, buf[:n]...)

	//SequentialCoins
	n = binary.PutUvarint(buf, c.nbSequentialCoins)
	msg = append(msg, buf[:n]...)

	//gazPrice
	n = binary.PutUvarint(buf, c.gazPrice)
	msg = append(msg, buf[:n]...)

	//target address
	msg = append(msg, c.address...)

	//target function
	n = binary.PutUvarint(buf, uint64(len([]byte(c.function))))
	msg = append(msg, buf[:n]...)
	msg = append(msg, []byte(c.function)...)

	//param
	n = binary.PutUvarint(buf, uint64(len(c.parameters)))
	msg = append(msg, buf[:n]...)
	msg = append(msg, c.parameters...)

	return msg
}
