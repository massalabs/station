package wallet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/massalabs/station/pkg/node/sendoperation/signer"
)

// AccountKeyPair account's key pair.
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

// Account object (V0).
//
//nolint:tagliatelle
type Account struct {
	// account's address.
	// Required: true
	Address string `json:"address"`

	// key pair
	// Required: true
	KeyPair KeyPair `json:"keyPair"`

	// account's nickname.
	// Required: true
	Nickname string `json:"nickname"`

	// account's balance.
	CandidateBalance string `json:"candidateBalance"`

	// account's status: ok or corrupted.
	Status string `json:"status,omitempty"`
}

func Fetch(nickname string) (*Account, error) {
	httpRawResponse, err := signer.ExecuteHTTPRequest(
		http.MethodGet,
		signer.WalletPluginURL+"accounts/"+nickname,
		bytes.NewReader([]byte("")),
	)
	if err != nil {
		res := signer.RespError{Code: "", Message: ""}
		_ = json.Unmarshal(httpRawResponse, &res)

		return nil, fmt.Errorf("calling executeHTTPRequest: %w, message: %s", err, res.Message)
	}

	account := Account{} //nolint:exhaustruct

	err = json.Unmarshal(httpRawResponse, &account)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling file 'wallet_%s.json': %w", nickname, err)
	}

	return &account, nil
}
