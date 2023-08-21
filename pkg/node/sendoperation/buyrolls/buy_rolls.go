package buyrolls

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const OpID = 1

//nolint:tagliatelle
type OperationDetails struct {
	CountRoll uint64 `json:"roll_count"`
}

//nolint:tagliatelle
type Operation struct {
	BuyRolls OperationDetails `json:"BuyRolls"`
}

type BuyRolls struct {
	countRoll uint64
}

func New(countRolls uint64) *BuyRolls {
	return &BuyRolls{
		countRoll: countRolls,
	}
}

func (b *BuyRolls) Content() (interface{}, error) {
	return &Operation{
		BuyRolls: OperationDetails{
			CountRoll: b.countRoll,
		},
	}, nil
}

func (b *BuyRolls) Message() []byte {
	msg := make([]byte, 0)
	buf := make([]byte, binary.MaxVarintLen64)

	// operationId
	nbBytes := binary.PutUvarint(buf, OpID)
	msg = append(msg, buf[:nbBytes]...)

	// count rolls
	nbBytes = binary.PutUvarint(buf, b.countRoll)

	return append(msg, buf[:nbBytes]...)
}

func DecodeMessage(data []byte) (*OperationDetails, error) {
	buf := bytes.NewReader(data)

	// Read operationId
	_, err := binary.ReadUvarint(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to read OpID: %w", err)
	}

	// Read count rolls
	countRoll, err := binary.ReadUvarint(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to read countRoll: %w", err)
	}

	return &OperationDetails{
		CountRoll: countRoll,
	}, nil
}
