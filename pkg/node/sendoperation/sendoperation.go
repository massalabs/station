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
	fmt.Print("here")
	// fee
	n := binary.PutUvarint(buf, fee)
	msg = append(msg, buf[:n]...)

	// expiration
	n = binary.PutUvarint(buf, expiry)
	msg = append(msg, buf[:n]...)

	// operation
	msg = append(msg, operation.Message()...)

	return msg
}

func sign(msg []byte, privKey []byte) ([]byte, error) {
	signature := ed25519.Sign(privKey, msg)
	return signature, nil
}

func Call(client *node.Client, expiry uint64, fee uint64, op Operation, pubKey []byte, privKey []byte) (string, error) {
	msg := message(expiry, fee, op)

	digest := blake3.Sum256(append(pubKey, msg...))

	signature, err := sign(
		digest[:],
		privKey,
	)
	if err != nil {
		return "", err
	}

	response, err := client.RPCClient.Call(
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
		return "", err
	}

	if response.Error != nil {
		return "", response.Error
	}

	var id []string

	err = response.GetObject(&id)
	if err != nil {
		return "", err
	}

	return id[0], nil
}
