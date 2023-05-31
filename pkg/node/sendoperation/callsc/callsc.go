package callsc

import (
	"encoding/binary"
	"fmt"

	"github.com/massalabs/thyra/pkg/node/base58"
	serializeAddress "github.com/massalabs/thyra/pkg/node/sendoperation/serializeaddress"
)

const CallSCOpID = uint64(4)

const versionByte = byte(1)

type OperationDetails struct {
	MaxGas     int64       `json:"max_gas"`
	Coins      string      `json:"coins"`
	TargetAddr string      `json:"target_addr"`
	TargetFunc string      `json:"target_func"`
	Param      interface{} `json:"param"`
}

//nolint:tagliatelle
type Operation struct {
	CallSC OperationDetails `json:"CallSC"`
}

type CallSC struct {
	address    []byte
	function   string
	parameters []byte
	gasLimit   uint64
	coins      uint64
}

func New(address string, function string, parameters []byte, gasLimit uint64, coins uint64,
) (*CallSC, error) {
	versionedAddress, err := serializeAddress.SerializeAddress(address)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare address: %w", err)
	}

	return &CallSC{
		address: versionedAddress, function: function, parameters: parameters,
		gasLimit: gasLimit, coins: coins,
	}, nil
}

func (c *CallSC) Content() interface{} {
	return &Operation{
		CallSC: OperationDetails{
			MaxGas:     int64(c.gasLimit),
			Coins:      fmt.Sprint(c.coins),
			TargetAddr: "AS" + base58.VersionedCheckEncode(c.address, versionByte),
			TargetFunc: c.function,
			Param:      c.parameters,
		},
	}
}

func (c *CallSC) Message() []byte {
	msg := make([]byte, 0)
	buf := make([]byte, binary.MaxVarintLen64)

	// operationId
	nbBytes := binary.PutUvarint(buf, CallSCOpID)
	msg = append(msg, buf[:nbBytes]...)

	// maxGas
	nbBytes = binary.PutUvarint(buf, c.gasLimit)
	msg = append(msg, buf[:nbBytes]...)

	// Coins
	nbBytes = binary.PutUvarint(buf, c.coins)
	msg = append(msg, buf[:nbBytes]...)

	// target address
	msg = append(msg, c.address...)

	// target function
	nbBytes = binary.PutUvarint(buf, uint64(len([]byte(c.function))))
	msg = append(msg, buf[:nbBytes]...)
	msg = append(msg, []byte(c.function)...)

	// param
	nbBytes = binary.PutUvarint(buf, uint64(len(c.parameters)))
	msg = append(msg, buf[:nbBytes]...)
	msg = append(msg, c.parameters...)

	return msg
}
