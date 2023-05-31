package sendoperation

import (
	"context"
	b64 "encoding/base64"
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/massalabs/station/pkg/node"
	"github.com/massalabs/station/pkg/node/base58"
	"github.com/massalabs/station/pkg/node/sendoperation/signer"
)

const DefaultGasLimit = 700_000_000

const DefaultSlotsDuration = 2

const NoFee = 0

const NoCoin = 0

const HundredMassa = 100_000_000_000

const OneMassa = 1_000_000_000

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

type OperationResponse struct {
	OperationID   string
	CorrelationID string
}

type OperationBatch struct {
	NewBatch      bool
	CorrelationID string
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

// Call uses the plugin wallet to sign an operation, then send the call to blockchain.
//
//nolint:funlen
func Call(client *node.Client,
	expiry uint64,
	fee uint64,
	operation Operation,
	nickname string,
	operationBatch OperationBatch,
	signer signer.Signer,
) (*OperationResponse, error) {
	msg, msgB64, err := MakeOperation(client, expiry, fee, operation)
	if err != nil {
		return nil, err
	}

	var content string

	switch {
	case operationBatch.NewBatch:
		content = `{
			"operation": "` + msgB64 + `",
			"batch": true
		}`
	case operationBatch.CorrelationID != "":
		content = `{
			"operation": "` + msgB64 + `",
			"correlationId": "` + operationBatch.CorrelationID + `"
		}`
	default:
		content = `{
			"operation": "` + msgB64 + `"
		}`
	}

	res, err := signer.Sign(nickname, []byte(content))
	if err != nil {
		return nil, fmt.Errorf("signing operation: %w", err)
	}

	signature, err := b64.StdEncoding.DecodeString(res.Signature)
	if err != nil {
		return nil, fmt.Errorf("decoding '%s' B64: %w", res.Signature, err)
	}

	resp, err := MakeRPCCall(msg, signature, res.PublicKey, client)
	if err != nil {
		return nil, err
	}

	return &OperationResponse{CorrelationID: res.CorrelationID, OperationID: resp[0]}, nil
}

func MakeRPCCall(msg []byte, signature []byte, publicKey string, client *node.Client) ([]string, error) {
	sendOpParams := [][]sendOperationsReq{
		{
			sendOperationsReq{
				SerializedContent: msg,
				Signature:         base58.CheckEncode(signature),
				PublicKey:         publicKey,
			},
		},
	}

	rawResponse, err := client.RPCClient.Call(
		context.Background(),
		"send_operations",
		sendOpParams,
	)
	if err != nil {
		return nil, fmt.Errorf("calling send_operations jsonrpc with '%+v': %w", sendOpParams, err)
	}

	if rawResponse.Error != nil {
		return nil, fmt.Errorf("receiving  send_operation response: %w", rawResponse.Error)
	}

	var resp []string

	err = rawResponse.GetObject(&resp)
	if err != nil {
		return nil, fmt.Errorf("parsing send_operations jsonrpc response '%+v': %w", rawResponse, err)
	}

	return resp, nil
}

func MakeOperation(client *node.Client, expiry uint64, fee uint64, operation Operation) ([]byte, string, error) {
	exp, err := node.NextSlot(client)
	if err != nil {
		return nil, "", fmt.Errorf("calling NextSlot: %w", err)
	}

	expiry += exp

	msg := message(expiry, fee, operation)

	msgB64 := b64.StdEncoding.EncodeToString(msg)

	return msg, msgB64, nil
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
