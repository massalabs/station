package sendoperation

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type MessageContent struct {
	OperationID uint64 `json:"operation_id"`
	RollCount   uint64 `json:"roll_count"`
}

func rollsDecodeMessage(data []byte) (*MessageContent, error) {
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
