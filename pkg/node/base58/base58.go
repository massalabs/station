package base58

import (
	"crypto/sha256"
	"errors"

	"github.com/btcsuite/btcutil/base58"
)

func checksum(input []byte) (cksum [4]byte) {
	h := sha256.Sum256(input)
	h2 := sha256.Sum256(h[:])
	copy(cksum[:], h2[:4])

	return
}

const checkHeaderLength = 4

// CheckEncode encodes a byte array to base58 with a checksum but without a version.
func CheckEncode(input []byte) string {
	b := make([]byte, 0, len(input)+checkHeaderLength)
	b = append(b, input...)
	cksum := checksum(b)
	b = append(b, cksum[:]...)

	return base58.Encode(b)
}

// CheckDecode decodes a string that was encoded with CheckEncode and verifies the checksum.
func CheckDecode(input string) (result []byte, err error) {
	decoded := base58.Decode(input)
	if len(decoded) <= checkHeaderLength {
		return nil, errors.New("ErrInvalidFormat")
	}

	var cksum [4]byte

	copy(cksum[:], decoded[len(decoded)-4:])

	if checksum(decoded[:len(decoded)-4]) != cksum {
		return nil, errors.New("ErrChecksum")
	}

	payload := decoded[:len(decoded)-4]
	result = append(result, payload...)

	return
}

// CheckEncode encodes a byte array  and a version to base58 with a checksum.
func VersionedCheckEncode(input []byte, version byte) string {
	return base58.CheckEncode(input, version)
}

// VersionedCheckDecode decodes a string that was encoded with VersionedCheckEncode and verifies the checksum.
//nolint:wrapcheck
func VersionedCheckDecode(input string) (result []byte, version byte, err error) {
	return base58.CheckDecode(input)
}
