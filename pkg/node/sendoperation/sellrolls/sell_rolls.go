package sellrolls

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const SellRollOpID = 2

//nolint:tagliatelle
type OperationDetails struct {
	CountRoll uint64 `json:"roll_count"`
}

//nolint:tagliatelle
type Operation struct {
	SellRolls OperationDetails `json:"SellRolls"`
}

type SellRolls struct {
	countRoll uint64
}

func New(countRolls uint64) *SellRolls {
	return &SellRolls{
		countRoll: countRolls,
	}
}

func (b *SellRolls) Content() (interface{}, error) {
	return &Operation{
		SellRolls: OperationDetails{
			CountRoll: b.countRoll,
		},
	}, nil
}

func (b *SellRolls) Message() []byte {
	msg := make([]byte, 0)
	buf := make([]byte, binary.MaxVarintLen64)

	// operationId
	nbBytes := binary.PutUvarint(buf, SellRollOpID)
	msg = append(msg, buf[:nbBytes]...)

	// count rolls
	nbBytes = binary.PutUvarint(buf, b.countRoll)
	msg = append(msg, buf[:nbBytes]...)

	return msg
}

func DecodeMessage(data []byte) (*OperationDetails, error) {
	buf := bytes.NewReader(data)

	// Skip operationId
	_, err := binary.ReadUvarint(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to read SellRollOpID: %w", err)
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
