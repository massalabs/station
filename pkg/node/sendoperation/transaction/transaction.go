package transaction

import (
	"encoding/binary"
	"fmt"

	"github.com/massalabs/thyra/pkg/node/base58"
	prepareAddress "github.com/massalabs/thyra/pkg/node/sendoperation/prepareaddress"
)

const TransactionOpID = 0

const versionByte = byte(0)

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

func New(recipientAddress string, amount uint64) (*Transaction, error) {
	versionedAddress, err := prepareAddress.PrepareAddress(recipientAddress)
	if err != nil {
		return nil, err
	}

	return &Transaction{
		recipientAddress: versionedAddress,
		amount:           amount,
	}, nil
}

func (t *Transaction) Content() interface{} {
	return &Operation{
		Transaction: OperationDetails{
			RecipientAddress: "AU" + base58.VersionedCheckEncode(t.recipientAddress, versionByte),
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
