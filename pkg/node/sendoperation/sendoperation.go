package sendoperation

import (
	"bytes"
	"context"
	b64 "encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/massalabs/thyra/pkg/node"
	"github.com/massalabs/thyra/pkg/node/base58"
)

const DefaultGazLimit = 700000000

const DefaultSlotsDuration = 2

const NoFee = 0

const NoCoin = 0

const HundredMassa = 100000000000

const OneMassa = 1000000000

const WalletPluginURL = "http://my.massa/thyra/plugin/Massalabs/Massa%20Wallet/rest/wallet/"

const HTTPRequestTimeout = 60 * time.Second

//nolint:tagliatelle
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

type operation interface {
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

func message(expiry uint64, fee uint64, operation operation) []byte {
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

// Call uses the plugin wallet instead of Thyra integrated wallet.
func Call(client *node.Client,
	expiry uint64,
	fee uint64,
	operation operation,
	nickname string,
) (string, error) {
	msg, msgB64, err := makeOperation(client, expiry, fee, operation)
	if err != nil {
		return "", err
	}

	httpRawResponse, err := executeHTTPRequest(http.MethodPost, WalletPluginURL+nickname+"/signOperation",
		bytes.NewBuffer([]byte(`{
		"operation": "`+msgB64+`"
		}`)))
	if err != nil {
		return "", fmt.Errorf("calling executeHTTPRequest: %w", err)
	}

	res := signOperationResponse{"", ""}
	err = json.Unmarshal(httpRawResponse, &res)

	if err != nil {
		return "", fmt.Errorf("unmarshalling '%s' JSON: %w", res, err)
	}

	signature, err := b64.StdEncoding.DecodeString(res.Signature)
	if err != nil {
		return "", fmt.Errorf("decoding '%s' B64: %w", res.Signature, err)
	}

	resp, err := makeRPCCall(msg, signature, res, client)
	if err != nil {
		return "", err
	}

	return resp[0], nil
}

func makeRPCCall(msg []byte, signature []byte, res signOperationResponse, client *node.Client) ([]string, error) {
	sendOpParams := [][]sendOperationsReq{
		{
			sendOperationsReq{
				SerializedContent: msg,
				Signature:         base58.CheckEncode(signature),
				PublicKey:         res.PublicKey,
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

func makeOperation(client *node.Client, expiry uint64, fee uint64, operation operation) ([]byte, string, error) {
	exp, err := node.NextSlot(client)
	if err != nil {
		return nil, "", fmt.Errorf("calling NextSlot: %w", err)
	}

	expiry += exp

	msg := message(expiry, fee, operation)

	msgB64 := b64.StdEncoding.EncodeToString(msg)

	return msg, msgB64, nil
}

func executeHTTPRequest(methodType string, url string, reader io.Reader) ([]byte, error) {
	request, err := http.NewRequestWithContext(
		context.Background(),
		methodType,
		url,
		reader)
	if err != nil {
		return nil, fmt.Errorf("creating HTTP request: %w", err)
	}

	request.Header.Set("Content-Type", "application/json;")

	HTTPClient := &http.Client{Timeout: HTTPRequestTimeout, Transport: nil, Jar: nil, CheckRedirect: nil}

	walletResp, err := HTTPClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("aborting during HTTP request: %w", err)
	}

	if walletResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status: %s", walletResp.Status)
	}

	body, err := io.ReadAll(walletResp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading request body: %w", err)
	}

	defer walletResp.Body.Close()

	return body, nil
}
