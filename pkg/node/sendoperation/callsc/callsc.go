package callsc

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"

	utils "github.com/massalabs/station/pkg/node/sendoperation/serializeaddress"
)

const (
	OpType = 4
)

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

type MessageContent struct {
	OperationType uint64
	MaxGas        uint64
	Coins         uint64
	Address       string
	Function      string
	Parameters    []byte
}

type CallSC struct {
	address    []byte
	function   string
	parameters []byte
	maxGas     uint64
	coins      uint64
}

func New(address string, function string, parameters []byte, maxGas uint64, coins uint64,
) (*CallSC, error) {
	versionedAddress, err := utils.SerializeAddress(address)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare address: %w", err)
	}

	return &CallSC{
		address: versionedAddress, function: function, parameters: parameters,
		maxGas: maxGas, coins: coins,
	}, nil
}

func (c *CallSC) Content() (interface{}, error) {
	addressString, err := utils.DeserializeAddress(c.address)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize address: %w", err)
	}

	return &Operation{
		CallSC: OperationDetails{
			MaxGas:     int64(c.maxGas),
			Coins:      strconv.FormatUint(c.coins, 10),
			TargetAddr: addressString,
			TargetFunc: c.function,
			Param:      c.parameters,
		},
	}, nil
}

func (c *CallSC) Message() []byte {
	msg := make([]byte, 0)
	buf := make([]byte, binary.MaxVarintLen64)

	// operationId
	nbBytes := binary.PutUvarint(buf, OpType)
	msg = append(msg, buf[:nbBytes]...)

	// maxGas
	nbBytes = binary.PutUvarint(buf, c.maxGas)
	msg = append(msg, buf[:nbBytes]...)

	// Coins
	nbBytes = binary.PutUvarint(buf, c.coins)
	msg = append(msg, buf[:nbBytes]...)

	// Extract type and version from the address (assuming format: type + version + 32 bytes)
	addressType := c.address[0]
	addressVersion := c.address[1]

	// Encode and append address type and version
	nbBytes = binary.PutUvarint(buf, uint64(addressType))
	msg = append(msg, buf[:nbBytes]...)

	nbBytes = binary.PutUvarint(buf, uint64(addressVersion))
	msg = append(msg, buf[:nbBytes]...)

	// Append the rest of the address (excluding type and version)
	msg = append(msg, c.address[2:]...)

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

//nolint:funlen
func DecodeMessage(data []byte) (*MessageContent, error) {
	callSCContent := &MessageContent{}
	buf := bytes.NewReader(data)

	// Read operation type
	callSCOpType, err := binary.ReadUvarint(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to read operation type: %w", err)
	}

	callSCContent.OperationType = callSCOpType

	// Read maxGas
	maxGas, err := binary.ReadUvarint(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to read maxGas: %w", err)
	}

	callSCContent.MaxGas = maxGas

	// Read Coins
	coins, err := binary.ReadUvarint(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to read coins: %w", err)
	}

	callSCContent.Coins = coins

	// Read recipient address
	addressString, err := utils.DecodeAddress(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to read Target address: %w", err)
	}

	callSCContent.Address = addressString

	// Read target function length
	functionLength, err := binary.ReadUvarint(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to read function length: %w", err)
	}

	functionBytes := make([]byte, functionLength)

	_, err = buf.Read(functionBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to read function: %w", err)
	}

	callSCContent.Function = string(functionBytes)

	// Read param length
	paramLength, err := binary.ReadUvarint(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to read param length: %w", err)
	}

	if paramLength == 0 {
		callSCContent.Parameters = []byte{}
	} else {
		parameters := make([]byte, paramLength)

		_, err = buf.Read(parameters)
		if err != nil {
			return nil, fmt.Errorf("failed to read parameters: %w", err)
		}

		callSCContent.Parameters = parameters
	}

	return callSCContent, nil
}
