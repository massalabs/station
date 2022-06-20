package sendoperation

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/massalabs/thyra/pkg/node"
	"lukechampine.com/blake3"

	"github.com/massalabs/thyra/pkg/node/base58"
)

type sendOperationReqContent struct {
	ExpirePeriod int64       `json:"expire_period"`
	Fee          string      `json:"fee"`
	Op           interface{} `json:"op"`
	SenderPK     string      `json:"sender_public_key"`
}

type sendOperationsReq struct {
	Content   sendOperationReqContent `json:"content"`
	Signature string                  `json:"signature"`
}

type Operation interface {
	Content() interface{}
	Message() []byte
}

func message(expiry uint64, fee uint64, senderPubKey []byte, operation Operation) []byte {
	msg := make([]byte, 0)
	buf := make([]byte, binary.MaxVarintLen64)

	//fee
	n := binary.PutUvarint(buf, fee)
	msg = append(msg, buf[:n]...)

	//expiration
	n = binary.PutUvarint(buf, expiry)
	msg = append(msg, buf[:n]...)

	//public key
	msg = append(msg, senderPubKey...)

	//operation
	msg = append(msg, operation.Message()...)

	return msg
}

func sign(msg []byte, privKey []byte) ([]byte, error) {
	digest := blake3.Sum256(msg)

	var auxBytes [32]byte

	if _, err := rand.Read(auxBytes[:]); err != nil {
		return nil, err
	}

	var signOpts = []schnorr.SignOption{schnorr.CustomNonce(auxBytes)}

	pk, _ := btcec.PrivKeyFromBytes(privKey)

	signature, err := schnorr.Sign(pk, digest[:], signOpts...)
	if err != nil {
		return nil, err
	}

	return signature.Serialize(), nil
}

func Call(client *node.Client, expiry uint64, fee uint64, op Operation, pubKey []byte, privKey []byte) (string, error) {
	signature, err := sign(
		message(expiry, fee, pubKey, op),
		privKey,
	)
	if err != nil {
		return "", err
	}

	response, err := client.RPCClient.Call(
		context.Background(),
		"send_operations",
		[][]sendOperationsReq{
			{sendOperationsReq{
				Content: sendOperationReqContent{
					ExpirePeriod: int64(expiry),
					Fee:          fmt.Sprint(fee),
					Op:           op.Content(),
					SenderPK:     base58.CheckEncode(pubKey)},
				Signature: base58.CheckEncode(signature)},
			},
		})
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
