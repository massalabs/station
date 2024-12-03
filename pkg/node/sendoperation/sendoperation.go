package sendoperation

import (
	"bytes"
	"context"
	b64 "encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/massalabs/station/pkg/convert"
	"github.com/massalabs/station/pkg/node"
	"github.com/massalabs/station/pkg/node/base58"
	"github.com/massalabs/station/pkg/node/sendoperation/signer"
)

const (
	MaxGasAllowedExecuteSC     = 3_980_167_295
	MaxGasAllowedCallSC        = 4_294_167_295
	DefaultExpiryInSlot        = 3
	DefaultFee                 = 0
	accountCreationStorageCost = 1_000_000
	StorageCostPerByte         = 100_000
	StorageEntryBaseBytes      = 4
	OneMassa                   = 1_000_000_000
)

//nolint:tagliatelle
type sendOperationsReq struct {
	SerializedContent JSONableSlice `json:"serialized_content"`
	PublicKey         string        `json:"creator_public_key"`
	Signature         string        `json:"signature"`
}

type Operation interface {
	Content() (interface{}, error)
	Message() []byte
}

type OperationResponse struct {
	OperationID   string
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
func Call(
	client *node.Client,
	chainID uint64,
	expiry uint64,
	fee uint64,
	operation Operation,
	nickname string,
	signer signer.Signer,
	description string,
) (*OperationResponse, error) {
	msg, msgB64, err := MakeOperation(client, expiry, fee, operation)
	if err != nil {
		return nil, err
	}

	content := createOperationContent(description, msgB64, chainID)

	res, err := signer.Sign(nickname, []byte(content))
	if err != nil {
		return nil, fmt.Errorf("signing operation: %w", err)
	}

	signature, err := b64.StdEncoding.DecodeString(res.Signature)
	if err != nil {
		return nil, fmt.Errorf("decoding '%s' B64: %w", res.Signature, err)
	}

	if res.Operation != "" {
		msg, err = b64.StdEncoding.DecodeString(res.Operation)
		if err != nil {
			return nil, fmt.Errorf("decoding '%s' B64: %w", res.Operation, err)
		}
	}

	resp, err := MakeRPCCall(msg, signature, res.PublicKey, client)
	if err != nil {
		return nil, err
	}

	return &OperationResponse{CorrelationID: res.CorrelationID, OperationID: resp[0]}, nil
}

func createOperationContent(description string, msgB64 string, chainID uint64) string {
	data := map[string]interface{}{
		"description": description,
		"operation":   msgB64,
		"chainId":     strconv.FormatUint(chainID, 10),
	}

	content, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}

	return convert.ToString(content)
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
		return nil, fmt.Errorf("receiving send_operations response: %w", rawResponse.Error)
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

// message constructs a message byte slice with the provided expiry, fee, and operation.
// It returns the composed message.
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

// DecodeMessage64 decodes a base64-encoded message and extracts the fee, expiry,
// and the actual message content(operation : CallSc, ExecuteSC, BuyRoll, SellRoll).
// It returns the decoded message, fee, expiry, and an error if any.
func DecodeMessage64(msgB64 string) ([]byte, uint64, uint64, error) {
	// Decode the base64-encoded message
	decodedMsg, err := b64.StdEncoding.DecodeString(msgB64)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("base64 decoding error: %w", err)
	}

	// Read the encoded fee from the decoded message and move the buffer index
	fee, bytesRead := binary.Uvarint(decodedMsg)
	if bytesRead <= 0 {
		return nil, 0, 0, fmt.Errorf("failed to read fee")
	}

	decodedMsg = decodedMsg[bytesRead:]

	// Read the encoded expiry from the decoded message and move the buffer index
	expiry, bytesRead := binary.Uvarint(decodedMsg)
	if bytesRead <= 0 {
		return nil, 0, 0, fmt.Errorf("failed to read expiry")
	}

	decodedMsg = decodedMsg[bytesRead:]

	return decodedMsg, fee, expiry, nil
}

// DecodeOperationType decodes a byte slice to retrieve the operation ID.
func DecodeOperationType(data []byte) (uint64, error) {
	buf := bytes.NewReader(data)

	// Read operation type
	opType, err := binary.ReadUvarint(buf)
	if err != nil {
		return 0, fmt.Errorf("failed to read operation type: %w", err)
	}

	return opType, nil
}

func StorageCostForEntry(keyByteLength, valueByteLength int) (int, error) {
	return (valueByteLength + keyByteLength + StorageEntryBaseBytes) * StorageCostPerByte, nil
}

func AccountCreationStorageCost() (int, error) {
	return accountCreationStorageCost, nil
}
