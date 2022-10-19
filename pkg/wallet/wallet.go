package wallet

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/massalabs/thyra/pkg/node/base58"
	"golang.org/x/crypto/pbkdf2"
	"gopkg.in/yaml.v3"
	"lukechampine.com/blake3"
)

var ErrUnprotectedSerialization = errors.New("private key must be protected before serialization")
var WalletDirectory = os.Getenv("HOME") + "/.config/thyra/"

const (
	SecretKeyLength           = 32
	PBKDF2NbRound             = 10000
	FileModeUserReadWriteOnly = 0o600
	MinAddressLength          = 49
)

//nolint:tagliatelle
type KeyPair struct {
	PrivateKey []byte   `yaml:",flow"`
	PublicKey  []byte   `yaml:",flow"`
	Salt       [16]byte `yaml:",flow"`
	Nonce      [12]byte `yaml:",flow"`
	Protected  bool
}

//nolint:tagliatelle
type Wallet struct {
	Version  uint8     `json:"version"`
	Nickname string    `json:"nickname"`
	Address  string    `json:"address"`
	KeyPairs []KeyPair `json:"keyPairs"`
}

type Config struct {
	// address
	Wallets []KeyPair `json:"wallets"`
}

func init() {
	if _, err := os.Stat(WalletDirectory); os.IsNotExist(err) {
		os.MkdirAll(WalletDirectory, 0755)
	}
}

func (w *Wallet) Protect(password string, keyPairIndex uint8) error {
	secretKey := pbkdf2.Key([]byte(password), w.KeyPairs[keyPairIndex].Salt[:], PBKDF2NbRound, SecretKeyLength, sha256.New)

	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return fmt.Errorf("intializing block ciphering: %w", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("intializing the AES block cipher wrapped in a Gallois Counter Mode ciphering: %w", err)
	}

	w.KeyPairs[keyPairIndex].PrivateKey = aesgcm.Seal(
		nil,
		w.KeyPairs[keyPairIndex].Nonce[:],
		w.KeyPairs[keyPairIndex].PrivateKey,
		nil)

	w.KeyPairs[keyPairIndex].Protected = true

	return nil
}

func (w *Wallet) Unprotect(password string, keyPairIndex uint8) error {
	secretKey := pbkdf2.Key([]byte(password), w.KeyPairs[keyPairIndex].Salt[:], PBKDF2NbRound, SecretKeyLength, sha256.New)

	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return fmt.Errorf("intializing block ciphering: %w", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("intializing the AES block cipher wrapped in a Gallois Counter Mode ciphering: %w", err)
	}

	pk, err := aesgcm.Open(nil, w.KeyPairs[keyPairIndex].Nonce[:], w.KeyPairs[keyPairIndex].PrivateKey, nil)
	if err != nil {
		return fmt.Errorf("opening the cipher seal: %w", err)
	}

	w.KeyPairs[keyPairIndex].PrivateKey = pk

	w.KeyPairs[keyPairIndex].Protected = false

	return nil
}

//nolint:wrapcheck
func (w *Wallet) YAML() ([]byte, error) {
	for _, v := range w.KeyPairs {
		if !v.Protected {
			return nil, ErrUnprotectedSerialization
		}
	}

	return yaml.Marshal(w)
}

func FromYAML(raw []byte) (w Wallet, err error) {
	err = yaml.Unmarshal(raw, &w)

	return
}

func LoadAll() (wallets []Wallet, e error) {
	wallets = []Wallet{}

	files, err := os.ReadDir(WalletDirectory)
	if err != nil {
		return nil, fmt.Errorf("reading wallet directory '%s': %w", WalletDirectory, err)
	}

	for _, f := range files {
		fileName := f.Name()

		if strings.HasPrefix(fileName, "wallet_") && strings.HasSuffix(fileName, ".json") {
			bytesInput, err := os.ReadFile(WalletDirectory + fileName)
			if err != nil {
				return nil, fmt.Errorf("reading file '%s': %w", fileName, err)
			}

			wallet := Wallet{} //nolint:exhaustruct

			err = json.Unmarshal(bytesInput, &wallet)
			if err != nil {
				return nil, fmt.Errorf("unmarshaling file '%s': %w", fileName, err)
			}

			wallets = append(wallets, wallet)
		}
	}

	return wallets, nil
}

func Load(nickname string) (*Wallet, error) {
	bytesInput, err := os.ReadFile(WalletDirectory + "wallet_" + nickname + ".json")
	if err != nil {
		return nil, fmt.Errorf("reading file 'wallet_%s.json': %w", nickname, err)
	}

	wallet := Wallet{} //nolint:exhaustruct

	err = json.Unmarshal(bytesInput, &wallet)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling file 'wallet_%s.json': %w", nickname, err)
	}

	return &wallet, nil
}

func New(nickname string) (*Wallet, error) {
	pubKey, privKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		return nil, fmt.Errorf("generating ed25519 keypair: %w", err)
	}

	addr := blake3.Sum256(pubKey)

	var salt [16]byte

	_, err = rand.Read(salt[:])
	if err != nil {
		return nil, fmt.Errorf("generating random salt: %w", err)
	}

	var nonce [12]byte

	_, err = rand.Read(nonce[:])
	if err != nil {
		return nil, fmt.Errorf("generating random nonce: %w", err)
	}

	wallet := Wallet{
		Version:  0,
		Nickname: nickname,
		Address:  "A" + base58.CheckEncode(append(make([]byte, 1), addr[:]...)),
		KeyPairs: []KeyPair{{
			PrivateKey: privKey,
			PublicKey:  pubKey,
			Salt:       salt,
			Nonce:      nonce,
		}},
	}

	bytesOutput, err := json.Marshal(wallet)
	if err != nil {
		return nil, fmt.Errorf("marshalling wallet: %w", err)
	}

	err = os.WriteFile(WalletDirectory+"wallet_"+nickname+".json", bytesOutput, FileModeUserReadWriteOnly)
	if err != nil {
		return nil, fmt.Errorf("writing wallet to 'wallet_%s.json': %w", nickname, err)
	}

	return &wallet, nil
}

func Delete(nickname string) (err error) {
	err = os.Remove(WalletDirectory + "wallet_" + nickname + ".json")
	if err != nil {
		return fmt.Errorf("deleting wallet 'wallet_%s.json': %w", nickname, err)
	}

	return nil
}

func AddressChecker(address string) bool {
	return len(address) > MinAddressLength
}
