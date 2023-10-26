package sendoperation

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type RollMessageContent struct {
	OperationType uint64 `json:"operation_type"`
	RollCount     uint64 `json:"roll_count"`
}

func RollDecodeMessage(data []byte) (*RollMessageContent, error) {
	buf := bytes.NewReader(data)

	// Read operation type
	opType, err := binary.ReadUvarint(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to read operation type: %w", err)
	}

	// Read count rolls
	countRoll, err := binary.ReadUvarint(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to read countRoll: %w", err)
	}

	operationDetails := &RollMessageContent{
		OperationType: opType,
		RollCount:     countRoll,
	}

	return operationDetails, nil
}
