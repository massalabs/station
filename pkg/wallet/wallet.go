package wallet

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"

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
	IV         [16]byte `yaml:",flow"`
	Protected  bool
}

type Wallet struct {
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

	mode := cipher.NewCBCEncrypter(block, w.KeyPairs[keyPairIndex].IV[:])
	mode.CryptBlocks(w.KeyPairs[keyPairIndex].PrivateKey, w.KeyPairs[keyPairIndex].PrivateKey)

	w.KeyPairs[keyPairIndex].Protected = true

	return nil
}

func (w *Wallet) Unprotect(password string, keyPairIndex uint8) error {
	secretKey := pbkdf2.Key([]byte(password), w.KeyPairs[keyPairIndex].Salt[:], PBKDF2NbRound, SecretKeyLength, sha256.New)

	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return err
	}

	mode := cipher.NewCBCDecrypter(block, w.KeyPairs[keyPairIndex].IV[:])
	mode.CryptBlocks(w.KeyPairs[0].PrivateKey, w.KeyPairs[keyPairIndex].PrivateKey)

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

	var initializationVector [16]byte

	_, err = rand.Read(initializationVector[:])
	if err != nil {
		return nil, err
	}

	return &Wallet{
		Nickname: nickname,
		Address:  "A" + base58.CheckEncode(append(make([]byte, 1), pubKeyBytes...)),
		KeyPairs: []KeyPair{{
			PrivateKey: privKeyBytes[:],
			PublicKey:  pubKeyBytes,
			Salt:       salt,
			IV:         initializationVector,
		}},
	}, nil
}
