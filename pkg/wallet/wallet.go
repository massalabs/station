package wallet

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"net/http"
	"strings"

	"github.com/btcsuite/btcd/btcec/v2"
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
	Version  uint8
	Nickname string
	Address  string
	KeyPairs []KeyPair
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

	return &Wallet{
		Version:  0,
		Nickname: nickname,
		Address:  "A" + base58.CheckEncode(append(make([]byte, 1), pubKeyBytes...)),
		KeyPairs: []KeyPair{{
			PrivateKey: privKeyBytes[:],
			PublicKey:  pubKeyBytes,
			Salt:       salt,
			Nonce:      nonce,
		}},
	}, nil
}

func HandleWalletManagementRequest(w http.ResponseWriter, r *http.Request) {

	target := r.URL.Path[1:]
	var fileText string
	if strings.Index(target, ".css") > 0 {
		fileText = IndexCss
		w.Header().Set("Content-Type", "text/css")
	} else if strings.Index(target, ".js") > 0 {
		fileText = IndexJs
		w.Header().Set("Content-Type", "application/json")
	} else if strings.Index(target, ".html") > 0 {
		fileText = IndexHtml
		w.Header().Set("Content-Type", "text/html")
	} else if strings.Index(target, ".webp") > 0 {
		fileText = Logo_massaWebp
		w.Header().Set("Content-Type", "image/webp")
	}

	w.Write([]byte(fileText))
}
