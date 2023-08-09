package serializeaddress

import (
	"errors"
	"fmt"

	utils "github.com/btcsuite/btcutil/base58"
	"github.com/massalabs/station/pkg/node/base58"
)

// SerializeAddress return the address in byte, ready to be used by the API
// It will add the prefix 1 or 0 depending on if it's a account address or a smart contract address
// It will also base58Check decode the address with version 0.
func SerializeAddress(addr string) ([]byte, error) {
	addressPrefix := addr[:2]
	addressWithoutPrefix := addr[2:]

	// New testnet20 addresses needs a byte 0 for AU addresses and byte 1 for AS addresses
	result := []byte{1}
	if addressPrefix == "AU" {
		result = []byte{0}
	}

	addressBytes, version, err := base58.VersionedCheckDecode(addressWithoutPrefix)
	if err != nil {
		return nil, fmt.Errorf("decoding address: %w", err)
	}

	// New testnet23 addresses needs a version byte
	result = append(result, version)

	return append(result, addressBytes...), nil
}

func DeserializeAddress(versionedAddress []byte) (string, error) {
	const minAddressLength = 2

	if len(versionedAddress) < minAddressLength {
		return "", errors.New("invalid versioned address length")
	}

	prefixByte := versionedAddress[0]
	addressVersion := versionedAddress[1]
	addressBytes := versionedAddress[2:]

	var addressPrefix string

	switch prefixByte {
	case 0:
		addressPrefix = "AU"
	case 1:
		addressPrefix = "AS"
	default:
		return "", errors.New("unknown address prefix")
	}

	addressWithoutPrefix := utils.CheckEncode(addressBytes, addressVersion)

	return addressPrefix + addressWithoutPrefix, nil
}
