package transaction

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/massalabs/station/pkg/node/base58"
	serializeAddress "github.com/massalabs/station/pkg/node/sendoperation/serializeaddress"

	utils "github.com/massalabs/station/pkg/node/sendoperation/serializeaddress"
)

const (
	TransactionOpID   = uint64(0)
	versionByte       = byte(0)
	publicKeyHashSize = 32
)

type OperationDetails struct {
	Amount           string `json:"amount"`
	RecipientAddress string `json:"recipient_address"`
}

//nolint:tagliatelle
type Operation struct {
	Transaction OperationDetails `json:"Transaction"`
}

// MessageContent stores essential fields extracted from the message during the sign operation.
type MessageContent struct {
	OperationID      uint64
	RecipientAddress string
	Amount           uint64
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

func (t *Transaction) Content() (interface{}, error) {
	return &Operation{
		Transaction: OperationDetails{
			RecipientAddress: "AU" + base58.VersionedCheckEncode(t.recipientAddress, versionByte),
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

// DecodeMessage decodes a byte slice for a transaction,
// It extracts the necessary fields: operationID, recipient address, and amount.
func DecodeMessage(data []byte) (*MessageContent, error) {
	transactionContent := &MessageContent{}
	buf := bytes.NewReader(data)

	// Read operationId
	opID, err := binary.ReadUvarint(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to read TransactionOpID: %w", err)
	}
	transactionContent.OperationID = opID

	// Read recipient address
	addressType, err := binary.ReadUvarint(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to read address type: %w", err)
	}

	addressVersion, err := binary.ReadUvarint(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to read address version: %w", err)
	}

	// Read fixed-size 32-byte portion left of the address
	addressBytes := make([]byte, publicKeyHashSize)
	_, err = buf.Read(addressBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to read address portion: %w", err)
	}

	// Concatenate address type, version, and addressBytes to form the full address
	fullAddressBytes := append([]byte{byte(addressType)}, byte(addressVersion))
	fullAddressBytes = append(fullAddressBytes, addressBytes...)

	addressString, err := utils.DeserializeAddress(fullAddressBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize address: %w", err)
	}

	transactionContent.RecipientAddress = addressString

	// Read amount
	amount, err := binary.ReadUvarint(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to read amount: %w", err)
	}

	transactionContent.Amount = amount

	return transactionContent, nil
}
