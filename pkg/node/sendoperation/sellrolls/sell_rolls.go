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

type MessageContent struct {
	OperationID uint64 `json:"operation_id"`
	RollCount   uint64 `json:"roll_count"`
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

func DecodeMessage(data []byte) (*MessageContent, error) {
	buf := bytes.NewReader(data)

	// Read operationId
	opID, err := binary.ReadUvarint(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to read BuyRollOpID: %w", err)
	}

	// Read count rolls
	countRoll, err := binary.ReadUvarint(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to read countRoll: %w", err)
	}

	operationDetails := &MessageContent{
		OperationID: opID,
		RollCount:   countRoll,
	}

	return operationDetails, nil
}
