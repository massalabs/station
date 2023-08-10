package callsc

import (
	"bytes"
	"encoding/binary"
	"fmt"

	utils "github.com/massalabs/station/pkg/node/sendoperation/serializeaddress"
)

const (
	CallSCOpID        = uint64(4)
	publicKeyHashSize = 32
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
	OperationID uint64
	GasLimit    uint64
	Coins       uint64
	Address     string
	Function    string
	Parameters  []byte
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
	versionedAddress, err := utils.SerializeAddress(address)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare address: %w", err)
	}

	return &CallSC{
		address: versionedAddress, function: function, parameters: parameters,
		gasLimit: gasLimit, coins: coins,
	}, nil
}

func (c *CallSC) Content() (interface{}, error) {
	addressString, err := utils.DeserializeAddress(c.address)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize address: %w", err)
	}

	return &Operation{
		CallSC: OperationDetails{
			MaxGas:     int64(c.gasLimit),
			Coins:      fmt.Sprint(c.coins),
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

//nolint:funlen, cyclop
func DecodeMessage(data []byte) (*MessageContent, error) {
	callSCContent := &MessageContent{}
	buf := bytes.NewReader(data)

	// Read operationId
	callSCOpID, err := binary.ReadUvarint(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to read CallSCOpID: %w", err)
	}

	callSCContent.OperationID = callSCOpID

	// Read maxGas
	gasLimit, err := binary.ReadUvarint(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to read gasLimit: %w", err)
	}

	callSCContent.GasLimit = gasLimit

	// Read Coins
	coins, err := binary.ReadUvarint(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to read coins: %w", err)
	}

	callSCContent.Coins = coins

	// Read address type
	addressType, err := binary.ReadUvarint(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to read address type: %w", err)
	}

	// Read address version
	addressVersion, err := binary.ReadUvarint(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to read address version: %w", err)
	}

	// Read fixed-size 32-byte portion left of the address
	addressBytes := make([]byte, publicKeyHashSize)

	_, err = buf.Read(addressBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to read address portion: %w", err)
	}

	// Concatenate address prefix, address version, and addressBytes to form the full address
	fullAddressBytes := append([]byte{byte(addressType)}, byte(addressVersion))
	fullAddressBytes = append(fullAddressBytes, addressBytes...)

	addressString, err := utils.DeserializeAddress(fullAddressBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize address: %w", err)
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

	parameters := make([]byte, paramLength)

	_, err = buf.Read(parameters)
	if err != nil {
		return nil, fmt.Errorf("failed to read parameters: %w", err)
	}

	callSCContent.Parameters = parameters

	return callSCContent, nil
}
