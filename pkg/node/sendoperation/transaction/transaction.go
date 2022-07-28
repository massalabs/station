package transaction

import (
	"encoding/binary"
	"fmt"

	"github.com/massalabs/thyra/pkg/node/base58"
)

type OperationDetails struct {
	Amount           string `json:"amount"`
	RecipientAddress string `json:"recipient_address"`
}

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
			RecipientAddress: "A" + base58.CheckEncode(append(make([]byte, 1), t.recipientAddress...)),
			Amount:           fmt.Sprint(t.amount),
		},
	}
}

func (t *Transaction) Message() []byte {
	msg := make([]byte, 0)
	buf := make([]byte, binary.MaxVarintLen64)

	TransactionOperationID := uint64(0)

	// operationId
	n := binary.PutUvarint(buf, TransactionOperationID)
	msg = append(msg, buf[:n]...)

	// receipient address
	msg = append(msg, t.recipientAddress...)

	// Amount
	n = binary.PutUvarint(buf, t.amount)
	msg = append(msg, buf[:n]...)

	return msg
}
