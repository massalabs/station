package sendoperation

import (
	"context"
	"crypto/ed25519"
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/massalabs/thyra/pkg/node"
	"github.com/massalabs/thyra/pkg/node/base58"
	"lukechampine.com/blake3"
)

const DefaultGazLimit = 700000000

const DefaultSlotsDuration = 2

const NoGazFee = 0

const NoFee = 0

const NoSequentialCoin = 0

const NoParallelCoin = 0

//nolint:tagliatelle
type sendOperationsReq struct {
	SerializedContent JSONableSlice `json:"serialized_content"`
	PublicKey         string        `json:"creator_public_key"`
	Signature         string        `json:"signature"`
}

type Operation interface {
	Content() interface{}
	Message() []byte
}

type JSONableSlice []uint8

func (u JSONableSlice) MarshalJSON() ([]byte, error) {
	var result string

	if u == nil {
		result = "null"
	} else {
		result = strings.Join(strings.Fields(fmt.Sprintf("%d", u)), ",")
	}

	return []byte(result), nil
}

func message(expiry uint64, fee uint64, operation Operation) []byte {
	msg := make([]byte, 0)
	buf := make([]byte, binary.MaxVarintLen64)
	// fee
	nbBytes := binary.PutUvarint(buf, fee)
	msg = append(msg, buf[:nbBytes]...)

	// expiration
	nbBytes = binary.PutUvarint(buf, expiry)
	msg = append(msg, buf[:nbBytes]...)

	// operation
	msg = append(msg, operation.Message()...)

	return msg
}

func Call(client *node.Client,
	expiry uint64, fee uint64,
	operation Operation,
	pubKey []byte, privKey []byte,
) (string, error) {
	exp, err := node.NextSlot(client)
	if err != nil {
		return "", fmt.Errorf("calling NextSlot: %w", err)
	}

	expiry += exp

	msg := message(expiry, fee, operation)

	digest := blake3.Sum256(append(pubKey, msg...))

	signature := ed25519.Sign(digest[:], privKey)

	rawResponse, err := client.RPCClient.Call(
		context.Background(),
		"send_operations",
		[][]sendOperationsReq{
			{
				sendOperationsReq{
					SerializedContent: msg,
					Signature:         base58.CheckEncode(signature),
					PublicKey:         "P" + base58.VersionedCheckEncode(pubKey, 0),
				},
			},
		},
	)
	if err != nil {
		return "", fmt.Errorf("calling send_operations jsonrpc with '%+v': %w",
			[][]sendOperationsReq{
				{
					sendOperationsReq{
						SerializedContent: msg,
						Signature:         base58.CheckEncode(signature),
						PublicKey:         "P" + base58.VersionedCheckEncode(pubKey, 0),
					},
				},
			},
			err)
	}

	if rawResponse.Error != nil {
		return "", rawResponse.Error
	}

	var resp []string

	err = rawResponse.GetObject(&resp)
	if err != nil {
		return "", fmt.Errorf("parsing send_operations jsonrpc response '%+v': %w", rawResponse, err)
	}

	return resp[0], nil
}
