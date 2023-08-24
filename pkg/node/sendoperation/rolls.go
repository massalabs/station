package sendoperation

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type RollMessageContent struct {
	OperationID uint64 `json:"operation_id"`
	RollCount   uint64 `json:"roll_count"`
}

func RollDecodeMessage(data []byte) (*RollMessageContent, error) {
	buf := bytes.NewReader(data)

	// Read operationId
	opID, err := binary.ReadUvarint(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to read Operation ID: %w", err)
	}

	// Read count rolls
	countRoll, err := binary.ReadUvarint(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to read countRoll: %w", err)
	}

	operationDetails := &RollMessageContent{
		OperationID: opID,
		RollCount:   countRoll,
	}

	return operationDetails, nil
}
