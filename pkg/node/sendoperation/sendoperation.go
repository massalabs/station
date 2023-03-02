package sendoperation

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"encoding/base64"
	b64 "encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/massalabs/thyra/pkg/node"
	"github.com/massalabs/thyra/pkg/node/base58"
	"lukechampine.com/blake3"
)

const DefaultGazLimit = 700000000

const DefaultSlotsDuration = 2

const SlotDurationBatch = 15

const NoGazFee = 0

const NoFee = 0

const NoCoin = 0

const HundredMassa = 100000000000

const OneMassa = 1000000000

const WalletPluginURL = "http://127.0.0.1:8080/rest/wallet/"

type signOperationResponse struct {
	PublicKey string `json:"publicKey"`
	Signature string `json:"signature"`
}

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

	fmt.Println("OPEARTION CONTENT : ", msg)

	digest := blake3.Sum256(append(pubKey, msg...))

	fmt.Println("DIGEST : ", digest)

	signature := ed25519.Sign(privKey, digest[:])

	fmt.Println("SIGNATURE : ", signature)

	fmt.Println("B58 SIGNATURE : ", base58.CheckEncode(signature))

	fmt.Println("B58 PKEY : ", "P"+base58.VersionedCheckEncode(pubKey, 0))

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

func CallV2(client *node.Client,
	expiry uint64, fee uint64,
	operation Operation,
	nickname string,
) (string, error) {
	exp, err := node.NextSlot(client)
	if err != nil {
		return "", fmt.Errorf("calling NextSlot: %w", err)
	}

	expiry += exp

	msg := message(expiry, fee, operation)

	b64EncodedMsg := b64.StdEncoding.EncodeToString(msg)

	request, err := http.NewRequest("POST", WalletPluginURL+nickname+"/signOperation",
		bytes.NewBuffer([]byte(`{
			"operation": "`+b64EncodedMsg+`"
	}`)))

	if err != nil {
		return "", fmt.Errorf("creating signOperation HTTP request: %w", err)
	}

	request.Header.Set("Content-Type", "application/json;")

	HTTPClient := &http.Client{Timeout: 30 * time.Second}
	walletResp, err := HTTPClient.Do(request)

	if err != nil {
		return "", fmt.Errorf("aborting during HTTP request: %w", err)
	}

	if walletResp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("signOperation failed with HTTP request status: %w", err)
	}

	body, err := ioutil.ReadAll(walletResp.Body)
	if err != nil {
		return "", fmt.Errorf("reading request body: %w", err)
	}

	defer walletResp.Body.Close()

	res := signOperationResponse{}
	err = json.Unmarshal(body, &res)

	if err != nil {
		return "", fmt.Errorf("unmarshalling json: %w", err)
	}

	signature, err := base64.StdEncoding.DecodeString(res.Signature)
	if err != nil {
		return "", fmt.Errorf("decoding b64: %w", err)
	}

	rawResponse, err := client.RPCClient.Call(
		context.Background(),
		"send_operations",
		[][]sendOperationsReq{
			{
				sendOperationsReq{
					SerializedContent: msg,
					Signature:         base58.CheckEncode(signature),
					PublicKey:         res.PublicKey,
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
						PublicKey:         res.PublicKey,
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
