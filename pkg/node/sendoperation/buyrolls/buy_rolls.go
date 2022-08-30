package buyrolls

import (
	"encoding/binary"
)

const BuyRollOpID = 1

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

	// operationId
	nbBytes := binary.PutUvarint(buf, BuyRollOpID)
	msg = append(msg, buf[:nbBytes]...)

	// count rolls
	nbBytes = binary.PutUvarint(buf, b.countRoll)

	return append(msg, buf[:nbBytes]...)
}
