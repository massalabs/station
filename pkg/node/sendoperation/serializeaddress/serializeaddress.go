package serializeaddress

import (
	"fmt"

	"github.com/massalabs/station/pkg/node/base58"
)

// SerializeAddress return the address in byte, ready to be used by the API
// It will add the prefix 1 or 0 depending on if it's a account address or a smart contract address
// It will also base58Check decode the address with version 0.
func SerializeAddress(addr string) ([]byte, error) {
	addressPrefix := addr[:2]
	addressWithoutPrefix := addr[2:]

	// New testnet20 addresses needs a byte 0 for AU addresses and byte 1 for AS addresses
	bytePrefix := []byte{1}
	if addressPrefix == "AU" {
		bytePrefix = []byte{0}
	}

	addressBytes, _, err := base58.VersionedCheckDecode(addressWithoutPrefix)
	if err != nil {
		return nil, fmt.Errorf("decoding address: %w", err)
	}

	return append(bytePrefix, addressBytes...), nil
}
