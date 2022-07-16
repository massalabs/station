package wallet

import (
	"crypto/ed25519"
	"errors"
	"fmt"

	"github.com/massalabs/thyra/pkg/node/base58"
	"lukechampine.com/blake3"
)

const VERSION byte = 0x0

type Wallet struct {
	PrivateKey ed25519.PrivateKey
	Address    string
}

func New() (*Wallet, error) {

	pubKey, privKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		return nil, errors.New("unable to create a new wallet")
	}
	//generate the wallet address from public key
	address := getAddressFromPubKey(pubKey)
	w := &Wallet{
		PrivateKey: privKey,
		Address:    address,
	}
	return w, nil
}

//NewFromSeed expects base58 encoded secreat key as per massa and returns the wallet
func NewFromSeed(key string) (*Wallet, error) {
	seed, _, err := base58.VersionedCheckDecode(key[1:])
	if err != nil {
		return nil, err
	}

	privKey := ed25519.NewKeyFromSeed(seed)
	pubKey := privKey.Public().(ed25519.PublicKey)

	//generate the wallet address from public key
	address := getAddressFromPubKey(pubKey)
	w := &Wallet{
		PrivateKey: privKey,
		Address:    address,
	}
	return w, nil
}

func (w *Wallet) GetPrivateKey() []byte {
	return w.PrivateKey
}
func (w *Wallet) GetEncodedPrivateKey() string {
	return "S" + base58.VersionedCheckEncode(w.PrivateKey.Seed(), VERSION)
}
func (w *Wallet) GetPublicKey() []byte {
	return w.PrivateKey.Public().(ed25519.PublicKey)
}
func (w *Wallet) GetEncodedPublicKey() string {
	pub := w.PrivateKey.Public().(ed25519.PublicKey)
	return "P" + base58.VersionedCheckEncode(pub, VERSION)
}

func (w *Wallet) Sign(data []byte) []byte {
	hash := blake3.Sum256(data)
	return ed25519.Sign(w.PrivateKey, hash[:])
}

func (w *Wallet) Print() {
	fmt.Printf("secret: %v\n", w.GetEncodedPrivateKey())
	fmt.Printf("public key: %v\n", w.GetEncodedPublicKey())
	fmt.Printf("address: %v\n", w.Address)
}

func getAddressFromPubKey(pubKey ed25519.PublicKey) string {
	hash := blake3.Sum256(pubKey)
	return "A" + base58.VersionedCheckEncode(hash[:], VERSION)
}
