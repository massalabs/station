package transaction

import (
	"encoding/binary"
	"fmt"

	"github.com/massalabs/thyra/pkg/node/base58"
)

type OperationDetails struct {
	Amount           string `json:"amount"`
	RecipientAddress []byte `json:"recipient_address"`
}

type Operation struct {
	Transaction OperationDetails `json:"Transaction"`
}

type Transaction struct {
	recepientAddress string
	amount           uint64
}

func New(recepientAddress string, amount uint64) *Transaction {
	return &Transaction{
		recepientAddress: recepientAddress,
		amount:           amount * 1e9,
	}
}

func (t *Transaction) Content() interface{} {
	addr, _, _ := base58.VersionedCheckDecode(t.recepientAddress[1:])
	return &Operation{
		Transaction: OperationDetails{
			RecipientAddress: addr, //"A" + base58.CheckEncode(append(make([]byte, 1), t.recepientAddress...)),
			Amount:           fmt.Sprint(t.amount * 1e9),
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
	addr, _, _ := base58.VersionedCheckDecode(t.recepientAddress[1:])

	msg = append(msg, addr...)

	// Amount
	n = binary.PutUvarint(buf, t.amount)
	msg = append(msg, buf[:n]...)

	return msg
}
