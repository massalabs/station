package send_operation

import (
	"context"
	"encoding/binary"
	"fmt"

	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/decred/dcrd/dcrec/secp256k1/v4/schnorr"
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

type sendOperationsRes struct {
	ID []string
}

type Operation interface {
	Content() interface{}
	Message() []byte
}

func message(expiry uint64, fee uint64, senderPubKey []byte, op Operation) []byte {
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
	msg = append(msg, op.Message()...)

	return msg
}

func sign(msg []byte, privKey []byte) ([]byte, error) {
	digest := blake3.Sum256(msg)

	sign, err := schnorr.Sign(
		secp256k1.PrivKeyFromBytes(privKey),
		digest[:])
	if err != nil {
		return nil, err
	}

	return sign.Serialize(), nil
}

func Call(c *node.Client, expiry uint64, fee uint64, op Operation, pubKey []byte, privKey []byte) (string, error) {

	signature, err := sign(
		message(expiry, fee, pubKey, op),
		privKey,
	)
	if err != nil {
		return "", err
	}

	r, err := c.RPCClient.Call(
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

	if r.Error != nil {
		return "", r.Error
	}

	var res *sendOperationsRes

	err = r.GetObject(&res)
	if err != nil {
		return "", err
	}

	return res.ID[0], nil
}
