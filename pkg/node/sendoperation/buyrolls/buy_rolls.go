package buyrolls

import (
	"encoding/binary"
)

type OperationDetails struct {
	CountRoll uint64 `json:"roll_count"`
}

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

func (b *BuyRolls) Content() interface{} {
	return &Operation{
		BuyRolls: OperationDetails{
			CountRoll: b.countRoll,
		},
	}
}

func (b *BuyRolls) Message() []byte {
	msg := make([]byte, 0)
	buf := make([]byte, binary.MaxVarintLen64)

	buyRollOperationID := uint64(1)

	// operationId
	n := binary.PutUvarint(buf, buyRollOperationID)
	msg = append(msg, buf[:n]...)

	// count rolls
	n = binary.PutUvarint(buf, b.countRoll)
	msg = append(msg, buf[:n]...)
	return msg
}
