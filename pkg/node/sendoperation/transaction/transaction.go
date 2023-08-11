package transaction

import (
	"encoding/binary"
	"fmt"

	"github.com/massalabs/station/pkg/node/base58"
	serializeAddress "github.com/massalabs/station/pkg/node/sendoperation/serializeaddress"
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
	versionedAddress, err := serializeAddress.SerializeAddress(recipientAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare address: %w", err)
	}

	return &Transaction{
		recipientAddress: versionedAddress,
		amount:           amount,
	}, nil
}

func (t *Transaction) Content() (interface{},error) {
	return &Operation{
		Transaction: OperationDetails{
			RecipientAddress: "AU" + base58.VersionedCheckEncode(t.recipientAddress, versionByte),
			Amount:           fmt.Sprint(t.amount),
		},
	},nil
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
