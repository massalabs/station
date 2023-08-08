package sendoperation

import (
	"bytes"
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

type CallSCFromSign struct {
	OperationID uint64
	GasLimit    uint64
	Coins       uint64
	Address     string
	Function    string
	Parameters  []byte
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
	/////////////////
	// Just For test
	// To Delete

	params, err := GetCallSCFromSign(msgB64)
	if err != nil {
		fmt.Println("🚀 ~ file: sendoperation.go:93 ~ err:", err)
	} else {
		if params != nil {
			// Print the fields of CallSCFromSign
			fmt.Println("🚀 ~ file: sendoperation.go:90 ~ Print the fields of CallSCFromSign:")
			fmt.Println("🚀 ~ file: sendoperation.go:92 ~ callSC.GasLimit:", params.GasLimit)
			fmt.Println("🚀 ~ file: sendoperation.go:94 ~ callSC.Coins:", params.Coins)
			fmt.Println("🚀 ~ file: sendoperation.go:96 ~ callSC.Address:", params.Address)
			fmt.Println("🚀 ~ file: sendoperation.go:98 ~ callSC.Function:", params.Function)
		} else {
			fmt.Println("🚀 ~ file: sendoperation.go:98 ~ params is nil")
		}
	}
	/////////////////
	//////////////////////
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

func DecodeMessage64(msgB64 string) ([]byte, error) {
	// Decode the base64-encoded message
	decodedMsg, err := b64.StdEncoding.DecodeString(msgB64)
	if err != nil {
		return nil, fmt.Errorf("base64 decoding error: %w", err)
	}

	// Read the encoded fee from the decoded message and move the buffer index
	_, bytesRead := binary.Uvarint(decodedMsg)
	if bytesRead <= 0 {
		return nil, fmt.Errorf("failed to read fee")
	}
	decodedMsg = decodedMsg[bytesRead:]

	// Read the encoded expiry from the decoded message and move the buffer index
	_, bytesRead = binary.Uvarint(decodedMsg)
	if bytesRead <= 0 {
		return nil, fmt.Errorf("failed to read expiry")
	}
	decodedMsg = decodedMsg[bytesRead:]

	// At this point, 'decodedMsg' contains the remaining part of the message
	// which is the 'operation.Message()' part
	return decodedMsg, nil
}

func DecodeCallSCMessage(data []byte) (*CallSCFromSign, error) {
	c := &CallSCFromSign{}
	buf := bytes.NewReader(data)

	// Read operationId
	callSCOpID, err := binary.ReadUvarint(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to read CallSCOpID: %w", err)
	}
	c.OperationID = callSCOpID

	// Read maxGas
	gasLimit, err := binary.ReadUvarint(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to read gasLimit: %w", err)
	}
	c.GasLimit = gasLimit

	// Read Coins
	coins, err := binary.ReadUvarint(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to read coins: %w", err)
	}
	c.Coins = coins

	// Read target address length
	addressLength, err := binary.ReadUvarint(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to read address length: %w", err)
	}
	address := make([]byte, addressLength)
	_, err = buf.Read(address)
	if err != nil {
		return nil, fmt.Errorf("failed to read address: %w", err)
	}
	const versionByte = byte(1)
	c.Address = "AS" + base58.VersionedCheckEncode(address, versionByte)

	// Read target function length
	functionLength, err := binary.ReadUvarint(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to read function length: %w", err)
	}
	functionBytes := make([]byte, functionLength)
	_, err = buf.Read(functionBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to read function: %w", err)
	}
	c.Function = string(functionBytes)

	// Read param length
	paramLength, err := binary.ReadUvarint(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to read param length: %w", err)
	}
	parameters := make([]byte, paramLength)
	_, err = buf.Read(parameters)
	if err != nil {
		return nil, fmt.Errorf("failed to read parameters: %w", err)
	}
	c.Parameters = parameters

	return c, nil
}

func GetCallSCFromSign(msgB64 string) (*CallSCFromSign, error) {
	decodedMsg, err := DecodeMessage64(msgB64)
	if err != nil {
		fmt.Println("Error decoding message:", err)
		return nil, err
	}

	callSC, err := DecodeCallSCMessage(decodedMsg)
	if err != nil {
		fmt.Println("Error decoding CallSC:", err)
		return nil, err
	}

	return callSC, nil
}
