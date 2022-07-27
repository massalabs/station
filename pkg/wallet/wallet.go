package wallet

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/massalabs/thyra/pkg/front"
	"github.com/massalabs/thyra/pkg/node/base58"
	"golang.org/x/crypto/pbkdf2"
	"gopkg.in/yaml.v3"
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

type WalletConfig struct {
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

func readWallets() (wallets []Wallet) {
	bytesInput, e := ioutil.ReadFile("wallet.json")
	wallets = []Wallet{}
	if e != nil {
		fmt.Print("No wallet dectected, new one created")
	} else {
		json.Unmarshal(bytesInput, &wallets)
	}
	return wallets

}

func GetWallets(w http.ResponseWriter, r *http.Request) ([]Wallet, error) {
	bytesInput, e := ioutil.ReadFile("wallet.json")

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(bytesInput))

	return readWallets(), e
}

func New(nickname string) (*Wallet, error) {
	privKey, err := btcec.NewPrivateKey()
	if err != nil {
		return nil, err
	}

	pubKeyBytes := privKey.PubKey().X().Bytes()
	privKeyBytes := privKey.Key.Bytes()

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
		Address:  "A" + base58.CheckEncode(append(make([]byte, 1), pubKeyBytes...)),
		KeyPairs: []KeyPair{{
			PrivateKey: privKeyBytes[:],
			PublicKey:  pubKeyBytes,
			Salt:       salt,
			Nonce:      nonce,
		}},
	}

	wallets := readWallets()
	wallets = append(wallets, wallet)

	bytesOutput, err := json.Marshal(wallets)

	os.WriteFile("wallet.json", bytesOutput, 0644)
	return &wallet, nil
}

func Update(wallet Wallet) (err error) {
	wallets := readWallets()

	wallets = append(wallets, wallet)
	bytesOutput, err := json.Marshal(wallets)
	if err != nil {
		panic(err)
	}

	os.WriteFile("wallet.json", bytesOutput, 0644)
	return err
}

func Delete(nickname string) (err error) {
	wallets := readWallets()

	for index, element := range wallets {
		if element.Nickname == nickname {
			wallets = append(wallets[:index], wallets[index+1:]...)
		}
	}

	bytesOutput, err := json.Marshal(wallets)
	os.WriteFile("wallet.json", bytesOutput, 0644)

	return err
}

func HandleWalletManagementRequest(w http.ResponseWriter, r *http.Request) {

	target := r.URL.Path[1:]
	var fileText string
	if strings.Index(target, ".css") > 0 {
		fileText = front.WalletCss
		w.Header().Set("Content-Type", "text/css")
	} else if strings.Index(target, ".js") > 0 {
		fileText = front.WalletJs
		w.Header().Set("Content-Type", "application/json")
	} else if strings.Index(target, ".html") > 0 {
		fileText = front.WalletHtml
		w.Header().Set("Content-Type", "text/html")
	} else if strings.Index(target, ".webp") > 0 {
		fileText = front.Logo_massaWebp
		w.Header().Set("Content-Type", "image/webp")
	}

	w.Write([]byte(fileText))
}
