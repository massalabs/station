package transaction

import (
	"encoding/binary"
	"fmt"

	"github.com/massalabs/thyra/pkg/node/base58"
)

const TransactionOpID = 0

type OperationDetails struct {
	Amount           string `json:"amount"`
	RecipientAddress string `json:"recipient_address"`
}

//nolint:tagliatelle
type Operation struct {
	Transaction OperationDetails `json:"Transaction"`
}

type Transaction struct {
	recipientAddress []byte
	amount           uint64
}

func New(recipientAddress []byte, amount uint64) *Transaction {
	return &Transaction{
		recipientAddress: recipientAddress,
		amount:           amount,
	}
}

func (t *Transaction) Content() interface{} {
	return &Operation{
		Transaction: OperationDetails{
			RecipientAddress: "AU" + base58.VersionedCheckEncode(t.recipientAddress, 0),
			Amount:           fmt.Sprint(t.amount),
		},
	}
}

func (t *Transaction) Message() []byte {
	msg := make([]byte, 0)
	buf := make([]byte, binary.MaxVarintLen64)

	// operationId
	nbBytes := binary.PutUvarint(buf, TransactionOpID)
	msg = append(msg, buf[:nbBytes]...)

	// recipient address
	msg = append(msg, t.recipientAddress...)

	// Amount
	nbBytes = binary.PutUvarint(buf, t.amount)
	msg = append(msg, buf[:nbBytes]...)

	return msg
}
