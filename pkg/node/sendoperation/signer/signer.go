package signer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/massalabs/station/int/config"
)

type CustomHeader struct {
	Origin  string
	Referer string
}

type WalletPlugin struct {
	headers CustomHeader
}

type RespError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

const WalletPluginURL = "http://" + config.MassaStationURL + "/plugin/massa-labs/massa-wallet/api/"

const HTTPRequestTimeout = 60 * time.Second

var _ Signer = &WalletPlugin{}

func NewWalletPlugin() *WalletPlugin {
	return &WalletPlugin{headers: CustomHeader{}}
}

func (s *WalletPlugin) SetCustomHeaders(headers CustomHeader) {
	s.headers = headers
}

func (s *WalletPlugin) Sign(nickname string, operation []byte) (*SignOperationResponse, error) {
	httpRawResponse, err := ExecuteHTTPRequest(
		http.MethodPost,
		WalletPluginURL+"accounts/"+nickname+"/sign?allow-fee-edition=true",
		bytes.NewBuffer(operation),
		&s.headers,
	)
	if err != nil {
		res := RespError{"", ""}
		_ = json.Unmarshal(httpRawResponse, &res)

		return nil, fmt.Errorf("calling executeHTTPRequest for call: %w, message: %s", err, res.Message)
	}

	res := SignOperationResponse{"", "", "", ""}

	err = json.Unmarshal(httpRawResponse, &res)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling '%s' JSON: %w", res, err)
	}

	return &res, nil
}

func ExecuteHTTPRequest(methodType string, url string, reader io.Reader, headers *CustomHeader) ([]byte, error) {
	request, err := http.NewRequestWithContext(
		context.Background(),
		methodType,
		url,
		reader)
	if err != nil {
		return nil, fmt.Errorf("creating HTTP request: %w", err)
	}

	request.Header.Set("Content-Type", "application/json;")

	if headers.Origin != "" {
		request.Header.Set("Origin", headers.Origin)
	}

	if headers.Referer != "" {
		request.Header.Set("Referer", headers.Referer)
	}

	HTTPClient := &http.Client{Timeout: HTTPRequestTimeout, Transport: nil, Jar: nil, CheckRedirect: nil}

	resp, err := HTTPClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("aborting during HTTP request: %w", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading request body: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return body, fmt.Errorf("request failed with status: %s", resp.Status)
	}

	return body, nil
}
