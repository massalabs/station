package sellrolls

import (
	"encoding/binary"
)

type OperationDetails struct {
	CountRoll uint64 `json:"roll_count"`
}

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

func (b *SellRolls) Content() interface{} {
	return &Operation{
		SellRolls: OperationDetails{
			CountRoll: b.countRoll,
		},
	}
}

func (b *SellRolls) Message() []byte {
	msg := make([]byte, 0)
	buf := make([]byte, binary.MaxVarintLen64)

	sellRollOperationID := uint64(2)

	// operationId
	n := binary.PutUvarint(buf, sellRollOperationID)
	msg = append(msg, buf[:n]...)

	// count rolls
	n = binary.PutUvarint(buf, b.countRoll)
	msg = append(msg, buf[:n]...)
	return msg
}
