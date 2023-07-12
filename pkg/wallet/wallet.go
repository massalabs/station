package wallet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/massalabs/station/pkg/node/sendoperation"
)

// WalletKeyPair wallet's key pair.
//
//nolint:tagliatelle
type KeyPair struct {
	// Nonce used by the AES-GCM algorithm used to protect the key pair's private key.
	// Required: true
	Nonce string `json:"nonce"`

	// Key pair's private key.
	// Required: true
	PrivateKey string `json:"privateKey"`

	// Key pair's public key.
	// Required: true
	PublicKey string `json:"publicKey"`

	// Salt used by the PBKDF that generates the secret key used to protect the key pair's private key.
	// Required: true
	Salt string `json:"salt"`
}

// Wallet object (V0).
//
//nolint:tagliatelle
type Wallet struct {
	// wallet's address.
	// Required: true
	Address string `json:"address"`

	// key pair
	// Required: true
	KeyPair KeyPair `json:"keyPair"`

	// wallet's nickname.
	// Required: true
	Nickname string `json:"nickname"`
}

func Fetch(nickname string) (*Wallet, error) {
	httpRawResponse, err := sendoperation.ExecuteHTTPRequest(
		http.MethodGet,
		sendoperation.WalletPluginURL+"accounts/"+nickname,
		bytes.NewReader([]byte("")),
	)
	if err != nil {
		res := sendoperation.RespError{Code: "", Message: ""}
		_ = json.Unmarshal(httpRawResponse, &res)

		return nil, fmt.Errorf("calling executeHTTPRequest: %w, message: %s", err, res.Message)
	}

	wallet := Wallet{} //nolint:exhaustruct

	err = json.Unmarshal(httpRawResponse, &wallet)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling file 'wallet_%s.json': %w", nickname, err)
	}

	return &wallet, nil
}
