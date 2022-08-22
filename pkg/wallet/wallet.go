package wallet

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"strings"

	"github.com/massalabs/thyra/pkg/node/base58"
	"golang.org/x/crypto/pbkdf2"
	"gopkg.in/yaml.v3"
	"lukechampine.com/blake3"
)

var ErrUnprotectedSerialization = errors.New("private key must be protected before serialization")

const (
	SecretKeyLength = 32
	PBKDF2NbRound   = 10000
)

type KeyPair struct {
	PrivateKey []byte   `yaml:",flow"`
	PublicKey  []byte   `yaml:",flow"`
	Salt       [16]byte `yaml:",flow"`
	Nonce      [12]byte `yaml:",flow"`
	Protected  bool
}

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

func (w *Wallet) Protect(password string, keyPairIndex uint8) error {
	secretKey := pbkdf2.Key([]byte(password), w.KeyPairs[keyPairIndex].Salt[:], PBKDF2NbRound, SecretKeyLength, sha256.New)

	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
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
		return err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	pk, err := aesgcm.Open(nil, w.KeyPairs[keyPairIndex].Nonce[:], w.KeyPairs[keyPairIndex].PrivateKey, nil)
	if err != nil {
		return err
	}

	w.KeyPairs[keyPairIndex].PrivateKey = pk

	w.KeyPairs[keyPairIndex].Protected = false

	return nil
}

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

func ReadWallets() ([]Wallet, error) {

	wd, err := os.Getwd()

	if err != nil {
		return nil, err
	}

	wallets := []Wallet{}
	files, err := ioutil.ReadDir(wd)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		fileName := f.Name()

		if strings.HasPrefix(fileName, "wallet_") && strings.HasSuffix(fileName, ".json") {
			bytesInput, err := ioutil.ReadFile(fileName)
			if err != nil {
				return nil, err
			}
			wallet := Wallet{}
			err = json.Unmarshal(bytesInput, &wallet)
			if err != nil {
				return nil, err
			}
			wallets = append(wallets, wallet)
		}
	}

	return wallets, nil
}

func GetWallet(nickname string) (*Wallet, error) {

	bytesInput, err := ioutil.ReadFile("wallet_" + nickname + ".json")
	if err != nil {
		return nil, err
	}
	wallet := Wallet{}
	err = json.Unmarshal(bytesInput, &wallet)
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

func New(nickname string) (*Wallet, error) {
	pubKey, privKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		return nil, err
	}

	addr := blake3.Sum256(pubKey)

	var salt [16]byte

	_, err = rand.Read(salt[:])
	if err != nil {
		return nil, err
	}

	var nonce [12]byte

	_, err = rand.Read(nonce[:])
	if err != nil {
		return nil, err
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
		return nil, err
	}

	err = os.WriteFile("wallet_"+nickname+".json", bytesOutput, 0o644)
	if err != nil {
		return nil, err
	}

	return &wallet, nil
}

func Update(wallet Wallet) (err error) {
	wallets, err := ReadWallets()
	if err != nil {
		return err
	}

	wallets = append(wallets, wallet)

	bytesOutput, err := json.Marshal(wallets)
	if err != nil {
		return err
	}

	err = os.WriteFile("wallet.json", bytesOutput, 0o644)
	if err != nil {
		return err
	}

	return nil
}

func Delete(nickname string) (err error) {
	wallets, err := ReadWallets()
	if err != nil {
		return err
	}

	for index, element := range wallets {
		if element.Nickname == nickname {
			wallets = append(wallets[:index], wallets[index+1:]...)
		}
	}

	bytesOutput, err := json.Marshal(wallets)
	if err != nil {
		return err
	}

	err = os.WriteFile("wallet.json", bytesOutput, 0o644)
	if err != nil {
		return err
	}

	return nil
}
