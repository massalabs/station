package transaction

import (
	"encoding/binary"
	"fmt"

	utils "github.com/massalabs/station/pkg/node/sendoperation/serializeaddress"
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

func New(recipientAddress string, amount uint64) (*Transaction, error) {
	versionedAddress, err := utils.SerializeAddress(recipientAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare address: %w", err)
	}

	return &Transaction{
		recipientAddress: versionedAddress,
		amount:           amount,
	}, nil
}

func (t *Transaction) Content() (interface{}, error) {
	addressString, err := utils.DeserializeAddress(t.recipientAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize address: %w", err)
	}

	return &Operation{
		Transaction: OperationDetails{
			RecipientAddress: addressString,
			Amount:           fmt.Sprint(t.amount),
		},
	}, nil
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
